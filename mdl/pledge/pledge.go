package pledge

import (
	"github.com/satori/go.uuid"
	"psychic-rat/mdl/item"
	"psychic-rat/mdl/user"
	"time"
)

type Id string

type Record interface {
	Id() Id
	UserId() user.Id
	ItemId() item.Id
	TimeStamp() time.Time
}

type ByTimeStamp []Record

func (b ByTimeStamp) Len() int           { return len(b) }
func (b ByTimeStamp) Less(i, j int) bool { return b[i].TimeStamp().Before(b[j].TimeStamp()) }
func (b ByTimeStamp) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

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
