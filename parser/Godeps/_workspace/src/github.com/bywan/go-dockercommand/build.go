package dockercommand

import (
	"archive/tar"
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

type BuildOptions struct {
	Dockerfile string
	Path       string
	Tag        string
}

func (dock *Docker) Build(options *BuildOptions) error {
	t := time.Now()

	inputbuf := bytes.NewBuffer(nil)

	bytearray, err := ioutil.ReadFile(options.Dockerfile)
	if err != nil {
		return err
	}

	tw := tar.NewWriter(inputbuf)
	if err = tw.WriteHeader(&tar.Header{Name: "/Dockerfile", Size: int64(len(bytearray)), ModTime: t, AccessTime: t, ChangeTime: t}); err != nil {
		return err
	}
	if _, err = tw.Write(bytearray); err != nil {
		return err
	}

	buildContextPath(options.Path, options.Path, tw)

	if err = tw.Close(); err != nil {
		return err
	}

	logsReader, outputbuf := io.Pipe()
	go func(reader io.Reader) {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			log.Printf("%s \n", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Println("There was an error with the scanner in attached container", err)
		}
	}(logsReader)

	opts := docker.BuildImageOptions{
		Name:         options.Tag,
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}
	err = dock.client.BuildImage(opts)

	if err != nil {
		return err
	}
	return nil
}

func buildContextPath(sourcePath, dirPath string, tw *tar.Writer) error {
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()
	fis, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		curPath := dirPath + "/" + fi.Name()
		if fi.IsDir() {
			buildContextPath(sourcePath, curPath, tw)
		} else {
			filePath := strings.Replace(curPath, sourcePath, "", 1)
			tarGzWrite(filePath, curPath, tw, fi)
		}
	}
	return nil
}

func tarGzWrite(tarPath string, _path string, tw *tar.Writer, fi os.FileInfo) error {
	h := new(tar.Header)
	h.Name = tarPath
	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	bytearray, err := ioutil.ReadFile(_path)
	if err != nil {
		return err
	}

	if err = tw.WriteHeader(h); err != nil {
		return err
	}

	if _, err = tw.Write(bytearray); err != nil {
		return err
	}

	return nil
}
