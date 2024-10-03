package mattermost

import (
	"fmt"
)

func (r *Repository) GetUsersWithSessions() (err error) {
	users, _, _ := r.client.GetUsers(0, 1500, "")
	_ = users

	for _, user := range users {
		if len(user.AuthService) > 0 {
			_ = user.AuthService
		}
		sessions, _, _ := r.client.GetSessions(user.Id, "")
		if len(sessions) > 0 {
			fmt.Printf("User ID: %s, Sessions: %d\n", user.Id, len(sessions))
		}
	}
	return nil
}
