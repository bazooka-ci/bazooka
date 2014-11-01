package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	docker "github.com/bywan/go-dockercommand"
)

type Builder struct {
	Options *BuildOptions
}

type BuildOptions struct {
	DockerfileFolder string
	SourceFolder     string
	JobID            string
	VariantID        string
}

func (b *Builder) Build() error {

	log.Printf("Starting building Dockerfiles\n")
	files, err := listBuildfiles(b.Options.DockerfileFolder)
	if err != nil {
		return err
	}

	client, err := docker.NewDocker(DockerEndpoint)
	if err != nil {
		return err
	}

	errChan := make(chan error)
	successChan := make(chan string)
	remainingBuilds := len(files)

	for i, file := range files {
		go buildContainer(client, i, b, file, successChan, errChan)
	}

	var buildImages []string
	for {
		select {
		case tag := <-successChan:
			buildImages = append(buildImages, tag)
			remainingBuilds--
		case err := <-errChan:
			return err
		}

		if remainingBuilds == 0 {
			break
		}
	}

	errChanRun := make(chan error)
	successChanRun := make(chan bool)
	remainingRuns := len(buildImages)
	for _, buildImage := range buildImages {
		go runContainer(client, buildImage, successChanRun, errChanRun)
	}

	for {
		select {
		case _ = <-successChanRun:
			remainingRuns--
		case err := <-errChanRun:
			return err
		}

		if remainingRuns == 0 {
			break
		}
	}

	log.Printf("Dockerfiles builds finished\n")
	return nil
}

func runContainer(client *docker.Docker, buildImage string, successChan chan bool, errChan chan error) {
	containerID, err := client.Run(&docker.RunOptions{
		Image: buildImage,
	})
	if err != nil {
		errChan <- err
		return
	}
	details, err := client.Inspect(containerID)
	if err != nil {
		errChan <- err
		return
	}
	if details.State.ExitCode != 0 {
		errChan <- fmt.Errorf("Build failed\n Check Docker container logs, id is %s\n", containerID)
		return
	}
	successChan <- true
}

func buildContainer(client *docker.Docker, i int, b *Builder, file *buildFiles, successChan chan string, errChan chan error) {
	for _, buildFile := range file.BuildFiles {
		splitString := strings.Split(buildFile, "/")
		CopyFile(buildFile, fmt.Sprintf("%s/%s", b.Options.SourceFolder, splitString[len(splitString)-1]))
	}

	tag := fmt.Sprintf("bazooka/build-%s-%s-%d", b.Options.JobID, b.Options.VariantID, i)
	err := client.Build(&docker.BuildOptions{
		Tag:        tag,
		Dockerfile: file.Dockerfile,
		Path:       b.Options.SourceFolder,
	})
	if err != nil {
		errChan <- err
	} else {
		successChan <- tag
	}
}

func listBuildfiles(source string) ([]*buildFiles, error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return nil, err
	}
	var output []*buildFiles
	for _, file := range files {
		if file.Mode().IsDir() {
			index, err := strconv.ParseInt(file.Name(), 10, 64)
			if err != nil {
				return nil, err
			}
			filesBuild, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", source, file.Name()))
			if err != nil {
				return nil, err
			}
			var result []string
			for _, fileBuild := range filesBuild {
				if fileBuild.Name() != "Dockerfile" {
					result = append(result, fmt.Sprintf("%s/%s/%s", source, file.Name(), fileBuild.Name()))
				}
			}
			output = append(output, &buildFiles{
				Dockerfile: fmt.Sprintf("%s/%s/Dockerfile", source, file.Name()),
				BuildFiles: result,
				JobIndex:   index,
			})

		}
	}
	return output, nil
}

type buildFiles struct {
	Dockerfile string
	BuildFiles []string
	JobIndex   int64
}

// TODO Extract this or replace it
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
