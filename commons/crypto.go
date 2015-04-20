package bazooka

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

// Encrypt encrypts some data with the key
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt decrypts some data with the key
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ReadCryptoKey(filePath string) ([]byte, error) {
	exists, err := FileExists(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to check existence of file: %s, reason is: %v\n", filePath, err)
	}

	if !exists {
		return nil, fmt.Errorf("Your bazooka config contains secured data but the keyfile can not be found at %s, reason is: %v\n", filePath, err)
	}

	key, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error while reading crypto key file: %s, reason is: %v\n", filePath, err)
	}
	return key, nil
}

func LoadCryptoKeyFromFile(filePath string) error {
	var err error
	PrivateKey, err = ReadCryptoKey(filePath)
	return err
}
