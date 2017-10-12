package api

import (
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

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
