package matrix

import (
	"fmt"
	"strings"
)

type Matrix map[string][]string

type Iterator func(permutation map[string]string, counter string)

func Iter(mx Matrix, it Iterator, vars ...string) {
	iter(mx, it, map[string]string{}, []string{}, vars...)
}

func IterAll(mx Matrix, it Iterator) {
	keys := make([]string, 0, len(mx))
	for key := range mx {
		keys = append(keys, key)
	}
	iter(mx, it, map[string]string{}, []string{}, keys...)
}

func iter(mx Matrix, it Iterator, permutation map[string]string, counter []string, vars ...string) {
	if len(vars) == 0 {
		//if no more variables, we reached a fixed permutation, call the iterator and return
		it(permutation, strings.Join(counter, ""))
		return
	}

	// handle the first variable
	v := vars[0]
	// pin the counter index of this variable
	ci := len(counter)
	// add a slot for this variable counter
	counter = append(counter, "")
	// prepare this variable iteration index format: calculate the index width based on the number of values
	// e.g. for a variable with 4 values, an index of width 1 is sufficient 0, 1, 2, ...
	// for a variable with 12 values, a width of 2 is needed: 00, 01, ..., 11, ...
	vw := len(mx[v])/10 + 1
	cf := fmt.Sprintf("%%0%dd", vw) // generate the format: %4d for example where 4 is the width

	// iterate over the fixed variable values
	for i, vv := range mx[v] {
		// fix a value
		permutation[v] = vv
		// set this variable counter by formatting the iteration index using the format computed above
		counter[ci] = fmt.Sprintf(cf, i)
		// recursively call _iter with a n-1 vars array (after removing self)
		iter(mx, it, permutation, counter, vars[1:]...)
	}
}
