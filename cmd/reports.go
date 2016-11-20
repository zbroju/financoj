// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/lib"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

func RepAccountBalance(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var bDate time.Time
	if td := c.String(OptDate); td != NotSetStringValue {
		if bDate, err = time.Parse(DateFormat, td); err != nil {
			printError.Fatalln(err)
		}
	} else {
		bDate = time.Now()
	}

	// Build formatting strings
	var getNextEntry func() *AccountBalanceEntry
	if getNextEntry, err = ReportAccountBalance(fh, bDate); err != nil {
		printError.Fatalln(err)
	}
	lA := utf8.RuneCountInString(HAName)
	lV := utf8.RuneCountInString(HTValue)
	lC := utf8.RuneCountInString(HACurrency)
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lA = MaxLen(e.Account.Name, lA)
		lV = MaxLen(strconv.FormatFloat(e.Value, 'f', 2, 64), lV)
		lC = MaxLen(e.Account.Currency, lC)
	}
	lineD := LineFor(NotSetStringValue, DFSForText(lA), DFSForValue(lV), DFSForText(lC))

	// Print report
	fmt.Fprintf(os.Stdout, "Accounts balance on %s:\n", bDate.Format(DateFormat))

	if getNextEntry, err = ReportAccountBalance(fh, bDate); err != nil {
		printError.Fatalln(err)
	}
	var currentType string
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.Account.AType.String() {
			currentType = e.Account.AType.String()
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)
		}
		fmt.Fprintf(os.Stdout, lineD, e.Account.Name, e.Value, e.Account.Currency)
	}

	return nil
}
