package entities

type UserId string

type CountryId string

type PublicUser struct {
	Id        UserId
	Country   CountryId
	FirstName string
}

type PrivateUser struct {
	Id         UserId
	Email      string
	FullName   string
	AuthMethod string
	AuthSecret string
}

type PublicUserRepo interface {
	Create(user PublicUser) (UserId, error)
	GetById(id UserId) (PublicUser, error)
}

type PrivateUserRepo interface {
	Create(user PrivateUser) (UserId, error)
	GetByEmail(email string) (PublicUser, error)
	GetById(id UserId) (PublicUser, error)
}
