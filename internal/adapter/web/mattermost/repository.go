package mattermost

type Repository struct {
	url   string
	token string
}

func New(url string, token string) *Repository {
	return &Repository{
		url:   url,
		token: token,
	}
}
