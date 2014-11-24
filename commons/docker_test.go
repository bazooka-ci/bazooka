package bazooka

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeDockerString(t *testing.T) {
	assert.Equal(t, SafeDockerString("dockerfile/mongodb:test"), "dockerfile_mongodb_test")
}
