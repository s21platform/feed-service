package api

import (
	"encoding/json"
	"net/http"
	"time"

	logger_lib "github.com/s21platform/logger-lib"

	api "github.com/s21platform/feed-service/internal/generated"
	"github.com/s21platform/feed-service/internal/service"
	feedproto "github.com/s21platform/feed-service/pkg/feed"
)

type Handler struct {
	feedService *service.Service
}

func New(feedService *service.Service) *Handler {
	return &Handler{
		feedService: feedService,
	}
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request, params api.GetFeedParams) {
	w.Header().Set("Content-Type", "application/json")
	ctx := logger_lib.WithField(r.Context(), "user_uuid", params.XUserUuid)

	// Вызываем gRPC метод GetFeed
	feedOut, err := h.feedService.GetFeed(ctx, &feedproto.GetFeedIn{})
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to get feed")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	// Конвертируем gRPC ответ в REST формат
	items := make([]api.FeedItem, 0, len(feedOut.Items))
	for _, item := range feedOut.Items {
		apiItem := api.FeedItem{}

		switch post := item.Post.(type) {
		case *feedproto.FeedItem_UserPost:
			var createdAt time.Time
			if post.UserPost.CreatedAt != nil {
				createdAt = post.UserPost.CreatedAt.AsTime()
			}
			fullName := post.UserPost.FullName
			avatarLink := post.UserPost.AvatarLink
			isEdited := post.UserPost.IsEdited
			apiItem.UserPost = &api.UserPost{
				PostUuid:   post.UserPost.PostUuid,
				Nickname:   post.UserPost.Nickname,
				FullName:   &fullName,
				AvatarLink: &avatarLink,
				Content:    post.UserPost.Content,
				CreatedAt:  createdAt,
				IsEdited:   &isEdited,
			}
			items = append(items, apiItem)
		}
	}

	response := api.GetFeedResponse{
		Items: items,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger_lib.Error(logger_lib.WithField(ctx, "error", err.Error()), "failed to marshal data")
		resolveError(&w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func resolveError(w *http.ResponseWriter, status int) {
	var message string
	switch status {
	case http.StatusBadRequest:
		message = "Произошла ошибка, попробуйте перезагрузить страницу"
	case http.StatusUnauthorized:
		message = "Вы не авторизованы для этого действия"
	case http.StatusInternalServerError:
		message = "У нас что-то сломалось, но мы уже чиним!"
	default:
		message = "У нас что-то сломалось, но мы уже чиним!"
	}

	body, _ := json.Marshal(api.Error{
		Message: message,
	})
	(*w).WriteHeader(status)
	_, _ = (*w).Write(body)
}
