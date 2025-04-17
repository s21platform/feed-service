package model

import (
	feed "github.com/s21platform/feed-proto/feed-proto"
)

type Post struct {
	OwnerUUID string `db:"owner_uuid"`
	Content   string `db:"content"`
}

func (a *Post) PostToDTO(UUID string, in *feed.CreateUserPostIn) (Post, error) {
	result := Post{
		OwnerUUID: UUID,
		Content:   in.Content,
	}

	return result, nil
}
