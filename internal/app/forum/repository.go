package forum

import (
	"database/sql"
	"tech-db-forum/internal/app/models"
)

type RepoPgx struct {
	DB *sql.DB
}

func NewPgxRepository(db *sql.DB) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (repo *RepoPgx) Create(newForum models.Forum) (models.Forum, error) {
	var forum models.Forum
	err := repo.DB.QueryRow(`INSERT INTO "forum" ("title", "user", "slug")
						VALUES ($1, $2, $3) RETURNING "title", "user", "slug", "posts", "threads";`,
		newForum.Title, newForum.User, newForum.Slug,
	).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}
