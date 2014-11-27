package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExcluded(t *testing.T) {
	exclusions := []map[string]interface{}{
		{
			"golang": "1.3.1",
			"env": []string{
				"TESTA=testa",
				"TESTB=testab",
			},
		},
		{
			"golang": "1.2.1",
			"env": []string{
				"TESTA=testa",
			},
		},
	}

	assert.True(t, IsExcluded(map[string]interface{}{
		"golang": "1.2.1",
		"env": []string{
			"TESTA=testa",
			"TESTB=testa",
			"TESTC=testa",
			"TESTB=testa",
		},
	}, exclusions))
	assert.False(t, IsExcluded(map[string]interface{}{
		"golang": "1.2",
		"env": []string{
			"TESTA=testa",
		},
	}, exclusions))
	assert.False(t, IsExcluded(map[string]interface{}{
		"golang": "1.3.1",
		"env": []string{
			"TESTA=testa",
		},
	}, exclusions))

}
