package repo

import (
	"psychic-rat/mdl/pledge"
	"psychic-rat/mdl/user"
)

type Pledges interface {
	Create(pledge pledge.Record) error
	GetById(id pledge.Id) (pledge.Record, error)
	GetByUser(id user.Id) []pledge.Id
}
