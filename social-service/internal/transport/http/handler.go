package http

import "social-service/internal/domain"

type Handlers struct {
}

func NewHandlers(repo domain.SocialRepository) *Handlers {
	return &Handlers{}
}
