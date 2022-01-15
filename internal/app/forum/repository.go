package forum

import (
	"database/sql"
	"fmt"
	"strconv"
	"tech-db-forum/internal/app/models"
	"time"
)

type RepoPgx struct {
	DB *sql.DB
}

func NewPgxRepository(db *sql.DB) *RepoPgx {
	return &RepoPgx{DB: db}
}

func (repo *RepoPgx) CreateForum(newForum models.Forum) (forum models.Forum, err error) {
	err = repo.DB.QueryRow(`INSERT INTO "forum" ("title", "user", "slug")
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
	err := repo.DB.QueryRow(`SELECT "id", "title", "user", "slug", "posts", "threads" FROM "forum"
						WHERE "slug" = $1;`, slug,
	).Scan(&forum.Id, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	if err != nil {
		return models.Forum{}, err
	}

	return forum, nil
}

func (repo *RepoPgx) CreateThread(newThread models.Thread) (thread models.Thread, err error) {
	if newThread.Slug != "" {
		err = repo.DB.QueryRow(`SELECT "id", "title", "author", "forum", "message", "votes", "slug", "created"
			FROM "thread" WHERE "slug" = $1;`, newThread.Slug,
		).Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created)
		if err == nil {
			return thread, models.ThreadAlreadyExistsError
		}
	}

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

	return thread, nil
}

func (repo *RepoPgx) GetThreads(slug, limit, since, desc string) ([]models.Thread, error) {
	threads := make([]models.Thread, 0)

	query := `SELECT "id", "title", "author", "forum", "message", "votes", "slug", "created" FROM "thread" WHERE "forum" = $1`
	if since != "" {
		sign := ">="
		if desc == "desc" {
			sign = "<="
		}
		query += fmt.Sprintf(` AND "created" %s '%s'`, sign, since)
	}
	query += fmt.Sprintf(` ORDER BY "created" %s LIMIT %s;`, desc, limit)

	rows, err := repo.DB.Query(query, slug)
	if err != nil {
		return []models.Thread{}, err
	}
	defer rows.Close()

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

func (repo *RepoPgx) GetUsers(forum models.Forum, limit, since, desc string) ([]models.User, error) {
	users := make([]models.User, 0)

	query := `SELECT "nickname", "about", "email", "fullname" FROM "user"
				WHERE "id" IN (SELECT "user" FROM "forum_user" WHERE forum = $1)`
	if since != "" {
		sign := ">"
		if desc == "desc" {
			sign = "<"
		}
		query += fmt.Sprintf(` AND "nickname" %s '%s'`, sign, since)
	}
	query += fmt.Sprintf(` ORDER BY "nickname" %s LIMIT %s;`, desc, limit)

	rows, err := repo.DB.Query(query, forum.Id)
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Nickname,
			&user.About,
			&user.Email,
			&user.Fullname,
		)
		if err != nil {
			return []models.User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repo *RepoPgx) CreatePosts(thread models.Thread, newPost []models.Post) (posts []models.Post, err error) {
	if len(newPost) == 0 {
		return []models.Post{}, nil
	}

	query := `INSERT INTO "post" ("parent", "author", "message", "forum", "thread", "created") VALUES `

	var newPostsData []interface{}
	created := time.Now()

	for i, post := range newPost {
		err := repo.DB.QueryRow(`SELECT "nickname" FROM "user" WHERE "nickname" = $1;`, post.Author).Scan(&post.Author)
		if err != nil {
			return []models.Post{}, models.NoAuthorError
		}

		if post.Parent != 0 {
			var id int64
			err := repo.DB.QueryRow(`SELECT "id" FROM "post" WHERE "thread" = $1 and "id" = $2;`, thread.Id, post.Parent).Scan(&id)
			if err != nil {
				return []models.Post{}, err
			}
		}

		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", 1+i*6, 2+i*6, 3+i*6, 4+i*6, 5+i*6, 6+i*6)
		newPostsData = append(newPostsData, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, created)
	}

	query += `RETURNING "id", "parent", "author", "message", "isEdited", "forum", "thread", "created";`

	rows, err := repo.DB.Query(query, newPostsData...)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return []models.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *RepoPgx) GetThreadBySlugOrId(slugOrId string) (thread models.Thread, err error) {
	id, _ := strconv.Atoi(slugOrId)
	err = repo.DB.QueryRow(`SELECT "id", "title", "author", "forum", "message", "votes", "slug", "created"
		FROM "thread" WHERE "id" = $1 OR "slug" = $2`, id, slugOrId).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	return
}

func (repo *RepoPgx) Vote(thread models.Thread, vote models.Vote) (models.Thread, error) {
	var user int64
	err := repo.DB.QueryRow(`SELECT "id" FROM "user" WHERE "nickname" = $1;`, vote.Nickname).Scan(&user)
	if err != nil {
		return models.Thread{}, err
	}

	var voteId int64
	err = repo.DB.QueryRow(`SELECT "id" FROM "vote" WHERE "user" = $1 AND "thread" = $2`,
		user, thread.Id).Scan(&voteId)

	if err == nil && voteId != 0 {
		err = repo.DB.QueryRow(`UPDATE "vote" SET "voice" = $1 WHERE "id" = $2;`,
			vote.Voice, user).Err()
		if err != nil {
			return models.Thread{}, err
		}
		thread.Votes += 2 * vote.Voice

		return thread, nil
	}

	err = repo.DB.QueryRow(`INSERT INTO "vote" ("user", "thread", "voice") VALUES ($1, $2, $3) RETURNING "user";`,
		user, thread.Id, vote.Voice).Scan(&user)
	if err != nil {
		return models.Thread{}, err
	}
	thread.Votes += vote.Voice

	return thread, nil
}

func (repo *RepoPgx) GetPosts(thread models.Thread, limit, since, sort, desc string) ([]models.Post, error) {
	posts := make([]models.Post, 0)

	query := `SELECT "id", "parent", "author", "message", "isEdited", "forum", "thread", "created"
				FROM "post" WHERE "thread" = $1 `

	sign := ">"
	if desc == "desc" {
		sign = "<"
	}

	switch sort {
	case "flat":
		if since != "" {
			query += fmt.Sprintf(`AND "id" %s %s `, sign, since)
		}
		query += fmt.Sprintf(`ORDER BY "created" %s, "id" %s LIMIT %s `, desc, desc, limit)
	case "tree":
		if since != "" {
			query += fmt.Sprintf(`AND "path" %s (SELECT "path" FROM "post" WHERE "id" = %s) `, sign, since)
		}
		query += fmt.Sprintf(`ORDER BY path[1] %s, path %s LIMIT %s `, desc, desc, limit)
	case "parent_tree":
		query += `AND "path" && (SELECT ARRAY (SELECT "id" FROM "post" WHERE "thread" = $1 AND "parent" = 0 `
		if since != "" {
			query += fmt.Sprintf(`AND "path" %s (SELECT path[1:1] FROM "post" WHERE "id" = %s) `, sign, since)
		}
		query += fmt.Sprintf(`ORDER BY path[1] %s, path LIMIT %s)) ORDER BY path[1] %s, path `, desc, limit, desc)
	default:
		return []models.Post{}, models.SortError
	}

	rows, err := repo.DB.Query(query, thread.Id)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return []models.Post{}, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repo *RepoPgx) UpdateThread(oldThread, newThread models.Thread) (thread models.Thread, err error) {
	err = repo.DB.QueryRow(`UPDATE "thread" SET "title" = $1, "message" = $2 WHERE "id" = $3
			RETURNING "id", "title", "author", "forum", "message", "votes", "slug", "created";`,
		newThread.Title, newThread.Message, oldThread.Id,
	).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Votes,
		&thread.Slug,
		&thread.Created,
	)
	return
}

func (repo *RepoPgx) GetPost(id int) (post models.Post, err error) {
	err = repo.DB.QueryRow(`SELECT "id", "parent", "author", "message", "isEdited", "forum", "thread", "created"
		FROM "post" WHERE "id" = $1;`, id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	return
}

func (repo *RepoPgx) UpdatePost(id int, newPost models.Post) (post models.Post, err error) {
	err = repo.DB.QueryRow(`UPDATE "post" SET "message" = $1, "isEdited" = true WHERE "id" = $2
			RETURNING "id", "parent", "author", "message", "isEdited", "forum", "thread", "created";`,
		newPost.Message, id,
	).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)
	return
}
