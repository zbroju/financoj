// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package lib

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// MainCategory represents the basic object for main category
type MainCategory struct {
	Id     int64
	MType  *MainCategoryType
	Name   string
	Status ItemStatus
}

// MainCategoryNew returns pointer to newly created MainCategory object
func MainCategoryNew() *MainCategory {
	m := new(MainCategory)
	m.MType = new(MainCategoryType)

	return m
}

// MainCategoryAdd adds new main category with type t and name n
func MainCategoryAdd(db *gsqlitehandler.SqliteDB, m *MainCategory) error {
	var err error
	var stmt *sql.Stmt

	if m.MType == nil {
		t := new(MainCategoryType)
		t.Id = MCTCost
	}

	if stmt, err = db.Handler.Prepare("INSERT INTO main_categories VALUES (NULL, ?, ?, ?);"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType.Id, m.Name, m.Status); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryForID returns pointer to the MainCategory for given id
func MainCategoryForID(db *gsqlitehandler.SqliteDB, i int) (m *MainCategory, err error) {
	var stmt *sql.Stmt

	sqlQuery := "SELECT m.id, m.name, m.status, t.id, t.name, t.factor " +
		"FROM main_categories m INNER JOIN main_categories_types t ON m.type_id=t.id " +
		"WHERE m.id=? AND m.status<>?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = MainCategoryNew()
	if err = stmt.QueryRow(i, ISClose).Scan(&m.Id, &m.Name, &m.Status, &m.MType.Id, &m.MType.Name, &m.MType.Factor); err != nil {
		return nil, errors.New(errMainCategoryWithIDNone)
	}

	return m, nil
	//TODO: add test
	//FIXME: move all sql string to separate variable
	//FIXME: replace all '*' in SELECT sql string to separated fields
}

// MainCategoryForName returns pointer to MainCategoryT for given (part of) name
func MainCategoryForName(db *gsqlitehandler.SqliteDB, n string) (m *MainCategory, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	n = "%" + n + "%"

	sqlQuery := "SELECT m.id, m.name, m.status, t.id, t.name, t.factor " +
		"FROM main_categories m INNER JOIN main_categories_types t ON m.type_id=t.id " +
		"WHERE m.name LIKE ? AND m.status<>?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = MainCategoryNew()
	if rows, err = stmt.Query(n, ISClose); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer rows.Close()

	var noOfMainCategories int
	for rows.Next() {
		noOfMainCategories++
		rows.Scan(&m.Id, &m.Name, &m.Status, &m.MType.Id, &m.MType.Name, &m.MType.Factor)
	}

	switch noOfMainCategories {
	case 0:
		return nil, errors.New(errMainCategoryWithNameNone)
	case 1:
		return m, nil
	default:
		return nil, errors.New(errMainCategoryNameAmbiguous)
	}

	//TODO: add test
}

// MainCategoryEdit updates main category with new values for type, name and status
// Both type, name and status is updated, so make sure you pass old values in argument 'm'
func MainCategoryEdit(db *gsqlitehandler.SqliteDB, m *MainCategory) error {
	var err error
	var stmt *sql.Stmt

	// Check if it is not a system object
	if m.Status == ISSystem {
		return errors.New(errSystemObject)
	}

	sqlQuery := "UPDATE main_categories SET type_id=?, name=?, status=? WHERE id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType.Id, m.Name, m.Status, m.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryRemove updates main category status with ISClose
func MainCategoryRemove(db *gsqlitehandler.SqliteDB, m *MainCategory) error {
	var err error
	var stmt *sql.Stmt

	// Check if it is not a system object
	if m.Status == ISSystem {
		return errors.New(errSystemObject)
	}

	// Set correct status (ISClose)
	if stmt, err = db.Handler.Prepare("UPDATE main_categories SET status=? WHERE id=?;"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(ISClose, m.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryList returns closure which generates a sequence of Main Category objects
func MainCategoryList(db *gsqlitehandler.SqliteDB, t *MainCategoryType, n string, s ItemStatus) (f func() *MainCategory, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if n == NotSetStringValue {
		n = noStringParamForSQL
	} else {
		n = "%" + n + "%"
	}
	var tId int
	if t == nil {
		tId = noIntParamForSQL
	} else {
		tId = t.Id
	}

	sqlQuery := "SELECT m.id, m.name, m.status, t.id, t.name, t.factor " +
		"FROM main_categories m INNER JOIN main_categories_types t ON m.type_id=t.id " +
		"WHERE (m.type_id=? OR ?=?) AND (m.name LIKE ? OR ?=?) AND (m.status=? or ?=?) ORDER BY t.id, m.name;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(tId, tId, noIntParamForSQL, n, n, noStringParamForSQL, s, s, ISUnset); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *MainCategory {
		if rows.Next() {
			m := MainCategoryNew()
			rows.Scan(&m.Id, &m.Name, &m.Status, &m.MType.Id, &m.MType.Name, &m.MType.Factor)
			return m
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}
