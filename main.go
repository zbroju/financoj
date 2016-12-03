// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"github.com/urfave/cli"
	. "github.com/zbroju/financoj/cmd"
	. "github.com/zbroju/financoj/lib"
	"os"
)

func main() {
	// Get error logger
	_, printError := GetLoggers()

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

	flagFile := cli.StringFlag{Name: OptFile + "," + OptFileAlias, Value: dataFile, Usage: "data file"}
	flagID := cli.IntFlag{Name: OptID + "," + OptIDAlias, Value: NotSetIntValue, Usage: "ID"}
	flagAll := cli.BoolFlag{Name: OptAll, Usage: "show all elements, including removed"}
	flagAccount := cli.StringFlag{Name: ObjAccount + "," + ObjAccountAlias, Value: NotSetStringValue, Usage: "account name"}
	flagDescription := cli.StringFlag{Name: OptDescription + "," + OptDescriptionAlias, Value: NotSetStringValue, Usage: "description of the object"}
	flagInstitution := cli.StringFlag{Name: OptInstitution + "," + OptInstitutionAlias, Value: NotSetStringValue, Usage: "institution (bank) where the account is located"}
	flagAccountType := cli.StringFlag{Name: OptAccountType + "," + OptAccountTypeAlias, Value: NotSetStringValue, Usage: "type of account: operational/o, savings/s, properties/p, investments/i, loans/l"}
	flagCategory := cli.StringFlag{Name: ObjCategory + "," + ObjCategoryAlias, Value: NotSetStringValue, Usage: "category name"}
	flagMainCategory := cli.StringFlag{Name: ObjMainCategory + "," + ObjMainCategoryAlias, Value: NotSetStringValue, Usage: "main category name"}
	flagMainCategoryType := cli.StringFlag{Name: OptMainCategoryType + "," + OptMainCategoryTypeAlias, Value: NotSetStringValue, Usage: "main category type (cost, transfer, income)"}
	flagCurrency := cli.StringFlag{Name: OptCurrency + "," + OptCurrencyAlias, Value: NotSetStringValue, Usage: "currency"}
	flagCurrencyWithDefault := cli.StringFlag{Name: OptCurrency + "," + OptCurrencyAlias, Value: defaultCurrency, Usage: "currency"}
	flagCurrencyTo := cli.StringFlag{Name: OptCurrencyTo + "," + OptCurrencyToAlias, Value: NotSetStringValue, Usage: "currency to"}
	flagExchangeRate := cli.Float64Flag{Name: ObjExchangeRate + "," + ObjExchangeRateAlias, Value: NotSetFloatValue, Usage: "currency exchange rate"}
	flagValue := cli.Float64Flag{Name: OptValue + "," + OptValueAlias, Value: NotSetFloatValue, Usage: "value"}
	flagDate := cli.StringFlag{Name: OptDate + "," + OptDateAlias, Value: NotSetStringValue, Usage: "date"}
	flagDateFrom := cli.StringFlag{Name: OptDateFrom, Value: NotSetStringValue, Usage: "date from"}
	flagDateTo := cli.StringFlag{Name: OptDateTo, Value: NotSetStringValue, Usage: "date to"}
	flagPeriod := cli.StringFlag{Name: OptPeriod + "," + OptPeriodAlias, Value: NotSetStringValue, Usage: "year-month period (yyyy-mm)"}

	app.Commands = []cli.Command{
		{Name: CmdInit,
			Aliases: []string{CmdInitAlias},
			Flags:   []cli.Flag{flagFile},
			Usage:   "Init a new data file specified by the user",
			Action:  CmdCreateNewDataFile},
		{Name: CmdAdd, Aliases: []string{CmdAddAlias}, Usage: "Add an object.",
			Subcommands: []cli.Command{
				{Name: ObjCategory,
					Aliases: []string{ObjCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagCategory, flagMainCategory},
					Usage:   "Add new category.",
					Action:  CmdCategoryAdd},
				{Name: ObjMainCategory,
					Aliases: []string{ObjMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType},
					Usage:   "Add new main category.",
					Action:  CmdMainCategoryAdd},
				{Name: ObjExchangeRate,
					Aliases: []string{ObjExchangeRateAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagCurrencyTo, flagExchangeRate},
					Usage:   "Add new currency exchange rate.",
					Action:  CmdExchangeRateAdd},
				{Name: ObjAccount,
					Aliases: []string{ObjAccountAlias},
					Flags:   []cli.Flag{flagFile, flagAccount, flagDescription, flagInstitution, flagCurrencyWithDefault, flagAccountType},
					Usage:   "Add new account",
					Action:  CmdAccountAdd},
				{Name: ObjTransaction,
					Aliases: []string{ObjTransactionAlias},
					Flags:   []cli.Flag{flagFile, flagDescription, flagValue, flagAccount, flagCategory, flagDate},
					Usage:   "Add new transaction.",
					Action:  CmdTransactionAdd},
				{Name: ObjBudget,
					Aliases: []string{ObjBudgetAlias},
					Flags:   []cli.Flag{flagFile, flagPeriod, flagCategory, flagValue, flagCurrencyWithDefault},
					Usage:   "Add new budget.",
					Action:  CmdBudgetAdd},
			},
		},
		{Name: CmdEdit, Aliases: []string{CmdEditAlias}, Usage: "Edit an object.",
			Subcommands: []cli.Command{
				{Name: ObjCategory,
					Aliases: []string{ObjCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagCategory, flagMainCategory},
					Usage:   "Edit category.",
					Action:  CmdCategoryEdit},
				{Name: ObjMainCategory,
					Aliases: []string{ObjMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagMainCategory, flagMainCategoryType},
					Usage:   "Edit main category.",
					Action:  CmdMainCategoryEdit},
				{Name: ObjExchangeRate,
					Aliases: []string{ObjExchangeRateAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagCurrencyTo, flagExchangeRate},
					Usage:   "Edit currency exchange rate.",
					Action:  CmdExchangeRateEdit},
				{Name: ObjAccount,
					Aliases: []string{ObjAccountAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagAccount, flagDescription, flagInstitution, flagCurrency, flagAccountType},
					Usage:   "Edit account.",
					Action:  CmdAccountEdit},
				{Name: ObjTransaction,
					Aliases: []string{ObjTransactionAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagDate, flagCategory, flagAccount, flagValue, flagDescription},
					Usage:   "Edit transaction.",
					Action:  CmdTransactionEdit},
				{Name: ObjBudget,
					Aliases: []string{ObjBudgetAlias},
					Flags:   []cli.Flag{flagFile, flagPeriod, flagCategory, flagValue, flagCurrency},
					Usage:   "Edit budget.",
					Action:  CmdBudgetEdit},
			},
		},
		{Name: CmdRemove, Aliases: []string{CmdRemoveAlias}, Usage: "Remove an object.",
			Subcommands: []cli.Command{
				{Name: ObjCategory,
					Aliases: []string{ObjCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove category.",
					Action:  CmdCategoryRemove},
				{Name: ObjMainCategory,
					Aliases: []string{ObjMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove main category.",
					Action:  CmdMainCategoryRemove},
				{Name: ObjExchangeRate,
					Aliases: []string{ObjExchangeRateAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagCurrencyTo},
					Usage:   "Remove currency exchange rate.",
					Action:  CmdExchangeRateRemove},
				{Name: ObjAccount,
					Aliases: []string{ObjAccountAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove account.",
					Action:  CmdAccountRemove},
				{Name: ObjTransaction,
					Aliases: []string{ObjTransactionAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove transaction.",
					Action:  CmdTransactionRemove},
				{Name: ObjBudget,
					Aliases: []string{ObjBudgetAlias},
					Flags:   []cli.Flag{flagFile, flagPeriod, flagCategory},
					Usage:   "Remove budget.",
					Action:  CmdBudgetRemove},
			},
		},
		{Name: CmdList, Aliases: []string{CmdListAlias}, Usage: "List objects on standard output.",
			Subcommands: []cli.Command{
				{Name: ObjMainCategory,
					Aliases: []string{ObjMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType, flagAll},
					Usage:   "List main categories.",
					Action:  CmdMainCategoryList},
				{Name: ObjCategory,
					Aliases: []string{ObjCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType, flagCategory, flagAll},
					Usage:   "List categories.",
					Action:  CmdCategoryList},
				{Name: ObjExchangeRate,
					Aliases: []string{ObjExchangeRateAlias},
					Flags:   []cli.Flag{flagFile},
					Usage:   "List currency exchange rates.",
					Action:  CmdExchangeRateList},
				{Name: ObjAccount,
					Aliases: []string{ObjAccountAlias},
					Flags:   []cli.Flag{flagFile, flagAccount, flagDescription, flagInstitution, flagCurrency, flagAccountType, flagAll},
					Usage:   "List accounts.",
					Action:  CmdAccountList},
				{Name: ObjTransaction,
					Aliases: []string{ObjTransactionAlias},
					Flags:   []cli.Flag{flagFile, flagDateFrom, flagDateTo, flagAccount, flagDescription, flagCategory, flagMainCategory},
					Usage:   "List transactions.",
					Action:  CmdTransactionList},
				{Name: ObjBudget,
					Aliases: []string{ObjBudgetAlias},
					Flags:   []cli.Flag{flagFile, flagPeriod, flagCategory},
					Usage:   "List budgets.",
					Action:  CmdBudgetList},
			},
		},
		{Name: CmdReport, Aliases: []string{CmdReportAlias}, Usage: "Prints report.",
			Subcommands: []cli.Command{
				{Name: ObjReportTransactionBalance,
					Aliases: []string{ObjReportTransactionBalanceAlias},
					Flags:   []cli.Flag{flagFile, flagCurrencyWithDefault, flagDateFrom, flagDateTo, flagAccount, flagCategory, flagMainCategory},
					Usage:   "Transactions balance for given criteria.",
					Action:  RepTransactionBalance},
				{Name: ObjReportAccountBalance,
					Aliases: []string{ObjReportAccountBalanceAlias},
					Flags:   []cli.Flag{flagFile, flagDate},
					Usage:   "Accounts balance on given date (or today if date flag missing).",
					Action:  RepAccountBalance},
				{Name: ObjReportBudgetCategories,
					Aliases: []string{ObjReportBudgetCategoriesAlias},
					Flags:   []cli.Flag{flagFile, flagPeriod, flagCurrencyWithDefault},
					Usage:   "Budget categories for given year or year-month (or current month if period flag is missing).",
					Action:  RepBudgetCategories},
				{Name: ObjReportBudgetMainCategories,
					Aliases: []string{ObjReportBudgetMainCategoriesAlias},
					Flags:   []cli.Flag{flagFile, flagPeriod, flagCurrencyWithDefault},
					Usage:   "Budget main categories for given year or year-month (or current month if period flag is missing).",
					Action:  RepBudgetMainCategories},
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
//DONE: transaction add
//DONE: transaction edit
//DONE: transaction remove
//DONE: transaction list
//DONE: budget add
//DONE: budget edit
//DONE: budget remove
//DONE: budget list
//DONE: report accounts balance
//DONE: report budget categories
//TODO: report assets summary
//DONE: report budget main categories
//TODO: report categories balance
//TODO: report main categories balance
//DONE: report transaction balance
//TODO: report net value
//
//DONE: 29/33 (87%)

// IDEAS
//TODO: add 'tag' or 'cost center' to transactions attribute (as a separate object)
//TODO: add condition to mainCategoryRemove checking if there are any transactions/categories connected and if not, remove it completely
//TODO: check all operations to see if there is checking if given object exists (e.g. before removing or updating an object)
//TODO: add automatic keeping number of backup copies (the number specified in config file)
//TODO: add export to csv any list and report
//TODO: add transaction split when transaction is equally split into two categories
//TODO: add transaction 'transfer' to have it in one shot (check currencies!)
//TODO: add transaction 'saving' which is cost (minus) on first account and transfer (plus) on the second one (it's for credit payment and amortization)

//FIXME: make all object private (requires 'ObjectNew' functions)
//FIXME: check if all 'list' functions respect flag --all
//TODO: review all comments inside function bodies and make them more verbose
//TODO: complete function descriptions for godoc
//TODO: add default account (especially to add transaction, but think about others)
//FIXME: change 'errMissing*Flag' to map and create function to easily check missing flags
//FIXME: for each function objectForID and objectForName, change returned error depending on the status of the object: if open -> return the object, if closed or system -> return respective error
//FIXME: move all sql queries to separate file and format them so that they are readable
