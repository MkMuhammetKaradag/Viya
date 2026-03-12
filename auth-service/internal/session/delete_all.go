package session

import (
	"context"
	"fmt"
	"log"
)

func (sr *SessionRepository) DeleteAllSession(ctx context.Context, userID string) error {
	// 1. Kullanıcının tüm session ID'lerini SET'ten çek
	userKey := sr.userSessionsKey(userID)
	sessionIDs, err := sr.client.SMembers(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	// Eğer hiç session yoksa zaten işlem tamamdır
	if len(sessionIDs) == 0 {
		return nil
	}

	// 2. İşlemleri toplu (Pipeline) yapalım ki performans artsın
	pipe := sr.client.TxPipeline()

	// Her bir session verisini (string olanları) tek tek sil
	for _, id := range sessionIDs {
		pipe.Del(ctx, "session:"+id)
	}

	// 3. Kullanıcının tüm listesini (SET) tek seferde sil
	pipe.Del(ctx, userKey)

	// Hepsini tek seferde Redis'e gönder
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute delete all sessions pipeline: %w", err)
	}

	log.Printf("🚫 User [%s] için tüm %d oturum sonlandırıldı.", userID, len(sessionIDs))
	return nil
}
