package grpc_client

import (
	"context"
	"fmt"
	"log"
	"time"

	// Kendi oluşturduğunuz proto paketini import edin
	pb "api-gateway/internal/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// İstemciyi uygulamada global olarak erişilebilir tutmak için.
var AuthServiceClient pb.AuthServiceClient

// User Servisine olan bağlantıyı temsil eder.
var conn *grpc.ClientConn

// Gateway uygulaması başlangıcında çağrılacak fonksiyon
func InitAuthClient(grpcAddress string) error {
	var err error

	// Güvenliksiz bağlantı (Genellikle internal mikroservisler için kabul edilebilir)
	conn, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	// Bağlantı üzerinden gRPC istemcisini oluşturun
	AuthServiceClient = pb.NewAuthServiceClient(conn)
	log.Printf("✅ Gateway, User Servisine gRPC ile bağlandı: %s", grpcAddress)
	return nil
}

// Uygulama kapanırken bağlantıyı kapatmak için
func CloseAuthClient() {
	if conn != nil {
		conn.Close()
	}
}

// AuthMiddleware'in çağıracağı ana doğrulama fonksiyonu
func ValidateAndRotateSession(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	// 3 saniyelik timeout ile bir context oluşturun
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	fmt.Println("istek atıldı buaya ")
	// User Servisindeki gRPC metodunu çağırın
	resp, err := AuthServiceClient.ValidateAndRotateSession(ctx, req)

	if err != nil {
		log.Printf("🔒 gRPC doğrulama çağrısı başarısız: %v", err)
		return nil, err
	}

	// Geri dönen cevabı kontrol edin
	return resp, nil
}
