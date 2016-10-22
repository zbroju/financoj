// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	. "github.com/zbroju/financoj/lib"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

// GetLoggers returns two loggers for standard formatting of messages and errors
func getLoggers() (messageLogger *log.Logger, errorLogger *log.Logger) {
	messageLogger = log.New(os.Stdout, fmt.Sprintf("%s: ", AppName), 0)
	errorLogger = log.New(os.Stderr, fmt.Sprintf("%s: ", AppName), 0)

	return
}

// mainCategoryTypeForString returns main category type for given string
func mainCategoryTypeForString(m string) (mct MainCategoryType) {
	switch m {
	case "c", "cost", NotSetStringValue: // If null string is given, then the default value is MCT_Cost
		mct = MCTCost
	case "i", "income":
		mct = MCTIncome
	case "t", "transfer":
		mct = MCTTransfer
	default:
		mct = MCTUnknown
	}

	return mct
}

// accountTypeForString returns account type for given string
func accountTypeForString(s string) (t AccountType) {
	switch s {
	case "o", "operational", NotSetStringValue: // If null string is given then the default value is ATTransactional
		t = ATTransactional
	case "s", "savings":
		t = ATSaving
	case "p", "properties":
		t = ATProperty
	case "i", "investment":
		t = ATInvestment
	case "l", "loans":
		t = ATLoan
	default:
		t = ATUnknown
	}

	return t
}

// getLineFor returns pre-formatted line formatting string for reporting
func lineFor(fs ...string) string {
	line := strings.Join(fs, FSSeparator) + "\n"
	return line
}

// getHFSForText return heading formatting string for string values
func hFSForText(l int) string {
	return fmt.Sprintf("%%-%ds", l)
}

// getHFSForNumeric return heading formatting string for numeric values
func hFSForNumeric(l int) string {
	return fmt.Sprintf("%%%ds", l)
}

// getDFSForText return data formatting string for string
func dFSForText(l int) string {
	return fmt.Sprintf("%%-%ds", l)
}

// getDFSForRates return data formatting string for rates
func getDFSForRates(l int) string {
	return fmt.Sprintf("%%%d.4f", l)
}

// getDFSForValue return data formatting string for values
func getDFSForValue(l int) string {
	return fmt.Sprintf("%%%d.2f", l)
}

// getDFSForID return data formatting string for id
func dFSForID(l int) string {
	return fmt.Sprintf("%%%dd", l)
}

// Return the bigger number out of the two given
func maxLen(s string, i int) int {
	if l := utf8.RuneCountInString(s); l > i {
		return l
	} else {
		return i
	}
}
