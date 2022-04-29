package main

import (
	"testing"
)

func Test_newFilterFuncs(t *testing.T) {
	testRow := []string{"tconst", "titleType", "primaryTitle", "originalTitle", "isAdult", "startYear", "endYear", "runtimeMinutes", "genres1,genres2"}
	type filterFuncTest struct {
		name    string
		filter  Filter
		len     int
		success bool
	}
	tests := []filterFuncTest{
		{
			name:    "Empty Filter",
			filter:  Filter{},
			len:     0,
			success: true,
		}, {
			name:    "FilterGenresDone",
			filter:  Filter{"", "", "", "", "", "", "genres1"},
			len:     1,
			success: true,
		}, {
			name:    "FilterGenresWrong",
			filter:  Filter{"", "", "", "", "", "", "wrong"},
			len:     1,
			success: false,
		}, {
			name:    "Full Filter Done",
			filter:  Filter{"titleType", "primaryTitle", "originalTitle", "startYear", "endYear", "runtimeMinutes", "genres1"},
			len:     7,
			success: true,
		}, {
			name:    "Full Filter One Wrong",
			filter:  Filter{"titleType", "primaryTitle", "wrong", "startYear", "endYear", "runtimeMinutes", "genres1"},
			len:     7,
			success: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			funcs := newFilterFuncs(tt.filter)
			//check len of filter funcs
			if len(funcs) != tt.len {
				t.Errorf("funcs len = %v, must %v", len(funcs), tt.len)
			}
			//if no filter all rows accepted
			result := true
			for _, f := range funcs {
				//if any filter fail return false
				if !f(testRow) {
					result = false
					break
				}
			}
			//check if result expected
			if result != tt.success {
				t.Errorf("filtering result = %v, must %v", result, tt.success)
			}
		})
	}
}
