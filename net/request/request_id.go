// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package request

// crypto/rand => http://blog.sgmansfield.com/2016/06/managing-syscall-overhead-with-crypto-rand/

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/corestoreio/csfw/log"
	"github.com/corestoreio/csfw/net/mw"
)

// RequestIDHeader defines the name of the header used to transmit the request ID.
const RequestIDHeader = "X-Request-Id"

// reqID is a global Counter used to create new request ids. This ID is not unique
// across multiple micro services.
var reqID = new(int64)

// IDGenerator defines the functions needed to generate a request
// prefix id.
type IDGenerator interface {
	// Init allows you to initialize a prefix which will be appended to
	// the NewID() function. Init is only called once.
	Init()
	// NewID returns an atomic ID. This function gets executed for every
	// request.
	NewID(*http.Request) string
}

// idService default prefix generator. Creates a prefix once the middleware
// is set up.
type idService struct {
	prefix string
}

// Prefix returns a unique prefix string for the current (micro) service.
// This id gets reset once you restart the service.
func (rp *idService) Init() {
	// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L40-L52
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		if _, err := rand.Read(buf[:]); err != nil {
			panic(err) // todo remove panic without giving up error reporting
		}
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	rp.prefix = fmt.Sprintf("%s/%s-", hostname, b64[0:10])
}

// NewID returns a new ID unique for the current compilation.
func (rp *idService) NewID(_ *http.Request) string {
	return rp.prefix + strconv.FormatInt(atomic.AddInt64(reqID, 1), 10)
}

// ID represents a middleware for request Id generation.
type ID struct {
	IDGenerator
	log.Logger
}

// With is a middleware that injects a request ID into the response header
// of each request. Retrieve it using:
// 		w.Header().Get(RequestIDHeader)
// If the incoming request has a RequestIDHeader header then that value is used
// otherwise a random value is generated. You can specify your own generator by
// providing the RequestPrefixGenerator in an option. No options uses the
// default request prefix generator.
// Supported options are: SetLogger() and SetRequestIDGenerator()
//
// Package store/storenet provides also a request ID generator containing
// the current requested store.
func (iw ID) With() mw.Middleware {
	if iw.Logger == nil {
		iw.Logger = log.BlackHole{}
	}
	if iw.IDGenerator == nil {
		iw.IDGenerator = &idService{}
	}
	iw.Init()

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(RequestIDHeader)
			if id == "" {
				id = iw.NewID(r)
			}
			if iw.IsDebug() {
				iw.Debug("request.ID.With", log.String("id", id), log.HTTPRequest("request", r))
			}
			w.Header().Set(RequestIDHeader, id)
			h.ServeHTTP(w, r)
		})
	}
}
