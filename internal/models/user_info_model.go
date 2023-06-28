package models

import (
	"time"

	"github.com/uptrace/bun"
)

//Информация о пользователе за конкретный день
type UserInfo struct {
	bun.BaseModel `bun:"table:reports,alias:u"`

	UserInfoID int `bun:"id,pk,autoincrement"`

	UserID string `bun:"user_id" `

	CreationDate time.Time `bun:"creation_date"`

	FinishDate time.Time `bun:"finish_date"`
}
