// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

// Application internal settings
const (
	AppName       = "financoj"
	ConfigFile    = ".financojrc"
	FSSeparator   = "  "
	NullDataValue = "-"

	NotSetIntValue    int     = -1
	NotSetFloatValue  float64 = -1
	NotSetStringValue         = ""
)

// DB Properties
var dataFileProperties = map[string]string{
	"applicationName": "financoj",
	"databaseVersion": "2.0",
}

// Config file settings
const (
	confDataFile = "DATA_FILE"
)

// Error messages
const (
	errMissingFileFlag = "missing information about data file. Specify it with --file or -f flag"
)
