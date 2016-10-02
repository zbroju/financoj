// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package financoj

// ItemStatus indicates the life cycle of an object
type ItemStatus int

const (
	ISClose ItemStatus = 0
	ISOpen  ItemStatus = 1
)

// MainCategoryStatusT describes the behaviour of categories and its descendants (transactions)
type MainCategoryTypeT int

const (
	MCTUnknown  MainCategoryTypeT = -2
	MCTCost     MainCategoryTypeT = -1
	MCTTransfer MainCategoryTypeT = 0
	MCTIncome   MainCategoryTypeT = 1
)

// String satisfies fmt.Stringer interface in order to get human readable names
func (mct MainCategoryTypeT) String() string {
	var mctName string

	switch mct {
	case MCTIncome:
		mctName = "income"
	case MCTCost:
		mctName = "cost"
	case MCTTransfer:
		mctName = "transfer"
	default:
		mctName = "not set"
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
