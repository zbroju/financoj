// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"github.com/urfave/cli"
	"os"
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
	app.Usage = "keeps track of you bike rides"
	app.Version = "2.0.0"
	app.Authors = []cli.Author{
		cli.Author{"Marcin 'Zbroju' Zbroinski", "marcin@zbroinski.net"},
	}

	flagFile := cli.StringFlag{Name: "file, f", Value: dataFile, Usage: "data file"}
	//flagAccount := cli.StringFlag{Name: "account, a", Value: NotSetStringValue, Usage: "account name"}

	app.Commands = []cli.Command{
		{Name: "init",
			Aliases: []string{"I"},
			Flags:   []cli.Flag{flagFile},
			Usage:   "Init a new data file specified by the user",
			Action:  cmdInit},
		/*{Name: "add", Aliases: []string{"A"}, Usage: "Add an object (account).",
			Subcommands: []cli.Command{
				{Name: objectBicycleType,
					Aliases: []string{objectBicycleTypeAlias},
					Flags:   []cli.Flag{flagFile, flagType},
					Usage:   "Add new bicycle type.",
					Action:  cmdTypeAdd},
			},
		},*/
	}

	app.Run(os.Args)

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
//TODO: main category add
//TODO: main category edit
//TODO: main category remove
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

//DONE: 1/33 (3%)