package models

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

func (u User) IsEmpty() bool {
	return u.Nickname == ""
}

type UserRep interface {
	CreateUser(newUser User) ([]User, error)
	GetUserByNickname(nickname string) (User, error)
	UpdateUser(newUserData User) (User, error)
}
