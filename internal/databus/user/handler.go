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
	dbR        DBRepo
	userClient UserClient
}

func New(dbR DBRepo, userClient UserClient) *Handler {
	return &Handler{
		dbR:        dbR,
		userClient: userClient,
	}
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

	postUUID, err := h.dbR.SaveNewEntity(ctx, msg.EntityUUID, model.User)
	if err != nil {
		m.Increment("save_user_post.error")
		log.Printf("failed to create post: %v", err)
		return err
	}

	resp, err := h.userClient.GetWhoFollowPeer(ctx, msg.UserUUID)
	if err != nil {
		m.Increment("save_user_post.error")
		log.Printf("failed to get followers: %v", err)
		return err
	}

	for _, peer := range resp {
		err = h.dbR.SaveNewEntitySuggestion(ctx, postUUID, peer.Uuid)
		if err != nil {
			m.Increment("save_user_post.error")
			log.Printf("failed to create suggestion: %v", err)
			return err
		}
	}

	m.Increment("save_user_post.success")

	return nil
}
