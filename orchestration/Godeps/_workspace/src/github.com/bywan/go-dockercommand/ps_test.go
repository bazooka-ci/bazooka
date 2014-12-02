package dockercommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerPs(t *testing.T) {
	docker, err := NewDocker("")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_, err = docker.Run(&RunOptions{
		Image: "ubuntu",
		Cmd:   []string{"ls", "/"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	containers, err := docker.Ps(&PsOptions{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	assert.NotEmpty(t, containers)

}
