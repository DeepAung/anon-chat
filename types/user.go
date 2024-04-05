package types

import "github.com/google/uuid"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func NewUser(username string) User {
	return User{
		Id:       uuid.NewString(),
		Username: username,
	}
}

var SystemUser = NewUser("system")
