// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

import (
	"database/sql"
	"errors"
	"github.com/zbroju/gsqlitehandler"
)

// Currency represents the object of currencies exchange rate
type ExchangeRate struct {
	CurrencyFrom string
	CurrencyTo   string
	ExchangeRate float64
}

// ExchangeRateAdd add new currency exchange rate
func ExchangeRateAdd(db *gsqlitehandler.SqliteDB, e *ExchangeRate) error {
	var err error
	var stmt *sql.Stmt

	// Check if such currency exchange rate exists
	if c, _ := ExchangeRateForCurrencies(db, e.CurrencyFrom, e.CurrencyTo); c != nil {
		return errors.New(errExchangeRateAlreadyExists)
	}

	// Add new currency exchange rate
	if stmt, err = db.Handler.Prepare("INSERT into currencies VALUES (?, ?, round(?,4));"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(e.CurrencyFrom, e.CurrencyTo, e.ExchangeRate); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// ExchangeRateEdit updates currency exchange rates for given currencies
func ExchangeRateEdit(db *gsqlitehandler.SqliteDB, e *ExchangeRate) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("UPDATE currencies SET exchange_rate=round(?,4) WHERE currency_from=? AND currency_to=?;"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(e.ExchangeRate, e.CurrencyFrom, e.CurrencyTo); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}

// ExchangeRateForCurrencies returns pointer to ExchangeRateT for given currency_from and currency_to
func ExchangeRateForCurrencies(db *gsqlitehandler.SqliteDB, cf string, ct string) (e *ExchangeRate, err error) {
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("SELECT currency_from, currency_to, exchange_rate FROM currencies WHERE currency_from=? AND currency_to=?;"); err != nil {
		errors.New(errReadingFromFile)
	}
	defer stmt.Close()

	e = new(ExchangeRate)
	if err = stmt.QueryRow(cf, ct).Scan(&e.CurrencyFrom, &e.CurrencyTo, &e.ExchangeRate); err != nil {
		return nil, errors.New(errExchangeRateNone)
	}

	return e, nil
	//TODO: add test

}

// ExchangeRateList returns all currency exchange rates as closure
func ExchangeRateList(db *gsqlitehandler.SqliteDB) (f func() *ExchangeRate, err error) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	if stmt, err = db.Handler.Prepare("SELECT currency_from, currency_to, exchange_rate FROM currencies ORDER BY currency_from, currency_to;"); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	if rows, err = stmt.Query(); err != nil {
		return nil, errors.New(errReadingFromFile)
	}

	f = func() *ExchangeRate {
		if rows.Next() {
			c := new(ExchangeRate)
			rows.Scan(&c.CurrencyFrom, &c.CurrencyTo, &c.ExchangeRate)
			return c
		}
		rows.Close()
		stmt.Close()

		return nil
	}

	return f, nil
	//TODO: add test
}

// ExchangeRateRemove removes given currency exchange rate
func ExchangeRateRemove(db *gsqlitehandler.SqliteDB, e *ExchangeRate) error {
	var err error
	var stmt *sql.Stmt

	if stmt, err = db.Handler.Prepare("DELETE FROM currencies WHERE currency_from=? AND currency_to=?;"); err != nil {
		return errors.New(errWritingToFile)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(e.CurrencyFrom, e.CurrencyTo); err != nil {
		return errors.New(errWritingToFile)
	}

	return nil
	//TODO: add test
}
