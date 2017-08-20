package privuser

type Id int

type Record struct {
	Id         Id
	Email      string
	FullName   string
	AuthMethod string
	AuthSecret string
}

type Repo interface {
	Create(user Record) (Id, error)
	GetByEmail(email string) (Record, error)
	GetById(id Id) (Record, error)
}
