package main

import "github.com/haklop/bazooka/commons/mongo"

func (c *Context) getVariant(params map[string]string, body BodyFunc) (*Response, error) {
	variant, err := c.Connector.GetVariantByID(params["id"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return NotFound("variant not found")
	}

	return Ok(&variant)
}

func (c *Context) getVariants(params map[string]string, body BodyFunc) (*Response, error) {
	variants, err := c.Connector.GetVariants(params["id"])
	if err != nil {
		return nil, err
	}

	return Ok(&variants)
}

func (c *Context) getVariantLog(params map[string]string, body BodyFunc) (*Response, error) {
	log, err := c.Connector.GetLog(&mongo.LogExample{
		VariantID: params["id"],
	})
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return NotFound("log not found")
	}

	return Ok(&log)
}
