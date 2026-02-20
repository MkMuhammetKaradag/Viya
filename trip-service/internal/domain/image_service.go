package domain

import (
	"context"
	"io"
)

type UploadOptions struct {
	Folder         string
	Width          int
	Height         int
	WayPointID     string
	Transformation string
}
type ImageService interface {
	UploadImage(ctx context.Context, reader io.Reader, opts UploadOptions) (string, error)
	UploadImageFromBytes(ctx context.Context, data []byte, opts UploadOptions) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
}
