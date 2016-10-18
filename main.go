// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/lib/financoj"
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
	optAllAlias              = "a"
	optMainCategoryType      = "main_category_type"
	optMainCategoryTypeAlias = "o"
	optID                    = "id"
	optIDAlias               = "i"
	optCurrency              = "currency"
	optCurrencyAlias         = "j"
	optCurrencyTo            = "currency_to"
	optCurrencyToAlias       = "k"

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
	flagAll := cli.BoolFlag{Name: optAll + "," + optAllAlias, Usage: "show all elements, including removed"}
	flagCategory := cli.StringFlag{Name: objCategory + "," + objCategoryAlias, Value: NotSetStringValue, Usage: "category name"}
	flagMainCategory := cli.StringFlag{Name: objMainCategory + "," + objMainCategoryAlias, Value: NotSetStringValue, Usage: "main category name"}
	flagMainCategoryType := cli.StringFlag{Name: optMainCategoryType + "," + optMainCategoryTypeAlias, Value: NotSetStringValue, Usage: "main category type (c/cost, t/transfer, i/income)"}
	flagCurrency := cli.StringFlag{Name: optCurrency + "," + optCurrencyAlias, Value: defaultCurrency, Usage: "currency (from)"}
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
					Flags:   []cli.Flag{flagFile, flagCurrency, flagCurrencyTo, flagExchangeRate},
					Usage:   "Add new currency exchange rate.",
					Action:  CmdExchangeRateAdd},
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
					Flags:   []cli.Flag{flagFile, flagCurrency, flagCurrencyTo, flagExchangeRate},
					Usage:   "Edit currency exchange rate.",
					Action:  CmdExchangeRateEdit},
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
					Flags:   []cli.Flag{flagFile, flagCurrency, flagCurrencyTo},
					Usage:   "Remove currency exchange rate.",
					Action:  CmdExchangeRateRemove},
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
	printUserMsg.Printf("changed details of category with id = %d\n", id)

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
		mct = mainCategoryTypeForString(t)
		if mct == MCTUnknown {
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
		lId = maxRune(strconv.FormatInt(ct.Id, 10), lId)
		lType = maxRune(ct.Main.MType.String(), lType)
		lMCat = maxRune(ct.Main.Name, lMCat)
		lCat = maxRune(ct.Name, lCat)
		lStatus = maxRune(ct.Status.String(), lStatus)
	}
	lineH := getLineFor(getHFSForNumeric(lId), getHFSForText(lType), getHFSForText(lMCat), getHFSForText(lCat), getHFSForText(lStatus))
	lineD := getLineFor(getDFSForID(lId), getDFSForText(lType), getDFSForText(lMCat), getDFSForText(lCat), getDFSForText(lStatus))

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
	printUserMsg.Printf("removed main category with id = %d\n", id)

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
		lId = maxRune(strconv.FormatInt(m.Id, 10), lId)
		lType = maxRune(m.MType.String(), lType)
		lName = maxRune(m.Name, lName)
		lStatus = maxRune(m.Status.String(), lStatus)
	}
	lineH := getLineFor(getHFSForNumeric(lId), getHFSForText(lType), getHFSForText(lName), getHFSForText(lStatus))
	lineD := getLineFor(getDFSForID(lId), getDFSForText(lType), getDFSForText(lName), getDFSForText(lStatus))

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
		lCurF = maxRune(cur.CurrencyFrom, lCurF)
		lCurT = maxRune(cur.CurrencyTo, lCurT)
		lRate = maxRune(strconv.FormatFloat(cur.Rate, 'f', -1, 64), lRate)
	}
	lineH := getLineFor(getHFSForText(lCurF), getHFSForText(lCurT), getHFSForNumeric(lRate))
	lineD := getLineFor(getDFSForText(lCurF), getDFSForText(lCurT), getDFSForRates(lRate))

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

// GetLoggers returns two loggers for standard formatting of messages and errors
func getLoggers() (messageLogger *log.Logger, errorLogger *log.Logger) {
	messageLogger = log.New(os.Stdout, fmt.Sprintf("%s: ", AppName), 0)
	errorLogger = log.New(os.Stderr, fmt.Sprintf("%s: ", AppName), 0)

	return
}

// ResolveMainCategoryType returns main category type for given string
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

// getLineFor returns pre-formatted line formatting string for reporting
func getLineFor(fs ...string) string {
	line := strings.Join(fs, FSSeparator) + "\n"
	return line
}

// getHFSForText return heading formatting string for string values
func getHFSForText(l int) string {
	return fmt.Sprintf("%%-%ds", l)
}

// getHFSForNumeric return heading formatting string for numeric values
func getHFSForNumeric(l int) string {
	return fmt.Sprintf("%%%ds", l)
}

// getDFSForText return data formatting string for string
func getDFSForText(l int) string {
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
func getDFSForID(l int) string {
	return fmt.Sprintf("%%%dd", l)
}

// Return the bigger number out of the two given
func maxRune(s string, i int) int {
	if l := utf8.RuneCountInString(s); i < l {
		return l
	} else {
		return i
	}
}

//DONE: init file
//TODO: account add
//TODO: account edit
//TODO: account close
//TODO: account list
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
//DONE: 13/33 (39%)

// IDEAS
//TODO: add 'tag' or 'cost center' to transactions attribute (as a separate object)
//TODO: add to main_categories column with 'coefficient', which will be used for calculations instead of signs in transactions (because of that we can have a real main category edit with correct type change)
//TODO: add condition to mainCategoryRemove checking if there are any transactions/categories connected and if not, remove it completely
//TODO: add automatic keeping number of backup copies (the number specified in config file)
//TODO: add export to csv any list and report
//FIXME: make all operations on currencies case insensitive
//FIXME: look through all functions if there is 'return error.new' and not only 'error.new'
