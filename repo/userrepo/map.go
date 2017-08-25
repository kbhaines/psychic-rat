package userrepo

import (
	"errors"
	"psychic-rat/repo"
	"fmt"
	"psychic-rat/mdl/user"
)

// declare that we implement Repo interface
var userRepo repo.Users = &userRepoMap{make(map[user.Id]user.Record)}

func GetUserRepoMapImpl() repo.Users {
	return userRepo
}

type userRepoMap struct {
	records map[user.Id]user.Record
}

func (repo *userRepoMap) Create(i user.Record) error {
	if _, found := repo.records[i.Id()]; found {
		return fmt.Errorf("user id %v already exists", i.Id())
	}
	repo.records[i.Id()] = i
	return nil
}

func (repo *userRepoMap) GetById(id user.Id) (user.Record, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return item, nil
}

func (repo *userRepoMap) List() []user.Id {
	ids := make([]user.Id, len(repo.records))
	i := 0
	for id := range repo.records {
		ids[i] = id
		i++
	}
	return ids
}
