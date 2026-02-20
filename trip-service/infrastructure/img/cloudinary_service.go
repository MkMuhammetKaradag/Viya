package img

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
	"trip-service/internal/domain"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudName, apiKey, apiSecret string) (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary connection failed: %w", err)
	}
	return &CloudinaryService{client: cld}, nil
}

// Ortak yükleme mantığı - io.Reader kullanarak hem file hem byte desteği sağlarız
func (s *CloudinaryService) upload(ctx context.Context, reader io.Reader, opts domain.UploadOptions) (string, error) {
	// 1. Benzersiz PublicID oluşturma
	publicID := fmt.Sprintf("%s_%d", opts.WayPointID, time.Now().UnixNano())

	// 2. Dinamik Transformation
	if opts.Transformation == "" {
		width := opts.Width
		height := opts.Height
		if width == 0 {
			width = 800
		}
		if height == 0 {
			height = 600
		}
		opts.Transformation = fmt.Sprintf("c_fill,g_auto,w_%d,h_%d,q_auto,f_auto", width, height)
	}

	// 3. Yükleme
	uploadRes, err := s.client.Upload.Upload(ctx, reader, uploader.UploadParams{
		Folder:         opts.Folder,
		PublicID:       publicID,
		Overwrite:      api.Bool(true),
		Invalidate:     api.Bool(true),
		Transformation: opts.Transformation,
	})

	if err != nil {
		return "", fmt.Errorf("cloud upload failed: %w", err)
	}

	// 4. Cloudinary iç hata kontrolü
	if uploadRes.Error.Message != "" {
		return "", fmt.Errorf("cloudinary api error: %s", uploadRes.Error.Message)
	}

	return uploadRes.SecureURL, nil
}

// UploadImage artık sadece dosyayı açıp ortak upload'a paslıyor
func (s *CloudinaryService) UploadImage(ctx context.Context, data io.Reader, opts domain.UploadOptions) (string, error) {
	return s.upload(ctx, data, opts)
}

func (s *CloudinaryService) UploadImageFromBytes(ctx context.Context, data []byte, opts domain.UploadOptions) (string, error) {
	return s.upload(ctx, bytes.NewReader(data), opts)
}

func (s *CloudinaryService) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}
