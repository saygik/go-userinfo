package app

type Integrations struct {
	AddCommentFromApi                  string   `json:"add-comment-from-api,omitempty"`
	DisableCalendarTaskNotificationApi string   `json:"disable-calendar-task-notification-api,omitempty"`
	AllowedHosts                       []string `json:"allowed-hosts,omitempty"`
}
type MattClient struct {
	url          string
	token        string
	integrations Integrations
}

func (a *App) newMattermostConnection(
	url string,
	token string,
	addCommentFromApi string,
	disableCalendarTaskNotificationApi string,
	allowedHosts []string,
) *MattClient {
	return &MattClient{
		url:          url,
		token:        token,
		integrations: Integrations{AddCommentFromApi: addCommentFromApi, DisableCalendarTaskNotificationApi: disableCalendarTaskNotificationApi, AllowedHosts: allowedHosts},
	}
}
