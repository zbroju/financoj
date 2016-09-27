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

	NotSetIntValue    int     = 0
	NotSetFloatValue  float64 = 0.0
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

// Commands, objects and options
const (
	cmdInit      = "init"
	cmdInitAlias = "I"
	cmdAdd       = "add"
	cmdAddAlias  = "A"
	cmdEdit      = "edit"
	cmdEditAlias = "E"

	optFile                  = "file"
	optFileAlias             = "f"
	optMainCategoryType      = "main_category_type"
	optMainCategoryTypeAlias = "o"
	optID                    = "id"
	optIDAlias               = "i"

	objMainCategory      = "main_category"
	objMainCategoryAlias = "m"
)

// Error messages
const (
	errMissingFileFlag           = "missing information about data file"
	errMissingIDFlag             = "missing ID"
	errMissingMainCategory       = "missing information about main category name"
	errIncorrectMainCategoryType = "incorrect main category type"
	errNoMainCategoryWithID      = "no main category with given ID"

	errWritingToFile   = "error writing to file"
	errReadingFromFile = "error reading from file"
)

// ItemStatus indicates the life cycle of an object
type ItemStatus int

const (
	IS_Close ItemStatus = 0
	IS_Open  ItemStatus = 1
)

// MainCategoryStatusT describes the behaviour of categories and its descendants (transactions)
type MainCategoryTypeT int

const (
	MCT_Unknown  MainCategoryTypeT = -2
	MCT_Cost     MainCategoryTypeT = -1
	MCT_Transfer MainCategoryTypeT = 0
	MCT_Income   MainCategoryTypeT = 1
)

// Satisfy fmt.Stringer interface in order to get human readable names
func (mct MainCategoryTypeT) String() string {
	var mctName string

	switch mct {
	case MCT_Income:
		mctName = "income"
	case MCT_Cost:
		mctName = "cost"
	case MCT_Transfer:
		mctName = "transfer"
	default:
		mctName = "not set"
	}

	return mctName
}

// ResolveMainCategoryType returns main category type for given string
func mainCategoryTypeForString(m string) (mct MainCategoryTypeT) {
	switch m {
	case "c", "cost", NotSetStringValue: // If null string is given, then the default value is MCT_Cost
		mct = MCT_Cost
	case "i", "income":
		mct = MCT_Income
	case "t", "transfer":
		mct = MCT_Transfer
	default:
		mct = MCT_Unknown
	}

	return mct
}
