package bazooka

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func FileExists(path string) (bool, error) {
	if _, err := os.Open(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func AppendToFile(path, text string, perm os.FileMode) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	b = append(b, []byte(text)...)
	return ioutil.WriteFile(path, b, perm)
}

func ListFilesWithPrefix(source, prefix string) ([]string, error) {
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return nil, err
	}
	var output []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) {
			output = append(output, fmt.Sprintf("%s/%s", source, file.Name()))
		}
	}
	return output, nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
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

	// This is an ugly hack to fix sporadic problems with boot2docker and virtualbox
	// where the copied file is sometimes empty
	exec.Command("sync").Run()

	return
}
