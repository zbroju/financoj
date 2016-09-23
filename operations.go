// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"github.com/urfave/cli"
	"github.com/zbroju/gsqlitehandler"
)

func cmdInit(c *cli.Context) error {
	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check the obligatory parameters and exit if missing
	if c.String("file") == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	// Create new file
	sqlCreateTables := "CREATE TABLE currencies (currency_from TEXT, currency_to TEXT, exchange_rate REAL, PRIMARY KEY (currency_from, currency_to));" +
		"CREATE TABLE accounts (id INTEGER PRIMARY KEY, name TEXT, description TEXT, institution TEXT, type INTEGER, currency TEXT, status INTEGER);" +
		"CREATE TABLE TRANSACTIONS (id INTEGER PRIMARY KEY, year INTEGER, month INTEGER, day INTEGER, account_id INTEGER, description TEXT, value REAL, category_id INTEGER);" +
		"CREATE TABLE BUDGETS (year INTEGER, month INTEGER, category_id INTEGER, value REAL, currency TEXT, PRIMARY KEY (YEAR, MONTH, CATEGORY_ID));" +
		"CREATE TABLE CATEGORIES (id INTEGER PRIMARY KEY, main_category_id INTEGER, name TEXT, status INTEGER);" +
		"CREATE TABLE MAIN_CATEGORIES (id INTEGER PRIMARY KEY, type INTEGER, name TEXT, status INTEGER);"

	f := gsqlitehandler.New(c.String("file"), dataFileProperties)

	err := f.CreateNew(sqlCreateTables)
	if err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("created file %s.\n", c.String("file"))

	return nil
}
