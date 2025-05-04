package user

import (
	"encoding/json"
	"log"

	"golang.org/x/net/context"

	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/model"
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

func (h *Handler) Handler(ctx context.Context, in []byte) error {
	m := pkg.FromContext(ctx, config.KeyMetrics)
	var msg model.NewEntityMessage
	
	err := convertMessage(in, &msg)
	if err != nil {
		m.Increment("save_user_post.error")
		log.Printf("failed to convert message: %v", err)
		return err
	}
	
	// TODO: функцию получения подписчиков пользователя

	err = h.dbR.SaveNewEntity(ctx, msg.EntityUUID, model.User)
	if err != nil {
		m.Increment("save_user_post.error")
		log.Printf("failed to create post: %v", err)
		return err
	}

	m.Increment("save_user_post.success")

	return nil
}
