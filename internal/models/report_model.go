package models

import (
	"github.com/uptrace/bun"
)

//Модель отчета по проекту
type Report struct {
	bun.BaseModel `bun:"table:reports,alias:u"`
	ReportID      int       `bun:"id,pk,autoincrement"`
	UserInfoID    int       `bun:"user_info_id"`
	UserInfo      *UserInfo `bun:"rel:belongs-to,join:UserInfo_UserInfoID=UserInfoID"`
	IssueID       string    `bun:"issue_id"`
	SpentOn       string    `bun:"spent_on"`
	Hours         string    `bun:"hours"`
	ActivityID    string    `bun:"activity_id"`
	Comments      string    `bun:"comments"`
}
