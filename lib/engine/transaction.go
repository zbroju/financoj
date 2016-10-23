// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
	"time"
)

// Transaction represents the basic object for transaction
type Transaction struct {
	Id          int64
	Date         time.Time
	Category    *Category
	Account     *Account
	Value       float64
	Description string
}

func TransactionNew() *Transaction {
	t := new(Transaction)
	t.Date=time.Now()
	t.Category = new(Category)
	t.Account = new(Account)

	return t
}

// TransactionAdd adds new transaction t
func TransactionAdd(db *gsqlitehandler.SqliteDB, t *Transaction) error {
	var err error
	var stmt *sql.Stmt

	sqlQuery := "INSERT INTO transactions VALUES(NULL, ?, ?, ?, ?, ?);"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.Date.Format(DateFormat), t.Account.Id, t.Description, t.Value, t.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil

	//TODO: add test
}
