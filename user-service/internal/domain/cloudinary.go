package domain

import (
	"context"
	"mime/multipart"
)

type UplodImageOptions struct {
	Folder         string
	Width          int
	Height         int
	PublicID       string
	Transformation string
}

type CloudinaryService interface {
	UplodImage(ctx context.Context, fileHeader *multipart.FileHeader, opts UplodImageOptions) (string, string, error)
	UploadAvatar(ctx context.Context, file multipart.File, userID string) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
}
