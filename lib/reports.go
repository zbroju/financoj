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
