// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/lib"
	"os"
	"strconv"
	"unicode/utf8"
)

// CmdCreateNewDataFile creates a new sqlite file and tables for financoj
func CmdCreateNewDataFile(c *cli.Context) error {
	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check the obligatory parameters and exit if missing
	f := c.String(optFile)
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
	printUserMsg, printError := getLoggers()

	// Check obligatory flags (file, name)
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(objCategory)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingCategory)
	}
	m := c.String(objMainCategory)
	if m == NotSetStringValue {
		printError.Fatalln(errMissingMainCategory)
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
	if m := c.String(objMainCategory); m != NotSetStringValue {
		var mcat *MainCategory
		if mcat, err = MainCategoryForName(fh, m); err != nil {
			printError.Fatalln(err)
		}
		cat.Main = mcat
	}
	if n := c.String(objCategory); n != NotSetStringValue {
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
	_, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	mcat := c.String(objMainCategory)
	var mct MainCategoryType
	if t := c.String(optMainCategoryType); t == NotSetStringValue {
		mct = MCTUnset
	} else {
		if mct = mainCategoryTypeForString(t); mct == MCTUnknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
	}
	cat := c.String(objCategory)
	s := ISOpen
	if a := c.Bool(optAll); a == true {
		s = ISUnset
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Build formatting strings
	var getNextCategory func() *Category
	if getNextCategory, err = CategoryList(fh, mcat, mct, cat, s); err != nil {
		printError.Fatalln(err)
	}
	lId, lType, lMCat, lCat, lStatus := utf8.RuneCountInString(HCId), utf8.RuneCountInString(HMCType), utf8.RuneCountInString(HMCName), utf8.RuneCountInString(HCName), utf8.RuneCountInString(HMCStatus)
	for ct := getNextCategory(); ct != nil; ct = getNextCategory() {
		lId = maxLen(strconv.FormatInt(ct.Id, 10), lId)
		lType = maxLen(ct.Main.MType.String(), lType)
		lMCat = maxLen(ct.Main.Name, lMCat)
		lCat = maxLen(ct.Name, lCat)
		lStatus = maxLen(ct.Status.String(), lStatus)
	}
	lineH := lineFor(hFSForNumeric(lId), hFSForText(lType), hFSForText(lMCat), hFSForText(lCat), hFSForText(lStatus))
	lineD := lineFor(dFSForID(lId), dFSForText(lType), dFSForText(lMCat), dFSForText(lCat), dFSForText(lStatus))

	// Print categories
	if getNextCategory, err = CategoryList(fh, mcat, mct, cat, s); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HCId, HMCType, HMCName, HCName, HMCStatus)
	for ct := getNextCategory(); ct != nil; ct = getNextCategory() {
		fmt.Fprintf(os.Stdout, lineD, ct.Id, ct.Main.MType, ct.Main.Name, ct.Name, ct.Status)
	}

	return nil
}

// CmdMainCategoryAdd adds new main category
func CmdMainCategoryAdd(c *cli.Context) error {
	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check obligatory flags (file, name)
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(objMainCategory)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingMainCategory)
	}
	t := mainCategoryTypeForString(c.String(optMainCategoryType))
	if t == MCTUnknown {
		printError.Fatalln(errIncorrectMainCategoryType)
	}

	// Add new main category
	fh := GetDataFileHandler(f)
	if err := fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	m := &MainCategory{MType: t, Name: n, Status: ISOpen}
	if err := MainCategoryAdd(fh, m); err != nil {
		printError.Fatalln(err)
	}

	// Show summary
	printUserMsg.Printf("added new main category: %s (type: %s)\n", n, t)

	return nil
}

// CmdMainCategoryEdit updates main category with new values
func CmdMainCategoryEdit(c *cli.Context) error {
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
	if t := c.String(optMainCategoryType); t != NotSetStringValue {
		mct := mainCategoryTypeForString(t)
		if mct == MCTUnknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
		mc.MType = mct
	}
	if n := c.String(objMainCategory); n != NotSetStringValue {
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
	_, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}
	var mct MainCategoryType
	if t := c.String(optMainCategoryType); t == NotSetStringValue {
		mct = MCTUnset
	} else {
		mct = mainCategoryTypeForString(t)
		if mct == MCTUnknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
	}
	n := c.String(objMainCategory)
	s := ISOpen
	if a := c.Bool(optAll); a == true {
		s = ISUnset
	}

	// Open data file
	fh := GetDataFileHandler(f)
	if err = fh.Open(); err != nil {
		printError.Fatalln(err)
	}
	defer fh.Close()

	// Build formatting strings
	var getNextMainCategory func() *MainCategory
	if getNextMainCategory, err = MainCategoryList(fh, mct, n, s); err != nil {
		printError.Fatalln(err)
	}
	lId, lType, lName, lStatus := utf8.RuneCountInString(HMCId), utf8.RuneCountInString(HMCType), utf8.RuneCountInString(HMCName), utf8.RuneCountInString(HMCStatus)
	for m := getNextMainCategory(); m != nil; m = getNextMainCategory() {
		lId = maxLen(strconv.FormatInt(m.Id, 10), lId)
		lType = maxLen(m.MType.String(), lType)
		lName = maxLen(m.Name, lName)
		lStatus = maxLen(m.Status.String(), lStatus)
	}
	lineH := lineFor(hFSForNumeric(lId), hFSForText(lType), hFSForText(lName), hFSForText(lStatus))
	lineD := lineFor(dFSForID(lId), dFSForText(lType), dFSForText(lName), dFSForText(lStatus))

	// Print main categories
	if getNextMainCategory, err = MainCategoryList(fh, mct, n, s); err != nil {
		printError.Fatalln(err)
	}
	fmt.Fprintf(os.Stdout, lineH, HMCId, HMCType, HMCName, HMCStatus)
	for m := getNextMainCategory(); m != nil; m = getNextMainCategory() {
		fmt.Fprintf(os.Stdout, lineD, m.Id, m.MType, m.Name, m.Status)
	}

	return nil
}

// CmdExchangeRateAdd adds new currency echchange rate
func CmdExchangeRateAdd(c *cli.Context) error {
	var err error

	// Get loggers
	printUserMsg, printError := getLoggers()

	// Check obligatory flags (file, name)
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	curFrom := c.String(optCurrency)
	if curFrom == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	curTo := c.String(optCurrencyTo)
	if curTo == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyToFlag)
	}
	rate := c.Float64(objExchangeRateAlias)
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
	printUserMsg, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	cf := c.String(optCurrency)
	if cf == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	ct := c.String(optCurrencyTo)
	if ct == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyToFlag)
	}
	r := c.Float64(objExchangeRate)
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
	_, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
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
		lCurF = maxLen(cur.CurrencyFrom, lCurF)
		lCurT = maxLen(cur.CurrencyTo, lCurT)
		lRate = maxLen(strconv.FormatFloat(cur.Rate, 'f', -1, 64), lRate)
	}
	lineH := lineFor(hFSForText(lCurF), hFSForText(lCurT), hFSForNumeric(lRate))
	lineD := lineFor(dFSForText(lCurF), dFSForText(lCurT), getDFSForRates(lRate))

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
	printUserMsg, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	j := c.String(optCurrency)
	if j == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	k := c.String(optCurrencyTo)
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
	printUserMsg, printError := getLoggers()

	// Check obligatory flags (file, name)
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(objAccount)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingAccountFlag)
	}
	j := c.String(optCurrency)
	if j == NotSetStringValue {
		printError.Fatalln(errMissingCurrencyFlag)
	}
	t := accountTypeForString(c.String(optAccountType))
	if t == ATUnknown {
		printError.Fatalln(errIncorrectAccountType)
	}

	// Other flags
	d := c.String(optDescription)
	i := c.String(optInstitution)

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
	_, printError := getLoggers()

	// Check obligatory flags
	f := c.String(optFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)
	}

	// Parse other flags
	name := c.String(objAccount)
	description := c.String(optDescription)
	institution := c.String(optInstitution)
	currency := c.String(optCurrency)
	var atype AccountType
	if t := c.String(optAccountType); t == NotSetStringValue {
		atype = ATUnset
	} else {
		if atype = accountTypeForString(t); atype == ATUnknown {
			printError.Fatalln(errIncorrectAccountType)
		}
	}
	status := ISOpen
	if s := c.Bool(optAll); s == true {
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
		lId = maxLen(strconv.FormatInt(a.Id, 10), lId)
		lN = maxLen(a.Name, lN)
		lD = maxLen(a.Description, lD)
		lI = maxLen(a.Institution, lI)
		lC = maxLen(a.Currency, lC)
		lT = maxLen(a.AType.String(), lT)
		lS = maxLen(a.Status.String(), lS)
	}
	lineH := lineFor(hFSForNumeric(lId), hFSForText(lN), hFSForText(lT), hFSForText(lC), hFSForText(lI), hFSForText(lS), hFSForText(lD))
	lineD := lineFor(dFSForID(lId), dFSForText(lN), dFSForText(lT), dFSForText(lC), dFSForText(lI), dFSForText(lS), dFSForText(lD))

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

	if n := c.String(objAccount); n != NotSetStringValue {
		a.Name = n
	}
	if d := c.String(optDescription); d != NotSetStringValue {
		a.Description = d
	}
	if i := c.String(optInstitution); i != NotSetStringValue {
		a.Institution = i
	}
	if j := c.String(optCurrency); j != NotSetStringValue {
		a.Currency = j
	}
	if ts := c.String(optAccountType); ts != NotSetStringValue {
		if at := accountTypeForString(ts); at == ATUnknown {
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
