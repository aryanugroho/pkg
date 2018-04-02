// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"sync"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
)

// MessageReceiver allows you to listen to write actions. The order of calling
// each subscriber is totally random. If a subscriber panics, it gets securely
// removed without crashing the whole system. This interface should be
// implemented in other packages. The Subscriber interface requires the
// MessageReceiver interface.
type MessageReceiver interface {
	// MessageConfig when a configuration value will be written this function
	// gets called to allow you to listen to changes. Path is never empty. Path
	// may contains up to three levels. For more details see the Subscriber
	// interface of this package. If an error will be returned, the subscriber
	// gets unsubscribed/removed.
	MessageConfig(Path) error
}

// Subscriber represents the overall service to receive subscriptions from
// MessageReceiver interfaces. This interface is at the moment only implemented
// by the config.Service.
type Subscriber interface {
	// Subscribe subscribes a MessageReceiver to a path. Path allows you to
	// filter to which path or part of a path you would like to listen. A path
	// can be e.g. "system/smtp/host" to receive messages by single host changes
	// or "system/smtp" to receive message from all smtp changes or "system" to
	// receive changes for all paths beginning with "system". A path is equal to
	// a topic in a PubSub system. Path cannot be empty means you cannot listen
	// to all changes. Returns a unique identifier for the Subscriber for later
	// removal, or an error.
	Subscribe(Path, MessageReceiver) (subscriptionID int, err error)
}

// pubSub embedded pointer struct into the Service
type pubSub struct {
	// subMap, subscribed writers are getting called when a write event
	// will happen. uint64 is the path/route (aka topic) and int the Subscriber ID for later
	// removal.
	subMap     map[uint32]map[int]MessageReceiver
	subAutoInc int // subAutoInc increased whenever a Subscriber has been added
	mu         sync.RWMutex
	pubPath    chan Path
	stop       chan struct{} // terminates the goroutine
	closeErr   chan error    // this one tells us that the go routine has really been terminated
	closed     bool          // if Close() has been called the config.Service can still Write() without panic
	log        log.Logger
}

// Close closes the internal channel for the pubsub Goroutine. Prevents a leaking
// Goroutine.
func (s *pubSub) Close() error {
	if s == nil {
		return nil
	}
	if s.closed {
		return errors.AlreadyClosed.Newf("[config] PubSub Service already closed")
	}
	defer func() { close(s.closeErr) }() // last close(s.closeErr) does not work and panics
	s.closed = true
	s.stop <- struct{}{}
	close(s.pubPath)
	close(s.stop)
	//close(s.closeErr)
	return <-s.closeErr
}

// Subscribe adds a Subscriber to be called when a write event happens. See
// interface Subscriber for a detailed description. Route can be any kind of
// level and can contain StrScope and Scope ID. Valid routes can be for example:
//		- StrScope/ID/currency/options/base
//		- StrScope/ID/currency/options
//		- StrScope/ID/currency
//		- currency/options/base
//		- currency/options
//		- currency
func (s *pubSub) Subscribe(p Path, mr MessageReceiver) (subscriptionID int, err error) {
	if p.IsEmpty() {
		return 0, errors.Empty.Newf("[config] pubSub.Subscribe %q", p)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subAutoInc++
	subscriptionID = s.subAutoInc

	hashPath := p.Hash32()

	if _, ok := s.subMap[hashPath]; !ok {
		s.subMap[hashPath] = make(map[int]MessageReceiver)
	}
	s.subMap[hashPath][subscriptionID] = mr

	return
}

// Unsubscribe removes a subscriber with a specific ID.
func (s *pubSub) Unsubscribe(subscriptionID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for path, subs := range s.subMap {
		if _, ok := subs[subscriptionID]; ok {
			delete(s.subMap[path], subscriptionID) // mem leaks?
			if len(s.subMap[path]) == 0 {
				delete(s.subMap, path)
			}
			return nil
		}
	}
	return nil
}

// sendMsg sends the arg into the channel
func (s *pubSub) sendMsg(p Path) {
	if false == s.closed {
		s.pubPath <- p
	}
}

// publish runs in a Goroutine and listens on the channel publishArg. Every time
// a message is coming in, it calls all subscribers. We must run asynchronously
// because we don't know how long each subscriber needs.
func (s *pubSub) publish() {

	// TODO: review this API and the running goroutine: http://www.jtolds.com/writing/2016/03/go-channels-are-bad-and-you-should-feel-bad/
	// but for know it works, someone can refactor it :-)

	for {
		select {
		case <-s.stop:
			s.closeErr <- nil
			return
		case p, ok := <-s.pubPath:
			if !ok {
				// channel closed
				return
			}

			if len(s.subMap) == 0 {
				break
			}

			var evict []int

			evict = append(evict, s.readMapAndSend(p, 1)...)  // e.g.: system and StrScope/ID/system
			evict = append(evict, s.readMapAndSend(p, 2)...)  // e.g.: system/smtp and StrScope/ID/system/smtp
			evict = append(evict, s.readMapAndSend(p, -1)...) // e.g.: system/smtp/host/... and StrScope/ID/system/smtp/host/...

			// remove all failed Subscribers
			if len(evict) > 0 {
				for _, e := range evict {
					if err := s.Unsubscribe(e); err != nil && s.log.IsDebug() {
						s.log.Debug("config.pubSub.publish.evict.Unsubscribe.err", log.Err(err), log.Int("subscriptionID", e))
					}
				}
			}
		}
	}
}

func (s *pubSub) readMapAndSend(p Path, level int) (evict []int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	h, err := p.Hash(level) // including scope and scopeID and the route
	if err != nil && s.log.IsDebug() {
		s.log.Debug("config.pubSub.publish.PathHash.err", log.Err(err), log.Stringer("path", p))
	}
	if subs, ok := s.subMap[h]; ok { // e.g.: strScope/ID/system/smtp/host/etc/pp
		evict = append(evict, s.sendMsgs(subs, p)...)
	}

	h, err = p.Hash(level) // without scope and scopeID and route only
	if err != nil && s.log.IsDebug() {
		s.log.Debug("config.pubSub.publish.RouteHash.err", log.Err(err), log.Stringer("path", p))
	}
	if subs, ok := s.subMap[h]; ok { // e.g.: system/smtp/host/etc/pp
		evict = append(evict, s.sendMsgs(subs, p)...)
	}

	return
}

func (s *pubSub) sendMsgs(subs map[int]MessageReceiver, p Path) (evict []int) {
	for id, sub := range subs {
		if err := s.sendMsgRecoverable(id, sub, p); err != nil {
			if s.log.IsDebug() {
				s.log.Debug("config.pubSub.publish.sendMessages", log.Err(err), log.Int("id", id), log.Stringer("path", p))
			}
			evict = append(evict, id) // mark Subscribers for removal which failed ...
		}
	}
	return
}

func (s *pubSub) sendMsgRecoverable(id int, sl MessageReceiver, p Path) (err error) {
	defer func() { // protect ... you'll never know
		if r := recover(); r != nil {
			if recErr, ok := r.(error); ok {
				s.log.Debug("config.pubSub.publish.recover.err", log.Err(recErr), log.Stringer("path", p))
				err = recErr
			} else {
				s.log.Debug("config.pubSub.publish.recover.r", log.Object("recover", r), log.Stringer("path", p))
				err = errors.Errorf("%#v", r)
			}
			// the overall trick here is, that defer will assign a new error to err
			// and therefore will overwrite the returned nil value!
		}
	}()
	err = sl.MessageConfig(p)
	return
}

func newPubSub(l log.Logger) *pubSub {
	return &pubSub{
		subMap:   make(map[uint32]map[int]MessageReceiver),
		pubPath:  make(chan Path),
		stop:     make(chan struct{}),
		closeErr: make(chan error),
		log:      l,
	}
}
