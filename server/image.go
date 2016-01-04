package main

import "strings"

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
	b := struct {
		Image string `json:"image" validate:"required"`
	}{}
	r.parseBody(&b)
	b.Image = strings.TrimSpace(b.Image)
	if len(b.Image) == 0 {
		return badRequest("image is required")
	}
	err := c.connector.SetImage(r.vars["name"], b.Image)
	if err != nil {
		return nil, err
	}

	return ok(b)
}
