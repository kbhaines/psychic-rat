package ctr

import (
	"psychic-rat/mdl/user"
	"fmt"
)

type UserController interface {
	AddUser(email string, country string, firstname string) (user.Id, error)
	GetUser(id user.Id) (user.Record, error)
	DeleteUser(id user.Id) error
}

type userController struct {}

var _ UserController = &userController{}

func (p *userController) AddUser(email string, country string, firstname string) (id user.Id, err error) {
	user := user.New(email, country, firstname)
	id, err = userRepo.Create(user)
	if err != nil {
		return id, fmt.Errorf("failed to create %v: %v", user, err)
	}
	return id, nil
}

func (p *userController) GetUser(id user.Id) (user.Record, error) {
	panic("implement me")
}

func (p *userController) DeleteUser(id user.Id) error {
	panic("implement me")
}


