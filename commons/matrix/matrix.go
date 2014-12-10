package matrix

import (
	"fmt"
	"strings"
)

type Matrix map[string][]string

type Iterator func(permutation map[string]string, counter string)

func (mx *Matrix) AddVar(name, value string) {
	if vs, ok := (*mx)[name]; ok {
		(*mx)[name] = append(vs, value)
		return
	}
	(*mx)[name] = []string{value}
}

func (mx *Matrix) Merge(with map[string][]string) {
	for k, vs := range with {
		for _, v := range vs {
			mx.AddVar(k, v)
		}
	}
}
func (mx *Matrix) Iter(it Iterator, exclusions []*Matrix, vars ...string) {
	mx.iter(it, map[string]string{}, exclusions, []string{}, vars...)
}

func (mx *Matrix) IterAll(it Iterator, exclusions []*Matrix) {
	keys := make([]string, 0, len(*mx))
	for key := range *mx {
		keys = append(keys, key)
	}
	mx.iter(it, map[string]string{}, exclusions, []string{}, keys...)
}

func isIn(needle *Matrix, haystack map[string]string) bool {
	ammo := len(*needle)
	for k, vs := range *needle {
		if w, ok := haystack[k]; ok {
			// the key also exists in the haystack
			found := false
			for _, v := range vs {
				if v == w {
					found = true
					ammo--
					break
				}
			}

			if !found {
				return false
			}
		} else {
			return false
		}
	}

	return ammo == 0
}

func (mx Matrix) iter(it Iterator, permutation map[string]string, exclusions []*Matrix, counter []string, vars ...string) {
	if len(vars) == 0 {
		// if no more variables, we reached a fixed permutation
		//check if this permutation is in the exclusions list
		for _, ex := range exclusions {
			if isIn(ex, permutation) {
				return
			}
		}

		// call the iterator and return
		it(copyMap(permutation), strings.Join(counter, ""))
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
		mx.iter(it, permutation, exclusions, counter, vars[1:]...)
	}
}

func copyMap(m map[string]string) map[string]string {
	res := map[string]string{}
	for k, v := range m {
		res[k] = v
	}
	return res
}
