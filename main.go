// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/lib"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Settings
const (
	FSSeparator   = "  "
	NullDataValue = "-"
)

// Headings for displaying data and reports
const (
	HCId   = "ID"
	HCName = "CATEGORY"

	HMCId     = "ID"
	HMCType   = "TYPE"
	HMCName   = "MAINCAT"
	HMCStatus = "STATUS"

	HCurF    = "CUR_FR"
	HCurT    = "CUR_TO"
	HCurRate = "EX.RATE"

	HAId          = "ID"
	HAName        = "ACCOUNT"
	HADescription = "DESCRIPTION"
	HAInstitution = "BANK"
	HACurrency    = "CUR"
	HAType        = "TYPE"
	HAStatus      = "STATUS"
)

// Errors
const (
	errMissingFileFlag           = "missing information about data file"
	errMissingIDFlag             = "missing ID"
	errMissingCategory           = "missing category name"
	errMissingMainCategory       = "missing main category name"
	errIncorrectMainCategoryType = "incorrect main category type"
	errMissingCurrencyFlag       = "missing currency (from) name"
	errMissingCurrencyToFlag     = "missing currency_to name"
	errMissingExchangeRateFlag   = "missing exchange rate"
	errMissingAccountFlag        = "missing account name"
	errIncorrectAccountType      = "incorrect account type"
)

// Commands, objects and options
const (
	cmdInit        = "init"
	cmdInitAlias   = "I"
	cmdAdd         = "add"
	cmdAddAlias    = "A"
	cmdEdit        = "edit"
	cmdEditAlias   = "E"
	cmdRemove      = "remove"
	cmdRemoveAlias = "R"
	cmdList        = "list"
	cmdListAlias   = "L"

	optFile                  = "file"
	optFileAlias             = "f"
	optAll                   = "all"
	optMainCategoryType      = "main-category-type"
	optMainCategoryTypeAlias = "o"
	optID                    = "id"
	optIDAlias               = "i"
	optCurrency              = "currency"
	optCurrencyAlias         = "j"
	optCurrencyTo            = "currency_to"
	optCurrencyToAlias       = "k"
	optDescription           = "description"
	optDescriptionAlias      = "s"
	optInstitution           = "bank"
	optInstitutionAlias      = "b"
	optAccountType           = "accout-type"
	optAccountTypeAlias      = "p"

	objAccount           = "account"
	objAccountAlias      = "a"
	objCategory          = "category"
	objCategoryAlias     = "c"
	objMainCategory      = "main_category"
	objMainCategoryAlias = "m"
	objExchangeRate      = "rate"
	objExchangeRateAlias = "r"
)

func main() {
	// Get error logger
	_, printError := getLoggers()

	// Get config settings
	dataFile, defaultCurrency, err := GetConfigSettings()
	if err != nil {
		printError.Fatalln(err)
	}

	// Parse user commands and flags
	cli.CommandHelpTemplate = `
NAME:
   {{.HelpName}} - {{.Usage}}
USAGE:
   {{.HelpName}}{{if .Subcommands}} [subcommand]{{end}}{{if .Flags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}{{if .Subcommands}}
SUBCOMMANDS:
    {{range .Subcommands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
{{end}}{{ end }}
`

	app := cli.NewApp()
	app.Name = AppName
	app.Usage = "keeps track of your finance"
	app.Version = "2.0.0"
	app.Authors = []cli.Author{
		cli.Author{"Marcin 'Zbroju' Zbroinski", "marcin@zbroinski.net"},
	}

	flagFile := cli.StringFlag{Name: optFile + "," + optFileAlias, Value: dataFile, Usage: "data file"}
	flagID := cli.IntFlag{Name: optID + "," + optIDAlias, Value: NotSetIntValue, Usage: "ID"}
	flagAll := cli.BoolFlag{Name: optAll, Usage: "show all elements, including removed"}
	flagAccount := cli.StringFlag{Name: objAccount + "," + objAccountAlias, Value: NotSetStringValue, Usage: "account name"}
	flagDescription := cli.StringFlag{Name: optDescription + "," + optDescriptionAlias, Value: NotSetStringValue, Usage: "description of the object"}
	flagInstitution := cli.StringFlag{Name: optInstitution + "," + optInstitutionAlias, Value: NotSetStringValue, Usage: "institution (bank) where the account is located"}
	flagAccountType := cli.StringFlag{Name: optAccountType + "," + optAccountTypeAlias, Value: NotSetStringValue, Usage: "type of account: operational/o, savings/s, properties/p, investments/i, loans/l"}
	flagCategory := cli.StringFlag{Name: objCategory + "," + objCategoryAlias, Value: NotSetStringValue, Usage: "category name"}
	flagMainCategory := cli.StringFlag{Name: objMainCategory + "," + objMainCategoryAlias, Value: NotSetStringValue, Usage: "main category name"}
	flagMainCategoryType := cli.StringFlag{Name: optMainCategoryType + "," + optMainCategoryTypeAlias, Value: NotSetStringValue, Usage: "main category type (c/cost, t/transfer, i/income)"}
	flagCurrency := cli.StringFlag{Name: optCurrency + "," + optCurrencyAlias, Value: NotSetStringValue, Usage: "currency"}
	flagCurrencyWithDefault := cli.StringFlag{Name: optCurrency + "," + optCurrencyAlias, Value: defaultCurrency, Usage: "currency"}
	flagCurrencyTo := cli.StringFlag{Name: optCurrencyTo + "," + optCurrencyToAlias, Value: NotSetStringValue, Usage: "currency to"}
	flagExchangeRate := cli.Float64Flag{Name: objExchangeRate + "," + objExchangeRateAlias, Value: NotSetFloatValue, Usage: "currency exchange rate"}

	app.Commands = []cli.Command{
		{Name: cmdInit,
			Aliases: []string{cmdInitAlias},
			Flags:   []cli.Flag{flagFile},
			Usage:   "Init a new data file specified by the user",
			Action:  CmdCreateNewDataFile},
		{Name: cmdAdd, Aliases: []string{cmdAddAlias}, Usage: "Add an object.",
			Subcommands: []cli.Command{
				{Name: objCategory,
					Aliases: []string{objCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagCategory, flagMainCategory},
					Usage:   "Add new category.",
					Action:  CmdCategoryAdd},
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType},
					Usage:   "Add new main category.",
					Action:  CmdMainCategoryAdd},
				{Name: objExchangeRate,
					Aliases: []string{objExchangeRateAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagCurrencyTo, flagExchangeRate},
					Usage:   "Add new currency exchange rate.",
					Action:  CmdExchangeRateAdd},
				{Name: objAccount,
					Aliases: []string{objAccountAlias},
					Flags:   []cli.Flag{flagFile, flagAccount, flagDescription, flagInstitution, flagCurrencyWithDefault, flagAccountType},
					Usage:   "Add new account",
					Action:  CmdAccountAdd},
			},
		},
		{Name: cmdEdit, Aliases: []string{cmdEditAlias}, Usage: "Edit an object.",
			Subcommands: []cli.Command{
				{Name: objCategory,
					Aliases: []string{objCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagCategory, flagMainCategory},
					Usage:   "Edit category.",
					Action:  CmdCategoryEdit},
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagMainCategory, flagMainCategoryType},
					Usage:   "Edit main category.",
					Action:  CmdMainCategoryEdit},
				{Name: objExchangeRate,
					Aliases: []string{objExchangeRateAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagCurrencyTo, flagExchangeRate},
					Usage:   "Edit currency exchange rate.",
					Action:  CmdExchangeRateEdit},
				{Name: objAccount,
					Aliases: []string{objAccountAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagAccount, flagDescription, flagInstitution, flagCurrency, flagAccountType},
					Usage:   "Edit account.",
					Action:  CmdAccountEdit},
			},
		},
		{Name: cmdRemove, Aliases: []string{cmdRemoveAlias}, Usage: "Remove an object.",
			Subcommands: []cli.Command{
				{Name: objCategory,
					Aliases: []string{objCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove category.",
					Action:  CmdCategoryRemove},
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove main category.",
					Action:  CmdMainCategoryRemove},
				{Name: objExchangeRate,
					Aliases: []string{objExchangeRateAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagCurrencyTo},
					Usage:   "Remove currency exchange rate.",
					Action:  CmdExchangeRateRemove},
				{Name: objAccount,
					Aliases: []string{objAccountAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove account.",
					Action:  CmdAccountRemove},
			},
		},
		{Name: cmdList, Aliases: []string{cmdListAlias}, Usage: "List objects on standard output.",
			Subcommands: []cli.Command{
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType, flagAll},
					Usage:   "List main categories.",
					Action:  CmdMainCategoryList},
				{Name: objCategory,
					Aliases: []string{objCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType, flagCategory, flagAll},
					Usage:   "List categories.",
					Action:  CmdCategoryList},
				{Name: objExchangeRate,
					Aliases: []string{objExchangeRateAlias},
					Flags:   []cli.Flag{flagFile},
					Usage:   "List currency exchange rates.",
					Action:  CmdExchangeRateList},
				{Name: objAccount,
					Aliases: []string{objAccountAlias},
					Flags:   []cli.Flag{flagFile, flagAccount, flagDescription, flagInstitution, flagCurrency, flagAccountType, flagAll},
					Usage:   "List accounts.",
					Action:  CmdAccountList},
			},
		},
	}

	app.Run(os.Args)
}

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

// GetLoggers returns two loggers for standard formatting of messages and errors
func getLoggers() (messageLogger *log.Logger, errorLogger *log.Logger) {
	messageLogger = log.New(os.Stdout, fmt.Sprintf("%s: ", AppName), 0)
	errorLogger = log.New(os.Stderr, fmt.Sprintf("%s: ", AppName), 0)

	return
}

// mainCategoryTypeForString returns main category type for given string
func mainCategoryTypeForString(m string) (mct MainCategoryType) {
	switch m {
	case "c", "cost", NotSetStringValue: // If null string is given, then the default value is MCT_Cost
		mct = MCTCost
	case "i", "income":
		mct = MCTIncome
	case "t", "transfer":
		mct = MCTTransfer
	default:
		mct = MCTUnknown
	}

	return mct
}

// accountTypeForString returns account type for given string
func accountTypeForString(s string) (t AccountType) {
	switch s {
	case "o", "operational", NotSetStringValue: // If null string is given then the default value is ATTransactional
		t = ATTransactional
	case "s", "savings":
		t = ATSaving
	case "p", "properties":
		t = ATProperty
	case "i", "investment":
		t = ATInvestment
	case "l", "loans":
		t = ATLoan
	default:
		t = ATUnknown
	}

	return t
}

// getLineFor returns pre-formatted line formatting string for reporting
func lineFor(fs ...string) string {
	line := strings.Join(fs, FSSeparator) + "\n"
	return line
}

// getHFSForText return heading formatting string for string values
func hFSForText(l int) string {
	return fmt.Sprintf("%%-%ds", l)
}

// getHFSForNumeric return heading formatting string for numeric values
func hFSForNumeric(l int) string {
	return fmt.Sprintf("%%%ds", l)
}

// getDFSForText return data formatting string for string
func dFSForText(l int) string {
	return fmt.Sprintf("%%-%ds", l)
}

// getDFSForRates return data formatting string for rates
func getDFSForRates(l int) string {
	return fmt.Sprintf("%%%d.4f", l)
}

// getDFSForValue return data formatting string for values
func getDFSForValue(l int) string {
	return fmt.Sprintf("%%%d.2f", l)
}

// getDFSForID return data formatting string for id
func dFSForID(l int) string {
	return fmt.Sprintf("%%%dd", l)
}

// Return the bigger number out of the two given
func maxLen(s string, i int) int {
	if l := utf8.RuneCountInString(s); l > i {
		return l
	} else {
		return i
	}
}

//DONE: init file
//DONE: account add
//DONE: account edit
//DONE: account close
//DONE: account list
//DONE: category add
//DONE: category edit
//DONE: category remove
//DONE: category list
//DONE: currency add
//DONE: currency edit
//DONE: currency remove
//DONE: currency list
//DONE: main category add
//DONE: main category edit
//DONE: main category remove
//DONE: main category list
//TODO: transaction add
//TODO: transaction edit
//TODO: transaction remove
//TODO: transaction list
//TODO: budget add
//TODO: budget edit
//TODO: budget remove
//TODO: budget list
//TODO: report accounts balance
//TODO: report budget categories
//TODO: report assets summary
//TODO: report budget main categories
//TODO: report categories balance
//TODO: report main categories balance
//TODO: report transaction balance
//TODO: report net value
//
//DONE: 17/33 (51%)

// IDEAS
//TODO: add 'tag' or 'cost center' to transactions attribute (as a separate object)
//TODO: add to main_categories column with 'coefficient', which will be used for calculations instead of signs in transactions (because of that we can have a real main category edit with correct type change)
//TODO: add condition to mainCategoryRemove checking if there are any transactions/categories connected and if not, remove it completely
//TODO: check all operations to see if there is checking if given object exists (e.g. before removing or updating an object)
//TODO: add automatic keeping number of backup copies (the number specified in config file)
//TODO: add export to csv any list and report

//FIXME: make all operations on currencies case insensitive
//FIXME: look through all functions if there is 'return error.new' and not only 'error.new'
//FIXME: make all object private (requires 'ObjectNew' functions)
//FIXME: check if all 'list' functions respect flag --all
//TODO: review all comments inside function bodies
//TODO: complete function descriptions for godoc
