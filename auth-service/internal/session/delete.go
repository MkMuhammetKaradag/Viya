package session

import (
	"context"
	"fmt"
	"log"
)

func (sr *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	// 1. Önce session verisini çekelim ki UserID'yi bulalım (Set'ten silmek için)
	session, err := sr.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if session == nil {
		// Zaten yoksa hata dönmeye gerek yok, "başarılı" sayılabilir (Idempotent)
		return nil
	}

	// Redis'teki gerçek anahtar adını oluşturduğundan emin ol
	// Eğer GetSession içinde zaten ekliyorsan burayı ona göre ayarla
	fullSessionKey := "session:" + sessionID

	pipe := sr.client.TxPipeline()

	// 2. Ana session verisini sil
	pipe.Del(ctx, fullSessionKey)

	// 3. Kullanıcının aktif session listesinden (Set) bu ID'yi çıkar
	pipe.SRem(ctx, sr.userSessionsKey(session.UserID), sessionID)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	log.Printf("🗑️ Session başarıyla silindi: %s", sessionID)
	return nil
}
