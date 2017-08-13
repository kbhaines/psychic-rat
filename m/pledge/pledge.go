package pledge

import (
	"time"
	"psychic-rat/m/pubuser"
	"psychic-rat/m/item"
)

type Id string

type Record struct {
	Id        Id
	UserId    pubuser.Id
	ItemId    item.Id
	Timestamp time.Time
}

type Repo interface {
	Create(pledge Record) (Id, error)
	GetById(id Id) (Record, error)
	GetByUser(id pubuser.Id) []Id
}
