// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// Errors
const (
	errMainCategoriesTypeWithNameNone = "there is no main category type with such name"
	errMainCategoryTypeNameAmbiguous  = "main category type name is ambiguous"
)

// MainCategoryStatusT describes the behaviour of categories and its descendants (transactions)
type MainCategoryType struct {
	Id     int
	Name   string
	Factor int
}

// MCT* constants identify particular type in database (used as id)
const (
	MCTUnknown = iota
	MCTUnset
	MCTCost
	MCTTransfer
	MCTIncome
)

// MainCategoryTypeForID returns pointer to the Main Category Type for given ID
func MainCategoryTypeForID(db *gsqlitehandler.SqliteDB, i int) (mt *MainCategoryType, err error) {
	var stmt *sql.Stmt

	sqlQuery := "SELECT id, name, factor " +
		"FROM main_categories_types " +
		"WHERE id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	mt = new(MainCategoryType)
	if err = stmt.QueryRow(i).Scan(&mt.Id, &mt.Name, &mt.Factor); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	return mt, nil
	//TODO: add test
}

// MainCategoryTYpeForName returns pointer to Main Category Type for given (part of) name
func MainCategoryTypeForName(db *gsqlitehandler.SqliteDB, n string) (mt *MainCategoryType, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	n = "%" + n + "%"
	sqlQuery := "SELECT id, name, factor " +
		"FROM main_categories_types " +
		"WHERE name LIKE ?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	mt = new(MainCategoryType)
	if rows, err = stmt.Query(n); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer rows.Close()

	var noOfTypes int
	for rows.Next() {
		noOfTypes++
		rows.Scan(&mt.Id, &mt.Name, &mt.Factor)
	}

	switch noOfTypes {
	case 0:
		return nil, errors.New(errMainCategoriesTypeWithNameNone)
	case 1:
		return mt, nil
	default:
		return nil, errors.New(errMainCategoryTypeNameAmbiguous)
	}

	//TODO: add test
}
