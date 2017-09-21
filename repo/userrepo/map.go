package userrepo

import (
	"errors"
	"fmt"
	"psychic-rat/mdl"
	"psychic-rat/repo"
)

// declare that we implement Repo interface
var userRepo repo.Users = &userRepoMap{make(map[mdl.Id]mdl.UserRecord)}

func GetUserRepoMapImpl() repo.Users {
	return userRepo
}

type userRepoMap struct {
	records map[mdl.Id]mdl.UserRecord
}

func (repo *userRepoMap) Create(i mdl.UserRecord) error {
	if _, found := repo.records[i.Id]; found {
		return fmt.Errorf("user id %v already exists", i.Id)
	}
	repo.records[i.Id] = i
	return nil
}

func (repo *userRepoMap) GetById(id mdl.Id) (*mdl.UserRecord, error) {
	item, found := repo.records[id]
	if !found {
		return nil, errors.New("not found")
	}
	return &item, nil
}

func (repo *userRepoMap) List() []mdl.Id {
	ids := make([]mdl.Id, len(repo.records))
	i := 0
	for id := range repo.records {
		ids[i] = id
		i++
	}
	return ids
}
