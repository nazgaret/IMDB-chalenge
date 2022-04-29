package main

import (
	"fmt"
	"strings"
)

// NewJob create job worked on butch of rows with provided filter funcs
func NewJob(filterFuncs FilterFuncs) func(strings []string) error {
	return func(butch []string) error {
		for _, s := range butch {
			for _, f := range filterFuncs {
				sliceString := strings.Split(s, "\t")
				if f(sliceString) {
					fmt.Println(s)
					//todo Do Job
					//too return err
				}
			}
		}
		return nil
	}
}
