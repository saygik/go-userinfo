package mattermost

import (
	"context"
	"fmt"
)

func (r *Repository) GetUsersWithSessions() (err error) {
	users, _, _ := r.client.GetUsers(context.Background(), 0, 1500, "")
	_ = users

	for _, user := range users {
		if len(user.AuthService) > 0 {
			_ = user.AuthService
		}
		sessions, _, _ := r.client.GetSessions(context.Background(), user.Id, "")
		if len(sessions) > 0 {
			fmt.Printf("User ID: %s, Sessions: %d\n", user.Id, len(sessions))
		}
	}
	return nil
}
