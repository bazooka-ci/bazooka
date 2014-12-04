package main

func (c *context) getImage(params map[string]string, body bodyFunc) (*response, error) {
	image, err := c.Connector.GetImage(params["name"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("image not found")
	}

	return ok(&image)
}

func (c *context) getImages(params map[string]string, body bodyFunc) (*response, error) {
	variants, err := c.Connector.GetImages()
	if err != nil {
		return nil, err
	}

	return ok(&variants)
}

func (c *context) setImage(params map[string]string, body bodyFunc) (*response, error) {
	b := map[string]string{}
	body(&b)
	image, ex := b["image"]
	if !ex {
		return badRequest("image is required")
	}
	err := c.Connector.SetImage(params["name"], image)
	if err != nil {
		return nil, err
	}

	return ok(b)
}
