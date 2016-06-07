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

package transcache_test

import (
	"net"
	"testing"

	"github.com/corestoreio/csfw/storage/transcache"
	"github.com/stretchr/testify/assert"
)

var _ transcache.Transcacher = (*transcache.Mock)(nil)

func TestMock_SetGet(t *testing.T) {

	mck := transcache.NewMock()

	key := net.ParseIP("192.168.100.0")
	val := "abc"
	if err := mck.Set(key, val); err != nil {
		t.Error(err)
	}

	var dstVal string
	if err := mck.Get(key, &dstVal); err != nil {
		t.Error(err)
	}

	assert.Exactly(t, val, dstVal)

	assert.Exactly(t, 1, mck.SetCount())
	assert.Exactly(t, 1, mck.GetCount())
}

//func TestMock_SetGet_Multi(t *testing.T) {
//	tests := []struct {
//		key []byte
//		val interface{}
//	}{
//		{net.ParseIP("192.168.100.0"), "a"},
//		{net.ParseIP("192.168.100.1"), 1},
//		{net.ParseIP("192.168.100.2"), 3.14152 * 2.7182},
//		{net.ParseIP("192.168.100.3"), t},
//	}
//	mck := transcache.NewMock()
//	for i, test := range tests {
//		if err := mck.Set(test.key, test.val); err != nil {
//			t.Error("Index", i, err)
//		}
//		var dst interface{}
//		switch test.val.(type) {
//		case string:
//			var dsts string
//			if err := mck.Get(test.key, &dsts); err != nil {
//				t.Error("Index", i, err)
//			}
//			dst = dsts
//		case int:
//			var dsti int
//			if err := mck.Get(test.key, &dsti); err != nil {
//				t.Error("Index", i, err)
//			}
//			dst = dsti
//		}
//		assert.Exactly(t, test.val, dst, "Index %d", i)
//	}
//	assert.Exactly(t, len(tests), mck.SetCount())
//	assert.Exactly(t, len(tests), mck.GetCount())
//}
