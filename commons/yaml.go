package bazooka

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Parse(file string, object interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// TODO Add validation
	return yaml.Unmarshal(b, object)
}

func Flush(object interface{}, outputFile string) error {
	d, err := yaml.Marshal(object)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFile, d, 0644)
}

func unmarshalOneOrMany(unmarshal func(interface{}) error, name string) ([]string, error) {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return nil, err
	}

	switch convCmd := raw.(type) {
	case string:
		return []string{convCmd}, nil
	case []interface{}:
		res := make([]string, len(convCmd))
		for i, rawCmd := range convCmd {
			cmd, ok := rawCmd.(string)
			if !ok {
				return nil, fmt.Errorf("%s list can only contain strings", name)
			}
			res[i] = cmd
		}
		return res, nil
	default:
		return nil, fmt.Errorf("%s can be either a string or a list of strings", name)
	}
}
