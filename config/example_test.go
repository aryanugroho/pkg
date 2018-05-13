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

package config

import (
	"fmt"

	"github.com/corestoreio/pkg/store/scope"
)

func Example() {

	fmt.Println(MustNewPath("system/smtp/host").String())
	fmt.Println(MustNewPath("system/smtp/host").BindWebsite(1).String())
	// alternative way
	fmt.Println(MustNewPath("system/smtp/host").BindWebsite(1).String())

	fmt.Println(MustNewPath("system/smtp/host").BindStore(3).String())
	// alternative way
	fmt.Println(MustNewPath("system/smtp/host").BindStore(3).String())
	// Group is not supported and falls back to default
	fmt.Println(MustNewPath("system/smtp/host").Bind(scope.Group.WithID(4)).String())

	p, err := NewPath("system/smtp/host")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(p.String())

	routes := MustNewPath("dev/css/merge_css_files")
	rs, err := routes.Split()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println("dev/css/merge_css_files => ", rs[0], rs[1], rs[2])

	// Output:
	//default/0/system/smtp/host
	//websites/1/system/smtp/host
	//websites/1/system/smtp/host
	//stores/3/system/smtp/host
	//stores/3/system/smtp/host
	//default/0/system/smtp/host
	//default/0/system/smtp/host
	//dev/css/merge_css_files =>  dev css merge_css_files
}
