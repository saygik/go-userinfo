package iutm

type Repository struct {
	url      string
	user     string
	password string
}

func New(url string, user string, password string) *Repository {
	return &Repository{
		url:      url,
		user:     user,
		password: password,
	}
}
