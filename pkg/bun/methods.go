package bun

import (
	"context"
	"mattermost/internal/models"

	"database/sql"
	"errors"
)

//Функции create будут объединены в одну, принимающую интерфейс
func (db *DB) CreateUserInfoRecord(ctx context.Context, UserInfo models.UserInfo) (err error) {

	var result sql.Result
	if result, err = db.DB.NewInsert().Model(UserInfo).Exec(ctx); err != nil {
		return
	}

	var rows int64
	if rows, err = result.RowsAffected(); err != nil {
		return
	}

	if rows == 0 {
		return errors.New("failed to insert new report to database")
	}

	return
}
func (db *DB) CreateReportRecord(ctx context.Context, Report models.Report) (err error) {

	var result sql.Result
	if result, err = db.DB.NewInsert().Model(Report).Exec(ctx); err != nil {
		return
	}

	var rows int64
	if rows, err = result.RowsAffected(); err != nil {
		return
	}

	if rows == 0 {
		return errors.New("failed to insert new report to database")
	}

	return
}
func(db *DB) CreateUserModelRecord(ctx context.Context,UserModel models.UserModel)(err error){
	var result sql.Result
	if result, err = db.DB.NewInsert().Model(UserModel).Exec(ctx); err != nil {
		return
	}

	var rows int64
	if rows, err = result.RowsAffected(); err != nil {
		return
	}

	if rows == 0 {
		return errors.New("failed to insert new report to database")
	}

	return
}

// func (db *DB) UpdateRecord(ctx context.Context, report *models.Report) (err error) {
// 	db.NewUpdate().Model(&report).Column("report_text").Where("user_name = ?", report.Username).Where("report_text = null").Exec(ctx)
// 	_, err = db.DB.NewInsert().Model(&report).Exec(ctx)
// 	if err != nil {
// 		println("UpdateRecord error")
// 	}
// 	return
// }

// func (db *DB) Select(ctx context.Context, report *models.Report) {
// 	db.NewSelect().Model(&report).Where("MMLogin = ?").Scan(ctx) //добавить путь к файлу
// }

func (db *DB)UserModelToStruct(ctx context.Context,u *models.UserModel,MattermostID string) error{
	err := db.DB.NewSelect().Model(u).Where("mm_id = ?",MattermostID).Scan(ctx)
	if err != nil {

	}
	return err
}