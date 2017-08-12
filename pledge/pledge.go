package pledge

import (
	"time"
	"psychic-rat/pubuser"
	"psychic-rat/item"
)

type Id string

type Pledge struct {
	Id        Id
	UserId    pubuser.Id
	ItemId    item.Id
	Timestamp time.Time
}
