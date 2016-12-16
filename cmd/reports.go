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
	"strings"
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
	var getNextEntry func() *AccountBalanceReportEntry
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

func RepTransactionBalance(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var df time.Time
	if ds := c.String(OptDateFrom); ds != NotSetStringValue {
		if df, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		df = time.Time{}
	}
	var dt time.Time
	if ds := c.String(OptDateTo); ds != NotSetStringValue {
		if dt, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dt = time.Time{}
	}
	var a *Account
	if as := c.String(ObjAccount); as != NotSetStringValue {
		if a, err = AccountForName(fh, as); err != nil {
			printError.Fatalln(err)
		}
	}
	var cat *Category
	if cs := c.String(ObjCategory); cs != NotSetStringValue {
		if cat, err = CategoryForName(fh, cs); err != nil {
			printError.Fatalln(err)
		}
	}
	var mcat *MainCategory
	if ms := c.String(ObjMainCategory); ms != NotSetStringValue {
		if mcat, err = MainCategoryForName(fh, ms); err != nil {
			printError.Fatalln(err)
		}
	}

	// Build formatting strings
	var getNextEntry func() *TransactionBalanceReportEntry
	if getNextEntry, err = ReportTransactionBalance(fh, cur, df, dt, a, cat, mcat); err != nil {
		printError.Fatalln(err)
	}
	lD := utf8.RuneCountInString(HTDate)
	lMC := utf8.RuneCountInString(HMCName)
	lC := utf8.RuneCountInString(HCName)
	lA := utf8.RuneCountInString(HAName)
	lV := utf8.RuneCountInString(HTValue)
	lDesc := utf8.RuneCountInString(HTDescription)
	var sumValue float64
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lD = MaxLen(e.Transaction.Date.Format(DateFormat), lD)
		lMC = MaxLen(e.Transaction.Category.Main.Name, lMC)
		lC = MaxLen(e.Transaction.Category.Name, lC)
		lA = MaxLen(e.Transaction.Account.Name, lA)
		sumValue += e.Balance
		lV = MaxLen(strconv.FormatFloat(sumValue, 'f', 2, 64), lV)
		lDesc = MaxLen(e.Transaction.Description, lDesc)
	}
	lineH := LineFor(HFSForText(lD), HFSForText(lMC), HFSForText(lC), HFSForText(lA), HFSForNumeric(lV), HFSForText(lDesc))
	lineD := LineFor(DFSForText(lD), DFSForText(lMC), DFSForText(lC), DFSForText(lA), DFSForValue(lV), DFSForText(lDesc))

	// Print report
	fmt.Fprintf(os.Stdout, "Transactions balance (in %s):\n", strings.ToUpper(cur))
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineH, HTDate, HMCName, HCName, HAName, HTValue, HTDescription)
	if getNextEntry, err = ReportTransactionBalance(fh, cur, df, dt, a, cat, mcat); err != nil {
		printError.Fatalln(err)
	}
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		fmt.Fprintf(os.Stdout, lineD, e.Transaction.Date.Format(DateFormat), e.Transaction.Category.Main.Name, e.Transaction.Category.Name, e.Transaction.Account.Name, e.Balance, e.Transaction.Description)
	}
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineD, "Total:", NotSetStringValue, NotSetStringValue, NotSetStringValue, sumValue, strings.ToUpper(cur))

	return nil
}

func RepCategoryBalance(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var df time.Time
	if ds := c.String(OptDateFrom); ds != NotSetStringValue {
		if df, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		df = time.Time{}
	}
	var dt time.Time
	if ds := c.String(OptDateTo); ds != NotSetStringValue {
		if dt, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dt = time.Time{}
	}
	var a *Account
	if as := c.String(ObjAccount); as != NotSetStringValue {
		if a, err = AccountForName(fh, as); err != nil {
			printError.Fatalln(err)
		}
	}
	var cat *Category
	if cs := c.String(ObjCategory); cs != NotSetStringValue {
		if cat, err = CategoryForName(fh, cs); err != nil {
			printError.Fatalln(err)
		}
	}
	var mcat *MainCategory
	if ms := c.String(ObjMainCategory); ms != NotSetStringValue {
		if mcat, err = MainCategoryForName(fh, ms); err != nil {
			printError.Fatalln(err)
		}
	}

	// Build formatting strings
	var getNextEntry func() *CategoryBalanceReportEntry
	if getNextEntry, err = ReportCategoryBalance(fh, cur, df, dt, a, cat, mcat); err != nil {
		printError.Fatalln(err)
	}
	lMT := utf8.RuneCountInString(HMCType)
	lM := utf8.RuneCountInString(HMCName)
	lC := utf8.RuneCountInString(HCName)
	lV := utf8.RuneCountInString(HTValue)
	var sumValue float64
	var currentType string
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lMT = MaxLen(e.Category.Main.MType.Name, lMT)
		lM = MaxLen(e.Category.Main.Name, lM)
		lC = MaxLen(e.Category.Name, lC)
		if currentType != e.Category.Main.MType.Name {
			lV = MaxLen(strconv.FormatFloat(sumValue, 'f', 2, 64), lV)
			sumValue = 0
			currentType = e.Category.Main.MType.Name
		}
		sumValue += e.Balance
	}
	lV = MaxLen(strconv.FormatFloat(sumValue, 'f', 2, 64), lV)

	lineH := LineFor(NotSetStringValue, HFSForText(lM), HFSForText(lC), HFSForNumeric(lV))
	lineD := LineFor(NotSetStringValue, DFSForText(lM), DFSForText(lC), DFSForValue(lV))
	lineS := LineFor(DFSForText(2*utf8.RuneCountInString(FSSeparator)+lM+lC), DFSForValue(lV))

	// Print report
	fmt.Fprintf(os.Stdout, "Categories balance (in %s):\n", strings.ToUpper(cur))

	if getNextEntry, err = ReportCategoryBalance(fh, cur, df, dt, a, cat, mcat); err != nil {
		printError.Fatalln(err)
	}
	currentType = NotSetStringValue
	var subtotalValue, totalValue float64
	beginning := true
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.Category.Main.MType.Name {
			if !beginning {
				fmt.Fprintf(os.Stdout, lineS, currentType, subtotalValue)
			}
			currentType = e.Category.Main.MType.Name
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)
			fmt.Fprintf(os.Stdout, lineH, HMCName, HCName, HTValue)

			beginning = false
			subtotalValue = NotSetFloatValue
		}
		fmt.Fprintf(os.Stdout, lineD, e.Category.Main.Name, e.Category.Name, e.Balance)
		subtotalValue += e.Balance
		totalValue += e.Balance
	}
	fmt.Fprintf(os.Stdout, lineS, currentType, subtotalValue)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineS, "TOTAL", totalValue)

	return nil
}

func RepCategoryBalanceMonthly(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	cs := c.String(ObjCategory)
	if cs == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var cat *Category
	if cat, err = CategoryForName(fh, cs); err != nil {
		printError.Fatalln(err)
	}
	var df time.Time
	if ds := c.String(OptDateFrom); ds != NotSetStringValue {
		if df, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		df = time.Time{}
	}
	var dt time.Time
	if ds := c.String(OptDateTo); ds != NotSetStringValue {
		if dt, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dt = time.Time{}
	}

	// Build formatting strings
	var getNextEntry func() *BalanceMonthlyReportEntry
	if getNextEntry, err = ReportCategoriesBalanceMonthly(fh, cur, cat, df, dt); err != nil {
		printError.Fatalln(err)
	}
	lP := utf8.RuneCountInString(HBPeriod)
	lV := utf8.RuneCountInString(HTValue)
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lP = MaxLen(e.Period.String(), lP)
		lV = MaxLen(strconv.FormatFloat(e.Value, 'f', 2, 64), lV)
	}
	lineH := LineFor(HFSForNumeric(lP), HFSForNumeric(lV))
	lineD := LineFor(DFSForText(lP), DFSForValue(lV))

	// Print report
	fmt.Fprintf(os.Stdout, "Category '%s' balance monthly (in %s):\n\n", cat.Name, strings.ToUpper(cur))
	fmt.Fprintf(os.Stdout, lineH, HBPeriod, HTValue)
	if getNextEntry, err = ReportCategoriesBalanceMonthly(fh, cur, cat, df, dt); err != nil {
		printError.Fatalln(err)
	}
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		fmt.Fprintf(os.Stdout, lineD, e.Period, e.Value)
	}

	return nil
}

func RepMainCategoryBalance(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var df time.Time
	if ds := c.String(OptDateFrom); ds != NotSetStringValue {
		if df, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		df = time.Time{}
	}
	var dt time.Time
	if ds := c.String(OptDateTo); ds != NotSetStringValue {
		if dt, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dt = time.Time{}
	}
	var a *Account
	if as := c.String(ObjAccount); as != NotSetStringValue {
		if a, err = AccountForName(fh, as); err != nil {
			printError.Fatalln(err)
		}
	}
	var mcat *MainCategory
	if ms := c.String(ObjMainCategory); ms != NotSetStringValue {
		if mcat, err = MainCategoryForName(fh, ms); err != nil {
			printError.Fatalln(err)
		}
	}

	// Build formatting strings
	var getNextEntry func() *MainCategoryBalanceReportEntry
	if getNextEntry, err = ReportMainCategoryBalance(fh, cur, df, dt, a, mcat); err != nil {
		printError.Fatalln(err)
	}
	lMT := utf8.RuneCountInString(HMCType)
	lM := utf8.RuneCountInString(HMCName)
	lV := utf8.RuneCountInString(HTValue)
	var subtotalValue, totalValue float64
	var currentType string
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lMT = MaxLen(e.MainCategory.MType.Name, lMT)
		lM = MaxLen(e.MainCategory.Name, lM)
		lV = MaxLen(strconv.FormatFloat(e.Balance, 'f', 2, 64), lV)
		if currentType != e.MainCategory.MType.Name {
			lV = MaxLen(strconv.FormatFloat(subtotalValue, 'f', 2, 64), lV)
			subtotalValue = 0
			currentType = e.MainCategory.MType.Name
		}
		subtotalValue += e.Balance
		totalValue += e.Balance
	}
	lV = MaxLen(strconv.FormatFloat(totalValue, 'f', 2, 64), lV)

	lineH := LineFor(NotSetStringValue, HFSForText(lM), HFSForNumeric(lV))
	lineD := LineFor(NotSetStringValue, DFSForText(lM), DFSForValue(lV))
	lineS := LineFor(DFSForText(utf8.RuneCountInString(FSSeparator)+lM), DFSForValue(lV))

	// Print report
	fmt.Fprintf(os.Stdout, "Categories balance (in %s):\n", strings.ToUpper(cur))

	if getNextEntry, err = ReportMainCategoryBalance(fh, cur, df, dt, a, mcat); err != nil {
		printError.Fatalln(err)
	}
	currentType = NotSetStringValue
	subtotalValue, totalValue = NotSetFloatValue, NotSetFloatValue
	beginning := true
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.MainCategory.MType.Name {
			if !beginning {
				fmt.Fprintf(os.Stdout, lineS, currentType, subtotalValue)
			}
			currentType = e.MainCategory.MType.Name
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)
			fmt.Fprintf(os.Stdout, lineH, HMCName, HTValue)

			beginning = false
			subtotalValue = NotSetFloatValue
		}
		fmt.Fprintf(os.Stdout, lineD, e.MainCategory.Name, e.Balance)
		subtotalValue += e.Balance
		totalValue += e.Balance
	}
	fmt.Fprintf(os.Stdout, lineS, currentType, subtotalValue)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineS, "TOTAL", totalValue)

	return nil
}

func RepAssetsSummary(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var onDate time.Time
	if ds := c.String(OptDate); ds != NotSetStringValue {
		if onDate, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		onDate = time.Now()
	}

	// Build formatting strings
	var getNextEntry func() *AssetsSummaryReportEntry
	if getNextEntry, err = ReportAssetsSummary(fh, cur, onDate); err != nil {
		printError.Fatalln(err)
	}
	lA := utf8.RuneCountInString(HAName)
	lV := utf8.RuneCountInString(HTValue)
	var subtotalValue, totalValue float64
	var currentType string
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lA = MaxLen(e.Account.Name, lA)
		lV = MaxLen(strconv.FormatFloat(e.Balance, 'f', 2, 64), lV)
		if currentType != e.Account.AType.String() {
			lV = MaxLen(strconv.FormatFloat(subtotalValue, 'f', 2, 64), lV)
			subtotalValue = 0
			currentType = e.Account.AType.String()
		}
		subtotalValue += e.Balance
		totalValue += e.Balance
	}
	lV = MaxLen(strconv.FormatFloat(totalValue, 'f', 2, 64), lV)
	lineD := LineFor(NotSetStringValue, DFSForText(lA), DFSForValue(lV))
	lineS := LineFor(DFSForText(utf8.RuneCountInString(FSSeparator)+lA), DFSForValue(lV))

	// Print report
	fmt.Fprintf(os.Stdout, "Assets summary on %s (in %s):\n", onDate.Format(DateFormat), strings.ToUpper(cur))

	if getNextEntry, err = ReportAssetsSummary(fh, cur, onDate); err != nil {
		printError.Fatalln(err)
	}
	currentType = NotSetStringValue
	subtotalValue, totalValue = NotSetFloatValue, NotSetFloatValue
	beginning := true
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.Account.AType.String() {
			if !beginning {
				fmt.Fprintf(os.Stdout, lineS, currentType, subtotalValue)
			}
			currentType = e.Account.AType.String()
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)

			beginning = false
			subtotalValue = NotSetFloatValue
		}
		fmt.Fprintf(os.Stdout, lineD, e.Account.Name, e.Balance)
		subtotalValue += e.Balance
		totalValue += e.Balance
	}
	fmt.Fprintf(os.Stdout, lineS, currentType, subtotalValue)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineS, "TOTAL", totalValue)

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
	var getNextEntry func() *BudgetCategoriesReportEntry
	if getNextEntry, err = ReportBudgetCategories(fh, p, currency); err != nil {
		printError.Fatalln(err)
	}
	lMT := utf8.RuneCountInString(HMCType)
	lMN := utf8.RuneCountInString(HMCName)
	lCN := utf8.RuneCountInString(HCName)
	lBL := utf8.RuneCountInString(HBLimit)
	lTV := utf8.RuneCountInString(HTValue)
	lD := utf8.RuneCountInString(HBDifference)
	var sumLimit, sumActual, sumDifference float64
	var currentType string
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lMT = MaxLen(e.Category.Main.MType.Name, lMT)
		lMN = MaxLen(e.Category.Main.Name, lMN)
		lCN = MaxLen(e.Category.Name, lCN)
		if currentType != e.Category.Main.MType.Name {
			lBL = MaxLen(strconv.FormatFloat(sumLimit, 'f', 2, 64), lBL)
			lTV = MaxLen(strconv.FormatFloat(sumActual, 'f', 2, 64), lTV)
			lD = MaxLen(strconv.FormatFloat(sumDifference, 'f', 2, 64), lD)
			sumLimit = 0
			sumActual = 0
			sumDifference = 0
			currentType = e.Category.Main.MType.Name
		}
		sumLimit += e.Limit
		sumActual += e.Actual
		sumDifference += e.Difference
	}
	lBL = MaxLen(strconv.FormatFloat(sumLimit, 'f', 2, 64), lBL)
	lTV = MaxLen(strconv.FormatFloat(sumActual, 'f', 2, 64), lTV)
	lD = MaxLen(strconv.FormatFloat(sumDifference, 'f', 2, 64), lD)

	lineH := LineFor(NotSetStringValue, HFSForText(lMN), HFSForText(lCN), HFSForNumeric(lBL), HFSForNumeric(lTV), HFSForNumeric(lD))
	lineD := LineFor(NotSetStringValue, DFSForText(lMN), DFSForText(lCN), DFSForValue(lBL), DFSForValue(lTV), DFSForValue(lD))
	lineS := LineFor(DFSForText(2*utf8.RuneCountInString(FSSeparator)+lMN+lCN), DFSForValue(lBL), DFSForValue(lTV), DFSForValue(lD))

	// Print report
	fmt.Fprintf(os.Stdout, "Budget report for %s (in %s):\n", p, strings.ToUpper(currency))

	if getNextEntry, err = ReportBudgetCategories(fh, p, currency); err != nil {
		printError.Fatalln(err)
	}
	currentType = NotSetStringValue
	var subtotalLimit, subtotalValue, subtotalDifference, totalLimit, totalValue, totalDifference float64
	beginning := true
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.Category.Main.MType.Name {
			if !beginning {
				fmt.Fprintf(os.Stdout, lineS, currentType, subtotalLimit, subtotalValue, subtotalDifference)
			}
			currentType = e.Category.Main.MType.Name
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)
			fmt.Fprintf(os.Stdout, lineH, HMCName, HCName, HBLimit, HTValue, HBDifference)

			beginning = false
			subtotalLimit = NotSetFloatValue
			subtotalValue = NotSetFloatValue
			subtotalDifference = NotSetFloatValue
		}
		fmt.Fprintf(os.Stdout, lineD, e.Category.Main.Name, e.Category.Name, e.Limit, e.Actual, e.Difference)
		subtotalLimit += e.Limit
		subtotalValue += e.Actual
		subtotalDifference += e.Difference
		totalLimit += e.Limit
		totalValue += e.Actual
		totalDifference += e.Difference
	}
	fmt.Fprintf(os.Stdout, lineS, currentType, subtotalLimit, subtotalValue, subtotalDifference)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineS, "TOTAL", totalLimit, totalValue, totalDifference)

	return nil
}

func RepBudgetMainCategories(c *cli.Context) error {
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
	var getNextEntry func() *BudgetMainCategoryReportEntry
	if getNextEntry, err = ReportBudgetMainCategories(fh, p, currency); err != nil {
		printError.Fatalln(err)
	}
	lMT := utf8.RuneCountInString(HMCType)
	lMN := utf8.RuneCountInString(HMCName)
	lBL := utf8.RuneCountInString(HBLimit)
	lTV := utf8.RuneCountInString(HTValue)
	lD := utf8.RuneCountInString(HBDifference)
	var sumLimit, sumActual, sumDifference float64
	var currentType string
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lMT = MaxLen(e.MainCategory.MType.Name, lMT)
		lMN = MaxLen(e.MainCategory.Name, lMN)
		if currentType != e.MainCategory.MType.Name {
			lBL = MaxLen(strconv.FormatFloat(sumLimit, 'f', 2, 64), lBL)
			lTV = MaxLen(strconv.FormatFloat(sumActual, 'f', 2, 64), lTV)
			lD = MaxLen(strconv.FormatFloat(sumDifference, 'f', 2, 64), lD)
			sumLimit = 0
			sumActual = 0
			sumDifference = 0
			currentType = e.MainCategory.MType.Name
		}
		sumLimit += e.Limit
		sumActual += e.Actual
		sumDifference += e.Difference
	}
	lBL = MaxLen(strconv.FormatFloat(sumLimit, 'f', 2, 64), lBL)
	lTV = MaxLen(strconv.FormatFloat(sumActual, 'f', 2, 64), lTV)
	lD = MaxLen(strconv.FormatFloat(sumDifference, 'f', 2, 64), lD)

	lineH := LineFor(NotSetStringValue, HFSForText(lMN), HFSForNumeric(lBL), HFSForNumeric(lTV), HFSForNumeric(lD))
	lineD := LineFor(NotSetStringValue, DFSForText(lMN), DFSForValue(lBL), DFSForValue(lTV), DFSForValue(lD))
	lineS := LineFor(DFSForText(utf8.RuneCountInString(FSSeparator)+lMN), DFSForValue(lBL), DFSForValue(lTV), DFSForValue(lD))

	// Print report
	fmt.Fprintf(os.Stdout, "Budget report for %s (in %s):\n", p, strings.ToUpper(currency))
	if getNextEntry, err = ReportBudgetMainCategories(fh, p, currency); err != nil {
		printError.Fatalln(err)
	}
	currentType = NotSetStringValue
	var subtotalLimit, subtotalValue, subtotalDifference, totalLimit, totalValue, totalDifference float64
	beginning := true
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		if currentType != e.MainCategory.MType.Name {
			if !beginning {
				fmt.Fprintf(os.Stdout, lineS, currentType, subtotalLimit, subtotalValue, subtotalDifference)
			}
			currentType = e.MainCategory.MType.Name
			fmt.Fprintf(os.Stdout, "\n%s\n", currentType)
			fmt.Fprintf(os.Stdout, lineH, HMCName, HBLimit, HTValue, HBDifference)
			beginning = false
			subtotalLimit = NotSetFloatValue
			subtotalValue = NotSetFloatValue
			subtotalDifference = NotSetFloatValue
		}
		fmt.Fprintf(os.Stdout, lineD, e.MainCategory.Name, e.Limit, e.Actual, e.Difference)
		subtotalLimit += e.Limit
		subtotalValue += e.Actual
		subtotalDifference += e.Difference
		totalLimit += e.Limit
		totalValue += e.Actual
		totalDifference += e.Difference
	}
	fmt.Fprintf(os.Stdout, lineS, currentType, subtotalLimit, subtotalValue, subtotalDifference)
	fmt.Fprint(os.Stdout, "\n")
	fmt.Fprintf(os.Stdout, lineS, "TOTAL", totalLimit, totalValue, totalDifference)

	return nil
}

func RepNetValueMonthly(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create filters
	var dateFrom, dateTo time.Time
	if ds := c.String(OptDateFrom); ds != NotSetStringValue {
		if dateFrom, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dateFrom = time.Time{}
	}
	if ds := c.String(OptDateTo); ds != NotSetStringValue {
		if dateTo, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dateTo = time.Now()
	}

	// Build formatting strings
	var getNextEntry func() *NetValueMonthlyReportEntry
	if getNextEntry, err = ReportNetValueMonthly(fh, cur, dateFrom, dateTo); err != nil {
		printError.Fatalln(err)
	}
	lP := utf8.RuneCountInString(HBPeriod)
	lV := utf8.RuneCountInString(HNV)
	var totalValue float64
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		lP = MaxLen(e.Period.String(), lP)
		totalValue += e.Value
		lV = MaxLen(strconv.FormatFloat(totalValue, 'f', 2, 64), lV)
	}
	LineH := LineFor(HFSForText(lP), HFSForNumeric(lV))
	LineD := LineFor(DFSForText(lP), DFSForValue(lV))

	// Print report
	if dateFrom.IsZero() {
		fmt.Fprintf(os.Stdout, "Net value until %s (in %s):\n", dateTo.Format(DateFormat), cur)
	} else {
		fmt.Fprintf(os.Stdout, "Net value between %s and %s (in %s):\n", dateFrom.Format(DateFormat), dateTo.Format(DateFormat), cur)
	}
	fmt.Fprintf(os.Stdout, LineH, HBPeriod, HNV)

	if getNextEntry, err = ReportNetValueMonthly(fh, cur, dateFrom, dateTo); err != nil {
		printError.Fatalln(err)
	}
	totalValue = NotSetFloatValue
	for e := getNextEntry(); e != nil; e = getNextEntry() {
		totalValue += e.Value
		fmt.Fprintf(os.Stdout, LineD, e.Period.String(), totalValue)
	}
	//FIXME: wrong net value if date_from is given
	return nil
}
