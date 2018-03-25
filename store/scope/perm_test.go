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

package scope_test

import (
	"testing"

	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
)

func TestPermAll(t *testing.T) {

	var p scope.Perm
	pa := p.All()
	assert.True(t, pa.Has(scope.Default))
	assert.True(t, pa.Has(scope.Website))
	assert.True(t, pa.Has(scope.Store))
}

func TestPermTop(t *testing.T) {

	assert.Exactly(t, scope.Website, scope.PermWebsite.Top())
	assert.Exactly(t, scope.Store, scope.PermStore.Top())
	assert.Exactly(t, scope.Default, scope.PermDefault.Top())
	assert.Exactly(t, scope.Website, scope.Perm(44).Top())
	assert.Exactly(t, scope.Store, scope.PermWebsiteReverse.Top())
	assert.Exactly(t, scope.Store, scope.PermStoreReverse.Top())
}

func TestPermMarshalJSONAll(t *testing.T) {

	var p scope.Perm
	pa := p.All()
	jd, err := pa.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "[\"Default\",\"Website\",\"Store\"]", string(jd))
}

func TestPermMarshalJSONNull(t *testing.T) {

	var p scope.Perm
	jd, err := p.MarshalJSON()
	assert.NoError(t, err)
	assert.Exactly(t, "null", string(jd))
}
