package user

import "golang.org/x/net/context"

type Handler struct {
	dbR DBRepo
}

func New(dbR DBRepo) *Handler {
	return &Handler{dbR: dbR}
}

func (h *Handler)Handler(ctx context.Context, in []byte) error {
	
}