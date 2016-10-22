// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"github.com/zbroju/gprops"
	"github.com/zbroju/gsqlitehandler"
	"os"
	"path"
)

// Config file settings
const configFile = ".financojrc"
const (
	confDataFile = "DATA_FILE"
	confCurrency = "DEFAULT_CURRENCY"
)

// DB Properties
var dataFileProperties = map[string]string{
	"applicationName": "financoj",
	"databaseVersion": "2.0",
}

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
		"CREATE TABLE accounts (id INTEGER PRIMARY KEY, name TEXT, description TEXT, institution TEXT, currency TEXT, type INTEGER, status INTEGER);" +
		"CREATE TABLE transactions (id INTEGER PRIMARY KEY, year INTEGER, month INTEGER, day INTEGER, account_id INTEGER, description TEXT, value REAL, category_id INTEGER);" +
		"CREATE TABLE budgets (year INTEGER, month INTEGER, category_id INTEGER, value REAL, currency TEXT, PRIMARY KEY (YEAR, MONTH, CATEGORY_ID));" +
		"CREATE TABLE categories (id INTEGER PRIMARY KEY, main_category_id INTEGER, name TEXT, status INTEGER);" +
		"CREATE TABLE main_categories (id INTEGER PRIMARY KEY, type INTEGER, name TEXT, status INTEGER);"

	err := db.CreateNew(sqlCreateTables)

	return err
	//TODO: add test
}
