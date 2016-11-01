// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

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
	b.Category = new(Category)

	return b
}

// BudgetAdd adds a new budget
func BudgetAdd(db *gsqlitehandler.SqliteDB, b *Budget) error {
	var err error
	var stmt *sql.Stmt

	sqlQuery := "INSERT INTO budgets VALUES (?, ?, ?, round(?,2), ?);"
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
