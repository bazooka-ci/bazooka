package main

import lib "github.com/bazooka-ci/bazooka/commons"

func (p *context) createUser(params map[string]string, body bodyFunc) (*response, error) {
	var user lib.User

	body(&user)

	switch {
	case len(user.Email) == 0:
		return badRequest("email is mandatory")
	case len(user.Password) == 0:
		return badRequest("password is mandatory")
	}

	exists, err := p.Connector.HasUser(user.Email)
	switch {
	case err != nil:
		return nil, err
	case exists:
		return conflict("email is already known")
	}

	if err = p.Connector.AddUser(&user); err != nil {
		return nil, err
	}
	return created(&user, "/user/"+user.ID)
}

func (p *context) getUser(params map[string]string, body bodyFunc) (*response, error) {
	user, err := p.Connector.GetUserByEmail(params["email"])
	if err != nil {
		if err.Error() != "not found" {
			return nil, err
		}
		return notFound("user not found")
	}

	return ok(&user)
}

func (p *context) getUsers(params map[string]string, body bodyFunc) (*response, error) {
	users, err := p.Connector.GetUsers()
	if err != nil {
		return nil, err
	}

	return ok(&users)
}
