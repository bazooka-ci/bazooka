package bazooka

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type BzkString struct {
	Name    string
	Value   string
	Secured bool
}

func FlattenEnvMap(mapp map[string][]BzkString) []BzkString {
	res := []BzkString{}
	for _, value := range mapp {
		res = append(res, value...)
	}
	return res
}

func GetEnvMap(envArray []BzkString) map[string][]BzkString {
	envKeyMap := make(map[string][]BzkString)
	for _, env := range envArray {
		envKeyMap[env.Name] = append(envKeyMap[env.Name], env)
	}
	return envKeyMap
}

func SplitNameValue(s string) (string, string) {
	split := strings.SplitN(string(s), "=", 2)
	value := ""
	if len(split) == 2 {
		value = split[1]
	}
	return split[0], value
}

func (c *BzkString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	bzkStr, err := extractBzkString(raw)
	if err != nil {
		return err
	}
	*c = bzkStr
	return nil
}

func (c BzkString) MarshalYAML() (interface{}, error) {
	merged := c.Name
	if len(c.Value) > 0 {
		merged = fmt.Sprintf("%s=%s", merged, c.Value)
	}
	if !c.Secured {
		return merged, nil
	}

	encrypted, err := encryptBzkString(merged)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"secure": encrypted,
	}, nil

}

func extractBzkString(raw interface{}) (BzkString, error) {
	switch convCmd := raw.(type) {
	case string:
		n, v := SplitNameValue(convCmd)
		return BzkString{n, v, false}, nil
	case map[interface{}]interface{}:
		if len(convCmd) > 1 {
			return BzkString{}, fmt.Errorf("BzkString should either be a go string or 'secure: <string>'")
		}
		if _, ok := convCmd["secure"]; !ok {
			return BzkString{}, fmt.Errorf("BzkString should either be a go string or 'secure: <string>'")
		}

		decrypted, err := decryptBzkString(convCmd["secure"].(string))
		if err != nil {
			return BzkString{}, fmt.Errorf("Error while trying to decrypt data, reason is: %v\n", err)
		}
		n, v := SplitNameValue(string(decrypted))
		return BzkString{n, v, true}, nil
	default:
		return BzkString{}, fmt.Errorf("BzkString should either be a go string or 'secure: <string>'")
	}
}

func decryptBzkString(str string) ([]byte, error) {
	if PrivateKey == nil {
		return nil, fmt.Errorf("PrivateKey is not set")
	}

	toDecryptDataAsHex, err := hex.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("Unable to decode string as hexa data, reason is: %v\n", err)
	}

	decrypted, err := Decrypt(PrivateKey, toDecryptDataAsHex)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to decrypt data, reason is: %v\n", err)
	}
	return decrypted, nil
}

func encryptBzkString(str string) (string, error) {
	if PrivateKey == nil {
		return "", fmt.Errorf("PrivateKey is not set")
	}

	encrypted, err := Encrypt(PrivateKey, []byte(str))
	if err != nil {
		return "", fmt.Errorf("Error while trying to encrypt data, reason is: %v\n", err)
	}

	return hex.EncodeToString(encrypted), nil
}
