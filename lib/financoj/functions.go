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
func GetConfigSettings() (dataFile string, err error) {
	// Read config file
	configSettings := gprops.New()
	configFile, err := os.Open(path.Join(os.Getenv("HOME"), configFile))
	if err == nil {
		err = configSettings.Load(configFile)
		if err != nil {
			return NotSetStringValue, err
		}
	}
	configFile.Close()
	dataFile = configSettings.GetOrDefault(confDataFile, NotSetStringValue)

	return dataFile, nil
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

// MainCategoryAdd adds new main category with type (t) and name (n)
func MainCategoryAdd(db *gsqlitehandler.SqliteDB, m MainCategoryT) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("INSERT INTO main_categories VALUES (NULL, ?, ?, ?);"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType, m.Name, m.Status); err != nil {
		return errors.New(errWritingToFile)
	}

	//TODO: add to the database schema coefficient so that transactions are always positive

	return nil
	//TODO: add test
}

// MainCategoryForID returns MainCategoryT for given id
func MainCategoryForID(db *gsqlitehandler.SqliteDB, i int) (m MainCategoryT, err error) {
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("SELECT * FROM main_categories WHERE id=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	m = MainCategoryT{}
	if err = stmt.QueryRow(i).Scan(&m.Id, &m.MType, &m.Name, &m.Status); err != nil {
		return m, errors.New(errNoMainCategoryWithID)
	}

	return m, nil
	//TODO: add test
}

// MainCategoryEdit updates main category with new values for type (t), name (n)
// Both type and name is updated, so make sure you pass old values in argument 'm'
func MainCategoryEdit(db *gsqlitehandler.SqliteDB, m MainCategoryT) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("UPDATE main_categories SET type=?, name=? WHERE id=?;"); err != nil {
		errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(m.MType, m.Name, m.Id); err != nil {
		errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// MainCategoryRemove updates main category status with isClose
func MainCategoryRemove(db *gsqlitehandler.SqliteDB, m MainCategoryT) error {
	var err error
	var stmt *sql.Stmt

	// Set correct status (IS_Close)
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
func MainCategoryList(db *gsqlitehandler.SqliteDB, t MainCategoryTypeT, n string) (f func() *MainCategoryT, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	const notSetName = "notSetName"
	if n == NotSetStringValue {
		n = notSetName
	} else {
		n = "%" + n + "%"
	}

	if stmt, err = db.Handler.Prepare("SELECT id, type, name, status FROM main_categories WHERE (type=? OR ?=?) AND (name LIKE ? OR ?=?) AND (status=?);"); err != nil {
		errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(t, t, MCTUnset, n, n, notSetName, ISOpen); err != nil {
		errors.New(errReadingFromFile)
	}

	f = func() *MainCategoryT {
		for rows.Next() {
			m := new(MainCategoryT)
			rows.Scan(&m.Id, &m.MType, &m.Name, &m.Status)
			return m
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
}
