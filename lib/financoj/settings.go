// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

// Application internal settings
const (
	AppName = "financoj"

	NotSetIntValue    int     = 0
	NotSetFloatValue  float64 = 0.0
	NotSetStringValue         = ""
)

// Config file settings
const (
	configFile = ".financojrc"

	confDataFile = "DATA_FILE"
	confCurrency = "DEFAULT_CURRENCY"
)

// DB Properties
var dataFileProperties = map[string]string{
	"applicationName": "financoj",
	"databaseVersion": "2.0",
}

// Error messages
const (
	errMainCategoryWithIDNone        = "no main category with given ID"
	errMainCategoryWithNameNone      = "no main category with given name"
	errMainCategoryWithNameAmbiguous = "given main category name is ambiguous"

	errCategoryWithIDNone = "no category with given ID"

	errWritingToFile   = "error writing to file"
	errReadingFromFile = "error reading from file"
)

// Other constants
const (
	noParameterValueForSQL = "NOT_SET_PARAMETER_VALUE"
)
