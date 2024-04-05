package app

type MattClient struct {
	url   string
	token string
}

func (a *App) newMattermostConnection(url string, token string) *MattClient {
	return &MattClient{
		url:   url,
		token: token,
	}
}
