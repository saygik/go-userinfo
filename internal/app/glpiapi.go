package app

type GLPIApiClient struct {
	url       string
	token     string
	usertoken string
}

func (a *App) newGLPIApiConnection(url string, token string, usertoken string) *GLPIApiClient {
	return &GLPIApiClient{
		url:       url,
		token:     token,
		usertoken: usertoken,
	}
}
