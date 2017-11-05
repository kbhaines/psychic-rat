package api

import (
	"fmt"
	"log"
	"psychic-rat/mdl"
)

func NewUserApi(repos Repos) *userApiRepoImpl {
	return &userApiRepoImpl{repos: repos}
}

type userApiRepoImpl struct {
	repos Repos
}

func (u *userApiRepoImpl) GetById(userId mdl.ID) (*mdl.User, error) {
	user, err := u.repos.User.GetById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userApiRepoImpl) Create(user mdl.User) error {
	if _, err := u.GetById(user.Id); err == nil {
		log.Printf("user %v already exists", user)
		return fmt.Errorf("user %s already exists", user.Id)
	}
	log.Printf("creating user %v", user)
	return u.repos.User.Create(user)
}
