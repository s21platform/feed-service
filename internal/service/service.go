package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userproto "github.com/s21platform/user-service/pkg/user"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/pkg/feed"
)

type Service struct {
	feed.UnimplementedFeedServiceServer
	dbR        DBRepo
	userClient UserClient
}

func New(dbR DBRepo, userClient UserClient) *Service {
	return &Service{
		dbR:        dbR,
		userClient: userClient,
	}
}

func (s *Service) GetFeed(ctx context.Context, _ *feed.GetFeedIn) (*feed.GetFeedOut, error) {
	userUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok || userUUID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get user UUID from context")
	}

	subscriptions, err := s.userClient.GetPeerFollow(ctx, userUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		return &feed.GetFeedOut{Items: []*feed.FeedItem{}}, nil
	}

	entityPosts, err := s.dbR.GetPostsByUserSubscriptions(ctx, userUUID, 50)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts from database: %v", err)
	}

	if len(entityPosts) == 0 {
		return &feed.GetFeedOut{Items: []*feed.FeedItem{}}, nil
	}

	postUUIDs := make([]string, 0, len(entityPosts))
	for _, entityPost := range entityPosts {
		if entityPost.Metadata == "user" {
			postUUIDs = append(postUUIDs, entityPost.ExternalUUID)
		}
	}

	if len(postUUIDs) == 0 {
		return &feed.GetFeedOut{Items: []*feed.FeedItem{}}, nil
	}

	postsInfo, err := s.userClient.GetPostsByIds(ctx, userUUID, postUUIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts info from user-service: %v", err)
	}

	postsMap := make(map[string]*userproto.PostInfo)
	for _, postInfo := range postsInfo {
		postsMap[postInfo.PostUuid] = postInfo
	}

	feedItems := make([]*feed.FeedItem, 0, len(entityPosts))
	for _, entityPost := range entityPosts {
		if entityPost.Metadata == "user" {
			postInfo, ok := postsMap[entityPost.ExternalUUID]
			if !ok {
				continue
			}

			feedItem := &feed.FeedItem{
				Post: &feed.FeedItem_UserPost{
					UserPost: &feed.UserPost{
						PostUuid:   postInfo.PostUuid,
						Nickname:   postInfo.Nickname,
						FullName:   postInfo.FullName,
						AvatarLink: postInfo.AvatarLink,
						Content:    postInfo.Content,
						CreatedAt:  postInfo.CreatedAt,
						IsEdited:   postInfo.IsEdited,
					},
				},
			}
			feedItems = append(feedItems, feedItem)
		}
	}

	return &feed.GetFeedOut{Items: feedItems}, nil
}
