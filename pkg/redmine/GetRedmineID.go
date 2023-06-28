package redmine

import (
	"context"
	"mattermost/internal/models"
	"net/http"
)
type RedmineUser struct{
	ID       int    `json:"id"`
}
func(c *Redmine)GetRedmineID(ctx context.Context,u *models.UserModel)(err error){
	out := new(RedmineUser)
	if err = c.Request(ctx, http.MethodGet, "/users/current.json", nil, nil, out); err != nil {
		return 
	}
	u.RedmineId = out.ID
	return err
}
