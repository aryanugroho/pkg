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

package cfgbigcache_test

import (
	"testing"

	"github.com/allegro/bigcache"
	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/storage/cfgbigcache"
	"github.com/corestoreio/pkg/util/conv"
	"github.com/stretchr/testify/assert"
)

var _ config.Storager = (*cfgbigcache.Storage)(nil)

func TestCacheGet(t *testing.T) {

	sc, err := cfgbigcache.New(bigcache.Config{
		Shards: 64,
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		key        cfgpath.Path
		val        interface{}
		wantSetErr error
		wantGetErr error
	}{
		{cfgpath.MustMakeByString("aa/bb/cc"), 12345, nil, nil},
	}
	for idx, test := range tests {

		haveSetErr := sc.Set(test.key, test.val)
		if test.wantSetErr != nil {
			assert.EqualError(t, haveSetErr, test.wantSetErr.Error(), "Index %d", idx)
		} else {
			assert.NoError(t, haveSetErr, "Index %d", idx)
		}

		haveVal, haveGetErr := sc.Value(test.key)
		if test.wantGetErr != nil {
			assert.EqualError(t, haveGetErr, test.wantGetErr.Error(), "Index %d", idx)
		} else {
			assert.NoError(t, haveGetErr, "Index %d", idx)
		}
		// don't do this 2x conv casting in production code
		assert.Exactly(t, test.val, conv.ToInt(conv.ToString(haveVal)), "Index %d => %v", idx, conv.ToString(haveVal))
	}
}

func TestCacheGetNotFound(t *testing.T) {
	sc, err := cfgbigcache.New(bigcache.Config{
		Shards: 64,
	})
	if err != nil {
		t.Fatal(err)
	}
	haveVal, haveGetErr := sc.Value(cfgpath.MustMakeByString("aa/bb/cc"))
	assert.True(t, errors.NotFound.Match(haveGetErr), "Error: %s", haveGetErr)
	assert.Empty(t, haveVal)
}

func TestCacheError(t *testing.T) {
	sc, err := cfgbigcache.New(bigcache.Config{
		Shards: 63,
	})
	assert.True(t, errors.IsFatal(err), "Error: %s", err)
	assert.Empty(t, sc)
}
