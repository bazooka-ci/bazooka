package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectBuildTool(t *testing.T) {
	tool, err := detectBuildTool("fixtures/buildtool/ant")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, tool, "ant")

	tool, err = detectBuildTool("fixtures/buildtool/gradle")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, tool, "gradle")

	tool, err = detectBuildTool("fixtures/buildtool/gradlew")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, tool, "gradlew")

	tool, err = detectBuildTool("fixtures/buildtool/maven")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.Equal(t, tool, "maven")
}
