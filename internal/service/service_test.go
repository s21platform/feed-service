package service

import (
	"context"
	"errors"
	"testing"

	feedproto "github.com/s21platform/feed-proto/feed-proto"
	"github.com/s21platform/feed-service/internal/config"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_CreateUserPosts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	t.Run("create_ok", func(t *testing.T) {
		mockRepo.EXPECT().CreateUserPost(ctx, uuid, gomock.Any()).Return(&feedproto.CreateUserPostOut{}, nil)

		s := New(mockRepo)
		_, err := s.CreateUserPost(ctx, &feedproto.CreateUserPostIn{})
		assert.NoError(t, err)
	})

	t.Run("create_no_uuid", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockDBRepo(ctrl)

		s := New(mockRepo)
		_, err := s.CreateUserPost(ctx, &feedproto.CreateUserPostIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to retrieve uuid")
	})

	t.Run("create_err", func(t *testing.T) {
		expectedErr := errors.New("get err")

		mockRepo.EXPECT().CreateUserPost(ctx, uuid, &feedproto.CreateUserPostIn{}).Return(&feedproto.CreateUserPostOut{}, expectedErr)

		s := New(mockRepo)
		_, err := s.CreateUserPost(ctx, &feedproto.CreateUserPostIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to create post: get err")
	})
}
