package service

import (
	"github.com/s21platform/feed-service/pkg/feed"
)

type Service struct {
	feed.UnimplementedFeedServiceServer
	dbR DBRepo
}

func New(dbR DBRepo) *Service {
	return &Service{dbR: dbR}
}
