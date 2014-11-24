package bazooka

import (
	"fmt"
	"io/ioutil"
	"os"
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
