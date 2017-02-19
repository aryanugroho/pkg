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

package problem

import (
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/corestoreio/errors"
)

//go:generate easyjson -snake_case -omit_empty detail.go

const (
	// MediaType specifies the default media type for a Detail response
	MediaType = "application/problem+json"

	// MediaTypeXML specifies the XML variant on the Detail Media type
	MediaTypeXML = "application/problem+xml"

	// DefaultURL is the default url to use for problem types
	DefaultURL = "about:blank"
)

// Detail is a default problem implementation. Problem details are not a
// debugging tool for the underlying implementation; rather, they are a way to
// expose greater detail about the HTTP interface itself.  Designers of new
// problem types need to carefully consider the Security Considerations, in
// particular, the risk of exposing attack vectors by exposing implementation
// internals through error messages.
//easyjson:json
type Detail struct {
	// Type - A URI reference [RFC3986] that identifies the problem type.  This
	// specification encourages that, when dereferenced, it provide
	// human-readable documentation for the problem type (e.g., using HTML
	// [W3C.REC-html5-20141028]).  When this member is not present, its value is
	// assumed to be "about:blank".
	Type string `json:"type" xml:"type"`

	// Title - A short, human-readable summary of the problem type.  It SHOULD
	// NOT change from occurrence to occurrence of the problem, except for
	// purposes of localization.
	Title string `json:"title" xml:"title"`

	// Status specifies the HTTP status code generated by the origin server for
	// this occurrence of the problem.
	Status int `json:"status,omitempty" xml:"status,omitempty"`

	// Detail states a human-readable explanation specific to this occurrence of the
	// problem.
	Detail string `json:"detail,omitempty" xml:"detail,omitempty"`

	// Instance states an URI that identifies the specific occurrence of the
	// problem. This URI may or may not yield further information if
	// dereferenced.
	Instance string `json:"instance,omitempty" xml:"instance,omitempty"`

	// Cause can reference to the underlying real cause of a problem detail.
	Cause *Detail `json:"cause,omitempty" xml:"cause,omitempty"`

	// Extension defines additional custom key/value pairs in a balanced slice.
	// i=key and i+1=value. They will be transformed into a key/value JSON type.
	Extension []string
}

// NewDetail returns a new instance of a Detail problem.
func NewDetail(title string, opts ...Option) (*Detail, error) {
	d := &Detail{Type: DefaultURL, Title: title}
	if err := d.Options(opts...); err != nil {
		return nil, errors.Wrap(err, "[problem] NewDetail failed")
	}
	if err := d.Validate(); err != nil {
		return nil, errors.Wrap(err, "[problem] NewDetail validation")
	}
	return d, nil
}

// MustNewDetail same as NewDetail but panics on error.
func MustNewDetail(title string, opts ...Option) *Detail {
	d, err := NewDetail(title, opts...)
	if err != nil {
		panic(err)
	}
	return d
}

// Options applies options to the details object.
func (d *Detail) Options(opts ...Option) error {
	for i, o := range opts {
		if err := o(d); err != nil {
			return errors.Wrapf(err, "[problem] Failed to apply option in index %d", i)
		}
	}
	return nil
}

// Validate returns an error if the title is empty, or type is not an URI or
// Extension field is imbalanced.
func (d *Detail) Validate() error {
	switch {
	case d.Title == "":
		return errors.NewEmptyf("[problem] Title cannot be empty")
	case !isURL(d.Type): // might be wrong validation
		return errors.NewNotValidf("[problem] Title cannot be empty")
	case len(d.Extension)%2 == 1:
		return errors.NewNotValidf("[problem] While creating a new detail the extensions are imbalanced: %v", d.Extension)
	}
	return nil
}

func isURL(str string) bool {
	const maxURLRuneCount = 2083
	const minURLRuneCount = 3

	if str == "" || utf8.RuneCountInString(str) >= maxURLRuneCount || len(str) <= minURLRuneCount || strings.HasPrefix(str, ".") {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	if strings.HasPrefix(u.Host, ".") {
		return false
	}
	if u.Host == "" && (u.Path != "" && !strings.Contains(u.Path, ".")) {
		return false
	}
	return true
}

// Option applies an option to the Detail object
type Option func(*Detail) error

// WithCause adds a detail as a new cause.
func WithCause(title string, opts ...Option) Option {
	return func(d *Detail) error {
		c, err := NewDetail(title, opts...)
		if err != nil {
			return errors.Wrap(err, "[problem] WithCause")
		}
		d.Cause = c
		return nil
	}
}

// WithExtensionString adds one or more key/value pairs to the Extension field.
func WithExtensionString(keyValues ...string) Option {
	return func(d *Detail) error {
		if len(keyValues)%2 == 1 {
			return errors.NewNotValidf("[problem] Imbalanced: %v", keyValues)
		}
		d.Extension = append(d.Extension, keyValues...)
		return nil
	}
}

// WithExtensionInt adds an int to the Extension field.
func WithExtensionInt(key string, value int) Option {
	return func(d *Detail) error {
		d.Extension = append(d.Extension, key, strconv.Itoa(value))
		return nil
	}
}

// WithExtensionUint adds an uint to the Extension field.
func WithExtensionUint(key string, value uint) Option {
	return func(d *Detail) error {
		d.Extension = append(d.Extension, key, strconv.FormatUint(uint64(value), 10))
		return nil
	}
}
