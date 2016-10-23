// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// Category represents the basic object for category
type Category struct {
	Id     int64
	Main   *MainCategory
	Name   string
	Status ItemStatus
}

func CategoryNew() *Category {
	c := new(Category)
	c.Main = new(MainCategory)

	return c
}

// CategoryAdd add new category with name n
func CategoryAdd(db *gsqlitehandler.SqliteDB, c *Category) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("INSERT INTO categories VALUES (NULL, ?, ?, ?);"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(c.Main.Id, c.Name, c.Status); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil

	//TODO: add test
}

// CategoryForID returns pointer to CategoryT for given id
func CategoryForID(db *gsqlitehandler.SqliteDB, i int) (c *Category, err error) {
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("SELECT c.id, c.name, c.status, m.id, m.type, m.name, m.status FROM categories c INNER JOIN main_categories m ON c.main_category_id=m.id WHERE c.id=? AND c.status=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	c = CategoryNew()
	if err = stmt.QueryRow(i, ISOpen).Scan(&c.Id, &c.Name, &c.Status, &c.Main.Id, &c.Main.MType, &c.Main.Name, &c.Main.Status); err != nil {
		return nil, errors.New(errCategoryWithIDNone)
	}
	return c, nil
	//TODO: add test
}

// CategoryForName returns pointer to Category for given (part of) name
func CategoryForName(db *gsqlitehandler.SqliteDB, n string) (c *Category, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	n = "%" + n + "%"
	sqlQuery := "SELECT c.id, c.name, c.status, m.id, m.type, m.name, m.status FROM categories c INNER JOIN main_categories m ON c.main_category_id=m.id WHERE c.name LIKE ? AND c.status=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	c = CategoryNew()
	if rows, err = stmt.Query(n, ISOpen); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer rows.Close()

	var noOfCategories int
	for rows.Next() {
		noOfCategories++
		rows.Scan(&c.Id, &c.Name, &c.Status, &c.Main.Id, &c.Main.MType, &c.Main.Name, &c.Main.Status)
	}

	switch noOfCategories {
	case 0:
		return nil, errors.New(errCategoryWithNameNone)
	case 1:
		return c, nil
	default:
		return nil, errors.New(errCategoryWithNameAmbiguous)
	}

	//TODO: add test
}

// CategoryEdit updates category with new values for name, main category and status
// All three fields are updated, so make sure you pass old values in argument 'c'
func CategoryEdit(db *gsqlitehandler.SqliteDB, c *Category) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("UPDATE categories SET main_category_id=?, name=?, status=? WHERE id=?;"); err != nil {
		errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(c.Main.Id, c.Name, c.Status, c.Id); err != nil {
		errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// CategoryRemove updates given category status with ISClose
func CategoryRemove(db *gsqlitehandler.SqliteDB, c *Category) error {
	var err error
	var stmt *sql.Stmt

	// Set correct status (ISClose)
	if stmt, err = db.Handler.Prepare("UPDATE categories SET status=? WHERE id=?;"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(ISClose, c.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// CategoryList returns all categories from file as closure
func CategoryList(db *gsqlitehandler.SqliteDB, m string, t MainCategoryType, c string, s ItemStatus) (f func() *Category, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if m == NotSetStringValue {
		m = noParameterValueForSQL
	} else {
		m = "%" + m + "%"
	}
	if c == NotSetStringValue {
		c = noParameterValueForSQL
	} else {
		c = "%" + c + "%"
	}

	if stmt, err = db.Handler.Prepare("SELECT c.id, c.name, c.status, m.id, m.type, m.name,m.status FROM categories c INNER JOIN main_categories m on c.main_category_id=m.id WHERE (m.name LIKE ? OR ?=?) AND (m.type=? OR ?=?) AND (c.name LIKE ? OR ?=?) AND (c.status=? or ?=?) ORDER BY m.type, m.name, c.name;"); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	if rows, err = stmt.Query(m, m, noParameterValueForSQL, t, t, MCTUnset, c, c, noParameterValueForSQL, s, s, ISUnset); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *Category {
		if rows.Next() {
			c := CategoryNew()
			rows.Scan(&c.Id, &c.Name, &c.Status, &c.Main.Id, &c.Main.MType, &c.Main.Name, &c.Main.Status)
			return c
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

//FIXME: rebuild SELECT queries in all functions so that the conditions are in separated lines
