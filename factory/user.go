package factory

import (
	"errors"
	"psychic-rat/mdl/user"
)

// declare that we implement Repo interface
var userRepo user.Repo = &userRepoMap{make(map[user.Id]user.Record)}

func GetUserRepo() user.Repo {
	return userRepo
}

type userRepoMap struct {
	records map[user.Id]user.Record
}

func (repo *userRepoMap) Create(i user.Record) (user.Id, error) {
	newId := user.Id(len(repo.records))
	repo.records[newId] = i
	return newId, nil
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
