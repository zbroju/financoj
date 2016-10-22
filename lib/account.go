// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// AccountType describes the type of an account
type AccountType int

const (
	ATUnknown       = -1
	ATUnset         = 0
	ATTransactional = 1
	ATSaving        = 2
	ATProperty      = 3
	ATInvestment    = 4
	ATLoan          = 5
)

// String satisfies fmt.Stringer interface in order to get human readambe names of account type
func (at AccountType) String() string {
	var name string

	switch at {
	case ATUnknown:
		name = "Unknown"
	case ATUnset:
		name = "Not set"
	case ATTransactional:
		name = "Operational"
	case ATSaving:
		name = "Savings"
	case ATProperty:
		name = "Properties"
	case ATInvestment:
		name = "Invenstments"
	case ATLoan:
		name = "Loans"
	}

	return name
}

// Account represents the basic object for account
type Account struct {
	Id          int64
	Name        string
	Description string
	Institution string
	Currency    string
	AType       AccountType
	Status      ItemStatus
}

// AccountAdd adds new account
func AccountAdd(db *gsqlitehandler.SqliteDB, a *Account) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("INSERT INTO accounts VALUES (NULL, ?, ?, ?, ?, ?, ?);"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(a.Name, a.Description, a.Institution, a.Currency, a.AType, a.Status); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// AccountList returns all accounts from file as closure
func AccountList(db *gsqlitehandler.SqliteDB, n string, d string, i string, c string, t AccountType, s ItemStatus) (f func() *Account, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	// Parse filtering criteria
	if n == NotSetStringValue {
		n = noParameterValueForSQL
	} else {
		n = "%" + n + "%"
	}
	if d == NotSetStringValue {
		d = noParameterValueForSQL
	} else {
		d = "%" + d + "%"
	}
	if i == NotSetStringValue {
		i = noParameterValueForSQL
	} else {
		i = "%" + i + "%"
	}
	if c == NotSetStringValue {
		c = noParameterValueForSQL
	} else {
		c = "%" + c + "%"
	}

	// Create and execute query
	sqlQuery := "SELECT id, name, description, institution, currency, type, status FROM accounts WHERE 1=1 " +
		"AND (name LIKE ? OR ?=?) " +
		"AND (description LIKE ? OR ?=?) " +
		"AND (institution LIKE ? OR ?=?) " +
		"AND (currency LIKE ? OR ?=?) " +
		"AND (type=? OR ?=?) " +
		"AND (status=? OR ?=?) " +
		"ORDER BY name ASC;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	if rows, err = stmt.Query(n, n, noParameterValueForSQL, d, d, noParameterValueForSQL, i, i, noParameterValueForSQL, c, c, noParameterValueForSQL, t, t, ATUnset, s, s, ISUnset); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *Account {
		if rows.Next() {
			a := new(Account)
			rows.Scan(&a.Id, &a.Name, &a.Description, &a.Institution, &a.Currency, &a.AType, &a.Status)
			return a
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

// AccountForID returns pointer to the Account for given id
func AccountForID(db *gsqlitehandler.SqliteDB, i int) (a *Account, err error) {
	var stmt *sql.Stmt

	sqlQuery := "SELECT id, name, description, institution, currency, type, status FROM accounts WHERE id=? AND status=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return nil, errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	a = new(Account)
	if err = stmt.QueryRow(i, ISOpen).Scan(&a.Id, &a.Name, &a.Description, &a.Institution, &a.Currency, &a.AType, &a.Status); err != nil {
		return nil, errors.New(errAccountWithIDNone)
	}

	return a, nil
	//TODO: add test
}

// AccountRemove updates given account status with ISClose
func AccountRemove(db *gsqlitehandler.SqliteDB, a *Account) error {
	var err error
	var stmt *sql.Stmt

	// Set correct status (ISClose)
	sqlQuery := "UPDATE accounts SET status=? WHERE id=?;"
	if stmt, err = db.Handler.Prepare(sqlQuery); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(ISClose, a.Id); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}
