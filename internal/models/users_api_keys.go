package models

import "github.com/uptrace/bun"

//Информация о пользователе
type UserModel struct {
	bun.BaseModel `bun:"table:UserModel,alias:u"`
	MattermostID  string `bun:"mm_id"`
	RedmineId     int    `bun:"redmine_id"`
	RedmineApiKey string `bun:"redmine_api_key"`
}
