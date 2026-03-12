package grpc_transport

import (
	"auth-service/internal/domain"
	"auth-service/internal/pb"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

// AuthGrpcHandler'ı app/application.go'dan buraya taşıyabilirsiniz veya
// sadece ValidateToken metodunu bu pakette uygulayabilirsiniz.
type AuthGrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	SessionRepo domain.SessionRepository
}

func NewAuthGrpcHandler(repo domain.SessionRepository) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		SessionRepo: repo,
	}
}

// GrpcServerRegistrar arayüzüne uyan Register metodu (Application.go'da kalabilir veya buraya taşınır)
func (h *AuthGrpcHandler) Register(gRPCServer *grpc.Server) {
	pb.RegisterAuthServiceServer(gRPCServer, h)
}

const SessionDuration = 7 * 24 * time.Hour
const MaxSessionLifetime = 52 * SessionDuration

// *** ASIL DOĞRULAMA MANTIĞI BURASI ***
func (h *AuthGrpcHandler) ValidateAndRotateSession(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	session, err := h.SessionRepo.GetSession(ctx, req.SessionId)
	if err != nil || session == nil {
		return &pb.ValidateResponse{Valid: false}, nil
	}

	fmt.Println("validate isteği")
	// 2. Güvenlik Kontrolü: IP veya UserAgent değişmiş mi? (Opsiyonel ama önerilir)
	// if session.Ip != req.Ip { return &pb.ValidateResponse{Valid: false}, nil }

	now := time.Now()

	// 3. Maksimum ömür kontrolü (365 gün geçtiyse oturumu kapat)
	if now.Sub(session.FirstCreatedAt) > MaxSessionLifetime {
		fmt.Println("1 yıllık süre  bitmiş")
		_ = h.SessionRepo.DeleteSession(ctx, req.SessionId)
		return &pb.ValidateResponse{Valid: false}, nil
	}

	res := &pb.ValidateResponse{
		Valid:  true,
		UserId: session.UserID,
	}

	// 4. ROTASYON KONTROLÜ (24 saat geçtiyse yeni bir SessionID üret)
	if now.Sub(session.CreatedAt) > SessionDuration {
		fmt.Println("7*24 saatlik süre  bitmiş")
		// Yeni bir SessionID üret ve eskisine Grace Period (1 dk) vererek yenisini kaydet
		newSessionID, err := h.SessionRepo.Rotate(ctx, req.SessionId, session)
		if err == nil {
			res.NewSessionId = newSessionID
		}
	}

	return res, nil
}
