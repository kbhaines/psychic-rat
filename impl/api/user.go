package api

import (
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

func NewUserApi(repos repo.Repos) *userApiRepoImpl {
	return &userApiRepoImpl{repos: repos}
}

type userApiRepoImpl struct {
	repos repo.Repos
}

func (u *userApiRepoImpl) GetById(userId mdl.Id) (*mdl.UserRecord, error) {
	user, err := u.repos.User.GetById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
