// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// MainCategoryStatusT describes the behaviour of categories and its descendants (transactions)
type MainCategoryType int

const (
	MCTUnknown  MainCategoryType = -1
	MCTUnset    MainCategoryType = 0
	MCTCost     MainCategoryType = 1
	MCTTransfer MainCategoryType = 2
	MCTIncome   MainCategoryType = 3
)

// String satisfies fmt.Stringer interface in order to get human readable names
func (mct MainCategoryType) String() string {
	var mctName string

	switch mct {
	case MCTUnknown:
		mctName = "unknown"
	case MCTUnset:
		mctName = "not set"
	case MCTIncome:
		mctName = "income"
	case MCTCost:
		mctName = "cost"
	case MCTTransfer:
		mctName = "transfer"
	}

	return mctName
}

// Factor returns coefficient for transaction values
func (mct MainCategoryType) Factor() float64 {
	var factor float64

	switch mct {
	case MCTIncome, MCTTransfer:
		factor = 1.0
	case MCTCost:
		factor = -1.0
	default:
		factor = 0.0
	}

	return factor
}

// MainCategory represents the basic object for main category
type MainCategory struct {
	Id     int64
	MType  MainCategoryType
	Name   string
	Status ItemStatus
}

// MainCategoryAdd adds new main category with type t and name n
func MainCategoryAdd(db *gsqlitehandler.SqliteDB, m *MainCategory) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("INSERT INTO main_categories VALUES (NULL, ?, ?, ?);"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType, m.Name, m.Status); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryForID returns pointer to the MainCategory for given id
func MainCategoryForID(db *gsqlitehandler.SqliteDB, i int) (m *MainCategory, err error) {
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("SELECT * FROM main_categories WHERE id=? AND status=?;"); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = new(MainCategory)
	if err = stmt.QueryRow(i, ISOpen).Scan(&m.Id, &m.MType, &m.Name, &m.Status); err != nil {
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
	if stmt, err = db.Handler.Prepare("SELECT * FROM main_categories WHERE name LIKE ? AND status=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = new(MainCategory)
	if rows, err = stmt.Query(n, ISOpen); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer rows.Close()

	var noOfMainCategories int
	for rows.Next() {
		noOfMainCategories++
		rows.Scan(&m.Id, &m.MType, &m.Name, &m.Status)
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

	if stmt, err = db.Handler.Prepare("UPDATE main_categories SET type=?, name=?, status=? WHERE id=?;"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType, m.Name, m.Status, m.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryRemove updates main category status with ISClose
func MainCategoryRemove(db *gsqlitehandler.SqliteDB, m *MainCategory) error {
	var err error
	var stmt *sql.Stmt

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
func MainCategoryList(db *gsqlitehandler.SqliteDB, t MainCategoryType, n string, s ItemStatus) (f func() *MainCategory, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if n == NotSetStringValue {
		n = noStringParamForSQL
	} else {
		n = "%" + n + "%"
	}

	if stmt, err = db.Handler.Prepare("SELECT id, type, name, status FROM main_categories WHERE (type=? OR ?=?) AND (name LIKE ? OR ?=?) AND (status=? or ?=?) ORDER BY type, name;"); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(t, t, MCTUnset, n, n, noStringParamForSQL, s, s, ISUnset); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *MainCategory {
		if rows.Next() {
			m := new(MainCategory)
			rows.Scan(&m.Id, &m.MType, &m.Name, &m.Status)
			return m
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}
