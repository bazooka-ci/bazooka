package dockercommand

import "testing"

func TestDockerBuild(t *testing.T) {
	docker, err := NewDocker("")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = docker.Build(&BuildOptions{
		Dockerfile: "Dockerfile_test",
		Path:       ".",
		Tag:        "testgodockercommand",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

}
