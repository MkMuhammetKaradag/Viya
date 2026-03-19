package img

import (
	"context"
	"fmt"
	"mime/multipart"
	"user-service/internal/domain"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	client *cloudinary.Cloudinary
}

func NewCloudinary(cloudName, apiKey, apiSecret string) (domain.CloudinaryService, error) {
	c, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &Cloudinary{
		client: c,
	}, nil
}
func (cld *Cloudinary) UploadAvatar(ctx context.Context, file multipart.File, userID string) (string, error) {

	uploadRes, err := cld.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "profile_pictures",
		PublicID:       userID,
		Overwrite:      api.Bool(true),
		Invalidate:     api.Bool(true),
		Transformation: "c_fill,g_face,h_500,w_500,q_auto,f_auto",
	})

	if err != nil {
		return "", fmt.Errorf("cloudinary upload: %w", err)
	}
	return uploadRes.SecureURL, nil
}
func (cld *Cloudinary) UplodImage(ctx context.Context, fileHeader *multipart.FileHeader, opts domain.UplodImageOptions) (string, string, error) {
	return "", "", nil
}
func (cld *Cloudinary) DeleteImage(ctx context.Context, publicID string) error {
	_, err := cld.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
func (cld *Cloudinary) Close() error {
	// Kapatılacak bir bağlantı yok
	return nil
}
