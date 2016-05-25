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

package backendauth

import (
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/source"
)

// Backend just exported for the sake of documentation. See fields
// for more information. The PkgBackend handles the reading and writing
// of configuration values within this package.
type Backend struct {
	cfgmodel.PkgBackend

	// NetAuthEnable indicates whether authentication has been enabled or not.
	//
	// Path: net/auth/enable
	NetAuthEnable cfgmodel.Bool

	// NetAuthRequireTLS indicates whether we require TLS.
	//
	// Path: net/auth/require_tls
	NetAuthRequireTLS cfgmodel.Bool

	// NetAuthAllowedIP indicates which IPs are allowed.
	// Separate via line break (\n).
	//
	// Path: net/auth/allowed_ips
	NetAuthAllowedIPs cfgmodel.StringCSV

	// NetAuthDeniedIPs indicates which IPs are denied.
	// Separate via line break (\n).
	//
	// Path: net/auth/denied_ips
	NetAuthDeniedIPs cfgmodel.StringCSV

	// NetAuthAllowedIPRange indicates which IP ranges are denied.
	// Separate via line break (\n).
	//
	// Path: net/auth/denied_ips
	NetAuthAllowedIPRange ConfigIPRange

	// NetAuthDeniedIPRange indicates which IP ranges are denied.
	// Separate via line break (\n).
	//
	// Path: net/auth/denied_ips
	NetAuthDeniedIPRange ConfigIPRange

	// and so on
	// range based allowances and denies
}

// New initializes the backend configuration models containing the
// cfgpath.Route variable to the appropriate entries.
// The function Load() will be executed to apply the SectionSlice
// to all models. See Load() for more details.
func New(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	return (&Backend{}).Load(cfgStruct, opts...)
}

// Load creates the configuration models for each PkgBackend field.
// Internal mutex will protect the fields during loading.
// The argument SectionSlice will be applied to all models.
func (pp *Backend) Load(cfgStruct element.SectionSlice, opts ...cfgmodel.Option) *Backend {
	pp.Lock()
	defer pp.Unlock()

	opts = append(opts, cfgmodel.WithFieldFromSectionSlice(cfgStruct))
	optsCSV := append([]cfgmodel.Option{}, opts...)
	optsCSV = append(optsCSV, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithCSVComma('\n'))
	optsYN := append([]cfgmodel.Option{}, opts...)
	optsYN = append(optsYN, cfgmodel.WithFieldFromSectionSlice(cfgStruct), cfgmodel.WithSource(source.YesNo))

	pp.NetAuthEnable = cfgmodel.NewBool(`net/auth/allow_credentials`, optsYN...)
	pp.NetAuthAllowedIPs = cfgmodel.NewStringCSV(`net/auth/exposed_headers`, optsCSV...)

	return pp
}
