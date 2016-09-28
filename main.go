// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/zbroju/gprops"
	"log"
	"os"
	"path"
)

func main() {
	// Get error logger
	_, printError := getLoggers()

	// Get config settings
	dataFile, err := getConfigSettings()
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

	flagFile := cli.StringFlag{Name: optFile + " ," + optFileAlias, Value: dataFile, Usage: "data file"}
	flagID := cli.IntFlag{Name: optID + " ," + optIDAlias, Value: NotSetIntValue, Usage: "ID"}
	flagMainCategory := cli.StringFlag{Name: objMainCategory + ", " + objMainCategoryAlias, Value: NotSetStringValue, Usage: "main category name"}
	flagMainCategoryType := cli.StringFlag{Name: optMainCategoryType + ", " + optMainCategoryTypeAlias, Value: NotSetStringValue, Usage: "main category type (c/cost, t/transfer, i/income)"}

	app.Commands = []cli.Command{
		{Name: cmdInit,
			Aliases: []string{cmdInitAlias},
			Flags:   []cli.Flag{flagFile},
			Usage:   "Init a new data file specified by the user",
			Action:  CreateNewDataFile},
		{Name: cmdAdd, Aliases: []string{cmdAddAlias}, Usage: "Add an object (main_category).",
			Subcommands: []cli.Command{
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagMainCategory, flagMainCategoryType},
					Usage:   "Add new main category.",
					Action:  MainCategoryAdd},
			},
		},
		{Name: cmdEdit, Aliases: []string{cmdEditAlias}, Usage: "Edit an object (main_category).",
			Subcommands: []cli.Command{
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID, flagMainCategory, flagMainCategoryType},
					Usage:   "Edit main category.",
					Action:  MainCategoryEdit},
			},
		},
		{Name: cmdRemove, Aliases: []string{cmdRemoveAlias}, Usage: "Remove an object (main_category).",
			Subcommands: []cli.Command{
				{Name: objMainCategory,
					Aliases: []string{objMainCategoryAlias},
					Flags:   []cli.Flag{flagFile, flagID},
					Usage:   "Remove main category.",
					Action:  MainCategoryRemove},
			},
		},
	}

	app.Run(os.Args)

}

// GetConfigSettings returns contents of settings file
func getConfigSettings() (dataFile string, err error) {
	// Read config file
	configSettings := gprops.New()
	configFile, err := os.Open(path.Join(os.Getenv("HOME"), ConfigFile))
	if err == nil {
		err = configSettings.Load(configFile)
		if err != nil {
			return NotSetStringValue, err
		}
	}
	configFile.Close()
	dataFile = configSettings.GetOrDefault(confDataFile, NotSetStringValue)

	return dataFile, nil
}

// GetLoggers returns two loggers for standard formatting of messages and errors
func getLoggers() (messageLogger *log.Logger, errorLogger *log.Logger) {
	messageLogger = log.New(os.Stdout, fmt.Sprintf("%s: ", AppName), 0)
	errorLogger = log.New(os.Stderr, fmt.Sprintf("%s: ", AppName), 0)

	return
}

//DONE: init file
//TODO: account add
//TODO: account edit
//TODO: account close
//TODO: account list
//TODO: category add
//TODO: category edit
//TODO: category remove
//TODO: category list
//TODO: currency add
//TODO: currency edit
//TODO: currency remove
//TODO: currency list
//DONE: main category add
//DONE: main category edit
//DONE: main category remove
//TODO: main category list
//TODO: transaction add
//TODO: transaction edit
//TODO: transaction remove
//TODO: transaction list
//TODO: budget add
//TODO: budget edit
//TODO: budget remove
//TODO: budget list
//TODO: report accounts balance
//TODO: report assets summary
//TODO: report budget categories
//TODO: report budget main categories
//TODO: report categories balance
//TODO: report main categories balance
//TODO: report transaction balance
//TODO: report net value
//
//DONE: 4/33 (12%)

// IDEAS
//TODO: add to main_categories column with 'coefficient', which will be used for calculations instead of signs in transactions (because of that we can have a real main category edit with correct type change)
//TODO: add condition to mainCategoryRemove checking if there are any transactions/categories connected and if not, remove it completely
