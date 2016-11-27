// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
	"strconv"
	"time"
)

// AccountBalanceEntry represents one line of the report
type AccountBalanceEntry struct {
	Account *Account
	Value   float64
}

func AccountBalanceEntryNew() *AccountBalanceEntry {
	e := new(AccountBalanceEntry)
	e.Account = new(Account)

	return e
}

func ReportAccountBalance(db *gsqlitehandler.SqliteDB, d time.Time) (f func() *AccountBalanceEntry, err error) {
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
	f = func() *AccountBalanceEntry {
		if rows.Next() {
			e := AccountBalanceEntryNew()
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
type BudgetCategoriesEntry struct {
	Category   *Category
	Limit      float64
	Actual     float64
	Difference float64
}

func BudgetCategoriesEntryNew() *BudgetCategoriesEntry {
	e := new(BudgetCategoriesEntry)
	e.Category = CategoryNew()

	return e
}

func ReportBudgetCategories(db *gsqlitehandler.SqliteDB, p *BPeriod, currency string) (f func() *BudgetCategoriesEntry, err error) {
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

	// Create closure
	f = func() *BudgetCategoriesEntry {
		if rows.Next() {
			e := BudgetCategoriesEntryNew()
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
