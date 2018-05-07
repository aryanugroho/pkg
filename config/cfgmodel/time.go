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

package cfgmodel

import (
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/corestoreio/pkg/util/conv"
)

// Time represents a path in config.Getter which handles time values.
type Time struct{ baseValue }

// NewTime creates a new Time cfgmodel with a given path.
func NewTime(path string, opts ...Option) Time {
	return Time{baseValue: newBaseValue(path, opts...)}
}

// Get returns a time value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
// Get is able to parse available time formats as defined in
// github.com/corestoreio/pkg/util/conv.StringToDate()
func (t Time) Value(sg config.Scoped) (time.Time, error) {
	// This code must be kept in sync with other Value() functions

	var v time.Time
	var scp = t.initScope().Top()
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		if d := t.Field.Default; d != nil {
			var err error
			v, err = conv.ToTimeE(d)
			if err != nil {
				return time.Time{}, errors.NotValid.Newf("[cfgmodel] ToTimeE: %v", err)
			}
		}
	}

	val, err := sg.Time(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a time value without validating it against the cfgsource.Slice.
func (t Time) Write(w config.Setter, v time.Time, h scope.TypeID) error {
	return t.baseValue.Write(w, v, h)
}

// Duration represents a path in config.Getter which handles duration values.
type Duration struct{ baseValue }

// NewDuration creates a new Duration cfgmodel with a given path.
func NewDuration(path string, opts ...Option) Duration {
	return Duration{baseValue: newBaseValue(path, opts...)}
}

// Get returns a duration value from ScopedGetter, if empty the
// *Field.Default value will be applied if provided.
// scope.DefaultID will be enforced if *Field.Scopes is empty.
// A duration string is a possibly signed sequence of
// decimal numbers, each with optional fraction and a unit suffix,
// such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// Error behaviour: NotValid
func (t Duration) Value(sg config.Scoped) (time.Duration, error) {
	// This code must be kept in sync with other Value() functions

	var v time.Duration
	var scp = t.initScope().Top()
	if t.Field != nil {
		scp = t.Field.Scopes.Top()
		if d := t.Field.Default; d != nil {
			var err error
			v, err = conv.ToDurationE(d)
			if err != nil {
				return 0, errors.NotValid.Newf("[cfgmodel] ToDurationE: %v", err)
			}
		}
	}

	val, err := sg.Duration(t.route, scp)
	switch {
	case err == nil: // we found the value in the config service
		v = val
	case !errors.NotFound.Match(err):
		err = errors.Wrapf(err, "[cfgmodel] Route %q", t.route)
	default:
		err = nil // a Err(Section|Group|Field)NotFound error and uninteresting, so reset
	}
	return v, err
}

// Write writes a duration value without validating it against the cfgsource.Slice.
func (t Duration) Write(w config.Setter, v time.Duration, h scope.TypeID) error {
	return t.baseValue.Write(w, v.String(), h)
}
