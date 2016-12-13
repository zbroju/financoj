// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Local errors
const (
	errPeriodIncorrect = "given period is not correct"
	errMonthIncorrect  = "month is not correct"
	errYearIncorrect   = "year is not correct"
)

// Basic type for keeping budget period (year-month).
type BPeriod struct {
	Year  int64
	Month int64
}

// incorrectYear returns error if y (year) is out of range.
func incorrectYear(y int64) error {
	switch {
	case y < 1900, y > 2100:
		return errors.New(errYearIncorrect)
	default:
		return nil
	}
}

// incorrectMonth returns error if m (month) is out of range.
func incorrectMonth(m int64) error {
	switch {
	case m < 1, m > 12:
		return errors.New(errMonthIncorrect)
	default:
		return nil
	}
}

// String satisfies fmt.Stringer interface in order to get human readable names.
func (p *BPeriod) String() string {
	if p.Month == int64(NotSetIntValue) {
		return fmt.Sprintf("%04d", p.Year)
	} else {
		return fmt.Sprintf("%04d%s%02d", p.Year, DateSeparator, p.Month)
	}

}

func (p *BPeriod) GetStrings() (y, m string) {
	y = fmt.Sprintf("%04d", p.Year)
	if p.Month == int64(NotSetIntValue) {
		m = NotSetStringValue
	} else {
		m = fmt.Sprintf("%02d", p.Month)
	}
	return y, m
}

// Set verifies if year y and month m are within their ranges and assigns it to the budgeting period fields.
func (p *BPeriod) Set(y, m int64) error {
	if err := incorrectYear(y); err != nil {
		return err
	}
	if err := incorrectMonth(m); err != nil {
		return err
	}
	p.Year = y
	p.Month = m

	return nil
}

// BPeriodParseYM converts string (expected format: yyyy-mm) to year and month and after verification if they are within
// their ranges assign them to the budget period fields.
func BPeriodParseYM(s string) (b *BPeriod, err error) {
	sarr := strings.SplitN(s, DateSeparator, 2)

	var y, m int64
	if y, err = strconv.ParseInt(sarr[0], 10, 64); err != nil {
		return nil, err
	}
	if m, err = strconv.ParseInt(sarr[1], 10, 64); err != nil {
		return nil, err
	}
	if err = incorrectYear(y); err != nil {
		return nil, err
	}
	if err = incorrectMonth(m); err != nil {
		return nil, err
	}

	b = new(BPeriod)
	b.Year = y
	b.Month = m

	return b, nil
}

// BPeriodParseYOrYM converts string (expected format: yyyy-mm or yyyy) to year and month or to year only and after verification if they are within
// their ranges assign them to the budget period fields.
func BPeriodParseYOrYM(s string) (b *BPeriod, err error) {
	switch utf8.RuneCountInString(s) {
	case 4:
		var y int64
		if y, err = strconv.ParseInt(s, 10, 64); err == nil {
			if err = incorrectYear(y); err == nil {
				b = new(BPeriod)
				b.Year = y
			}
		}
	case 6, 7:
		b, err = BPeriodParseYM(s)
	default:
		err = errors.New(errPeriodIncorrect)
	}

	return b, nil
}

func BPeriodCurrent() (b *BPeriod, err error) {
	d := time.Now()
	y := int64(d.Year())
	m := int64(d.Month())

	if err = incorrectYear(y); err != nil {
		return nil, err
	}
	if err = incorrectMonth(m); err != nil {
		return nil, err
	}

	b = new(BPeriod)
	b.Year = y
	b.Month = m

	return b, nil
}
