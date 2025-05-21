package user

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/context"

	logger_lib "github.com/s21platform/logger-lib"
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
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("Handler")

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
		logger.Error(fmt.Sprintf("failed to create post: %v", err))
		return err
	}

	followers, err := h.userClient.GetWhoFollowPeer(ctx, msg.UserUUID)
	if err != nil {
		m.Increment("save_user_post.error")
		logger.Error(fmt.Sprintf("failed to get followers: %v", err))
		return err
	}

	for _, follower := range followers {
		err = h.dbR.SaveNewEntitySuggestion(ctx, postUUID, follower.Uuid)
		if err != nil {
			m.Increment("save_user_post.error")
			logger.Error(fmt.Sprintf("failed to create suggestion: %v", err))
			return err
		}
	}

	m.Increment("save_user_post.success")

	return nil
}
