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

	HAId          = "ID"
	HAName        = "ACCOUNT"
	HADescription = "DESCRIPTION"
	HAInstitution = "BANK"
	HACurrency    = "CUR"
	HAType        = "TYPE"
	HAStatus      = "STATUS"

	HTId          = "ID"
	HTDate        = "DATE"
	HTValue       = "VALUE"
	HTDescription = "DESCRIPTION"

	HBPeriod     = "PERIOD"
	HBLimit      = "LIMIT"
	HBCurrency   = "CUR"
	HBDifference = "DIFFERENCE"

	HNV = "NET VALUE"
)

// Errors
const (
	errMissingFileFlag           = "missing information about data file"
	errMissingIDFlag             = "missing ID"
	errMissingCategoryFlag       = "missing category name"
	errMissingMainCategoryFlag   = "missing main category name"
	errIncorrectMainCategoryType = "incorrect main category type"
	errMissingCurrencyFlag       = "missing currency (from) name"
	errMissingCurrencyToFlag     = "missing currency_to name"
	errMissingExchangeRateFlag   = "missing exchange rate"
	errMissingAccountFlag        = "missing account name"
	errIncorrectAccountType      = "incorrect account type"
	errMissingDescriptionFlag    = "missing description"
	errMissingValueFlag          = "missing value"
	errMissingPeriodFlag         = "missing period"
)

// Commands, objects and options
const (
	CmdInit        = "init"
	CmdInitAlias   = "I"
	CmdAdd         = "add"
	CmdAddAlias    = "A"
	CmdEdit        = "edit"
	CmdEditAlias   = "E"
	CmdRemove      = "delete"
	CmdRemoveAlias = "D"
	CmdList        = "list"
	CmdListAlias   = "L"
	CmdReport      = "report"
	CmdReportAlias = "R"

	OptFile                  = "file"
	OptFileAlias             = "f"
	OptAll                   = "all"
	OptMainCategoryType      = "main-category-type"
	OptMainCategoryTypeAlias = "o"
	OptID                    = "id"
	OptIDAlias               = "i"
	OptCurrency              = "currency"
	OptCurrencyAlias         = "j"
	OptCurrencyTo            = "currency-to"
	OptCurrencyToAlias       = "k"
	OptDescription           = "description"
	OptDescriptionAlias      = "s"
	OptInstitution           = "bank"
	OptInstitutionAlias      = "b"
	OptAccountTo             = "account-to"
	OptAccountType           = "accout-type"
	OptAccountTypeAlias      = "p"
	OptValue                 = "value"
	OptValueAlias            = "v"
	OptDate                  = "date"
	OptDateAlias             = "d"
	OptDateFrom              = "date-from"
	OptDateTo                = "date-to"
	OptPeriod                = "period"
	OptPeriodAlias           = "e"

	ObjAccount           = "account"
	ObjAccountAlias      = "a"
	ObjCategory          = "category"
	ObjCategoryAlias     = "c"
	ObjMainCategory      = "main_category"
	ObjMainCategoryAlias = "m"
	ObjExchangeRate      = "rate"
	ObjExchangeRateAlias = "r"
	ObjTransaction       = "transaction"
	ObjTransactionAlias  = "t"
	ObjBudget            = "budget"
	ObjBudgetAlias       = "b"

	ObjCompoundTransfer      = "transfer"
	ObjCompoundTransferAlias = "T"

	ObjReportAccountBalance            = "account-balance"
	ObjReportAccountBalanceAlias       = "ab"
	ObjReportBudgetCategories          = "budget-categories"
	ObjReportBudgetCategoriesAlias     = "bc"
	ObjReportBudgetMainCategories      = "budget-main-categories"
	ObjReportBudgetMainCategoriesAlias = "bmc"
	ObjReportTransactionBalance        = "transaction-balance"
	ObjReportTransactionBalanceAlias   = "tb"
	ObjReportCategoryBalance           = "category-balance"
	ObjReportCategoryBalanceAlias      = "cb"
	ObjReportMainCategoryBalance       = "main-category-balance"
	ObjReportMainCategoryBalanceAlias  = "mcb"
	ObjReportAssetsSummary             = "assets-summary"
	ObjReportAssetsSummaryAlias        = "as"
	ObjReportNetValueMonthly           = "net-value"
	ObjReportNetValueMonthlyAlias      = "nv"
)
