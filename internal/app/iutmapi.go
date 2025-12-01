package app

type IUTMApiClient struct {
	url      string
	user     string
	password string
}

func (a *App) newIUTMApiConnection(url string, user string, password string) *IUTMApiClient {
	return &IUTMApiClient{
		url:      url,
		user:     user,
		password: password,
	}
}
