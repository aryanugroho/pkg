// Copyright 2015 CoreStore Authors
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
	"strings"
	"testing"

	"github.com/corestoreio/csfw/utils/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.Set(log.NewStdLogger())
	log.SetLevel(log.StdLevelError)
}

func TestScopeKey(t *testing.T) {
	tests := []struct {
		haveArg []ScopeOption
		want    string
	}{
		{[]ScopeOption{Path("a/b/c")}, StringScopeDefault + "/0/a/b/c"},
		{[]ScopeOption{Path("")}, ""},
		{[]ScopeOption{Path()}, ""},
		{[]ScopeOption{Scope(IDScopeDefault, nil)}, ""},
		{[]ScopeOption{Scope(IDScopeWebsite, nil)}, ""},
		{[]ScopeOption{Scope(IDScopeStore, nil)}, ""},
		{[]ScopeOption{Path("a/b/c"), Scope(IDScopeWebsite, nil)}, StringScopeDefault + "/0/a/b/c"},
		{[]ScopeOption{Path("a/b/c"), Scope(IDScopeWebsite, scopeID(2))}, StringScopeWebsites + "/2/a/b/c"},
		{[]ScopeOption{Path("a", "b", "c"), Scope(IDScopeWebsite, scopeID(200))}, StringScopeWebsites + "/200/a/b/c"},
		{[]ScopeOption{Path("a", "b", "c"), Scope(IDScopeStore, scopeID(4))}, StringScopeStores + "/4/a/b/c"},
		{[]ScopeOption{Path("a", "b"), Scope(IDScopeStore, scopeID(4))}, StringScopeStores + "/4/a"},
		{[]ScopeOption{nil, Scope(IDScopeStore, scopeID(4))}, ""},
		{[]ScopeOption{Path("a", "b", "c"), ScopeStore(scopeID(5))}, StringScopeStores + "/5/a/b/c"},
		{[]ScopeOption{Path("a", "b", "c"), ScopeStore(nil)}, StringScopeDefault + "/0/a/b/c"},
		{[]ScopeOption{Path("a", "b", "c"), ScopeWebsite(scopeID(50))}, StringScopeWebsites + "/50/a/b/c"},
		{[]ScopeOption{Path("a", "b", "c"), ScopeWebsite(nil)}, StringScopeDefault + "/0/a/b/c"},
		{nil, ""},
	}

	for i, test := range tests {
		arg := scopeKey(test.haveArg...)

		if arg == nil {
			t.Errorf("\narg is nil at index %d => %#v\n", i, test)
			return
		}

		actualPath := arg.scopePath()
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
	}
}

func TestScopeKeyValue(t *testing.T) {
	tests := []struct {
		haveArg []ScopeOption
		want    string
	}{
		{[]ScopeOption{Value(1), Path("a/b/c")}, StringScopeDefault + "/0/a/b/c"},
		{[]ScopeOption{Value("1"), Path("")}, ""},
		{[]ScopeOption{Value(1.1), Path()}, ""},
		{[]ScopeOption{Value(1), Scope(IDScopeDefault, nil)}, ""},
		{[]ScopeOption{Value(1), Scope(IDScopeWebsite, nil)}, ""},
		{[]ScopeOption{Value(1), Scope(IDScopeStore, nil)}, ""},
		{[]ScopeOption{Value(1), Path("a/b/c"), Scope(IDScopeWebsite, nil)}, StringScopeDefault + "/0/a/b/c"},
		{[]ScopeOption{Value(1), Path("a/b/c"), Scope(IDScopeWebsite, scopeID(2))}, StringScopeWebsites + "/2/a/b/c"},
		{[]ScopeOption{Value(1), Path("a", "b", "c"), Scope(IDScopeWebsite, scopeID(200))}, StringScopeWebsites + "/200/a/b/c"},
		{[]ScopeOption{Value(1), Path("a", "b", "c"), Scope(IDScopeStore, scopeID(4))}, StringScopeStores + "/4/a/b/c"},
		{[]ScopeOption{Value(1), Path("a", "b"), Scope(IDScopeStore, scopeID(4))}, StringScopeStores + "/4/a"},
		{[]ScopeOption{Value(1), nil, Scope(IDScopeStore, scopeID(4))}, ""},
		{[]ScopeOption{Value(1), Path("a", "b", "c"), ScopeStore(scopeID(5))}, StringScopeStores + "/5/a/b/c"},
		{[]ScopeOption{Value(1.2), Path("a", "b", "c"), ScopeStore(nil)}, StringScopeDefault + "/0/a/b/c"},
		{[]ScopeOption{Value(1.3), Path("a", "b", "c"), ScopeWebsite(scopeID(50))}, StringScopeWebsites + "/50/a/b/c"},
		{[]ScopeOption{ValueReader(strings.NewReader("a config value")), Path("a", "b", "c"), ScopeWebsite(nil)}, StringScopeDefault + "/0/a/b/c"},
		{nil, ""},
	}

	for _, test := range tests {
		arg := scopeKeyValue(test.haveArg...)
		actualPath, actualVal := arg.scopePath(), arg.v
		assert.EqualValues(t, test.want, actualPath, "Test: %#v", test)
		if test.haveArg != nil {
			assert.NotEmpty(t, actualVal, "Test: %#v", test)
		} else {
			assert.Empty(t, actualVal, "Test: %#v", test)
		}
	}
}

// All benchmarks MacBook Air (13-inch, Mid 2012); 1.8 GHz Intel Core i5; 8 GB 1600 MHz DDR3

var benchmarkScopeKey string

// BenchmarkScopeKey____InMap	 2000000	       936 ns/op	     176 B/op	       9 allocs/op
func BenchmarkScopeKey____InMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg := scopeKey(Path("a", "b", "c"), Scope(IDScopeWebsite, scopeID(4)))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey_NotInMap	 2000000	       992 ns/op	     200 B/op	      10 allocs/op
func BenchmarkScopeKey_NotInMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg := scopeKey(Path("a", "b", "c"), Scope(IDScopeWebsite, scopeID(40)))
		benchmarkScopeKey = arg.scopePath()
	}
}

// BenchmarkScopeKey____InMapNoJoin	 2000000	       824 ns/op	     176 B/op	       8 allocs/op
func BenchmarkScopeKey____InMapNoJoin(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		arg := scopeKey(Path("a/b/c"), Scope(IDScopeWebsite, scopeID(3)))
		benchmarkScopeKey = arg.scopePath()
	}
}
