// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"github.com/zbroju/gsqlitehandler"
)

// cmdInit creates a new sqlite file and tables for financoj
func cmdInit(c *cli.Context) error {
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
		"CREATE TABLE TRANSACTIONS (id INTEGER PRIMARY KEY, year INTEGER, month INTEGER, day INTEGER, account_id INTEGER, description TEXT, value REAL, category_id INTEGER);" +
		"CREATE TABLE BUDGETS (year INTEGER, month INTEGER, category_id INTEGER, value REAL, currency TEXT, PRIMARY KEY (YEAR, MONTH, CATEGORY_ID));" +
		"CREATE TABLE CATEGORIES (id INTEGER PRIMARY KEY, main_category_id INTEGER, name TEXT, status INTEGER);" +
		"CREATE TABLE MAIN_CATEGORIES (id INTEGER PRIMARY KEY, type INTEGER, name TEXT, status INTEGER);"

	df := gsqlitehandler.New(f, dataFileProperties)

	err := df.CreateNew(sqlCreateTables)
	if err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("created file %s\n", f)

	return nil
}

// cmdMainCategoryAdd adds new main category
func cmdMainCategoryAdd(c *cli.Context) error {
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
	var mcType MainCategoryTypeT
	if mcType, err = resolveMainCategoryType(c.String(optMainCategoryType)); err != nil {
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
	printUserMsg.Printf("added new main category: %s of type: %s\n", mcName, mcType)

	return nil
}

// ResolveMainCategoryType returns main category type for given string
func resolveMainCategoryType(m string) (mc MainCategoryTypeT, err error) {
	switch m {
	case "c", "cost", NotSetStringValue:
		mc = MCT_Cost
	case "i", "income":
		mc = MCT_Income
	case "t", "transfer":
		mc = MCT_Transfer
	default:
		err = errors.New("unknown type of main category")
		mc = MCT_NotSet
	}

	return mc, err
}
