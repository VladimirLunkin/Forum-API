package models

type User struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (u User) IsEmpty() bool {
	return u.Nickname == ""
}

type UserRep interface {
	CreateUser(newUser User) ([]User, error)
	GetUserByNickname(nickname string) (User, error)
	UpdateUser(newUserData User) (User, error)
}
