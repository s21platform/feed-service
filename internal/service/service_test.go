package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/pkg/feed"
)

func TestServer_CreateUserPosts(t *testing.T) {
	t.Parallel()

	content := "test-content"
	userUUID := uuid.New().String()
	expUUID := uuid.New().String()

	ctx := context.Background()
	ctx = context.WithValue(ctx, config.KeyUUID, userUUID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)
	s := New(mockRepo)

	t.Run("create_ok", func(t *testing.T) {
		mockRepo.EXPECT().Post(ctx, userUUID, content).Return(expUUID, nil)

		_, err := s.CreateUserPost(ctx, &feed.CreateUserPostIn{Content: content})

		assert.NoError(t, err)
	})

	t.Run("create_no_uuid", func(t *testing.T) {
		ctx := context.Background()

		_, err := s.CreateUserPost(ctx, &feed.CreateUserPostIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to retrieve uuid")
	})

	t.Run("create_err", func(t *testing.T) {
		expectedErr := errors.New("get err")

		mockRepo.EXPECT().Post(ctx, userUUID, content).Return("", expectedErr)

		_, err := s.CreateUserPost(ctx, &feed.CreateUserPostIn{Content: content})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to create post: get err")
	})
}
