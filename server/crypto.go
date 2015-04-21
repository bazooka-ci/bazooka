package main

import (
	"encoding/hex"

	log "github.com/Sirupsen/logrus"
	"github.com/bazooka-ci/bazooka/commons"
)

func (c *context) encryptData(params map[string]string, body bodyFunc) (*response, error) {
	var v bazooka.StringValue

	body(&v)

	if len(v.Value) == 0 {
		return badRequest("value is mandatory")
	}

	_, err := c.Connector.GetProjectById(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("project not found")
	}

	keys, err := c.Connector.GetCryptoKeys(params["id"])
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return notFound("Crypto Key not found")
	}

	keyContent := keys[0].Content

	encrypted, err := bazooka.Encrypt(keyContent, []byte(v.Value))
	if err != nil {
		log.Fatal(err)
	}

	encryptedData := &bazooka.StringValue{
		Value: hex.EncodeToString(encrypted),
	}

	return ok(&encryptedData)
}
