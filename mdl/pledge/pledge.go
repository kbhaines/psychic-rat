package pledge

import (
	"time"
	"psychic-rat/mdl/user"
	"github.com/satori/go.uuid"
	"psychic-rat/mdl/item"
)

type Id string

type Record interface {
	Id() Id
	UserId() user.Id
	ItemId() item.Id
	TimeStamp() time.Time
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
