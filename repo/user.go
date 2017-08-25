package repo

import "psychic-rat/mdl/user"

type Users interface {
	Create(user user.Record) error
	GetById(id user.Id) (user.Record, error)
}
