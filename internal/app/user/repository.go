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
	users := repo.GetUsersByUser(newUser)
	if len(users) != 0 {
		return users, models.UserExistsError
	}

	var user models.User
	err := repo.DB.QueryRow(`INSERT INTO "user" ("nickname", "fullname", "about", "email")
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

func (repo *RepoPgx) GetUsersByUser(user models.User) (users []models.User) {
	userByNickname, err := repo.GetUserByNickname(user.Nickname)
	if !userByNickname.IsEmpty() && err == nil {
		users = append(users, userByNickname)
	}

	userByEmail, err := repo.GetUserByEmail(user.Email)
	if !userByEmail.IsEmpty() && err == nil && userByEmail != userByNickname {
		users = append(users, userByEmail)
	}

	return
}

func (repo *RepoPgx) UpdateUser(newUserData models.User) (models.User, error) {
    userExists, err := repo.GetUserByEmail(newUserData.Email)
	if !userExists.IsEmpty() {
		return models.User{}, models.NewUserDataError
	}

	var user models.User
	err = repo.DB.QueryRow(`UPDATE "user" SET "fullname" = $2, "about" = $3, "email" = $4
									WHERE "nickname" = $1 RETURNING "nickname", "fullname", "about", "email";`,
		newUserData.Nickname, newUserData.Fullname, newUserData.About, newUserData.Email,
	).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return user, err
}
