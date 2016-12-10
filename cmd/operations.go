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

// CmdCreateNewDataFile creates a new sqlite file and tables for financoj
func CmdCreateNewDataFile(c *cli.Context) error {
	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check the obligatory parameters and exit if missing
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	// Create new data file
	fh := GetDataFileHandler(f)
	if err := CreateNewDataFile(fh); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("created file %s\n", f)

	return nil
}

// CmdCategoryAdd adds new category
func CmdCategoryAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags (file, name)
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(ObjCategory)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}
	m := c.String(ObjMainCategory)
	if m == NotSetStringValue {
		printError.Fatalln(errMissingMainCategoryFlag)
	}

	// Add new category
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var mc *MainCategory
	if mc, err = MainCategoryForName(fh, m); err != nil {
		printError.Fatalln(err)
	}

	newCategory := &Category{Main: mc, Name: n, Status: ISOpen}
	if err = CategoryAdd(fh, newCategory); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("added new category: %s\n", n)

	return nil
}

// CmdCategoryEdit updates category with new values
func CmdCategoryEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Prepare new values based on old ones
	var cat *Category
	if cat, err = CategoryForID(fh, id); err != nil {
		printError.Fatalln(err)
	}
	if m := c.String(ObjMainCategory); m != NotSetStringValue {
		var mcat *MainCategory
		if mcat, err = MainCategoryForName(fh, m); err != nil {
			printError.Fatalln(err)
		}
		cat.Main = mcat
	}
	if n := c.String(ObjCategory); n != NotSetStringValue {
		cat.Name = n
	}

	// Execute the changes
	if err = CategoryEdit(fh, cat); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("changed details of category with id = %d\n", cat.Id)

	return nil
}

// CmdCategoryRemove sets category status to ISClose
func CmdCategoryRemove(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var cat *Category
	if cat, err = CategoryForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	// Remove the category
	if err = CategoryRemove(fh, cat); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("removed  category with id = %d\n", id)

	return nil

}

// CmdCategoryList prints categories on standard output
func CmdCategoryList(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	mn := c.String(ObjMainCategory)
	cat := c.String(ObjCategory)
	s := ISOpen
	if a := c.Bool(OptAll); a == true {
		s = ISUnset
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var mcat *MainCategory
	if mn != NotSetStringValue {
		if mcat, err = MainCategoryForName(fh, mn); err != nil {
			printError.Fatalln(err)
		}
	}

	// Build formatting strings
	var getNextCategory func() *Category
	if getNextCategory, err = CategoryList(fh, mcat, cat, s); err != nil {
		printError.Fatalln(err)
	}
	lId, lType, lMCat, lCat, lStatus := utf8.RuneCountInString(HCId), utf8.RuneCountInString(HMCType), utf8.RuneCountInString(HMCName), utf8.RuneCountInString(HCName), utf8.RuneCountInString(HMCStatus)
	for ct := getNextCategory(); ct != nil; ct = getNextCategory() {
		lId = MaxLen(strconv.FormatInt(ct.Id, 10), lId)
		lType = MaxLen(ct.Main.MType.Name, lType)
		lMCat = MaxLen(ct.Main.Name, lMCat)
		lCat = MaxLen(ct.Name, lCat)
		lStatus = MaxLen(ct.Status.String(), lStatus)
	}
	lineH := LineFor(HFSForNumeric(lId), HFSForText(lType), HFSForText(lMCat), HFSForText(lCat), HFSForText(lStatus))
	lineD := LineFor(DFSForID(lId), DFSForText(lType), DFSForText(lMCat), DFSForText(lCat), DFSForText(lStatus))

	// Print categories
	if getNextCategory, err = CategoryList(fh, mcat, cat, s); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HCId, HMCType, HMCName, HCName, HMCStatus)
	for ct := getNextCategory(); ct != nil; ct = getNextCategory() {
		fmt.Fprintf(os.Stdout, lineD, ct.Id, ct.Main.MType.Name, ct.Main.Name, ct.Name, ct.Status)
	}

	return nil
}

// CmdMainCategoryAdd adds new main category
func CmdMainCategoryAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags (file, name)
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(ObjMainCategory)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingMainCategoryFlag)
	}
	tn := c.String(OptMainCategoryType)

	// Add new main category
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var t *MainCategoryType
	if tn != NotSetStringValue {
		if t, err = MainCategoryTypeForName(fh, tn); err != nil {
			printError.Fatalln(err)
		}
	} else {
		if t, err = MainCategoryTypeForID(fh, MCTCost); err != nil {
			printError.Fatalln(err)
		}
	}

	m := &MainCategory{MType: t, Name: n, Status: ISOpen}
	if err = MainCategoryAdd(fh, m); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("added new main category: %s (type: %s)\n", n, t.Name)

	return nil
}

// CmdMainCategoryEdit updates main category with new values
func CmdMainCategoryEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var mc *MainCategory
	if mc, err = MainCategoryForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	// Edit main category
	if t := c.String(OptMainCategoryType); t != NotSetStringValue {
		if mc.MType, err = MainCategoryTypeForName(fh, t); err != nil {
			printError.Fatalln(err)
		}
	}
	if n := c.String(ObjMainCategory); n != NotSetStringValue {
		mc.Name = n
	}
	if err = MainCategoryEdit(fh, mc); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("changed details of main category with id = %d\n", id)

	return nil
}

// CmdMainCategoryRemove sets main category status to ISClose
func CmdMainCategoryRemove(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var mc *MainCategory
	if mc, err = MainCategoryForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	// Remove the main category
	if err = MainCategoryRemove(fh, mc); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("removed main category with id = %d\n", mc.Id)

	return nil
}

// CmdMainCategoryList prints main categories on standard output
func CmdMainCategoryList(c *cli.Context) error {
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
	var mct *MainCategoryType
	if t := c.String(OptMainCategoryType); t != NotSetStringValue {
		if mct, err = MainCategoryTypeForName(fh, t); err != nil {
			printError.Fatalln(err)
		}
	}
	n := c.String(ObjMainCategory)
	s := ISOpen
	if a := c.Bool(OptAll); a == true {
		s = ISUnset
	}

	// Build formatting strings
	var getNextMainCategory func() *MainCategory
	if getNextMainCategory, err = MainCategoryList(fh, mct, n, s); err != nil {
		printError.Fatalln(err)
	}
	lId, lType, lName, lStatus := utf8.RuneCountInString(HMCId), utf8.RuneCountInString(HMCType), utf8.RuneCountInString(HMCName), utf8.RuneCountInString(HMCStatus)
	for m := getNextMainCategory(); m != nil; m = getNextMainCategory() {
		lId = MaxLen(strconv.FormatInt(m.Id, 10), lId)
		lType = MaxLen(m.MType.Name, lType)
		lName = MaxLen(m.Name, lName)
		lStatus = MaxLen(m.Status.String(), lStatus)
	}
	lineH := LineFor(HFSForNumeric(lId), HFSForText(lType), HFSForText(lName), HFSForText(lStatus))
	lineD := LineFor(DFSForID(lId), DFSForText(lType), DFSForText(lName), DFSForText(lStatus))

	// Print main categories
	if getNextMainCategory, err = MainCategoryList(fh, mct, n, s); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HMCId, HMCType, HMCName, HMCStatus)
	for m := getNextMainCategory(); m != nil; m = getNextMainCategory() {
		fmt.Fprintf(os.Stdout, lineD, m.Id, m.MType.Name, m.Name, m.Status)
	}

	return nil
}

// CmdExchangeRateAdd adds new currency exchange rate
func CmdExchangeRateAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags (file, name)
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	curFrom := c.String(OptCurrency)
	if curFrom == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	curTo := c.String(OptCurrencyTo)
	if curTo == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyToFlag)
	}
	rate := c.Float64(ObjExchangeRate)
	if rate == NotSetFloatValue {
		printError.Fatalln(errMissingExchangeRateFlag)
	}

	// Add currency exchange rate
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	newCurrency := &ExchangeRate{CurrencyFrom: curFrom, CurrencyTo: curTo, Rate: rate}
	if err = ExchangeRateAdd(fh, newCurrency); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("added new currency exchange rate: %s-%s\n", curFrom, curTo)
	return nil
}

// CmdExchangeRateEdit edits currency exchange rate
func CmdExchangeRateEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	cf := c.String(OptCurrency)
	if cf == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	ct := c.String(OptCurrencyTo)
	if ct == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyToFlag)
	}
	r := c.Float64(ObjExchangeRate)
	if r == NotSetFloatValue {
		printError.Fatalln(errMissingExchangeRateFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var e *ExchangeRate
	if e, err = ExchangeRateForCurrencies(fh, cf, ct); err != nil {
		printError.Fatalln(err)
	}

	// Edit exchange rate
	e.Rate = r
	if err = ExchangeRateEdit(fh, e); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("changed exchange rate for %s-%s\n", e.CurrencyFrom, e.CurrencyTo)

	return nil
}

// CmdExchangeRateList lists currency exchange rates
func CmdExchangeRateList(c *cli.Context) error {
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

	// Build formatting strings
	var getNextCurrency func() *ExchangeRate
	if getNextCurrency, err = ExchangeRateList(fh); err != nil {
		printError.Fatalln(err)
	}
	lCurF, lCurT, lRate := utf8.RuneCountInString(HCurF), utf8.RuneCountInString(HCurT), utf8.RuneCountInString(HCurRate)
	for cur := getNextCurrency(); cur != nil; cur = getNextCurrency() {
		lCurF = MaxLen(cur.CurrencyFrom, lCurF)
		lCurT = MaxLen(cur.CurrencyTo, lCurT)
		lRate = MaxLen(strconv.FormatFloat(cur.Rate, 'f', -1, 64), lRate)
	}
	lineH := LineFor(HFSForText(lCurF), HFSForText(lCurT), HFSForNumeric(lRate))
	lineD := LineFor(DFSForText(lCurF), DFSForText(lCurT), DFSForRates(lRate))

	// Print currencies
	if getNextCurrency, err = ExchangeRateList(fh); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HCurF, HCurT, HCurRate)
	for cur := getNextCurrency(); cur != nil; cur = getNextCurrency() {
		fmt.Fprintf(os.Stdout, lineD, cur.CurrencyFrom, cur.CurrencyTo, cur.Rate)
	}

	return nil
}

// CmdExchangeRateRemove removes exchange rates for given currencies
func CmdExchangeRateRemove(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	j := c.String(OptCurrency)
	if j == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	k := c.String(OptCurrencyTo)
	if k == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyToFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var cur *ExchangeRate
	if cur, err = ExchangeRateForCurrencies(fh, j, k); err != nil {
		printError.Fatalln(err)
	}

	// Remove the exchange rate
	if err = ExchangeRateRemove(fh, cur); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("removed currency exchange rate for %s and %s\n", cur.CurrencyFrom, cur.CurrencyTo)

	return nil

}

// CmdAccountAdd adds new account
func CmdAccountAdd(c *cli.Context) error {
	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags (file, name)
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(ObjAccount)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	j := c.String(OptCurrency)
	if j == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	t := AccountTypeForString(c.String(OptAccountType))
	if t == ATUnknown {
		printError.Fatalln(errIncorrectAccountType)
	}

	// Other flags
	d := c.String(OptDescription)
	i := c.String(OptInstitution)

	// Add new account
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	a := &Account{Name: n, Description: d, Institution: i, Currency: j, AType: t, Status: ISOpen}
	if err := AccountAdd(fh, a); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("added new account: %s (type: %s)\n", a.Name, a.AType)

	return nil

}

// CmdAccountList lists account fitting given criteria
func CmdAccountList(c *cli.Context) error {
	var err error

	// Get loggers
	_, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	// Parse other flags
	name := c.String(ObjAccount)
	description := c.String(OptDescription)
	institution := c.String(OptInstitution)
	currency := c.String(OptCurrency)
	var atype AccountType
	if t := c.String(OptAccountType); t == NotSetStringValue {
		atype = ATUnset
	} else {
		if atype = AccountTypeForString(t); atype == ATUnknown {
			printError.Fatalln(errIncorrectAccountType)
		}
	}
	status := ISOpen
	if s := c.Bool(OptAll); s == true {
		status = ISUnset
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Build formatting strings
	var getNextAccount func() *Account
	if getNextAccount, err = AccountList(fh, name, description, institution, currency, atype, status); err != nil {
		printError.Fatalln(err)
	}
	lId := utf8.RuneCountInString(HAId)
	lN := utf8.RuneCountInString(HAName)
	lD := utf8.RuneCountInString(HADescription)
	lI := utf8.RuneCountInString(HAInstitution)
	lC := utf8.RuneCountInString(HACurrency)
	lT := utf8.RuneCountInString(HAType)
	lS := utf8.RuneCountInString(HAStatus)
	for a := getNextAccount(); a != nil; a = getNextAccount() {
		lId = MaxLen(strconv.FormatInt(a.Id, 10), lId)
		lN = MaxLen(a.Name, lN)
		lD = MaxLen(a.Description, lD)
		lI = MaxLen(a.Institution, lI)
		lC = MaxLen(a.Currency, lC)
		lT = MaxLen(a.AType.String(), lT)
		lS = MaxLen(a.Status.String(), lS)
	}
	lineH := LineFor(HFSForNumeric(lId), HFSForText(lN), HFSForText(lT), HFSForText(lC), HFSForText(lI), HFSForText(lS), HFSForText(lD))
	lineD := LineFor(DFSForID(lId), DFSForText(lN), DFSForText(lT), DFSForText(lC), DFSForText(lI), DFSForText(lS), DFSForText(lD))

	// Print accounts
	if getNextAccount, err = AccountList(fh, name, description, institution, currency, atype, status); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HAId, HAName, HAType, HACurrency, HAInstitution, HAStatus, HADescription)
	for a := getNextAccount(); a != nil; a = getNextAccount() {
		fmt.Fprintf(os.Stdout, lineD, a.Id, a.Name, a.AType, a.Currency, a.Institution, a.Status, a.Description)
	}

	return nil
}

// CmdAccountEdit updates account with new values
func CmdAccountEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Prepare new values based on old ones
	var a *Account
	if a, err = AccountForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	if n := c.String(ObjAccount); n != NotSetStringValue {
		a.Name = n
	}
	if d := c.String(OptDescription); d != NotSetStringValue {
		a.Description = d
	}
	if i := c.String(OptInstitution); i != NotSetStringValue {
		a.Institution = i
	}
	if j := c.String(OptCurrency); j != NotSetStringValue {
		a.Currency = j
	}
	if ts := c.String(OptAccountType); ts != NotSetStringValue {
		if at := AccountTypeForString(ts); at == ATUnknown {
			printError.Fatalln(errIncorrectAccountType)
		} else {
			a.AType = at
		}
	}

	// Execute the changes
	if err = AccountEdit(fh, a); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("changed details of account with id = %d\n", a.Id)

	return nil
}

// CmdAccountRemove removes account with given id
func CmdAccountRemove(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var a *Account
	if a, err = AccountForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	// Remove the account
	if err = AccountRemove(fh, a); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("removed account with id = %d\n", a.Id)

	return nil
}

// CmdTransactionAdd adds new transaction
func CmdTransactionAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	d := c.String(OptDescription)
	if d == NotSetStringValue {
		printError.Fatalln(errMissingDescriptionFlag)
	}
	v := c.Float64(OptValue)
	if v == NotSetFloatValue {
		printError.Fatalln(errMissingValueFlag)
	}
	an := c.String(ObjAccount)
	if an == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	cn := c.String(ObjCategory)
	if cn == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Create the transaction object
	t := TransactionNew()
	if td := c.String(OptDate); td != NotSetStringValue {
		if t.Date, err = time.Parse(DateFormat, td); err != nil {
			printError.Fatalln(err)
		}
	}
	if t.Category, err = CategoryForName(fh, cn); err != nil {
		printError.Fatalln(err)
	}
	if t.Account, err = AccountForName(fh, an); err != nil {
		printError.Fatalln(err)
	}
	t.Value = v
	t.Description = d

	// Add transaction
	if err = TransactionAdd(fh, t); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("add new transaction\n")

	return nil
}

// CmdTransactionList prints transactions on standard output
func CmdTransactionList(c *cli.Context) error {
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

	// Get filtering criteria
	var dateFrom time.Time
	if ds := c.String(OptDateFrom); ds != NotSetStringValue {
		if dateFrom, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dateFrom = time.Time{}
	}
	var dateTo time.Time
	if ds := c.String(OptDateTo); ds != NotSetStringValue {
		if dateTo, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	} else {
		dateTo = time.Time{}
	}
	var account *Account
	if as := c.String(ObjAccount); as != NotSetStringValue {
		if account, err = AccountForName(fh, as); err != nil {
			printError.Fatalln(err)
		}
	}
	description := c.String(OptDescription)
	var category *Category
	if cs := c.String(ObjCategory); cs != NotSetStringValue {
		if category, err = CategoryForName(fh, cs); err != nil {
			printError.Fatalln(err)
		}
	}
	var mainCategory *MainCategory
	if ms := c.String(ObjMainCategory); ms != NotSetStringValue {
		if mainCategory, err = MainCategoryForName(fh, ms); err != nil {
			printError.Fatalln(err)
		}
	}

	// Build formatting strings
	var getNextTransaction func() *Transaction
	if getNextTransaction, err = TransactionList(fh, dateFrom, dateTo, account, description, category, mainCategory); err != nil {
		printError.Fatalln(err)
	}
	lId := utf8.RuneCountInString(HTId)
	lDate := utf8.RuneCountInString(HTDate)
	lAccount := utf8.RuneCountInString(HAName)
	lType := utf8.RuneCountInString(HMCType)
	lMCat := utf8.RuneCountInString(HMCName)
	lCat := utf8.RuneCountInString(HCName)
	lValue := utf8.RuneCountInString(HTValue)
	lCur := utf8.RuneCountInString(HACurrency)
	lDesc := utf8.RuneCountInString(HTDescription)

	for t := getNextTransaction(); t != nil; t = getNextTransaction() {
		lId = MaxLen(strconv.FormatInt(t.Id, 10), lId)
		lDate = MaxLen(t.Date.Format(DateFormat), lDate)
		lAccount = MaxLen(t.Account.Name, lAccount)
		lType = MaxLen(t.Category.Main.MType.Name, lType)
		lMCat = MaxLen(t.Category.Main.Name, lMCat)
		lCat = MaxLen(t.Category.Name, lCat)
		lValue = MaxLen(strconv.FormatFloat(t.GetSValue(), 'f', 2, 64), lValue)
		lCur = MaxLen(t.Account.Currency, lCur)
		lDesc = MaxLen(t.Description, lDesc)
	}
	lineH := LineFor(HFSForNumeric(lId), HFSForText(lDate), HFSForText(lAccount), HFSForText(lType), HFSForText(lMCat), HFSForText(lCat), HFSForNumeric(lValue), HFSForText(lCur), HFSForText(lDesc))
	lineD := LineFor(DFSForID(lId), DFSForText(lDate), DFSForText(lAccount), DFSForText(lType), DFSForText(lMCat), DFSForText(lCat), DFSForValue(lValue), DFSForText(lCur), DFSForText(lDesc))

	// Print transactions
	if getNextTransaction, err = TransactionList(fh, dateFrom, dateTo, account, description, category, mainCategory); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HTId, HTDate, HAName, HMCType, HMCName, HCName, HTValue, HACurrency, HTDescription)
	for t := getNextTransaction(); t != nil; t = getNextTransaction() {
		fmt.Fprintf(os.Stdout, lineD, t.Id, t.Date.Format(DateFormat), t.Account.Name, t.Category.Main.MType.Name, t.Category.Main.Name, t.Category.Name, t.GetSValue(), t.Account.Currency, t.Description)
	}

	return nil
}

// CmdTransactionEdit updates transaction with new values
func CmdTransactionEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var t *Transaction
	if t, err = TransactionForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	// Edit transaction
	if ds := c.String(OptDate); ds != NotSetStringValue {
		if t.Date, err = time.Parse(DateFormat, ds); err != nil {
			printError.Fatalln(err)
		}
	}
	if cs := c.String(ObjCategory); cs != NotSetStringValue {
		if t.Category, err = CategoryForName(fh, cs); err != nil {
			printError.Fatalln(err)
		}
	}
	if as := c.String(ObjAccount); as != NotSetStringValue {
		if t.Account, err = AccountForName(fh, as); err != nil {
			printError.Fatalln(err)
		}
	}
	if vf := c.Float64(OptValue); vf != NotSetFloatValue {
		t.Value = vf
	}
	if descr := c.String(OptDescription); descr != NotSetStringValue {
		t.Description = descr
	}

	if err = TransactionEdit(fh, t); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("changed details of transaction with id = %d\n", id)

	return nil
}

// CmdTransactionRemove removes transaction with given id
func CmdTransactionRemove(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	id := c.Int(OptID)
	if id == NotSetIntValue {
		printError.Fatalln(errMissingIDFlag)
	}

	// Open data file and get original main category
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var t *Transaction
	if t, err = TransactionForID(fh, id); err != nil {
		printError.Fatalln(err)
	}

	// Remove the transaction
	if err = TransactionRemove(fh, t); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("removed transaction with id = %d\n", t.Id)

	return nil
}

// CmdBudgetAdd adds new budget
func CmdBudgetAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags (file, name)
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	p := c.String(OptPeriod)
	if p == NotSetStringValue {
		printError.Fatalln(errMissingPeriodFlag)
	}
	cat := c.String(ObjCategory)
	if cat == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}
	v := c.Float64(OptValue)
	if v == NotSetFloatValue {
		printError.Fatalln(errMissingValueFlag)
	}
	cur := c.String(OptCurrency)
	if cur == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}

	// Open data file and validate parameters
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	b := BudgetNew()
	if b.Period, err = BPeriodParseYM(p); err != nil {
		printError.Fatalln(err)
	}
	if b.Category, err = CategoryForName(fh, cat); err != nil {
		printError.Fatalln(err)
	}
	b.Value = v
	b.Currency = cur

	// Add new budget
	if err = BudgetAdd(fh, b); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("added new budget\n")

	return nil
}

// CmdBudgetRemove removes budget
func CmdBudgetRemove(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	ps := c.String(OptPeriod)
	if ps == NotSetStringValue {
		printError.Fatalln(errMissingPeriodFlag)
	}
	cs := c.String(ObjCategory)
	if cs == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}

	// Open data file and validate parameters
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var p *BPeriod
	if p, err = BPeriodParseYM(ps); err != nil {
		printError.Fatalln(err)
	}
	var cat *Category
	if cat, err = CategoryForName(fh, cs); err != nil {
		printError.Fatalln(err)
	}

	// Find the budget and remove it
	var b *Budget
	if b, err = BudgetGet(fh, p, cat); err != nil {
		printError.Fatalln(err)
	}
	if err = BudgetRemove(fh, b); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("removed budget for period: %s and category: %s\n", p, cat.Name)

	return nil
}

// CmdBudgetEdit updates budget with new values
func CmdBudgetEdit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	ps := c.String(OptPeriod)
	if ps == NotSetStringValue {
		printError.Fatalln(errMissingPeriodFlag)
	}
	cs := c.String(ObjCategory)
	if cs == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}

	// Open data file and validate parameters
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	var p *BPeriod
	if p, err = BPeriodParseYM(ps); err != nil {
		printError.Fatalln(err)
	}
	var cat *Category
	if cat, err = CategoryForName(fh, cs); err != nil {
		printError.Fatalln(err)
	}

	// Find the budget and remove it
	var b *Budget
	if b, err = BudgetGet(fh, p, cat); err != nil {
		printError.Fatalln(err)
	}
	if v := c.Float64(OptValue); v != NotSetFloatValue {
		b.Value = v
	}
	if cur := c.String(OptCurrency); cur != NotSetStringValue {
		b.Currency = cur
	}

	// Edit budget
	if err = BudgetEdit(fh, b); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("updated budget with new values")

	return nil
}

// CmdBudgetList prints budgets on standard output
func CmdBudgetList(c *cli.Context) error {
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

	// Get filtering criteria
	var p *BPeriod
	if p, err = BPeriodParseYOrYM(OptPeriod); err != nil {
		printError.Fatalln(err)
	}
	var ct *Category
	if cs := c.String(ObjCategory); cs != NotSetStringValue {
		if ct, err = CategoryForName(fh, cs); err != nil {
			printError.Fatalln(err)
		}
	}

	// Build formatting strings
	var getNextBudget func() *Budget
	if getNextBudget, err = BudgetList(fh, p, ct); err != nil {
		printError.Fatalln(err)
	}
	lP := utf8.RuneCountInString(HBPeriod)
	lT := utf8.RuneCountInString(HMCType)
	lMC := utf8.RuneCountInString(HMCName)
	lC := utf8.RuneCountInString(HCName)
	lL := utf8.RuneCountInString(HBLimit)
	lCur := utf8.RuneCountInString(HBCurrency)

	for b := getNextBudget(); b != nil; b = getNextBudget() {
		lP = MaxLen(b.Period.String(), lP)
		lT = MaxLen(b.Category.Main.MType.Name, lT)
		lMC = MaxLen(b.Category.Main.Name, lMC)
		lC = MaxLen(b.Category.Name, lC)
		lL = MaxLen(strconv.FormatFloat(b.Value, 'f', 2, 64), lL)
		lCur = MaxLen(b.Currency, lCur)
	}
	lineH := LineFor(HFSForText(lP), HFSForText(lT), HFSForText(lMC), HFSForText(lC), HFSForNumeric(lL), HFSForText(lCur))
	LineD := LineFor(DFSForText(lP), DFSForText(lT), DFSForText(lMC), DFSForText(lC), DFSForValue(lL), DFSForText(lCur))

	// Print budgets
	if getNextBudget, err = BudgetList(fh, p, ct); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HBPeriod, HMCType, HMCName, HCName, HBLimit, HBCurrency)
	for b := getNextBudget(); b != nil; b = getNextBudget() {
		fmt.Fprintf(os.Stdout, LineD, b.Period, b.Category.Main.MType.Name, b.Category.Main.Name, b.Category.Name, b.Value, b.Currency)
	}

	return nil
}

// CmdCompoundTransferAdd adds two transactions with non-budgetable category 'Transfer'
func CmdCompoundTransferAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	af := c.String(ObjAccount)
	if af == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	at := c.String(OptAccountTo)
	if at == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	v := c.Float64(OptValue)
	if v == NotSetFloatValue {
		printError.Fatalln(errMissingValueFlag)
	}
	desc := c.String(OptDescription)
	if desc == NotSetStringValue {
		printError.Fatalln(errMissingDescriptionFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Parse necessary parameters
	d := time.Now()
	if td := c.String(OptDate); td != NotSetStringValue {
		if d, err = time.Parse(DateFormat, td); err != nil {
			printError.Fatalln(err)
		}
	}
	var accFrom, accTo *Account
	if accFrom, err = AccountForName(fh, af); err != nil {
		printError.Fatalln(err)
	}
	if accTo, err = AccountForName(fh, at); err != nil {
		printError.Fatalln(err)
	}
	var er *ExchangeRate
	if r := c.Float64(ObjExchangeRate); r == NotSetFloatValue {
		if er, err = ExchangeRateForCurrencies(fh, accFrom.Currency, accTo.Currency); err != nil {
			printError.Fatalln(err)
		}
	} else {
		er = new(ExchangeRate)
		er.CurrencyFrom = accFrom.Currency
		er.CurrencyTo = accTo.Currency
		er.Rate = r
	}

	// Add transaction
	if err = CompoundTransferAdd(fh, d, accFrom, accTo, v, desc, er); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("add new transfer\n")

	return nil
}

// CmdCompoundInternalCostAdd adds two transactions: cost in first account and non-budgetable category 'Transfer' into the other
func CmdCompoundInternalCostAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	cs := c.String(ObjCategory)
	if cs == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}
	ac := c.String(ObjAccount)
	if ac == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	at := c.String(OptAccountTo)
	if at == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	v := c.Float64(OptValue)
	if v == NotSetFloatValue {
		printError.Fatalln(errMissingValueFlag)
	}
	desc := c.String(OptDescription)
	if desc == NotSetStringValue {
		printError.Fatalln(errMissingDescriptionFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Parse necessary parameters
	d := time.Now()
	if td := c.String(OptDate); td != NotSetStringValue {
		if d, err = time.Parse(DateFormat, td); err != nil {
			printError.Fatalln(err)
		}
	}
	var cat *Category
	if cat, err = CategoryForName(fh, cs); err != nil {
		printError.Fatalln(err)
	}
	var accCost, accTransfer *Account
	if accCost, err = AccountForName(fh, ac); err != nil {
		printError.Fatalln(err)
	}
	if accTransfer, err = AccountForName(fh, at); err != nil {
		printError.Fatalln(err)
	}
	var er *ExchangeRate
	if r := c.Float64(ObjExchangeRate); r == NotSetFloatValue {
		if er, err = ExchangeRateForCurrencies(fh, accCost.Currency, accTransfer.Currency); err != nil {
			printError.Fatalln(err)
		}
	} else {
		er = new(ExchangeRate)
		er.CurrencyFrom = accCost.Currency
		er.CurrencyTo = accTransfer.Currency
		er.Rate = r
	}

	// Add transaction
	if err = CompoundInternalCostAdd(fh, d, cat, accCost, accTransfer, v, desc, er); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("add new internal cost\n")

	return nil
}

// CmdCompoundTransactionSplit adds two transactions for two different categories with half of the original value
func CmdCompoundTransactionSplit(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	a := c.String(ObjAccount)
	if a == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	v := c.Float64(OptValue)
	if v == NotSetFloatValue {
		printError.Fatalln(errMissingValueFlag)
	}
	desc := c.String(OptDescription)
	if desc == NotSetStringValue {
		printError.Fatalln(errMissingDescriptionFlag)
	}
	c1 := c.String(ObjCategory)
	if c1 == NotSetStringValue {
		printError.Fatalln(errMissingCategoryFlag)
	}
	c2 := c.String(OptCategorySplit)
	if c2 == NotSetStringValue {
		printError.Fatalln(errmissingCategorySplitFlag)
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Parse necessary parameters
	d := time.Now()
	if td := c.String(OptDate); td != NotSetStringValue {
		if d, err = time.Parse(DateFormat, td); err != nil {
			printError.Fatalln(err)
		}
	}
	var acc *Account
	if acc, err = AccountForName(fh, a); err != nil {
		printError.Fatalln(err)
	}
	var cat1, cat2 *Category
	if cat1, err = CategoryForName(fh, c1); err != nil {
		printError.Fatalln(err)
	}
	if cat2, err = CategoryForName(fh, c2); err != nil {
		printError.Fatalln(err)
	}

	// Add transaction
	if err = CompoundSplitAdd(fh, d, acc, v, desc, cat1, cat2); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("add split transaction\n")

	return nil
}

//FIXME: make user messages more verbose (good example: BudgetRemove)
//FIXME: split operations file into separate files, one for each object
