package parallel

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	type taskCase struct {
		tag      interface{}
		err      error
		duration time.Duration
	}
	cases := [][]taskCase{
		{{"1", nil, 100 * time.Millisecond}, {"2", nil, 0 * time.Millisecond}},
		{{"1", nil, 100 * time.Millisecond}, {"2", nil, 200 * time.Millisecond}},
		{{"1", nil, 100 * time.Millisecond}, {"2", fmt.Errorf("lol wut"), 200 * time.Millisecond}},
		{{"1", fmt.Errorf("oh snap"), 0 * time.Millisecond}, {"2", fmt.Errorf("lol wut"), 0 * time.Millisecond}},
	}

	for _, cas := range cases {
		par := New()
		called := map[interface{}]bool{}
		expectedErrs := map[interface{}]error{}
		for _, itc := range cas {
			tc := itc
			expectedErrs[tc.tag] = tc.err
			par.Submit(func() error {
				time.Sleep(tc.duration)
				called[tc.tag] = true
				return tc.err
			}, tc.tag)

		}
		par.Exec(func(tag interface{}, err error) {
			expectedErr, found := expectedErrs[tag]
			if !found {
				t.Fatalf("callback called with an unknown tag %#v", tag)
			}
			require.Equal(t, expectedErr, err, "The task associated with the tag %#v didn't return the expected error", tag)
		})

		for _, tc := range cas {
			c, found := called[tc.tag]
			require.True(t, found, "The task associated with the tag %#v wasn't called", tc.tag)
			require.True(t, c, "The task associated with the tag %#v wasn't called", tc.tag)
		}
	}

}
