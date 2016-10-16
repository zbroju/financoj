// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package cli

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
