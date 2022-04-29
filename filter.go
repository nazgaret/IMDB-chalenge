package main

import (
	"strings"
)

const (
	//filter fields indexes
	FilterTitleTypeIndex = iota
	FilterPrimaryTitleIndex
	FilterOriginalTitleIndex
	FilterStartYearIndex
	FilterEndYearIndex
	FilterRuntimeMinutesIndex
	FilterGenresIndex

	//row fields indexes
	RowTitleTypeIndex      = 1
	RowPrimaryTitleIndex   = 2
	RowOriginalTitleIndex  = 3
	RowStartYearIndex      = 5
	RowEndYearIndex        = 6
	RowRuntimeMinutesIndex = 7
	RowGenresIndex         = 8
)

type (
	Filter      [7]string
	FilterFuncs []func(lines []string) bool
)

//newFilterFuncs create list of filter funcs to check row fields for equality
func newFilterFuncs(filter Filter) []func(lines []string) bool {
	funcs := make(FilterFuncs, 0, len(filter))
	for i, v := range filter {
		if v != "" {
			//we need local var for closure because i change value through range
			i := i
			rowIndex := 0
			switch i {
			case FilterTitleTypeIndex:
				rowIndex = RowTitleTypeIndex
			case FilterPrimaryTitleIndex:
				rowIndex = RowPrimaryTitleIndex
			case FilterOriginalTitleIndex:
				rowIndex = RowOriginalTitleIndex
			case FilterStartYearIndex:
				rowIndex = RowStartYearIndex
			case FilterEndYearIndex:
				rowIndex = RowEndYearIndex
			case FilterRuntimeMinutesIndex:
				rowIndex = RowRuntimeMinutesIndex
			case FilterGenresIndex:
				rowIndex = RowGenresIndex
			}

			//special filter for multi genres field
			if i == FilterGenresIndex {
				funcs = append(funcs, func(s []string) bool {
					for _, s2 := range strings.Split(s[rowIndex], ",") {
						if filter[i] == s2 {
							return true
						}
					}
					return false
				})
				continue
			}
			//usual function to check field by existed filter
			funcs = append(funcs, func(s []string) bool {
				return filter[i] == s[rowIndex]
			})
		}
	}
	return funcs
}
