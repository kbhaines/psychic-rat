package api

import (
	"psychic-rat/mdl"
)

func NewUserApi(repos Repos) *userApiRepoImpl {
	return &userApiRepoImpl{repos: repos}
}

type userApiRepoImpl struct {
	repos Repos
}

func (u *userApiRepoImpl) GetById(userId mdl.Id) (*mdl.UserRecord, error) {
	user, err := u.repos.User.GetById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
