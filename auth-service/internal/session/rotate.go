package session

import (
	"auth-service/internal/domain"
	"context"
	"encoding/json"
	"time"
)

func (sr *SessionRepository) Rotate(ctx context.Context, oldSessionID string, session *domain.SessionData) (string, error) {
	// 1. Yeni bir UUID oluştur
	newSessionID := generateSessionID()

	// Eski session verilerini kopyala ve güncelleme zamanını yenile
	session.CreatedAt = time.Now()

	jsonData, _ := json.Marshal(session)

	// Redis Pipeline kullanarak atomik işlem yapalım
	pipe := sr.client.TxPipeline()

	// 1. Yeni session'ı kaydet (265 günlük genel ömür)
	pipe.Set(ctx, "session:"+newSessionID, jsonData, 365*24*time.Hour)

	// 2. Kullanıcının session listesine yenisini ekle
	pipe.SAdd(ctx, sr.userSessionsKey(session.UserID), newSessionID)

	// 3. ESKİ SESSION'A GRACE PERIOD VER (Önemli!)
	// Anında silmiyoruz, 1 dakika daha tanıyalım ki havada kalan istekler patlamasın
	pipe.Expire(ctx, "session:"+oldSessionID, 1*time.Minute)
	pipe.SRem(ctx, sr.userSessionsKey(session.UserID), oldSessionID)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", err
	}

	return newSessionID, nil
}
