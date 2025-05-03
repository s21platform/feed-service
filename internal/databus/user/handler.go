package user

import (
	"encoding/json"
	"log"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/model"
	"github.com/s21platform/metrics-lib/pkg"
	"golang.org/x/net/context"
)

type Handler struct {
	dbR DBRepo
}

func New(dbR DBRepo) *Handler {
	return &Handler{dbR: dbR}
}

func convertMessage(bMessage []byte, targer interface{}) error {
	err := json.Unmarshal(bMessage, targer)
	if err != nil {
		return err
	}
	return nil
}

// TO DO: функцию получения подписчиков пользователя
func (h *Handler) Handler(ctx context.Context, in []byte) error {
	m := pkg.FromContext(ctx, config.KeyMetrics)
	var msg model.CreateUserPostMessage

	log.Printf("Received message: %s", string(in))

	err := convertMessage(in, &msg)
	if err != nil {
		m.Increment("create_user_post.error")
		log.Printf("failed to convert message: %v", err)
		return err
	}

	log.Printf("Parsed message: UserUUID=%s, PostUUID=%s", msg.UserUUID, msg.PostUUID)

	err = h.dbR.CreatePost(ctx, msg.PostUUID, "user")
	if err != nil {
		m.Increment("create_user_post.error")
		log.Printf("failed to create post: %v", err)
		return err
	}

	log.Printf("Successfully created post for user %s", msg.UserUUID)
	m.Increment("create_user_post.success")

	return nil
}
