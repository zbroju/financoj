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

func RepBudgetCategories(c *cli.Context) error {
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
	var p *BPeriod
	if ps := c.String(OptPeriod); ps != NotSetStringValue {
		if p, err = BPeriodParseYOrYM(ps); err != nil {
			printError.Fatalln(err)
		}
	} else {
		if p, err = BPeriodCurrent(); err != nil {
			printError.Fatalln(err)
		}
	}
	currency := c.String(OptCurrency)
	if currency == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Build formatting strings
	var getNextEntry func() *BudgetCategoriesEntry
	if getNextEntry, err = ReportBudgetCategories(fh, p, currency); err != nil {
		printError.Fatalln(err)
	}
	lMT := utf8.RuneCountInString(HMCType)
	lMN := utf8.RuneCountInString(HMCName)
	lCN := utf8.RuneCountInString(HCName)
	lBL := utf8.RuneCountInString(HBLimit)
	lTV := utf8.RuneCountInString(HTValue)
	lD := utf8.RuneCountInString(HBDifference)
	lCur := utf8.RuneCountInString(HBCurrency)
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lMT = MaxLen(e.Category.Main.MType.Name, lMT)
		lMN = MaxLen(e.Category.Main.Name, lMN)
		lCN = MaxLen(e.Category.Name, lCN)
		lBL = MaxLen(strconv.FormatFloat(e.Limit, 'f', 2, 64), lBL)
		lTV = MaxLen(strconv.FormatFloat(e.Actual, 'f', 2, 64), lTV)
		lD = MaxLen(strconv.FormatFloat(e.Difference, 'f', 2, 64), lD)
		lCur = MaxLen(e.Currency, lCur)
	}
	lineH := LineFor(NotSetStringValue, HFSForText(lMN), HFSForText(lCN), HFSForNumeric(lBL), HFSForText(lCur), HFSForNumeric(lTV), HFSForText(lCur), HFSForNumeric(lD), HFSForText(lCur))
	lineD := LineFor(NotSetStringValue, DFSForText(lMN), DFSForText(lCN), DFSForValue(lBL), DFSForText(lCur), DFSForValue(lTV), DFSForText(lCur), DFSForValue(lD), DFSForText(lCur))
	lineS := LineFor(DFSForText(2*utf8.RuneCountInString(FSSeparator)+lMN+lCN), DFSForValue(lBL), DFSForText(lCur), DFSForValue(lTV), DFSForText(lCur), DFSForValue(lD), DFSForText(lCur))

	// Print report
	fmt.Fprintf(os.Stdout, "Budget report for %s:\n", p)

	if getNextEntry, err = ReportBudgetCategories(fh, p, currency); err != nil {
		printError.Fatalln(err)
	}
	var currentType string
	var subtotalLimit, subtotalValue, subtotalDifference, totalLimit, totalValue, totalDifference float64
	beginning := true
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.Category.Main.MType.Name {
			if !beginning {
				fmt.Fprintf(os.Stdout, lineS, currentType, subtotalLimit, e.Currency, subtotalValue, e.Currency, subtotalDifference, e.Currency)
			}
			currentType = e.Category.Main.MType.Name
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)
			fmt.Fprintf(os.Stdout, lineH, HMCName, HCName, HBLimit, HBCurrency, HTValue, HBCurrency, HBDifference, HBCurrency)

			beginning = false
			subtotalLimit = NotSetFloatValue
			subtotalValue = NotSetFloatValue
			subtotalDifference = NotSetFloatValue
		}
		fmt.Fprintf(os.Stdout, lineD, e.Category.Main.Name, e.Category.Name, e.Limit, e.Currency, e.Actual, e.Currency, e.Difference, e.Currency)
		subtotalLimit += e.Limit
		subtotalValue += e.Actual
		subtotalDifference += e.Difference
		totalLimit += e.Limit
		totalValue += e.Actual
		totalDifference += e.Difference
	}
	fmt.Fprintf(os.Stdout, lineS, currentType, subtotalLimit, currency, subtotalValue, currency, subtotalDifference, currency)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineS, "TOTAL", totalLimit, currency, totalValue, currency, totalDifference, currency)

	return nil
}
