package user

type Id int

type Record interface {
	Id() Id
	Country() string
	FirstName() string

	Email() string
	AuthToken() string
	Fullname() string
}

type record struct {
	id        Id
	country   string
	firstName string
	email     string
	authToken string
	fullname  string
}

func New(email string, country string, firstname string) Record {
	return &record{email: email, country: country, firstName: firstname}
}

func (r *record) Id() Id {
	return r.id
}

func (r *record) Country() string {
	return r.country
}

func (r *record) FirstName() string {
	return r.firstName
}

func (r *record) Email() string {
	return r.email
}

func (r *record) AuthToken() string {
	return r.authToken
}

func (r *record) Fullname() string {
	return r.fullname
}

