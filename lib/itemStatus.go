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
