package matrix

import (
	"reflect"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestAddVar(t *testing.T) {
	mx := Matrix{}
	mx.AddVar("a", "42")
	v, ok := mx["a"]
	require.True(t, ok, "AddVar should have set the key on the matrix")

	require.Equal(t, 1, len(v), "AddVar on an empty matrix should have added exactly one value")

	require.Equal(t, "42", v[0], "AddVar should have added the passed value")

	// Now a second add with the same variable name
	mx.AddVar("a", "666")
	v, ok = mx["a"]
	require.True(t, ok, "A second AddVar should have kept the key in the matrix")

	require.Equal(t, 2, len(v), "A second AddVar should have added a second value")

	require.Equal(t, "42", v[0], "A second AddVar shouldn't have touched the first value")
	require.Equal(t, "666", v[1], "A second AddVar shouldn have added a second value")
}

func TestMerge(t *testing.T) {
	mx1 := Matrix{}
	mx1.AddVar("a", "42")

	mx2 := Matrix{}
	mx2.AddVar("a", "666")
	mx2.AddVar("b", "bzk")

	mx1.Merge(mx2)

	v, ok := mx1["a"]
	require.True(t, ok, "The merged matrix should contain the a variable")

	require.Equal(t, 2, len(v), "The merged matrix should contain the merged values of the original ones")

	require.Equal(t, "42", v[0], "A second AddVar shouldn't have touched the first value")
	require.Equal(t, "666", v[1], "A second AddVar shouldn have added a second value")

	v, ok = mx1["b"]
	require.True(t, ok, "The merged matrix should contain the b variable")

	require.Equal(t, 1, len(v), "The merged matrix should contain the original values of the the b variable")

	require.Equal(t, "bzk", v[0], "The merged matrix should contain the original values of the the b variable")
}

type expectedPerm struct {
	visited     bool
	permutation map[string]string
}
type iterCapturer struct {
	captured map[string]map[string]string // permutations index by the counter
	expected []*expectedPerm
}

func newCapturer(permutations []map[string]string) *iterCapturer {
	res := &iterCapturer{map[string]map[string]string{}, []*expectedPerm{}}
	for _, p := range permutations {
		res.expected = append(res.expected, &expectedPerm{false, p})
	}
	return res
}

func (i *iterCapturer) iter(t *testing.T) Iterator {
	return func(permutation map[string]string, counter string) {
		if _, ok := i.captured[counter]; ok {
			t.Fatalf("Duplicate counter value %s", counter)
		}
		i.captured[counter] = permutation
	}
}

func (i *iterCapturer) check(t *testing.T) {
	require.Equal(t, len(i.expected), len(i.captured), "Different permutations count")
	if len(i.expected) == 0 {
		return
	}

	for _, perm := range i.captured {
		found := false
		for _, exPerm := range i.expected {
			if reflect.DeepEqual(perm, exPerm.permutation) {
				if exPerm.visited {
					t.Fatalf("Duplicate permutation %v", perm)
				}
				exPerm.visited = true
				found = true
			}
			if found {
				break
			}
		}
		if !found {
			t.Fatalf("The permutation %v wasn't in the expected permutations list", perm)
		}
	}

}

func TestIter(t *testing.T) {
	cases := []struct {
		mx           *Matrix
		iterVars     []string
		exclusions   []*Matrix
		permutations []map[string]string
	}{
		{&Matrix{"a": {"42", "666"}}, []string{"a"},
			nil,
			[]map[string]string{{"a": "42"}, {"a": "666"}}},

		{&Matrix{"a": {"42", "666"}, "b": {"bzk"}}, []string{"a"},
			nil,
			[]map[string]string{{"a": "42"}, {"a": "666"}}},

		{&Matrix{"a": {"42", "666"}, "b": {"bzk"}}, []string{"a", "b"},
			nil,
			[]map[string]string{{"a": "42", "b": "bzk"}, {"a": "666", "b": "bzk"}}},

		{&Matrix{"a": {"42", "666"}, "b": {"bzk", "see"}}, []string{"a", "b"},
			nil,
			[]map[string]string{
				{"a": "42", "b": "bzk"}, {"a": "42", "b": "see"},
				{"a": "666", "b": "bzk"}, {"a": "666", "b": "see"}}},

		{&Matrix{"a": {"42", "666"}, "b": {"bzk", "see"}}, []string{"a", "b"},
			[]*Matrix{{"a": {"666"}}},
			[]map[string]string{
				{"a": "42", "b": "bzk"}, {"a": "42", "b": "see"}}},

		{&Matrix{"a": {"42", "666"}, "b": {"bzk", "see"}}, []string{"a", "b"},
			[]*Matrix{{"a": {"666"}}, {"b": {"see"}}},
			[]map[string]string{
				{"a": "42", "b": "bzk"}}},

		{&Matrix{"a": {"42", "666", "1024"}, "b": {"bzk", "see"}}, []string{"a", "b"},
			[]*Matrix{{"a": {"666", "1024"}}},
			[]map[string]string{
				{"a": "42", "b": "bzk"}, {"a": "42", "b": "see"}}},
	}

	for _, cas := range cases {
		t.Logf("mx: %v", *cas.mx)
		cap := newCapturer(cas.permutations)
		cas.mx.Iter(cap.iter(t), cas.exclusions, cas.iterVars...)
		cap.check(t)

		if len(cas.iterVars) == len(*cas.mx) {
			cap := newCapturer(cas.permutations)
			cas.mx.IterAll(cap.iter(t), cas.exclusions)
			cap.check(t)
		}
	}

}
