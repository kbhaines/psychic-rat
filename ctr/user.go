package ctr

import (
	"psychic-rat/mdl/user"
	"fmt"
	"psychic-rat/repo"
)

type UserController interface {
	AddUser(email string, country string, firstname string) (user.Id, error)
	GetUser(id user.Id) (user.Record, error)
	DeleteUser(id user.Id) error
}

type userController struct {}

var _ UserController = &userController{}

var userRepo repo.Users

func (p *userController) AddUser(email string, country string, firstname string) (id user.Id, err error) {
	user := user.New(email, country, firstname)
	err = userRepo.Create(user)
	if err != nil {
		return id, fmt.Errorf("failed to create %v: %v", user, err)
	}
	return id, nil
}

func (p *userController) GetUser(id user.Id) (user user.Record, err error) {
	user, err = userRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("user %v not available: %v", id, err)
	}
	return user, nil
}

func (p *userController) DeleteUser(id user.Id) error {
	panic("implement me")
}


