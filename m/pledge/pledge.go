package pledge

import (
	"time"
	"psychic-rat/m/pubuser"
	"psychic-rat/m/item"
	"github.com/satori/go.uuid"
)

type Id string

type Record interface {
	Id() Id
	UserId() pubuser.Id
	ItemId() item.Id
	TimeStamp() time.Time
}

type Repo interface {
	Create(pledge Record) (Id, error)
	GetById(id Id) (Record, error)
	GetByUser(id pubuser.Id) []Id
}

type record struct {
	id        Id
	userId    pubuser.Id
	itemId    item.Id
	timestamp time.Time
}

func New(userId pubuser.Id, itemId item.Id, timestamp time.Time) Record {
	return &record{id: Id(uuid.NewV4().String()), userId: userId, itemId: itemId, timestamp: timestamp}
}

func (r *record) Id() Id {
	return r.id
}

func (r *record) UserId() pubuser.Id {
	return r.userId
}

func (r *record) ItemId() item.Id {
	return r.itemId
}

func (r *record) TimeStamp() time.Time {
	return r.timestamp
}
