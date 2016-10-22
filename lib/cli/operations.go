// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/lib/engine"
	"os"
	"strconv"
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
		printError.Fatalln(errMissingCategory)
	}
	m := c.String(ObjMainCategory)
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

	mcat := c.String(ObjMainCategory)
	var mct MainCategoryType
	if t := c.String(OptMainCategoryType); t == NotSetStringValue {
		mct = MCTUnset
	} else {
		if mct = MainCategoryTypeForString(t); mct == MCTUnknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
	}
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

	// Build formatting strings
	var getNextCategory func() *Category
	if getNextCategory, err = CategoryList(fh, mcat, mct, cat, s); err != nil {
		printError.Fatalln(err)
	}
	lId, lType, lMCat, lCat, lStatus := utf8.RuneCountInString(HCId), utf8.RuneCountInString(HMCType), utf8.RuneCountInString(HMCName), utf8.RuneCountInString(HCName), utf8.RuneCountInString(HMCStatus)
	for ct := getNextCategory(); ct != nil; ct = getNextCategory() {
		lId = MaxLen(strconv.FormatInt(ct.Id, 10), lId)
		lType = MaxLen(ct.Main.MType.String(), lType)
		lMCat = MaxLen(ct.Main.Name, lMCat)
		lCat = MaxLen(ct.Name, lCat)
		lStatus = MaxLen(ct.Status.String(), lStatus)
	}
	lineH := LineFor(HFSForNumeric(lId), HFSForText(lType), HFSForText(lMCat), HFSForText(lCat), HFSForText(lStatus))
	lineD := LineFor(DFSForID(lId), DFSForText(lType), DFSForText(lMCat), DFSForText(lCat), DFSForText(lStatus))

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
	printUserMsg, printError := GetLoggers()

	// Check obligatory flags (file, name)
	f := c.String(OptFile)
	if f == NotSetStringValue {
		printError.Fatalln(errMissingFileFlag)

	}
	n := c.String(ObjMainCategory)
	if n == NotSetStringValue {
		printError.Fatalln(errMissingMainCategory)
	}
	t := MainCategoryTypeForString(c.String(OptMainCategoryType))
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
		mct := MainCategoryTypeForString(t)
		if mct == MCTUnknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
		mc.MType = mct
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
	var mct MainCategoryType
	if t := c.String(OptMainCategoryType); t == NotSetStringValue {
		mct = MCTUnset
	} else {
		mct = MainCategoryTypeForString(t)
		if mct == MCTUnknown {
			printError.Fatalln(errIncorrectMainCategoryType)
		}
	}
	n := c.String(ObjMainCategory)
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

	// Build formatting strings
	var getNextMainCategory func() *MainCategory
	if getNextMainCategory, err = MainCategoryList(fh, mct, n, s); err != nil {
		printError.Fatalln(err)
	}
	lId, lType, lName, lStatus := utf8.RuneCountInString(HMCId), utf8.RuneCountInString(HMCType), utf8.RuneCountInString(HMCName), utf8.RuneCountInString(HMCStatus)
	for m := getNextMainCategory(); m != nil; m = getNextMainCategory() {
		lId = MaxLen(strconv.FormatInt(m.Id, 10), lId)
		lType = MaxLen(m.MType.String(), lType)
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
		fmt.Fprintf(os.Stdout, lineD, m.Id, m.MType, m.Name, m.Status)
	}

	return nil
}

// CmdExchangeRateAdd adds new currency echchange rate
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
	rate := c.Float64(ObjExchangeRateAlias)
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
