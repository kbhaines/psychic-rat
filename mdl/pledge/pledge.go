package pledge

import (
	"time"
	"psychic-rat/mdl/user"
	"psychic-rat/mdl/item"
	"github.com/satori/go.uuid"
)

type Id string

type Record interface {
	Id() Id
	UserId() user.Id
	ItemId() item.Id
	TimeStamp() time.Time
}

type Repo interface {
	Create(pledge Record) (Id, error)
	GetById(id Id) (Record, error)
	GetByUser(id user.Id) []Id
}

type record struct {
	id        Id
	userId    user.Id
	itemId    item.Id
	timestamp time.Time
}

func New(userId user.Id, itemId item.Id, timestamp time.Time) Record {
	return &record{id: Id(uuid.NewV4().String()), userId: userId, itemId: itemId, timestamp: timestamp}
}

func (r *record) Id() Id {
	return r.id
}

func (r *record) UserId() user.Id {
	return r.userId
}

func (r *record) ItemId() item.Id {
	return r.itemId
}

func (r *record) TimeStamp() time.Time {
	return r.timestamp
}
