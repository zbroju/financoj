// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package lib

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// Budget is primary structure for budget entity
type Budget struct {
	Period   *BPeriod
	Category *Category
	Value    float64
	Currency string
}

// BudgetNew returns pointer to new instance of Budget object
func BudgetNew() *Budget {
	b := new(Budget)
	b.Period = new(BPeriod)
	b.Category = CategoryNew()

	return b
}

// BudgetAdd adds a new budget
func BudgetAdd(db *gsqlitehandler.SqliteDB, b *Budget) error {
	var err error
	var stmt *sql.Stmt

	sqlQuery := "INSERT INTO budgets VALUES (?, ?, ?, round(?,2), upper(?));"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(b.Period.Year, b.Period.Month, b.Category.Id, b.Value, b.Currency); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// BudgetGet returns pointer to Budget for given period and category
func BudgetGet(db *gsqlitehandler.SqliteDB, p *BPeriod, c *Category) (b *Budget, err error) {
	var stmt *sql.Stmt

	sqlQuery := "SELECT b.year, b.month, b.value, b.currency, c.id, c.name, c.status, m.id, m.name, m.status, t.id, t.name, t.factor " +
		"FROM budgets b INNER JOIN categories c ON b.category_id=c.id INNER JOIN main_categories m ON c.main_category_id=m.id INNER JOIN main_categories_types t ON m.type_id=t.id " +
		"WHERE b.year=? AND b.month=? AND b.category_id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	b = BudgetNew()
	if err = stmt.QueryRow(p.Year, p.Month, c.Id).Scan(&b.Period.Year, &b.Period.Month, &b.Value, &b.Currency, &b.Category.Id, &b.Category.Name, &b.Category.Status, &b.Category.Main.Id, &b.Category.Main.Name, &b.Category.Main.Status, &b.Category.Main.MType.Id, &b.Category.Main.MType.Name, &b.Category.Main.MType.Factor); err != nil {
		return nil, errors.New(errBudgetNone)
	}

	return b, nil
	//TODO: add test
}

// BudgetRemove removes given Budget from file
func BudgetRemove(db *gsqlitehandler.SqliteDB, b *Budget) error {
	var err error
	var stmt *sql.Stmt

	// Remove budget
	sqlQuery := "DELETE FROM budgets WHERE year=? AND month=? AND category_id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(b.Period.Year, b.Period.Month, b.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO add test
}

// BudgetEdit updates budget with new values.
// All fields are updated, so make sure you pass old values in argument 'b'.
func BudgetEdit(db *gsqlitehandler.SqliteDB, b *Budget) error {
	var err error
	var stmt *sql.Stmt

	sqlQuery := "UPDATE budgets SET value=?, currency=upper(?) WHERE year=? AND month=? AND category_id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(b.Value, b.Currency, b.Period.Year, b.Period.Month, b.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
}

// BudgetList returns budgets from file as closure
func BudgetList(db *gsqlitehandler.SqliteDB, p *BPeriod, c *Category) (f func() *Budget, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	var y, m int64
	if p == nil {
		y, m = noIntParamForSQL, noIntParamForSQL
	} else {
		y = p.Year
		if p.Month == int64(NotSetIntValue) {
			m = noIntParamForSQL
		} else {
			m = p.Month
		}
	}
	var cId int64
	if c == nil {
		cId = noIntParamForSQL
	} else {
		cId = c.Id
	}

	sqlQuery := "SELECT b.year, b.month, b.value * t.factor as value, b.currency, c.id, c.name, c.status, m.id, m.name, m.status, t.id, t.name, t.factor " +
		"FROM budgets b INNER JOIN categories c ON b.category_id=c.id INNER JOIN main_categories m ON c.main_category_id=m.id INNER JOIN main_categories_types t ON m.type_id=t.id " +
		"WHERE (b.year=? OR ?=?) AND (b.month=? OR ?=?) AND (c.id=? OR ?=?) " +
		"ORDER BY b.year, b.month, t.name, m.name, c.name;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	if rows, err = stmt.Query(y, y, noIntParamForSQL, m, m, noIntParamForSQL, cId, cId, noIntParamForSQL); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *Budget {
		if rows.Next() {
			b := BudgetNew()
			rows.Scan(&b.Period.Year, &b.Period.Month, &b.Value, &b.Currency, &b.Category.Id, &b.Category.Name, &b.Category.Status, &b.Category.Main.Id, &b.Category.Main.Name, &b.Category.Main.Status, &b.Category.Main.MType.Id, &b.Category.Main.MType.Name, &b.Category.Main.MType.Factor)
			return b
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}
