package bazooka

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveConfigFile(t *testing.T) {
	res, err := ResolveConfigFile("fixtures/config/resolveconfigfile/0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, res, "fixtures/config/resolveconfigfile/0/.bazooka.yml")

	res, err = ResolveConfigFile("fixtures/config/resolveconfigfile/1")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, res, "fixtures/config/resolveconfigfile/1/.travis.yml")

	res, err = ResolveConfigFile("fixtures/config/resolveconfigfile/2")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, res, "fixtures/config/resolveconfigfile/2/.bazooka.yml")

	res, err = ResolveConfigFile("fixtures/config/resolveconfigfile/3")
	if err == nil {
		t.Fatalf("Error should have been raised %s", err)
	}
	assert.Equal(t, err.Error(), "Unable to find either .bazooka.yml or .travis.yml at the root of the project")
}

func TestParser(t *testing.T) {
	parsetest := &parse{}
	err := Parse("fixtures/config/parse/test.yml", parsetest)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, parsetest.Type1, "TestABC")
	assert.Len(t, parsetest.Type2, 2)
	assert.Contains(t, parsetest.Type2, "testDEF1")
	assert.Contains(t, parsetest.Type2, "testDEF2")
}

func TestFlush(t *testing.T) {
	parsetest := &parse{
		Type1: "TestABC",
		Type2: []string{"testDEF1", "testDEF2"},
	}
	err := Flush(parsetest, "fixtures/config/parse/test-out.yml")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	breal, err := ioutil.ReadFile("fixtures/config/parse/test-out.yml")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	bexpected, err := ioutil.ReadFile("fixtures/config/parse/test.yml")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, breal, bexpected)
}

type parse struct {
	Type1 string   "abc"
	Type2 []string "def"
}
