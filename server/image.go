package main

func (c *context) getImage(r *request) (*response, error) {
	image, err := c.Connector.GetImage(r.vars["name"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("image not found")
	}

	return ok(&image)
}

func (c *context) getImages(r *request) (*response, error) {
	variants, err := c.Connector.GetImages()
	if err != nil {
		return nil, err
	}

	return ok(&variants)
}

func (c *context) setImage(r *request) (*response, error) {
	b := map[string]string{}
	r.parseBody(&b)
	image, ex := b["image"]
	if !ex {
		return badRequest("image is required")
	}
	err := c.Connector.SetImage(r.vars["name"], image)
	if err != nil {
		return nil, err
	}

	return ok(b)
}
