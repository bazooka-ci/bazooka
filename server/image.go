package main

func (c *context) getImage(r *request) (*response, error) {
	image, err := c.connector.GetImage(r.vars["name"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("image not found")
	}

	return ok(&image)
}

func (c *context) getImages(r *request) (*response, error) {
	images, err := c.connector.GetImages()
	if err != nil {
		return nil, err
	}

	return ok(&images)
}

func (c *context) setImage(r *request) (*response, error) {
	b := map[string]string{}
	r.parseBody(&b)
	image, ex := b["image"]
	if !ex {
		return badRequest("image is required")
	}
	err := c.connector.SetImage(r.vars["name"], image)
	if err != nil {
		return nil, err
	}

	return ok(b)
}
