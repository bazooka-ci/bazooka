package dockercommand

import "testing"

func TestDockerRun(t *testing.T) {
	docker, err := NewDocker("")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_, err = docker.Run(&RunOptions{
		Image: "ubuntu",
		Cmd:   []string{"/bin/sh", "-c", "while true; do echo hello world; sleep 1; done"},
		Env: map[string]string{
			"TOTO1": "toto1",
			"TOTO2": "toto2",
		},
		VolumeBinds: []string{
			"/:/volumes",
		},
		Detach: true,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
