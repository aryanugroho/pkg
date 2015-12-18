// Copyright 2015, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package model_test

import (
	"testing"

	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/configsource"
	"github.com/corestoreio/csfw/config/model"
	"github.com/corestoreio/csfw/config/scope"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {

	wantPath := scope.StrStores.FQPathInt64(3, "web/cors/allow_credentials")
	b := model.NewBool("web/cors/allow_credentials", configsource.YesNo...)

	assert.Exactly(t, configsource.YesNo, b.Options)
	// because default value in packageConfiguration is "true"
	assert.True(t, b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.False(t, b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: 0,
		}),
	).NewScoped(0, 0, 3)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, true, scope.StoreID, 3))
	assert.Exactly(t, wantPath, mw.ArgPath)
}

func TestString(t *testing.T) {

	wantPath := scope.StrDefault.FQPathInt64(0, "web/cors/exposed_headers")
	b := model.NewString("web/cors/exposed_headers")

	assert.Empty(t, b.Options)

	assert.Exactly(t, "Content-Type,X-CoreStore-ID", b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, "X-Gopher", b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: "X-Gopher",
		}),
	).NewScoped(0, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, "dude", scope.DefaultID, 0))
	assert.Exactly(t, wantPath, mw.ArgPath)

}

func TestInt(t *testing.T) {

	wantPath := scope.StrWebsites.FQPathInt64(10, "web/cors/int")
	b := model.NewInt("web/cors/int")

	assert.Empty(t, b.Options)

	assert.Exactly(t, 2015, b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, 2016, b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: 2016,
		}),
	).NewScoped(10, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, 1, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath, mw.ArgPath)

}

func TestFloat64(t *testing.T) {
	wantPath := scope.StrWebsites.FQPathInt64(10, "web/cors/float64")
	b := model.NewFloat64("web/cors/float64")

	assert.Empty(t, b.Options)

	assert.Exactly(t, 2015.1000001, b.Get(packageConfiguration, config.NewMockGetter().NewScoped(0, 0, 0)))

	assert.Exactly(t, 2016.1000001, b.Get(packageConfiguration, config.NewMockGetter(
		config.WithMockValues(config.MockPV{
			wantPath: 2016.1000001,
		}),
	).NewScoped(10, 0, 0)))

	mw := &config.MockWrite{}
	assert.NoError(t, b.Write(mw, 1, scope.WebsiteID, 10))
	assert.Exactly(t, wantPath, mw.ArgPath)

}