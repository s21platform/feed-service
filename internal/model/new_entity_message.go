package model

const User string = "user"

type NewEntityMessage struct {
	UserUUID   string `json:"user_uuid"`
	EntityUUID string `json:"entity_uuid"`
}
