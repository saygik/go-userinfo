package iutm

type Repository struct {
	url      string
	User     string
	Password string
}

func New(url string, user string, password string) *Repository {
	return &Repository{
		url:      url,
		User:     user,
		Password: password,
	}
}
