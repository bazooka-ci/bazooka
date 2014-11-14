package bazooka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFilesWithPrefix(t *testing.T) {
	files, err := ListFilesWithPrefix("fixtures/commons/listfileswithprefix/0", "test")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Len(t, files, 3)
	assert.Contains(t, files, "fixtures/commons/listfileswithprefix/0/test0")
	assert.Contains(t, files, "fixtures/commons/listfileswithprefix/0/test12")
	assert.Contains(t, files, "fixtures/commons/listfileswithprefix/0/test1234test")

	files, err = ListFilesWithPrefix("fixtures/commons/listfileswithprefix/0", "abc")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Len(t, files, 1)
	assert.Contains(t, files, "fixtures/commons/listfileswithprefix/0/abc")
}
