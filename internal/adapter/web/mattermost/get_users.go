package mattermost

import (
	"context"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/saygik/go-userinfo/internal/entity"
)

// func (r *Repository) GetUsers1() ([]entity.MattermostUser, error) {

// 	users := []entity.MattermostUser{}
// 	url_api := r.url
// 	urlUsers := url_api + "/users/search"

// 	var jsonStr = []byte(`{ "term": "*", "allow_inactive": false, "limit": 1000 }`)
// 	req, err := http.NewRequest("POST", urlUsers, bytes.NewBuffer(jsonStr))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+r.token)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return users, err
// 	}

// 	defer resp.Body.Close()
// 	body, _ := io.ReadAll(resp.Body)
// 	err = json.Unmarshal(body, &users)
// 	if err != nil {
// 		return users, err
// 	}

// 	return users, err
// }

func (r *Repository) GetUsers() ([]entity.MattermostUserWithSessions, error) {
	usersWithSessions := []entity.MattermostUserWithSessions{}

	users, _, err := r.client.SearchUsers(context.Background(), &model.UserSearch{
		Term:  "*",
		Limit: 1000,
	})
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		status, err := r.GetUserStatus(user.Id)
		if err != nil {
			_ = status
		}
		sessions, _, _ := r.client.GetSessions(context.Background(), user.Id, "")
		entitySessions := []entity.MattermostSession{}

		if len(sessions) > 0 {
			for _, session := range sessions {
				sessionsn := entity.MattermostSession{
					Id:             session.Id,
					CreateAt:       session.CreateAt,
					ExpiresAt:      session.ExpiresAt,
					LastActivityAt: session.LastActivityAt,
					UserId:         session.UserId,
					DeviceId:       session.DeviceId,
					Roles:          session.Roles,
					IsOAuth:        session.IsOAuth,
					ExpiredNotify:  session.ExpiredNotify,
					Props:          entity.StringMap(session.Props),
					Local:          session.Local,
				}
				entitySessions = append(entitySessions, sessionsn)
			}
		}

		userws := entity.MattermostUserWithSessions{
			Id:             user.Id,
			Username:       user.Username,
			Email:          user.Email,
			EmailVerified:  user.EmailVerified,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Nickname:       user.Nickname,
			Position:       user.Position,
			CreateAt:       user.CreateAt,
			UpdateAt:       user.UpdateAt,
			DeleteAt:       user.DeleteAt,
			Roles:          user.Roles,
			AuthService:    user.AuthService,
			LastActivityAt: user.LastActivityAt,
		}
		userws.Sessions = entitySessions
		userws.Status = status.Status
		userws.LastActivityAt = status.LastActivityAt
		usersWithSessions = append(usersWithSessions, userws)
	}
	return usersWithSessions, err
}
