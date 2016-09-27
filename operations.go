// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/zbroju/gsqlitehandler"
)

// CreateNewDataFile creates a new sqlite file and tables for financoj
func CreateNewDataFile(c *cli.Context) error {
	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check the obligatory parameters and exit if missing
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	// Create new file
	sqlCreateTables := "CREATE TABLE currencies (currency_from TEXT, currency_to TEXT, exchange_rate REAL, PRIMARY KEY (currency_from, currency_to));" +
		"CREATE TABLE accounts (id INTEGER PRIMARY KEY, name TEXT, description TEXT, institution TEXT, type INTEGER, currency TEXT, status INTEGER);" +
		"CREATE TABLE transactions (id INTEGER PRIMARY KEY, year INTEGER, month INTEGER, day INTEGER, account_id INTEGER, description TEXT, value REAL, category_id INTEGER);" +
		"CREATE TABLE budgets (year INTEGER, month INTEGER, category_id INTEGER, value REAL, currency TEXT, PRIMARY KEY (YEAR, MONTH, CATEGORY_ID));" +
		"CREATE TABLE categories (id INTEGER PRIMARY KEY, main_category_id INTEGER, name TEXT, status INTEGER);" +
		"CREATE TABLE main_categories (id INTEGER PRIMARY KEY, type INTEGER, name TEXT, status INTEGER);"

	df := gsqlitehandler.New(f, dataFileProperties)

	err := df.CreateNew(sqlCreateTables)
	if err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("created file %s\n", f)

	return nil
}

// MainCategoryAdd adds new main category
func MainCategoryAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check obligatory flags (file, name)
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	mcName := c.String(objMainCategory)
	if mcName == NotSetStringValue {
		printError.Fatalln(errMissingMainCategory)
	}
	mcType := mainCategoryTypeForString(c.String(optMainCategoryType))
	if mcType == MCT_Unknown {
		printError.Fatalln(errIncorrectMainCategoryType)
	}

	// Open data file
	df := gsqlitehandler.New(f, dataFileProperties)
	if err = df.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer df.Close()

	// Add new type
	sqlAddType := fmt.Sprintf("INSERT INTO main_categories VALUES (NULL, %d,'%s', %d);", mcType, mcName, IS_Open)
	if _, err = df.Handler.Exec(sqlAddType); err != nil {
		printError.Fatalln(errWritingToFile)
	}

	// Show summary
	printUserMsg.Printf("added new main category: %s (type: %s)\n", mcName, mcType)

	return nil
}

// MainCategoryEdit updates main category with new values
func MainCategoryEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	id := c.Int(optID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file
	df := gsqlitehandler.New(f, dataFileProperties)
	if err = df.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer df.Close()

	// Edit main category
	sqlQuery := fmt.Sprintf("BEGIN TRANSACTION;")
	if t := c.String(optMainCategoryType); t != NotSetStringValue {
		mct := mainCategoryTypeForString(t)
		if mct == MCT_Unknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
		sqlQuery = sqlQuery + fmt.Sprintf("UPDATE main_categories SET type=%d WHERE id=%d;", mct, id)
	}
	if n := c.String(objMainCategory); n != NotSetStringValue {
		sqlQuery = sqlQuery + fmt.Sprintf("UPDATE main_categories SET name='%s' WHERE id=%d;", n, id)
	}
	sqlQuery = sqlQuery + "COMMIT;"

	r, err := df.Handler.Exec(sqlQuery)
	if err != nil {
		printError.Fatalln(errWritingToFile)
	}
	if i, _ := r.RowsAffected(); i == 0 {
		printError.Fatalln(errNoMainCategoryWithID)
	}

	// Show summary
	printUserMsg.Printf("changed details of main category with id = %d\n", id)

	return nil
}
