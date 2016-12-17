// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package lib

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
	"math"
	"time"
)

// SQL queries
const (
	sqlTransactionAdd string = "INSERT INTO transactions VALUES(NULL, ?, ?, ?, ?, ?);"
)

// Transaction represents the basic object for transaction
type Transaction struct {
	Id       int64
	Date     time.Time
	Category *Category
	Account  *Account

	// Value holds the values with the same sign as user typed it. The same goes to database file.
	// To classify transaction as cost, income (i.e. negative or positive value) you need to check
	// the transactions main category (attribute of category) type factor.
	// To get transaction value with sign use the method GetSValue().
	// This principle is important in order to correctly re-classify transactions in case their
	// categories and/or main categories change.
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

	if stmt, err = db.Handler.Prepare(sqlTransactionAdd); err != nil {
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
		"WHERE (t.date>=? OR ?=?) AND (t.date<=? OR ?=?) AND (a.id=? OR ?=?) AND (t.description LIKE ? OR ?=?) AND (c.id=? OR ?=?) AND (m.id=? OR ?=?) " +
		"ORDER BY t.date, t.id;"
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

// TransactionEdit updates transaction with new values.
// All fields are updated, so make sure you pass old values in argument t.
func TransactionEdit(db *gsqlitehandler.SqliteDB, t *Transaction) error {
	var err error
	var stmt *sql.Stmt

	sqlQuery := "UPDATE transactions " +
		"SET date=?, account_id=?, description=?, value=?, category_id=? " +
		"WHERE id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t.Date.Format(DateFormat), t.Account.Id, t.Description, t.Value, t.Category.Id, t.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
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

// CompoundTransferAdd adds two transactions with NonBudgetary category 'transfer'.
// It should be used to transfer money between accounts.
func CompoundTransferAdd(db *gsqlitehandler.SqliteDB, date time.Time, accFrom, accTo *Account, value float64, description string, e *ExchangeRate) error {
	var err error
	var tx *sql.Tx
	var stmt *sql.Stmt

	// Create two separate transactions
	tMinus, tPlus := TransactionNew(), TransactionNew()

	tMinus.Date, tPlus.Date = date, date
	tMinus.Category.Id, tPlus.Category.Id = SOCategoryTransferID, SOCategoryTransferID
	tMinus.Account, tPlus.Account = accFrom, accTo
	tMinus.Value, tPlus.Value = -value, value*e.Rate
	tMinus.Description, tPlus.Description = description, description

	// Save transactions to DB
	if tx, err = db.Handler.Begin(); err != nil {
		return errors.New(errWritingToFile)
	}

	if stmt, err = tx.Prepare(sqlTransactionAdd); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(tMinus.Date.Format(DateFormat), tMinus.Account.Id, tMinus.Description, tMinus.Value, tMinus.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}
	if _, err = stmt.Exec(tPlus.Date.Format(DateFormat), tPlus.Account.Id, tPlus.Description, tPlus.Value, tPlus.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	tx.Commit()

	return nil

	//TODO: add test
}

// CompoundInternalCostAdd adds two transactions with NonBudgetary category 'transfer'.
// It should be used to transfer money between accounts.
func CompoundInternalCostAdd(db *gsqlitehandler.SqliteDB, date time.Time, c *Category, accCost, accTransfer *Account, value float64, description string, e *ExchangeRate) error {
	var err error
	var tx *sql.Tx
	var stmt *sql.Stmt

	// Create two separate transactions
	tCost, tTransfer := TransactionNew(), TransactionNew()

	tCost.Date, tTransfer.Date = date, date
	tCost.Category.Id, tTransfer.Category.Id = c.Id, SOCategoryTransferID
	tCost.Account, tTransfer.Account = accCost, accTransfer
	tCost.Value, tTransfer.Value = value, float64(-1)*float64(c.Main.MType.Factor)*value*e.Rate
	tCost.Description, tTransfer.Description = description, description

	// Save transactions to DB
	if tx, err = db.Handler.Begin(); err != nil {
		return errors.New(errWritingToFile)
	}

	if stmt, err = tx.Prepare(sqlTransactionAdd); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(tCost.Date.Format(DateFormat), tCost.Account.Id, tCost.Description, tCost.Value, tCost.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}
	if _, err = stmt.Exec(tTransfer.Date.Format(DateFormat), tTransfer.Account.Id, tTransfer.Description, tTransfer.Value, tTransfer.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	tx.Commit()

	return nil

	//TODO: add test
}

// CompountSplitAdd adds two transactions for two different categories with half of the value each.
func CompoundSplitAdd(db *gsqlitehandler.SqliteDB, d time.Time, a *Account, value float64, description string, c1, c2 *Category) error {
	var err error
	var tx *sql.Tx
	var stmt *sql.Stmt

	// Create two separate transactions
	t1, t2 := TransactionNew(), TransactionNew()
	if !d.IsZero() {
		t1.Date, t2.Date = d, d
	}
	t1.Account, t2.Account = a, a
	t1.Description, t2.Description = description, description
	t1.Value, t2.Value = splitValue(value)
	t1.Category, t2.Category = c1, c2

	// Save transactions to DB
	if tx, err = db.Handler.Begin(); err != nil {
		return errors.New(errWritingToFile)
	}
	if stmt, err = tx.Prepare(sqlTransactionAdd); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(t1.Date.Format(DateFormat), t1.Account.Id, t1.Description, t1.Value, t1.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}
	if _, err = stmt.Exec(t2.Date.Format(DateFormat), t2.Account.Id, t2.Description, t2.Value, t2.Category.Id); err != nil {
		return errors.New(errWritingToFile)
	}
	tx.Commit()

	return nil

	//TODO: add test
}

// splitValue splits given value into two values without reminder
func splitValue(value float64) (v1, v2 float64) {
	var tmpTotV, tmpV1, tmpV2 int64

	tmpTotV = int64(math.Ceil(value * 100))
	if tmpTotV%2 == 0 {
		tmpV1 = tmpTotV / 2
		tmpV2 = tmpV1
	} else {
		tmpTotV -= 1
		tmpV1 = tmpTotV / 2
		tmpV2 = tmpV1 + 1
	}

	v1 = float64(tmpV1) / 100
	v2 = float64(tmpV2) / 100

	return v1, v2

	//TODO: add test
}
