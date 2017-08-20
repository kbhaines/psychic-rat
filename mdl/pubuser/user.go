package pubuser

type Id int

type Record struct {
	Id        Id
	Country   string
	FirstName string
}

type Repo interface {
	Create(user Record) (Id, error)
	GetById(id Id) (Record, error)
}
