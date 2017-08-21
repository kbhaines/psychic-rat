package company

type Id string
type Record struct {
	Id   Id
	Name string

}

type Repo interface {
	Create(company Record) (Id, error)
	GetById(id Id) (Record, error)
	List() []Id
}

