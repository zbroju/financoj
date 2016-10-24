// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/zbroju/gsqlitehandler"
	"time"
)

// Transaction represents the basic object for transaction
type Transaction struct {
	Id          int64
	Date        time.Time
	Category    *Category
	Account     *Account
	Value       float64
	Description string
}

func TransactionNew() *Transaction {
	t := new(Transaction)
	t.Date = time.Now()
	t.Category = CategoryNew()
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

// TransactionForID returns pointer to Transaction for given id
func TransactionForID(db *gsqlitehandler.SqliteDB, i int) (t *Transaction, err error) {
	var stmt *sql.Stmt

	sqlQuery := "SELECT t.id, t.date, t.description, t.value, a.id, a.name, a.description, a.institution, a.currency, a.type, a.status, c.id, c.name, c.status, m.id, m.type, m.name, m.status FROM transactions t INNER JOIN accounts a ON t.account_id=a.id INNER JOIN categories c ON t.category_id=c.id INNER JOIN main_categories m ON c.main_category_id=m.id WHERE t.id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	t = TransactionNew()
	var tmpDate string
	if err = stmt.QueryRow(i).Scan(&t.Id, &tmpDate, &t.Description, &t.Value, &t.Account.Id, &t.Account.Name, &t.Account.Description, &t.Account.Institution, &t.Account.Currency, &t.Account.AType, &t.Account.Status, &t.Category.Id, &t.Category.Name, &t.Category.Status, &t.Category.Main.Id, &t.Category.Main.MType, &t.Category.Main.Name, &t.Category.Main.Status); err != nil {
		return nil, errors.New(errTransactionWithIDNone)
	}
	if t.Date, err = time.Parse(DateFormat, tmpDate); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	return t, nil
	//TODO: add test
}

// TransactionRemove removes given transaction completely from data file
func TransactionRemove(db *gsqlitehandler.SqliteDB, t *Transaction) error {
	var err error
	var stmt *sql.Stmt

	// Remove transaction
	sqlQuery := "DELETE FROM transactions WHERE id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}
