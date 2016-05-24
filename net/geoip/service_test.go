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

package geoip_test

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"bytes"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/net/geoip"
	"github.com/corestoreio/csfw/store"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/store/storemock"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/corestoreio/csfw/util/log"
	"github.com/stretchr/testify/assert"
)

func mustGetTestService(opts ...geoip.Option) *geoip.Service {
	maxMindDB := filepath.Join("testdata", "GeoIP2-Country-Test.mmdb")
	return geoip.MustNew(append(opts, geoip.WithGeoIP2File(maxMindDB))...)
}

func finalHandlerFinland(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.NotNil(t, ipc)
		assert.NoError(t, err)
		assert.Exactly(t, "2a02:d200::", ipc.IP.String())
		assert.Exactly(t, "FI", ipc.Country.IsoCode)
	}
}

func mustGetRequestFinland() *http.Request {
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	return req
}

func TestNewServiceErrorWithoutOptions(t *testing.T) {
	s, err := geoip.New()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Nil(t, s.GeoIP)
}

func TestNewServiceErrorWithGeoIP2Reader(t *testing.T) {
	s, err := geoip.New(geoip.WithGeoIP2File("Walhalla/GeoIP2-Country-Test.mmdb"))
	assert.Nil(t, s)
	assert.True(t, errors.IsNotFound(err), "Error: %s", err)
}

func deferClose(t *testing.T, c io.Closer) {
	assert.NoError(t, c.Close())
}

func TestNewServiceWithGeoIP2Reader(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)
	ip, _, err := net.ParseCIDR("2a02:d200::/29") // IP range for Finland

	assert.NoError(t, err)
	haveCty, err := s.GeoIP.Country(ip)
	assert.NoError(t, err)
	assert.Exactly(t, "FI", haveCty.Country.IsoCode)
}

type geoReaderMock struct{}

func (geoReaderMock) Country(ipAddress net.IP) (*geoip.Country, error) {
	return nil, errors.New("Failed to read country from MMDB")
}
func (geoReaderMock) Close() error { return nil }

func TestWithCountryByIPErrorGetCountryByIP(t *testing.T) {
	s := mustGetTestService()
	s.GeoIP = geoReaderMock{}
	defer deferClose(t, s.GeoIP)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.Nil(t, ipc)
		assert.True(t, errors.IsFatal(err), "Error: %s", err)
	})

	countryHandler := s.WithCountryByIP()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2a02:d200::")
	countryHandler.ServeHTTP(rec, req)
}

func TestWithCountryByIPSuccess(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithCountryByIP()(finalHandlerFinland(t))
	rec := httptest.NewRecorder()

	countryHandler.ServeHTTP(rec, mustGetRequestFinland())
}

func TestWithIsCountryAllowedByIPErrorStoreManager(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.Nil(t, ipc)
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	})

	countryHandler := s.WithIsCountryAllowedByIP()(finalHandler)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	countryHandler.ServeHTTP(rec, req)
}

func ipErrorFinalHandler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ipc, err := geoip.FromContextCountry(r.Context())
		assert.Nil(t, ipc)
		assert.True(t, errors.IsNotFound(err), "Error: %s", err)
	}
}

func TestWithCountryByIPErrorRemoteAddr(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithCountryByIP()(ipErrorFinalHandler(t))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2324.2334.432.534")
	countryHandler.ServeHTTP(rec, req)
}

func TestWithIsCountryAllowedByIPErrorWithContextCountryByIP(t *testing.T) {
	s := mustGetTestService()
	defer deferClose(t, s.GeoIP)

	countryHandler := s.WithIsCountryAllowedByIP()(ipErrorFinalHandler(t))
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://corestore.io", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "2R02:d2'0.:")

	countryHandler.ServeHTTP(rec, req) // managerStoreSimpleTest,  ?
}

func TestWithIsCountryAllowedByIP_MultiScopes(t *testing.T) {
	logBuf := &bytes.Buffer{}
	s := mustGetTestService(
		geoip.WithLogger(log.NewStdLog(log.WithStdWriter(logBuf), log.WithStdLevel(log.StdLevelDebug))),
	)
	defer deferClose(t, s.GeoIP)

	o, err := scope.SetByCode(scope.Website, "euro")
	if err != nil {
		t.Fatal(err)
	}
	storeSrv := storemock.NewEurozzyService(o)

	var finalTestHandler = func(i int, wantCountryISO string, wantErrorBhf errors.BehaviourFunc, wantAltHandler bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if wantAltHandler && i < 900 {
				t.Fatalf("Expecting alternative Handler, Got %d", i) // bit hacky but ok
			}

			ipc, err := geoip.FromContextCountry(r.Context())
			if wantErrorBhf != nil {
				assert.Nil(t, ipc, "Index %d", i)
				assert.True(t, wantErrorBhf(err), "Index %d Error: %s", i, err)
				// t.Log(err)
				return
			}
			assert.NoError(t, err, "Index %d", i)
			if assert.NotNil(t, ipc, "Index %d", i) {
				assert.Exactly(t, wantCountryISO, ipc.Country.IsoCode, "Index %d", i)
			}
		})
	}

	if err := s.Options(
		geoip.WithAlternativeHandler(scope.Store, 2, finalTestHandler(999, "AT", nil, false)), // for test case 2
		geoip.WithAllowedCountryCodes(scope.Store, 2, "DE", "CH", "FI"),
	); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		req            func() *http.Request
		wantCountryISO string
		wantErrorBhf   errors.BehaviourFunc
		wantAltHandler bool
	}{
		// requested store not found
		0: {func() *http.Request {
			req, _ := http.NewRequest("GET", "http://corestore.io", nil)
			return req
		},
			"", errors.IsNotFound, false},

		// IP detected as origin from Finland
		1: {func() *http.Request {
			req, _ := http.NewRequest("GET", "http://corestore.io", nil)
			req.Header.Set("X-Forwarded-For", "2a02:d200::")
			st, err := storeSrv.Store(scope.MockID(1)) // German Store
			if err != nil {
				t.Fatal(errors.PrintLoc(err))
			}
			return req.WithContext(store.WithContextRequestedStore(req.Context(), st))
		},
			"FI", nil, false},

		// IP detected as origin from AT and alternative handler for scope Store == 2 gets called but AT not allowed
		2: {func() *http.Request {
			req, _ := http.NewRequest("GET", "http://corestore.io", nil)
			req.RemoteAddr = "2a02:da80::"
			st, err := storeSrv.Store(scope.MockID(2)) // Austria Store
			if err != nil {
				t.Fatal(errors.PrintLoc(err))
			}
			st.Config = cfgmock.NewService().NewScoped(1, 2)
			return req.WithContext(store.WithContextRequestedStore(req.Context(), st))
		},
			"AT", nil, true},

		// IP detection errors and an error gets attached to the context
		3: {func() *http.Request {
			req, _ := http.NewRequest("GET", "http://corestore.io", nil)
			req.RemoteAddr = "Er00r"
			st, err := storeSrv.Store(scope.MockID(2)) // Austria Store
			if err != nil {
				t.Fatal(errors.PrintLoc(err))
			}
			st.Config = cfgmock.NewService().NewScoped(1, 2)
			return req.WithContext(store.WithContextRequestedStore(req.Context(), st))
		},
			"XX", errors.IsNotFound, false},
	}
	for i, test := range tests {
		h := s.WithIsCountryAllowedByIP()(finalTestHandler(i, test.wantCountryISO, test.wantErrorBhf, test.wantAltHandler))
		h.ServeHTTP(nil, test.req())
	}
	//	println("\n", logBuf.String())
}
