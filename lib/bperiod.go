// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Local errors
const (
	errMonthIncorrect = "month is not correct"
	errYearIncorrect  = "year is not correct"
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
	return fmt.Sprintf("%04d-%02d", p.Year, p.Month)
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

// Parse converts string (expected format: yyyy-mm) to year and month and after verification if they are within
// their ranges assign them to the budgeting period fields.
func BPeriodParse(s string) (b *BPeriod, err error) {

	sarr := strings.SplitN(s, "-", 2)

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
