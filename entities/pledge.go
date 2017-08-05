package entities

import "time"

type PledgeId string

type Pledge struct {
	Id        PledgeId
	UserId    UserId
	ItemId    ItemId
	Timestamp time.Time
}
