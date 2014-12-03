package dockercommand

import "testing"

func TestContainerRemove(t *testing.T) {
	docker, err := NewDocker("")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	container, err := docker.Run(&RunOptions{
		Image: "ubuntu",
		Cmd:   []string{"ls", "/"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = container.Remove(&RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
