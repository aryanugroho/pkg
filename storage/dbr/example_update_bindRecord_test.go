// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr_test

import (
	"fmt"
	"strings"

	"github.com/corestoreio/csfw/storage/dbr"
	"github.com/corestoreio/errors"
)

// Make sure that type categoryEntity implements interface
var _ dbr.ArgumentsAppender = (*categoryEntity)(nil)

// categoryEntity represents just a demo record.
type categoryEntity struct {
	EntityID       int64 // Auto Increment
	AttributeSetID int64
	ParentID       string
	Path           dbr.NullString
	// TeaserIDs contain a list of foreign primary keys which identifies special
	// teaser to be shown on the category page. Each teaser ID gets joined by a
	// | and stored as a long string in the database.
	TeaserIDs []string
}

func (pe categoryEntity) appendBind(args dbr.Arguments, column string) (_ dbr.Arguments, err error) {
	switch column {
	case "entity_id":
		args = args.Int64(pe.EntityID)
	case "attribute_set_id":
		args = args.Int64(pe.AttributeSetID)
	case "parent_id":
		args = args.Str(pe.ParentID)
	case "path":
		args = args.NullString(pe.Path)
	case "teaser_id_s":
		if pe.TeaserIDs == nil {
			args = args.Null()
		} else {
			args = args.Str(strings.Join(pe.TeaserIDs, "|"))
		}
	case "fk_teaser_id_s": // TODO ...
		args = args.Strs(pe.TeaserIDs...)
	default:
		return nil, errors.NewNotFoundf("[dbr_test] Column %q not found", column)
	}
	return args, nil
}

// AppendArgs implements dbr.ArgumentsAppender interface
func (pe categoryEntity) AppendArgs(args dbr.Arguments, columns []string) (dbr.Arguments, error) {
	l := len(columns)
	if l == 1 {
		// Most commonly used case
		return pe.appendBind(args, columns[0])
	}
	if l == 0 {
		// This case gets executed when an INSERT statement doesn't contain any
		// columns.
		return args.Int64(pe.EntityID).Int64(pe.AttributeSetID).Str(pe.ParentID), nil
	}
	// This case gets executed when an INSERT statement requests specific columns.
	for _, col := range columns {
		var err error
		if args, err = pe.appendBind(args, col); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return args, nil
}

// ExampleUpdate_BindRecord performs an UPDATE query in the table `catalog_category_entity` with the
// fix specified columns. The Go type categoryEntity implements the dbr.ArgumentsAppender interface and can
// append the required arguments.
func ExampleUpdate_BindRecord() {

	ce := categoryEntity{345, 6, "p123", dbr.MakeNullString("4/5/6/7"), []string{"saleAutumn", "saleShoe"}}

	// Updates all rows in the table because of missing WHERE statement.
	u := dbr.NewUpdate("catalog_category_entity").
		AddColumns("attribute_set_id", "parent_id", "path", "teaser_id_s").
		BindRecord(
			dbr.Qualify("", ce), // qualifier can be empty because no alias and no additional tables.
		)
	writeToSQLAndInterpolate(u)

	fmt.Print("\n\n")

	ce = categoryEntity{678, 6, "p456", dbr.NullString{}, nil}

	// Updates only one row in the table because of the WHERE. You can call
	// AddArgumentsAppender and Exec as often as you like. Each call to Exec will
	// reassemble the arguments from AddArgumentsAppender, means you can exchange AddArgumentsAppender
	// with different objects.
	u.
		BindRecord(dbr.Qualify("", ce)). // overwrites previously set default qualifier
		Where(dbr.Column("entity_id").PlaceHolder())
	writeToSQLAndInterpolate(u)

	// Output:
	//Prepared Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?,
	//`path`=?, `teaser_id_s`=?
	//Arguments: [6 p123 4/5/6/7 saleAutumn|saleShoe]
	//
	//Interpolated Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p123',
	//`path`='4/5/6/7', `teaser_id_s`='saleAutumn|saleShoe'
	//
	//Prepared Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=?, `parent_id`=?,
	//`path`=?, `teaser_id_s`=? WHERE (`entity_id` = ?)
	//Arguments: [6 p456 <nil> <nil> 678]
	//
	//Interpolated Statement:
	//UPDATE `catalog_category_entity` SET `attribute_set_id`=6, `parent_id`='p456',
	//`path`=NULL, `teaser_id_s`=NULL WHERE (`entity_id` = 678)
}