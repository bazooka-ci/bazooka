package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	assert.Equal(t, "bazooka", safeDockerAlias("bazooka"))
	assert.Equal(t, "bazooka_test", safeDockerAlias("bazooka/test"))
	assert.Equal(t, "bazooka_test", safeDockerAlias("bazooka.test"))
	assert.Equal(t, "bazooka_test", safeDockerAlias("bazooka:test"))
	assert.Equal(t, "baz_o_oka_te_st_bzk", safeDockerAlias("baz.o-oka/te:st-bzk"))
}
