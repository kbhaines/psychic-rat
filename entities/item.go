package entities

type ItemId string
type Item struct {

}

type CompanyId string
type Company string

type ItemRepo interface {
	Create(item Item) (ItemId, error)
	GetById(id ItemId) (Item, error)
	List() []ItemId
}

type CompanyRepo interface {
	Create(company Company) (CompanyId, error)
	GetById(id CompanyId) (Company, error)
}
