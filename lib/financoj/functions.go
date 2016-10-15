// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gprops"
	"github.com/zbroju/gsqlitehandler"
	"os"
	"path"
)

// GetConfigSettings returns contents of settings file
func GetConfigSettings() (dataFile string, currency string, err error) {
	// Read config file
	configSettings := gprops.New()
	configFile, err := os.Open(path.Join(os.Getenv("HOME"), configFile))
	if err == nil {
		err = configSettings.Load(configFile)
		if err != nil {
			return NotSetStringValue, NotSetStringValue, err
		}
	}
	configFile.Close()
	dataFile = configSettings.GetOrDefault(confDataFile, NotSetStringValue)
	currency = configSettings.GetOrDefault(confCurrency, NotSetStringValue)

	return dataFile, currency, nil
	//TODO: add test
}

// GetDataFileHandler returns new file handler for given path
func GetDataFileHandler(filePath string) *gsqlitehandler.SqliteDB {
	return gsqlitehandler.New(filePath, dataFileProperties)
	//TODO: add test
}

// CreateNewDataFile creates new data file for given data file handler
func CreateNewDataFile(db *gsqlitehandler.SqliteDB) error {
	// Create new file
	sqlCreateTables := "CREATE TABLE currencies (currency_from TEXT, currency_to TEXT, exchange_rate REAL, PRIMARY KEY (currency_from, currency_to));" +
		"CREATE TABLE accounts (id INTEGER PRIMARY KEY, name TEXT, description TEXT, institution TEXT, type INTEGER, currency TEXT, status INTEGER);" +
		"CREATE TABLE transactions (id INTEGER PRIMARY KEY, year INTEGER, month INTEGER, day INTEGER, account_id INTEGER, description TEXT, value REAL, category_id INTEGER);" +
		"CREATE TABLE budgets (year INTEGER, month INTEGER, category_id INTEGER, value REAL, currency TEXT, PRIMARY KEY (YEAR, MONTH, CATEGORY_ID));" +
		"CREATE TABLE categories (id INTEGER PRIMARY KEY, main_category_id INTEGER, name TEXT, status INTEGER);" +
		"CREATE TABLE main_categories (id INTEGER PRIMARY KEY, type INTEGER, name TEXT, status INTEGER);"

	err := db.CreateNew(sqlCreateTables)

	return err
	//TODO: add test
}

// CategoryAdd add new category with name n
func CategoryAdd(db *gsqlitehandler.SqliteDB, c *CategoryT) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("INSERT INTO categories VALUES (NULL, ?, ?, ?);"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(c.MainCategory.Id, c.Name, c.Status); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil

	//TODO: add test
}

// CategoryForID returns pointer to CategoryT for given id
func CategoryForID(db *gsqlitehandler.SqliteDB, i int) (c *CategoryT, err error) {
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("SELECT c.id, c.name, c.status, m.id, m.type, m.name, m.status FROM categories c INNER JOIN main_categories m ON c.main_category_id=m.id WHERE c.id=? AND c.status=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	c = CategoryNew()
	if err = stmt.QueryRow(i, ISOpen).Scan(&c.Id, &c.Name, &c.Status, &c.MainCategory.Id, &c.MainCategory.MType, &c.MainCategory.Name, &c.MainCategory.Status); err != nil {
		return nil, errors.New(errCategoryWithIDNone)
	}
	return c, nil
	//TODO: add test
}

// CategoryEdit updates category with new values for name, main category and status
// All three fields are updated, so make sure you pass old values in argument 'c'
func CategoryEdit(db *gsqlitehandler.SqliteDB, c *CategoryT) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("UPDATE categories SET main_category_id=?, name=?, status=? WHERE id=?;"); err != nil {
		errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(c.MainCategory.Id, c.Name, c.Status, c.Id); err != nil {
		errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// CategoryRemove updates given category status with ISClose
func CategoryRemove(db *gsqlitehandler.SqliteDB, c *CategoryT) error {
	var err error
	var stmt *sql.Stmt

	// Set correct status (ISClose
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
func CategoryList(db *gsqlitehandler.SqliteDB, m string, t MainCategoryTypeT, c string, s ItemStatus) (f func() *CategoryT, err error) {
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

	f = func() *CategoryT {
		if rows.Next() {
			c := CategoryNew()
			rows.Scan(&c.Id, &c.Name, &c.Status, &c.MainCategory.Id, &c.MainCategory.MType, &c.MainCategory.Name, &c.MainCategory.Status)
			return c
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

// MainCategoryAdd adds new main category with type t and name n
func MainCategoryAdd(db *gsqlitehandler.SqliteDB, m *MainCategoryT) error {
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

// MainCategoryForID returns pointer to MainCategoryT for given id
func MainCategoryForID(db *gsqlitehandler.SqliteDB, i int) (m *MainCategoryT, err error) {
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("SELECT * FROM main_categories WHERE id=? AND status=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = new(MainCategoryT)
	if err = stmt.QueryRow(i, ISOpen).Scan(&m.Id, &m.MType, &m.Name, &m.Status); err != nil {
		return m, errors.New(errMainCategoryWithIDNone)
	}

	return m, nil
	//TODO: add test
}

// MainCategoryForName returns pointer to MainCategoryT for given (part of) name
func MainCategoryForName(db *gsqlitehandler.SqliteDB, n string) (m *MainCategoryT, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	n = "%" + n + "%"
	if stmt, err = db.Handler.Prepare("SELECT * FROM main_categories WHERE name LIKE ? AND status=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = new(MainCategoryT)
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
		return nil, errors.New(errMainCategoryWithNameAmbiguous)
	}

	//TODO: add test
}

// MainCategoryEdit updates main category with new values for type, name and status
// Both type, name and status is updated, so make sure you pass old values in argument 'm'
func MainCategoryEdit(db *gsqlitehandler.SqliteDB, m *MainCategoryT) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("UPDATE main_categories SET type=?, name=?, status=? WHERE id=?;"); err != nil {
		errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType, m.Name, m.Status, m.Id); err != nil {
		errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryRemove updates main category status with ISClose
func MainCategoryRemove(db *gsqlitehandler.SqliteDB, m *MainCategoryT) error {
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
func MainCategoryList(db *gsqlitehandler.SqliteDB, t MainCategoryTypeT, n string, s ItemStatus) (f func() *MainCategoryT, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if n == NotSetStringValue {
		n = noParameterValueForSQL
	} else {
		n = "%" + n + "%"
	}

	if stmt, err = db.Handler.Prepare("SELECT id, type, name, status FROM main_categories WHERE (type=? OR ?=?) AND (name LIKE ? OR ?=?) AND (status=? or ?=?) ORDER BY type, name;"); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(t, t, MCTUnset, n, n, noParameterValueForSQL, s, s, ISUnset); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *MainCategoryT {
		if rows.Next() {
			m := new(MainCategoryT)
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

// CurrencyAdd add new currency exchange rate
func CurrencyAdd(db *gsqlitehandler.SqliteDB, c *CurrencyT) error {
	var err error
	var stmt *sql.Stmt
	//TODO: check if such currency exchange rate exists
	if stmt, err = db.Handler.Prepare("INSERT into currencies VALUES (?, ?, round(?,4));"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(c.CurrencyFrom, c.CurrencyTo, c.ExchangeRate); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// CurrencyList returns all currency exchange rates as closure
func CurrencyList(db *gsqlitehandler.SqliteDB) (f func() *CurrencyT, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if stmt, err = db.Handler.Prepare("SELECT currency_from, currency_to, exchange_rate FROM currencies ORDER BY currency_from, currency_to;"); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	if rows, err = stmt.Query(); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *CurrencyT {
		if rows.Next() {
			c := new(CurrencyT)
			rows.Scan(&c.CurrencyFrom, &c.CurrencyTo, &c.ExchangeRate)
			return c
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}
