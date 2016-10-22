// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

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

	HAId          = "ID"
	HAName        = "ACCOUNT"
	HADescription = "DESCRIPTION"
	HAInstitution = "BANK"
	HACurrency    = "CUR"
	HAType        = "TYPE"
	HAStatus      = "STATUS"
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
	errMissingAccountFlag        = "missing account name"
	errIncorrectAccountType      = "incorrect account type"
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
	optMainCategoryType      = "main-category-type"
	optMainCategoryTypeAlias = "o"
	optID                    = "id"
	optIDAlias               = "i"
	optCurrency              = "currency"
	optCurrencyAlias         = "j"
	optCurrencyTo            = "currency_to"
	optCurrencyToAlias       = "k"
	optDescription           = "description"
	optDescriptionAlias      = "s"
	optInstitution           = "bank"
	optInstitutionAlias      = "b"
	optAccountType           = "accout-type"
	optAccountTypeAlias      = "p"

	objAccount           = "account"
	objAccountAlias      = "a"
	objCategory          = "category"
	objCategoryAlias     = "c"
	objMainCategory      = "main_category"
	objMainCategoryAlias = "m"
	objExchangeRate      = "rate"
	objExchangeRateAlias = "r"
)
