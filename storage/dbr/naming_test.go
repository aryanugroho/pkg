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

package dbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeAlias(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "`table1`", MakeAlias("table1").String())
	assert.Exactly(t, "`table0` AS `table1`", MakeAlias("table0", "table1").String())
	assert.Exactly(t, "`(table1)`", MakeAlias("(table1)").String())
	assert.Exactly(t, "`(table1)` AS `table2`", MakeAlias("(table1)", "table2").String())
	assert.Exactly(t, "`(table1)`", MakeAlias("(table1)", "").String())
	assert.Exactly(t, "`table1`", MakeAlias("table1", "").String())
}

func TestMakeAliasExpr(t *testing.T) {
	t.Parallel()
	assert.Exactly(t, "(table1)", MakeAliasExpr("(table1)", "").String())
	assert.Exactly(t, "(table1) AS `x`", MakeAliasExpr("(table1)", "x").String())
	assert.Exactly(t, "(table1)", MakeAliasExpr("(table1)").String())
}

func TestQuoteAs(t *testing.T) {
	tests := []struct {
		have []string
		want string
	}{
		0: {[]string{"a"}, "`a`"},
		1: {[]string{"a", "b"}, "`a` AS `b`"},
		2: {[]string{"a", ""}, "`a`"},
		3: {[]string{"`c`"}, "`c`"},
		4: {[]string{"d.e"}, "`d`.`e`"},
		5: {[]string{"`d`.`e`"}, "`d`.`e`"},
		6: {[]string{"f", "g", "h"}, "`f` AS `g_h`"},
		7: {[]string{"f", "g", "h`h"}, "`f` AS `g_hh`"},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, Quoter.QuoteAs(test.have...), "Index %d", i)
	}
}

func BenchmarkQuoteAs(b *testing.B) {
	const want = "`e`.`entity_id` AS `ee`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.QuoteAs("e.entity_id", "ee"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

func BenchmarkQuoteAlias(b *testing.B) {
	const want = "(e.price * a.tax * e.weee) AS `final_price`"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.exprAlias("(e.price * a.tax * e.weee)", "final_price"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

func BenchmarkQuoteQuote(b *testing.B) {
	const want = "`databaseName`.`tableName`"

	b.ReportAllocs()
	b.ResetTimer()
	b.Run("Worse Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := Quoter.Quote("database`Name", "table`Name"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
	b.Run("Best Case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if have := Quoter.Quote("databaseName", "tableName"); have != want {
				b.Fatalf("Have %s\nWant %s\n", have, want)
			}
		}
	})
}

func TestMysqlQuoter_Quote(t *testing.T) {
	assert.Exactly(t, "`tableName`", Quoter.Quote("tableName"))
	assert.Exactly(t, "`databaseName`.`tableName`", Quoter.Quote("databaseName", "tableName"))
	assert.Exactly(t, "`tableName`", Quoter.Quote("", "tableName")) // qualifier is empty
	assert.Exactly(t, "`databaseName`.`tableName`", Quoter.Quote("database`Name", "table`Name"))
}

func TestIsValidIdentifier(t *testing.T) {
	tests := []struct {
		have string
		want int8
	}{
		{"*", 0},
		{"table.*", 0},
		{"*.*", 2},
		{"table.p*", 2},
		{"`table`.*", 2},     // not valid because of backticks
		{"`table`.`col`", 2}, // not valid because of backticks
		{"", 1},
		{"a", 0},
		{"a.", 1},
		{"a.b", 0},
		{".b", 1},
		{"", 2},
		{"花间一壶酒，独酌无相亲。", 2}, // no idea what this means but found it in x/text pkg
		{"独酌无相", 2},         // no idea what this means but found it in x/text pkg
		{"Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 1},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 0},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher", 0},
		{"Gooooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooo0opher", 1},
		{"Goooooooooooooooooooooooooooooooooooooooooooooooooooooooooopher.Gooooooooooooooooooooooooooooooooooooooooooooooooooooooo0oph€r", 2},
		{"DATE_FORMAT(t3.period, '%Y-%m-01')", 2},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, isValidIdentifier(test.have), "Index %d with %q", i, test.have)
	}
}

var benchmarkIsValidIdentifier int8

// BenchmarkIsValidIdentifier-4   	20000000	        92.0 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIsValidIdentifier-4   	 5000000	       280 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIsValidIdentifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = isValidIdentifier(`store_owner.catalog_product_entity_varchar`)
	}
	if benchmarkIsValidIdentifier != 0 {
		b.Fatalf("Should be zero but got %d", benchmarkIsValidIdentifier)
	}
}
