package model

type CreateUserPostMessage struct {
	UserUUID string `json:"user_uuid"`
	PostUUID string `json:"post_uuid"`
}
