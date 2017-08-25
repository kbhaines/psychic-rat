package company

type Id string

type Record interface {
	Id() Id
	Name() string
}

type record struct {
	id   Id
	name string
}

func New(id Id, name string) Record {
	return &record{id, name}
}

func (r *record) Id() Id {
	return r.id
}

func (r *record) Name() string {
	return r.name
}
