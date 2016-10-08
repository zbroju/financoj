// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

// ItemStatus indicates the life cycle of an object
type ItemStatus int

const (
	ISUnset ItemStatus = -1
	ISClose ItemStatus = 0
	ISOpen  ItemStatus = 1
)

// String satisfies fmt.Stringer interface in order to get human readable names
func (is ItemStatus) String() string {
	var s string

	switch is {
	case ISUnset:
		s = "unset"
	case ISOpen:
		s = "active"
	case ISClose:
		s = "inactive"
	}

	return s
}

// Category represents the basic object for category
type CategoryT struct {
	Id           int
	MainCategory *MainCategoryT
	Name         string
	Status       ItemStatus
}

// MainCategoryStatusT describes the behaviour of categories and its descendants (transactions)
type MainCategoryTypeT int

const (
	MCTUnknown  MainCategoryTypeT = -1
	MCTUnset    MainCategoryTypeT = 0
	MCTCost     MainCategoryTypeT = 1
	MCTTransfer MainCategoryTypeT = 2
	MCTIncome   MainCategoryTypeT = 3
)

// String satisfies fmt.Stringer interface in order to get human readable names
func (mct MainCategoryTypeT) String() string {
	var mctName string

	switch mct {
	case MCTUnknown:
		mctName = "unknown"
	case MCTUnset:
		mctName = "not set"
	case MCTIncome:
		mctName = "income"
	case MCTCost:
		mctName = "cost"
	case MCTTransfer:
		mctName = "transfer"
	}

	return mctName
}

// MainCategory represents the basic object for main category
type MainCategoryT struct {
	Id     int
	MType  MainCategoryTypeT
	Name   string
	Status ItemStatus
}

// MainCategory list represents the list of main categories
type MainCategoryListT struct {
	MainCategories map[*int]MainCategoryT
}

// MainCategoryListNew returns initialized MainCategoryList
func MainCategoryListNew() *MainCategoryListT {
	mcList := new(MainCategoryListT)
	mcList.MainCategories = make(map[*int]MainCategoryT)
	return mcList
}