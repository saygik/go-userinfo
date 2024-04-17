package glpiapi

type Repository struct {
	url       string
	token     string
	usertoken string
}

func New(url string, token string, usertoken string) *Repository {
	return &Repository{
		url:       url,
		token:     token,
		usertoken: usertoken,
	}
}
