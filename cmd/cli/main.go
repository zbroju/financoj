// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/lib"
	"os"
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
