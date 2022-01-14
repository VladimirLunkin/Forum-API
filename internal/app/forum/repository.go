package forum

import (
	"database/sql"
	"fmt"
	"tech-db-forum/internal/app/models"
	"time"
)

type RepoPgx struct {
	DB *sql.DB
}

func NewPgxRepository(db *sql.DB) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (repo *RepoPgx) CreateForum(newForum models.Forum) (models.Forum, error) {
	var forum models.Forum
	err := repo.DB.QueryRow(`INSERT INTO "forum" ("title", "user", "slug")
						VALUES ($1, $2, $3) RETURNING "title", "user", "slug", "posts", "threads";`,
		newForum.Title, newForum.User, newForum.Slug,
	).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, models.SlugAlreadyExistsError
	}

	return forum, nil
}

func (repo *RepoPgx) GetForumBySlug(slug string) (models.Forum, error) {
	var forum models.Forum
	err := repo.DB.QueryRow(`SELECT "title", "user", "slug", "posts", "threads" FROM "forum"
						WHERE "slug" = $1;`, slug,
	).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}

func (repo *RepoPgx) CreateThread(newThread models.Thread) (thread models.Thread, err error) {
	if newThread.Created.String() == "" {
		newThread.Created = time.Now()
	}

	err = repo.DB.QueryRow(`INSERT INTO "thread" ("title", "author", "forum", "message", "slug", "created")
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING "id", "title", "author", "forum", "message", "votes", "slug", "created";`,
		newThread.Title, newThread.Author, newThread.Forum, newThread.Message, newThread.Slug, newThread.Created,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created)
	if err != nil {
		return models.Thread{}, models.SlugAlreadyExistsError
	}

	err = repo.DB.QueryRow(`UPDATE "forum" SET threads = threads + 1 WHERE slug = $1 RETURNING id;`, newThread.Slug).Err()
	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func (repo *RepoPgx) GetThreads(slug, limit, since, desc string) ([]models.Thread, error) {
	threads := make([]models.Thread, 0)

	query := `SELECT "id", "title", "author", "forum", "message", "votes", "slug", "created" FROM "thread" WHERE "forum" = $1`
	//if since != "" {
	//	query += fmt.Sprintf(" AND created %s '%s'", comparisonSign, since)
	//}
	query += fmt.Sprintf(`ORDER BY created %s LIMIT %s;`, desc, limit)

	rows, err := repo.DB.Query(query, slug)
	if err != nil {
		return []models.Thread{}, err
	}

	for rows.Next() {
		var thread models.Thread
		err := rows.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
		if err != nil {
			return []models.Thread{}, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}
