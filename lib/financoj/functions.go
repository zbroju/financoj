// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

import (
	"errors"
	"fmt"
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
}

// GetDataFileHandler returns new file handler for given path
func GetDataFileHandler(filePath string) *gsqlitehandler.SqliteDB {
	return gsqlitehandler.New(filePath, dataFileProperties)
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
}

// MainCategoryAdd adds new main category with type (t) and name (n)
func MainCategoryAdd(db *gsqlitehandler.SqliteDB, t MainCategoryTypeT, n string) error {
	sqlAddType := fmt.Sprintf("INSERT INTO main_categories VALUES (NULL, %d,'%s', %d);", t, n, isOpen)
	if _, err := db.Handler.Exec(sqlAddType); err != nil {
		return errors.New(errWritingToFile)
	}

	//TODO: add to the database schema coefficient so that transactions are always positive

	return nil
}

// MainCategoryForID returns MainCategoryT for given id
func MainCategoryForID(db *gsqlitehandler.SqliteDB, i int) (m MainCategoryT, err error) {
	m = MainCategoryT{}
	sqlQuery := fmt.Sprintf("SELECT * FROM main_categories WHERE id=%d;", i)
	if err = db.Handler.QueryRow(sqlQuery).Scan(&m.Id, &m.MCType, &m.Name, &m.Status); err != nil {
		return m, errors.New(errNoMainCategoryWithID)
	}

	return m, err

}

// MainCategoryEdit updates main category with new values for type (t), name (n)
func MainCategoryEdit(db *gsqlitehandler.SqliteDB, m MainCategoryT) error {
	sqlQuery := fmt.Sprintf("UPDATE main_categories SET type=%d, name='%s' WHERE id=%d;", m.MCType, m.Name, m.Id)

	if _, err := db.Handler.Exec(sqlQuery); err != nil {
		errors.New(errWritingToFile)
	}

	return nil

}

// MainCategoryRemove updates main category status with isClose
func MainCategoryRemove(db *gsqlitehandler.SqliteDB, m MainCategoryT) error {
	// Set correct status (IS_Close)
	sqlQuery := fmt.Sprintf("UPDATE main_categories SET status=%d WHERE id=%d;", isClose, m.Id)

	if _, err := db.Handler.Exec(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
}
