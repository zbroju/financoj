// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package lib

// Application internal settings
const (
	AppName = "financoj"

	NotSetIntValue    int     = 0
	NotSetFloatValue  float64 = 0.0
	NotSetStringValue         = ""

	DateFormat    = "2006-01-02"
	DateSeparator = "-"
)

// Special objects
const (
	// SOMCNonBudgetaryID is a special main category for transfers
	SOMCNonBudgetaryID int64 = 10

	// SOCategoryTransferID is a special category for transfers
	SOCategoryTransferID int64 = 10
)

// Error messages
const (
	errMainCategoryWithIDNone    = "no main category with given ID"
	errMainCategoryWithNameNone  = "no main category with given name"
	errMainCategoryNameAmbiguous = "given main category name is ambiguous"
	errMainCategoryMissing       = "main category missing"

	errCategoryWithIDNone        = "no category with given ID"
	errCategoryWithNameNone      = "no category with given name"
	errCategoryWithNameAmbiguous = "given category name is ambiguous"
	errCategoryMissing           = "category missing"

	errExchangeRateNone          = "no exchange rate for given currencies"
	errExchangeRateAlreadyExists = "exchange rate for given currencies already exists"

	errAccountWithIDNone    = "no account with given ID"
	errAccountForNameNone   = "no account with given name"
	errAccountNameAmbiguous = "given account name is ambiguous"

	errTransactionWithIDNone = "no transaction with given ID"

	errBudgetNone = "no budget"

	errWritingToFile   = "error writing to file"
	errReadingFromFile = "error reading from file"

	errReportMissingCurrencies string = "missing currency exchange rate(s) for: "

	errSystemObject = "this is system object and cannot be changed or removed"
)

// Other constants
const (
	noStringParamForSQL = "NOT_SET_PARAMETER_VALUE"
	noIntParamForSQL    = 0
)
