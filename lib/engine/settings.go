// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

// Application internal settings
const (
	AppName = "financoj"

	NotSetIntValue    int     = 0
	NotSetFloatValue  float64 = 0.0
	NotSetStringValue         = ""

	DateFormat = "2006-01-02"
)

// Error messages
const (
	errMainCategoryWithIDNone    = "no main category with given ID"
	errMainCategoryWithNameNone  = "no main category with given name"
	errMainCategoryNameAmbiguous = "given main category name is ambiguous"

	errCategoryWithIDNone        = "no category with given ID"
	errCategoryWithNameNone      = "no category with given name"
	errCategoryWithNameAmbiguous = "given category name is ambiguous"

	errExchangeRateNone          = "no exchange rate for given currencies"
	errExchangeRateAlreadyExists = "exchange rate for given currencies already exists"

	errAccountWithIDNone    = "no account with given ID"
	errAccountForNameNone   = "no account with given name"
	errAccountNameAmbiguous = "given account name is ambiguous"

	errWritingToFile   = "error writing to file"
	errReadingFromFile = "error reading from file"
)

// Other constants
const (
	noParameterValueForSQL = "NOT_SET_PARAMETER_VALUE"
)
