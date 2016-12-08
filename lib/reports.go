// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
	"strconv"
	"strings"
	"time"
)

// AccountBalanceEntry represents one line of the report
type AccountBalanceReportEntry struct {
	Account *Account
	Value   float64
}

func AccountBalanceReportEntryNew() *AccountBalanceReportEntry {
	e := new(AccountBalanceReportEntry)
	e.Account = new(Account)

	return e
}

func ReportAccountBalance(db *gsqlitehandler.SqliteDB, d time.Time) (f func() *AccountBalanceReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Prepare query
	sqlQuery := `
SELECT
	a.id
	, a.name
	, a.description
	, a.institution
	, a.currency
	, a.type
	, a.status
	, sum(t.value * mct.factor) as value
FROM
	transactions t
	INNER JOIN accounts a ON t.account_id = a.id
	INNER JOIN categories c ON t.category_id = c.id
	INNER JOIN main_categories mc ON c.main_category_id = mc.id
	INNER JOIN main_categories_types mct ON mc.type_id = mct.id
WHERE
	t.date<=?
	AND a.status=?
GROUP BY
	a.id
	, a.name
	, a.description
	, a.institution
	, a.currency
	, a.type
	, a.status
ORDER BY
	a.type
	, a.name
;
`
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	if rows, err = stmt.Query(d.Format(DateFormat), ISOpen); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *AccountBalanceReportEntry {
		if rows.Next() {
			e := AccountBalanceReportEntryNew()
			rows.Scan(&e.Account.Id, &e.Account.Name, &e.Account.Description, &e.Account.Institution, &e.Account.Currency, &e.Account.AType, &e.Account.Status, &e.Value)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

type TransactionBalanceReportEntry struct {
	Transaction *Transaction
	Balance     float64
}

func TransactionBalanceReportEntryNew() *TransactionBalanceReportEntry {
	e := new(TransactionBalanceReportEntry)
	e.Transaction = TransactionNew()

	return e
}

func ReportTransactionBalance(db *gsqlitehandler.SqliteDB, currency string, dateFrom, dateTo time.Time, a *Account, c *Category, m *MainCategory) (f func() *TransactionBalanceReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Check input parameters
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	df := noStringParamForSQL
	if !dateFrom.IsZero() {
		df = dateFrom.Format(DateFormat)
	}
	dt := noStringParamForSQL
	if !dateTo.IsZero() {
		dt = dateTo.Format(DateFormat)
	}
	aId := int64(noIntParamForSQL)
	if a != nil {
		aId = a.Id
	}
	cId := int64(noIntParamForSQL)
	if c != nil {
		cId = c.Id
	}
	mId := int64(noIntParamForSQL)
	if m != nil {
		mId = m.Id
	}

	// Execute main query
	if stmt, err = db.Handler.Prepare(sqlReportTransactionsBalance); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(currency, currency, df, df, noStringParamForSQL, dt, dt, noStringParamForSQL, aId, aId, noIntParamForSQL, cId, cId, noIntParamForSQL, mId, mId, noIntParamForSQL); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *TransactionBalanceReportEntry {
		if rows.Next() {
			e := TransactionBalanceReportEntryNew()
			var tDate string
			rows.Scan(&e.Transaction.Id, &tDate, &e.Transaction.Description, &e.Transaction.Value, &e.Transaction.Account.Id, &e.Transaction.Account.Name, &e.Transaction.Account.Description, &e.Transaction.Account.Institution, &e.Transaction.Account.Currency, &e.Transaction.Account.AType, &e.Transaction.Account.Status, &e.Transaction.Category.Id, &e.Transaction.Category.Name, &e.Transaction.Category.Status, &e.Transaction.Category.Main.Id, &e.Transaction.Category.Main.Name, &e.Transaction.Category.Main.Status, &e.Transaction.Category.Main.MType.Id, &e.Transaction.Category.Main.MType.Name, &e.Transaction.Category.Main.MType.Factor, &e.Balance)
			if e.Transaction.Date, err = time.Parse(DateFormat, tDate); err != nil {
				e.Transaction.Date = time.Time{}
			}
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

type CategoryBalanceReportEntry struct {
	Category *Category
	Balance  float64
}

func CategoryBalanceReportBalanceNew() *CategoryBalanceReportEntry {
	e := new(CategoryBalanceReportEntry)
	e.Category = CategoryNew()

	return e
}

func ReportCategoryBalance(db *gsqlitehandler.SqliteDB, currency string, dateFrom, dateTo time.Time, a *Account, c *Category, m *MainCategory) (f func() *CategoryBalanceReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Check input parameters
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	df := noStringParamForSQL
	if !dateFrom.IsZero() {
		df = dateFrom.Format(DateFormat)
	}
	dt := noStringParamForSQL
	if !dateTo.IsZero() {
		dt = dateTo.Format(DateFormat)
	}
	aId := int64(noIntParamForSQL)
	if a != nil {
		aId = a.Id
	}
	cId := int64(noIntParamForSQL)
	if c != nil {
		cId = c.Id
	}
	mId := int64(noIntParamForSQL)
	if m != nil {
		mId = m.Id
	}

	// Execute main query
	if stmt, err = db.Handler.Prepare(sqlReportCategoriesBalance); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(MCTTransfer, currency, currency, df, df, noStringParamForSQL, dt, dt, noStringParamForSQL, aId, aId, noIntParamForSQL, cId, cId, noIntParamForSQL, mId, mId, noIntParamForSQL); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *CategoryBalanceReportEntry {
		if rows.Next() {
			e := CategoryBalanceReportBalanceNew()
			rows.Scan(&e.Category.Main.Id, &e.Category.Main.Name, &e.Category.Main.Status, &e.Category.Main.MType.Id, &e.Category.Main.MType.Name, &e.Category.Main.MType.Factor, &e.Category.Id, &e.Category.Name, &e.Category.Status, &e.Balance)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

type MainCategoryBalanceReportEntry struct {
	MainCategory *MainCategory
	Balance      float64
}

func MainCategoryBalanceReportEntryNew() *MainCategoryBalanceReportEntry {
	e := new(MainCategoryBalanceReportEntry)
	e.MainCategory = MainCategoryNew()

	return e
}

func ReportMainCategoryBalance(db *gsqlitehandler.SqliteDB, currency string, dateFrom, dateTo time.Time, a *Account, m *MainCategory) (f func() *MainCategoryBalanceReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Check input parameters
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	df := noStringParamForSQL
	if !dateFrom.IsZero() {
		df = dateFrom.Format(DateFormat)
	}
	dt := noStringParamForSQL
	if !dateTo.IsZero() {
		dt = dateTo.Format(DateFormat)
	}
	aId := int64(noIntParamForSQL)
	if a != nil {
		aId = a.Id
	}
	mId := int64(noIntParamForSQL)
	if m != nil {
		mId = m.Id
	}

	// Execute main query
	if stmt, err = db.Handler.Prepare(sqlReportMainCategoriesBalance); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(MCTTransfer, currency, currency, df, df, noStringParamForSQL, dt, dt, noStringParamForSQL, aId, aId, noIntParamForSQL, mId, mId, noIntParamForSQL); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *MainCategoryBalanceReportEntry {
		if rows.Next() {
			e := MainCategoryBalanceReportEntryNew()
			rows.Scan(&e.MainCategory.Id, &e.MainCategory.Name, &e.MainCategory.Status, &e.MainCategory.MType.Id, &e.MainCategory.MType.Name, &e.MainCategory.MType.Factor, &e.Balance)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

type AssetsSummaryReportEntry struct {
	Account *Account
	Balance float64
}

func AssetsSummaryReportEntryNew() *AssetsSummaryReportEntry {
	e := new(AssetsSummaryReportEntry)
	e.Account = new(Account)

	return e
}

func ReportAssetsSummary(db *gsqlitehandler.SqliteDB, currency string, onDate time.Time) (f func() *AssetsSummaryReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Check input parameters
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	dt := noStringParamForSQL
	if onDate.IsZero() {
		dt = time.Now().Format(DateFormat)
	} else {
		dt = onDate.Format(DateFormat)
	}

	// Execute main query
	if stmt, err = db.Handler.Prepare(sqlReportAssetsSummary); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(currency, currency, dt, ISClose); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *AssetsSummaryReportEntry {
		if rows.Next() {
			e := AssetsSummaryReportEntryNew()
			rows.Scan(&e.Account.Id, &e.Account.Name, &e.Account.Description, &e.Account.Institution, &e.Account.Currency, &e.Account.AType, &e.Account.Status, &e.Balance)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

// BudgetCategoryEntry represents one line of the report
type BudgetCategoriesReportEntry struct {
	Category   *Category
	Limit      float64
	Actual     float64
	Difference float64
}

func BudgetCategoriesReportEntryNew() *BudgetCategoriesReportEntry {
	e := new(BudgetCategoriesReportEntry)
	e.Category = CategoryNew()

	return e
}

func ReportBudgetCategories(db *gsqlitehandler.SqliteDB, p *BPeriod, currency string) (f func() *BudgetCategoriesReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	y := int(p.Year)
	m := int(p.Month)
	if m == NotSetIntValue {
		if stmt, err = db.Handler.Prepare(sqlReportBudgetCategoriesYearly); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
		if rows, err = stmt.Query(strconv.Itoa(y), y, MCTTransfer, currency, currency, y, currency, currency, strconv.Itoa(y)); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
	} else {
		if stmt, err = db.Handler.Prepare(sqlReportBudgetCategoriesMonthly); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
		if rows, err = stmt.Query(strconv.Itoa(y), strconv.Itoa(m), y, m, MCTTransfer, currency, currency, y, m, currency, currency, strconv.Itoa(y), strconv.Itoa(m)); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
	}

	// Check if we have all necessary currency exchange rates
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	if s, err := missingCurrenciesForBudgets(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}

	// Create closure
	f = func() *BudgetCategoriesReportEntry {
		if rows.Next() {
			e := BudgetCategoriesReportEntryNew()
			rows.Scan(&e.Category.Main.Id, &e.Category.Main.Name, &e.Category.Main.Status, &e.Category.Main.MType.Id, &e.Category.Main.MType.Name, &e.Category.Main.MType.Factor, &e.Category.Id, &e.Category.Name, &e.Category.Status, &e.Limit, &e.Actual, &e.Difference)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

// BudgetMainCategoryEntry represents one line of the report
type BudgetMainCategoryReportEntry struct {
	MainCategory *MainCategory
	Limit        float64
	Actual       float64
	Difference   float64
}

func BudgetMainCategoryReportEntryNew() *BudgetMainCategoryReportEntry {
	e := new(BudgetMainCategoryReportEntry)
	e.MainCategory = MainCategoryNew()

	return e
}

func ReportBudgetMainCategories(db *gsqlitehandler.SqliteDB, p *BPeriod, currency string) (f func() *BudgetMainCategoryReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	y := int(p.Year)
	m := int(p.Month)
	if m == NotSetIntValue {
		if stmt, err = db.Handler.Prepare(sqlReportBudgetMainCategoriesYearly); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
		if rows, err = stmt.Query(strconv.Itoa(y), y, MCTTransfer, currency, currency, y, strconv.Itoa(y), currency, currency); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
	} else {
		if stmt, err = db.Handler.Prepare(sqlReportBudgetMainCategoriesMonthly); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
		if rows, err = stmt.Query(strconv.Itoa(y), strconv.Itoa(m), y, m, MCTTransfer, currency, currency, y, m, strconv.Itoa(y), strconv.Itoa(m), currency, currency); err != nil {
			return nil, errors.New(errReadingFromFile)
		}
	}

	// Check if we have all necessary currency exchange rates
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	if s, err := missingCurrenciesForBudgets(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}

	// Create closure
	f = func() *BudgetMainCategoryReportEntry {
		if rows.Next() {
			e := BudgetMainCategoryReportEntryNew()
			rows.Scan(&e.MainCategory.Id, &e.MainCategory.Name, &e.MainCategory.Status, &e.MainCategory.MType.Id, &e.MainCategory.MType.Name, &e.MainCategory.MType.Factor, &e.Limit, &e.Actual, &e.Difference)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

type NetValueMonthlyReportEntry struct {
	Period *BPeriod
	Value  float64
}

func NetValueMonthlyReportEntryNew() *NetValueMonthlyReportEntry {
	e := new(NetValueMonthlyReportEntry)
	e.Period = new(BPeriod)

	return e
}

func ReportNetValueMonthly(db *gsqlitehandler.SqliteDB, currency string, dateFrom, dateTo time.Time) (f func() *NetValueMonthlyReportEntry, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Check input parameters
	if s, err := missingCurrenciesForTransactions(db, currency); err == nil {
		if s != nil {
			return nil, errors.New(errReportMissingCurrencies + strings.Join(s, ", "))
		}
	} else {
		return nil, err
	}
	df := noStringParamForSQL
	if !dateFrom.IsZero() {
		df = dateFrom.Format(DateFormat)
	}
	dt := time.Now().Format(DateFormat)
	if !dateTo.IsZero() {
		dt = dateTo.Format(DateFormat)
	}

	// Execute main query
	if stmt, err = db.Handler.Prepare(sqlReportNetValueMonthly); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(currency, currency, df, df, noStringParamForSQL, dt, dt, noStringParamForSQL); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *NetValueMonthlyReportEntry {
		if rows.Next() {
			e := NetValueMonthlyReportEntryNew()
			rows.Scan(&e.Period.Year, &e.Period.Month, &e.Value)
			return e
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

// missingCurrenciesForTransactions returns list of missing currency exchange rates for transactions
// or empty slice if all the currencies exist
func missingCurrenciesForTransactions(db *gsqlitehandler.SqliteDB, c string) (l []string, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if stmt, err = db.Handler.Prepare(sqlReportMissingCurrenciesForTransactions); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()
	if rows, err = stmt.Query(c, c, c); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		rows.Scan(&s)
		l = append(l, s)
	}

	return l, nil
	//TODO: add test
}

// missingCurrenciesForBudgets returns list of missing currency exchange rates for budgets
// or empty slice if all the currencies exist
func missingCurrenciesForBudgets(db *gsqlitehandler.SqliteDB, c string) (l []string, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if stmt, err = db.Handler.Prepare(sqlReportMissingCurrenciesForBudgets); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()
	if rows, err = stmt.Query(c, c, c); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		rows.Scan(&s)
		l = append(l, s)
	}

	return l, nil
	//TODO: add test
}
