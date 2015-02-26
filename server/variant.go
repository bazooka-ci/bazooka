package main

import "github.com/bazooka-ci/bazooka/commons/mongo"

func (c *context) getVariant(params map[string]string, body bodyFunc) (*response, error) {
	variant, err := c.Connector.GetVariantByID(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("variant not found")
	}

	return ok(&variant)
}

func (c *context) getVariants(params map[string]string, body bodyFunc) (*response, error) {
	variants, err := c.Connector.GetVariants(params["id"])
	if err != nil {
		return nil, err
	}

	return ok(&variants)
}

func (c *context) getVariantLog(params map[string]string, body bodyFunc) (*response, error) {
	log, err := c.Connector.GetLog(&mongo.LogExample{
		VariantID: params["id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("log not found")
	}

	return ok(&log)
}
