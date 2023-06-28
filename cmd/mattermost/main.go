package main

import (
	"context"

	// "fmt"
	"mattermost/config"

	"mattermost/internal/models"
	"mattermost/pkg/bun"
	"mattermost/pkg/logger"
	"mattermost/pkg/redmine"

	"mattermost/pkg/mattermost"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	config := config.New().Load()

	log := logger.New(config.Logger)

	var report models.Report

	postgres, err := bun.New(config.DB, log)
	_ = redmine.TimeEntry{
		IssueID:    9984,
		SpentOn:    "2022-06-24",
		Hours:      "1",
		ActivityID: "8",
		Comments:   "TestTrack",
		//	UserId:     "266",
		CustomField: redmine.CustomFieldObject{
			ID:    33,
			Name:  "Project manager (timelog)",
			Value: "Evgeny Dubetsky",
		},
	}

	if err != nil {

		log.Panic(ctx, "connecting to DB", zap.Error(err))
	}
	defer postgres.Close()
	err = postgres.CreateRecord(ctx, report)
	if err != nil {
		print("model creation failed")
	}

	mattermostClient, err := mattermost.New(config.Mattermost, postgres)
	if err != nil {
		log.Panic(ctx, "mattermost connection", zap.Error(err))
		print("mattermost conection failed")
	}
	err = mattermostClient.GetTeamByName()
	if err != nil {
		print("mattermost team not found")
	}
	mattermostClient.DialWebSockets(ctx)

	redmineClient := redmine.New(log, config.Redmine)
	err = redmine.Retrieve

	/*
		var a redmine.IssueSingleGetRequest
		x, _ := redmineClient.IssueSingleGet(ctx, 9984, a)
		fmt.Printf("%+v\n", x)

		err = redmineClient.IssuesToStopped(ctx, 266)
		if err != nil {
			print("ошибка")
		}
		err = redmineClient.IssueToInprogress(ctx, 266, 9984)
		if err != nil {
			print("Ошибка")
		}
		err = redmineClient.TimeTracker(ctx, SpendTime)
		if err != nil {
			print("TimeTrack successful failed")
		}
	*/
}
