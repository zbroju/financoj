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
	Id       int64
	Date     time.Time
	Category *Category
	Account  *Account
	// Value holds the values with the same sign as user typed it.
	// The same goes to database file.
	// To classify transaction as cost, income (i.e. negative or positive value) you need to check
	// the transactions main category (attribute of category) type factor.
	// To get transaction value with sign use the method GetSValue().
	// This principle is important in order to correctly re-classify transactions in case their
	// categories or main categories change.
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

func (t *Transaction) GetSValue() float64 {
	return t.Value * float64(t.Category.Main.MType.Factor)
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

	sqlQuery := "SELECT t.id, t.date, t.description, t.value, a.id, a.name, a.description, a.institution, a.currency, a.type, a.status, c.id, c.name, c.status, m.id, m.name, m.status, mt.id, mt.name, mt.factor " +
		"FROM transactions t INNER JOIN accounts a ON t.account_id=a.id INNER JOIN categories c ON t.category_id=c.id INNER JOIN main_categories m ON c.main_category_id=m.id INNER JOIN main_categories_types mt on m.type_id=mt.id " +
		"WHERE t.id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	t = TransactionNew()
	var tmpDate string
	if err = stmt.QueryRow(i).Scan(&t.Id, &tmpDate, &t.Description, &t.Value, &t.Account.Id, &t.Account.Name, &t.Account.Description, &t.Account.Institution, &t.Account.Currency, &t.Account.AType, &t.Account.Status, &t.Category.Id, &t.Category.Name, &t.Category.Status, &t.Category.Main.Id, &t.Category.Main.Name, &t.Category.Main.Status, &t.Category.Main.MType.Id, &t.Category.Main.MType.Name, &t.Category.Main.MType.Factor); err != nil {
		return nil, errors.New(errTransactionWithIDNone)
	}
	if t.Date, err = time.Parse(DateFormat, tmpDate); err != nil {
		return nil, err
	}

	return t, nil
	//TODO: add test
}

// TransactionList returns all transactions from file as closure
func TransactionList(db *gsqlitehandler.SqliteDB, dateF, dateT time.Time, a *Account, description string, c *Category, m *MainCategory) (f func() *Transaction, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Prepare filtering parameters
	df, dt := noStringParamForSQL, noStringParamForSQL
	if !dateF.IsZero() {
		df = dateF.Format(DateFormat)
	}
	if !dateT.IsZero() {
		dt = dateT.Format(DateFormat)
	}
	if description == NotSetStringValue {
		description = noStringParamForSQL
	} else {
		description = "%" + description + "%"
	}
	var aId int64
	if a == nil {
		aId = noIntParamForSQL
	} else {
		aId = a.Id
	}
	var cId int64
	if c == nil {
		cId = noIntParamForSQL
	} else {
		cId = c.Id
	}
	var mId int64
	if m == nil {
		mId = noIntParamForSQL
	} else {
		mId = m.Id
	}

	// Prepare query
	sqlQuery := "SELECT t.id, t.date, t.description, t.value, a.id, a.name, a.description, a.institution, a.currency, a.type, a.status, c.id, c.name, c.status, m.id, m.name, m.status, mt.id, mt.name, mt.factor " +
		"FROM transactions t INNER JOIN accounts a ON t.account_id=a.id INNER JOIN categories c ON t.category_id=c.id INNER JOIN main_categories m ON c.main_category_id=m.id INNER JOIN main_categories_types mt on m.type_id=mt.id " +
		"WHERE (t.date>=? OR ?=?) AND (t.date<=? OR ?=?) AND (a.id=? OR ?=?) AND (t.description LIKE ? OR ?=?) AND (c.id=? OR ?=?) AND (m.id=? OR ?=?);"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	if rows, err = stmt.Query(df, df, noStringParamForSQL, dt, dt, noStringParamForSQL, aId, aId, NotSetIntValue, description, description, noStringParamForSQL, cId, cId, NotSetIntValue, mId, mId, NotSetIntValue); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	// Create closure
	f = func() *Transaction {
		if rows.Next() {
			t := TransactionNew()
			var tmpDate string
			rows.Scan(&t.Id, &tmpDate, &t.Description, &t.Value, &t.Account.Id, &t.Account.Name, &t.Account.Description, &t.Account.Institution, &t.Account.Currency, &t.Account.AType, &t.Account.Status, &t.Category.Id, &t.Category.Name, &t.Category.Status, &t.Category.Main.Id, &t.Category.Main.Name, &t.Category.Main.Status, &t.Category.Main.MType.Id, &t.Category.Main.MType.Name, &t.Category.Main.MType.Factor)
			if t.Date, err = time.Parse(DateFormat, tmpDate); err != nil {
				t.Date = time.Time{}
			}
			return t
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
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
