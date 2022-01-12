package user

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

func (repo *RepoPgx) CreateUser(newUser models.User) ([]models.User, error) {
	var users []models.User
	userByNickname, err := repo.GetUserByNickname(newUser.Nickname)
	if !userByNickname.IsEmpty() && err == nil {
		users = append(users, userByNickname)
	}

	userByEmail, err := repo.GetUserByEmail(newUser.Email)
	if !userByEmail.IsEmpty() && err == nil && userByEmail != userByNickname {
		users = append(users, userByEmail)
	}

	if len(users) != 0 {
		return users, models.UserExistsError
	}

	var user models.User
	err = repo.DB.QueryRow(`INSERT INTO "user" ("nickname", "fullname", "about", "email")
									VALUES ($1, $2, $3, $4) RETURNING "nickname", "fullname", "about", "email";`,
		newUser.Nickname, newUser.Fullname, newUser.About, newUser.Email,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return users, err
}

func (repo *RepoPgx) GetUserByNickname(nickname string) (user models.User, err error) {
	err = repo.DB.QueryRow(`SELECT "nickname", "fullname", "about", "email" FROM "user"
									WHERE "nickname" = $1;`, nickname,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return
}

func (repo *RepoPgx) GetUserByEmail(email string) (user models.User, err error) {
	err = repo.DB.QueryRow(`SELECT "nickname", "fullname", "about", "email" FROM "user"
									WHERE "email" = $1;`, email,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return
}
