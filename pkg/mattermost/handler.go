package mattermost

import (
	"context"
	"encoding/json"
	"mattermost/internal/models"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

type state int

const (
	LISTENING state = iota
	REPORT_ISSUE
	REPORT_HOURS
	REPORT_ACTIVITY
	REPORT_TEXT
	API_KEY_WAITING
)

//находится ли пользователь в процессе написания отчета
var UsersState = map[string]state{}

func (c *Client) handlMsg(ctx context.Context, event *model.WebSocketEvent) error {
	var post model.Post
	var UserInfo models.UserInfo
	var Report models.Report
	var UserModel models.UserModel
	
	
	err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post)
	if err != nil {
		print(err)
		return err
	}
	switch UsersState[post.UserId] {
	case LISTENING:
		switch post.Message {
		case "!help":
			c.sendMsgToUser(c.config.Help, post.ChannelId)
		case "!апиключ":
			//внесение апи ключа в бд
			//добавить проверку валидности апи ключа
			UsersState[post.UserId] = API_KEY_WAITING
			c.sendMsgToUser("введите апи-ключ",post.ChannelId)
			
		case "!старт":
			//проверка наличия отчета
			count, err := c.db.NewSelect().Model(UserInfo).Where("finish_date is NULL").Count(ctx)
			if err != nil {

			}
			if count == 1 {
				c.sendMsgToUser(c.config.NoReport, post.ChannelId)
				return err
			}
			//таймер
			var t time.Duration
			t, err = time.ParseDuration("9h")
			time.AfterFunc(t, func() {
				reportStatus, _ := c.db.NewSelect().Model(UserInfo).Where("finish_date is NULL").Count(ctx)
				if reportStatus == 1 {
					c.sendMsgToUser(c.config.ReportRemind, post.ChannelId)
				}
			})

			c.sendMsgToUser(c.config.Start, post.ChannelId)
			UserInfo.CreationDate = time.Now()
			UserInfo.UserID = post.UserId
			err = c.db.CreateUserInfoRecord(ctx, UserInfo)
			if err != nil {

			}

		case "!отчет":
			//проверка наличия апи ключа
			apiKeyExistence, _ := c.db.NewSelect().Model(UserModel).Where("mm_id=?", post.UserId).Where("redmine_api_key NOT NULL").Count(ctx)
			if apiKeyExistence == 0 {
				c.sendMsgToUser(c.config.NoApiKey, post.ChannelId)
			}

			if _, err := c.db.NewSelect().Model(UserInfo).Column("id").Where("user_id=?", post.UserId).Order("DESC").Limit(1).Exec(ctx); err != nil {

			}

			UsersState[post.UserId] = REPORT_ISSUE

			c.sendMsgToUser("Введите номер задачи", post.ChannelId)
			Report.UserInfoID = UserInfo.UserInfoID
			err = c.db.CreateReportRecord(ctx, Report)
			if err != nil {

			}

		}
	case API_KEY_WAITING:
		UsersState[post.UserId] = LISTENING
		UserModel.MattermostID = post.UserId
		UserModel.RedmineApiKey = post.Message
		c.redmine.GetRedmineID(ctx,&UserModel) 
		err = c.db.CreateUserModelRecord(ctx, UserModel)
			if err != nil {

			}
	case REPORT_ISSUE:
		UsersState[post.UserId] = REPORT_HOURS
		Report.IssueID = post.Message
		c.sendMsgToUser("Введите время", post.ChannelId)
		c.db.NewUpdate().Model(Report).Where("issue_id is NULL").Exec(ctx)
	case REPORT_HOURS:
		UsersState[post.UserId] = REPORT_ACTIVITY
		Report.Hours = post.Message
		c.sendMsgToUser("Введите тип трека:", post.ChannelId) //дополнить описания треков
		c.db.NewUpdate().Model(Report).Where("hours is NULL").Exec(ctx)
	case REPORT_ACTIVITY:
		UsersState[post.UserId] = REPORT_TEXT
		Report.ActivityID = post.Message
		c.sendMsgToUser("Введите отчет", post.ChannelId)
		c.db.NewUpdate().Model(Report).Where("activity_id is NULL").Exec(ctx)
		//внесение в бд
	case REPORT_TEXT:
		UsersState[post.UserId] = LISTENING
		Report.Comments = post.Message
		c.sendMsgToUser(c.config.Finish, post.ChannelId)
		c.db.NewUpdate().Model(Report).Where("comments is NULL").Exec(ctx)

	}
	return err
}
