package dockercommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertEnvMapToSlice(t *testing.T) {
	envMap := map[string]string{
		"TOTO1": "toto-1",
		"TOTO2": "toto-2",
	}
	slice := convertEnvMapToSlice(envMap)
	assert.Contains(t, slice, "TOTO1=toto-1")
	assert.Contains(t, slice, "TOTO2=toto-2")
}
