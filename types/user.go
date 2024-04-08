package types

import "github.com/google/uuid"

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

var SystemUser = User{
	Id:       uuid.NewString(),
	Username: "system",
}
