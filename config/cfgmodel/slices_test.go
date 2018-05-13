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

package cfgmodel_test

import (
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/config/cfgmock"
	"github.com/corestoreio/pkg/config/cfgmodel"
	"github.com/corestoreio/pkg/config/cfgpath"
	"github.com/corestoreio/pkg/config/cfgsource"
	"github.com/corestoreio/pkg/store/scope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringCSVGet(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := cfgpath.MustMakeByString(pathWebCorsHeaders).String()
	b := cfgmodel.NewStringCSV(
		"web/cors/exposed_headers",
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithSourceByString(
			"Content-Type", "Content Type", "X-CoreStore-ID", "CoreStore Microservice ID",
		),
	)
	assert.NotEmpty(t, b.Options())

	sm := cfgmock.NewService()
	sl, err := b.Value(sm.NewScoped(0, 0))
	assert.NoError(t, err)
	assert.Exactly(t, []string{"Content-Type", "X-CoreStore-ID"}, sl) // default values from variable configStructure
	assert.Exactly(t, typeIDsDefault, sm.StringInvokes().ScopeIDs())

	tests := []struct {
		have    string
		want    []string
		wantIDs scope.TypeIDs
		wantErr error
	}{
		{"Content-Type,X-CoreStore-ID", []string{"Content-Type", "X-CoreStore-ID"}, scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(1)}, nil},
		{"", nil, scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(1)}, nil},
		{"X-CoreStore-ID", []string{"X-CoreStore-ID"}, scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(1)}, nil},
		{"Content-Type,X-CS", []string{"Content-Type", "X-CS"}, scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(1)}, nil},
		// todo add errors
	}
	for i, test := range tests {
		sm := cfgmock.NewService(cfgmock.PathValue{
			wantPath: test.have,
		})
		haveSL, haveErr := b.Value(sm.NewScoped(1, 0)) // 1,0 because scope of pathWebCorsHeaders is default,website

		assert.Exactly(t, test.want, haveSL, "Index %d", i)
		assert.Exactly(t, test.wantIDs, sm.StringInvokes().ScopeIDs(), "Index %d", i)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestStringCSVWrite(t *testing.T) {

	const pathWebCorsHeaders = "web/cors/exposed_headers"
	wantPath := cfgpath.MustMakeByString(pathWebCorsHeaders).String()
	b := cfgmodel.NewStringCSV(
		"web/cors/exposed_headers",
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithSourceByString(
			"Content-Type", "Content Type", "X-CoreStore-ID", "CoreStore Microservice ID",
		),
	)

	mw := &cfgmock.Write{}
	b.Source.Merge(cfgsource.MustNewByString("a", "a", "b", "b", "c", "c"))

	assert.NoError(t, b.Write(mw, []string{"a", "b", "c"}, scope.DefaultTypeID))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "a,b,c", mw.ArgValue.(string))
	err := b.Write(mw, []string{"abc"}, scope.DefaultTypeID)
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

func TestStringCSVCustomSeparator(t *testing.T) {

	const cfgPath = "aa/bb/cc"

	b := cfgmodel.NewStringCSV(
		cfgPath,
		cfgmodel.WithSourceByString(
			"2014", "Year 2014",
			"2015", "Year 2015",
			"2016", "Year 2016",
			"2017", "Year 2017",
		),
		cfgmodel.WithCSVComma(''),
	)
	wantPath := cfgpath.MustMakeByString(cfgPath).String() // Default Scope

	sm := cfgmock.NewService(cfgmock.PathValue{
		wantPath: `20152016`,
	})
	haveSL, haveErr := b.Value(sm.NewScoped(34, 4))
	if haveErr != nil {
		t.Fatal(haveErr)
	}
	assert.Exactly(t, []string{"2015", "2016"}, haveSL)
	assert.Exactly(t, typeIDsDefault, sm.StringInvokes().ScopeIDs())
}

func TestIntCSV(t *testing.T) {

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := cfgmodel.NewIntCSV(
		pathWebCorsIntSlice,
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithSourceByInt(cfgsource.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
	)
	assert.Len(t, b.Options(), 4)
	assert.Exactly(t, pathWebCorsIntSlice, b.String())
	// default values:
	sl, err := b.Value(cfgmock.NewService().NewScoped(0, 4))
	assert.NoError(t, err)
	assert.Exactly(t, []int{2014, 2015, 2016}, sl) // three years are defined in variable configStructure
	//assert.Exactly(t, scope.DefaultTypeID.String(), h.String())

	wantPath := cfgpath.MustMakeByString(pathWebCorsIntSlice).BindStore(4).String()

	tests := []struct {
		lenient bool
		have    string
		want    []int
		wantIDs scope.TypeIDs
		wantBhf errors.BehaviourFunc
	}{
		{false, "3015,3016", []int{3015, 3016}, scope.TypeIDs{scope.Store.WithID(4)}, nil},
		{false, "2015,2017", []int{2015, 2017}, scope.TypeIDs{scope.Store.WithID(4)}, nil},
		{false, "", nil, scope.TypeIDs{scope.Store.WithID(4)}, nil},
		{false, "2015,,20x17", []int{2015}, scope.TypeIDs{scope.Store.WithID(4)}, errors.IsNotValid},
		{true, "2015,,2017", []int{2015, 2017}, scope.TypeIDs{scope.Store.WithID(4)}, nil},
	}
	for i, test := range tests {
		b.Lenient = test.lenient
		sm := cfgmock.NewService(cfgmock.PathValue{
			wantPath: test.have,
		})
		haveSL, haveErr := b.Value(sm.NewScoped(0, 4))

		assert.Exactly(t, test.want, haveSL, "Index %d", i)
		assert.Exactly(t, test.wantIDs, sm.StringInvokes().ScopeIDs(), "Index %d", i)
		if test.wantBhf != nil {
			assert.True(t, test.wantBhf(haveErr), "Index %d => %+v", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestIntCSVWrite(t *testing.T) {

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := cfgmodel.NewIntCSV(
		pathWebCorsIntSlice,
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithSourceByInt(cfgsource.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
	)
	wantPath := cfgpath.MustMakeByString(pathWebCorsIntSlice).BindStore(4).String()

	mw := &cfgmock.Write{}
	b.Source.Merge(cfgsource.NewByInt(cfgsource.Ints{
		{2018, "Year 2018"},
	}))
	assert.NoError(t, b.Write(mw, []int{2016, 2017, 2018}, scope.Store.WithID(4)))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "2016,2017,2018", mw.ArgValue.(string))
	err := b.Write(mw, []int{2019}, scope.Store.WithID(4))
	assert.True(t, errors.IsNotValid(err), "Error: %s", err)
}

func TestIntCSVCustomSeparator(t *testing.T) {

	const pathWebCorsIntSlice = "web/cors/int_slice"

	b := cfgmodel.NewIntCSV(
		pathWebCorsIntSlice,
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithSourceByInt(cfgsource.Ints{
			{2014, "Year 2014"},
			{2015, "Year 2015"},
			{2016, "Year 2016"},
			{2017, "Year 2017"},
		}),
		cfgmodel.WithCSVComma('|'),
	)
	wantPath := cfgpath.MustMakeByString(pathWebCorsIntSlice).BindWebsite(34).String()

	sm := cfgmock.NewService(cfgmock.PathValue{
		wantPath: `2015|2016|`,
	})
	haveSL, haveErr := b.Value(sm.NewScoped(34, 4))
	if haveErr != nil {
		t.Fatal(haveErr)
	}
	assert.Exactly(t, []int{2015, 2016}, haveSL)
	assert.Exactly(t, scope.TypeIDs{scope.Website.WithID(34), scope.Store.WithID(4)}, sm.StringInvokes().ScopeIDs())
}

func TestCSVGet(t *testing.T) {

	const pathWebCorsCSV = "web/cors/csv"
	wantPath := cfgpath.MustMakeByString(pathWebCorsCSV).String()
	b := cfgmodel.NewCSV(
		"web/cors/csv",
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithCSVComma('|'),
	)
	assert.Empty(t, b.Options())

	sl, err := b.Value(cfgmock.NewService().NewScoped(0, 0))
	require.NoError(t, err, "%+v", err)
	assert.Exactly(t,
		[][]string{{"0", "\"Did you mean...\" Suggestions", "\"meinten Sie...?\""}, {"1", "Accuracy for Suggestions", "Genauigkeit der Vorschläge"}, {"2", "After switching please reindex the<br /><em>Catalog Search Index</em>.", "Nach dem Umschalten reindexieren Sie bitte den <br /><em>Katalog Suchindex</em>."}, {"3", "CATALOG", "KATALOG"}},
		sl) // default values from variable configStructure

	tests := []struct {
		have       string
		want       [][]string
		wantErrBhf errors.BehaviourFunc
	}{
		{"Content-Type|X-CoreStore-ID", [][]string{{"Content-Type", "X-CoreStore-ID"}}, nil},
		{"", nil, nil},
		{"X-CoreStore-ID", [][]string{{"X-CoreStore-ID"}}, nil},
		{"Content-Type|X-CS", [][]string{{"Content-Type", "X-CS"}}, nil},
		{"Content-Type|X-CS\nApplication", nil, errors.IsNotValid},
	}
	for i, test := range tests {
		sm := cfgmock.NewService(cfgmock.PathValue{
			wantPath: test.have,
		})
		haveSL, haveErr := b.Value(sm.NewScoped(1, 2)) // because scope of pathWebCorsHeaders is default,website

		assert.Exactly(t, test.want, haveSL, "Index %d", i)
		assert.Exactly(t, scope.TypeIDs{scope.DefaultTypeID, scope.Website.WithID(1)}, sm.StringInvokes().ScopeIDs(), "Index %d", i)
		if test.wantErrBhf != nil {
			assert.True(t, test.wantErrBhf(haveErr), "Index %d Error: %s", i, haveErr)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
	}
}

func TestCSVWrite(t *testing.T) {

	const pathWebCorsCsv = "web/cors/csv"
	wantPath := cfgpath.MustMakeByString(pathWebCorsCsv).String()
	b := cfgmodel.NewCSV(
		"web/cors/csv",
		cfgmodel.WithFieldFromSectionSlice(configStructure),
		cfgmodel.WithCSVComma('!'),
	)

	mw := &cfgmock.Write{}

	assert.NoError(t, b.Write(mw, [][]string{{"a", "b", "c"}, {"d", "e", "f"}}, scope.DefaultTypeID))
	assert.Exactly(t, wantPath, mw.ArgPath)
	assert.Exactly(t, "a!b!c\nd!e!f\n", mw.ArgValue.(string))
}
