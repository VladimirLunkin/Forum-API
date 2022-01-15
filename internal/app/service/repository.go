package service

import (
	"database/sql"
	"fmt"
	"tech-db-forum/internal/app/models"
)

type RepoPgx struct {
	DB *sql.DB
}

func NewPgxRepository(db *sql.DB) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (repo *RepoPgx) GetStatus() (models.Service, error) {
	var status models.Service
	err := repo.DB.QueryRow(`SELECT
									(SELECT COUNT(*) FROM "forum"),
									(SELECT COUNT(*) FROM "post"),
									(SELECT COUNT(*) FROM "thread"), 
									(SELECT COUNT(*) FROM "user");`,
	).Scan(&status.Forum, &status.Post, &status.Thread, &status.User)
	fmt.Println(status, err)
	return status, err
}

func (repo *RepoPgx) Clear() error {
	err := repo.DB.QueryRow(`TRUNCATE "user", "forum", "thread", "post", "vote", "forum_user" CASCADE;`).Scan()
	return err
}
