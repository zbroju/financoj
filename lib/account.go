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
