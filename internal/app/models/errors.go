package models

import "errors"

var (
	UserExistsError          = errors.New("user exists")
	NewUserDataError         = errors.New("user exists")
	SlugAlreadyExistsError   = errors.New("slug already exists")
	ThreadAlreadyExistsError = errors.New("thread already exists")
	SortError                = errors.New("undefined sort type")
	NoAuthorError            = errors.New("posts create no author")
)
