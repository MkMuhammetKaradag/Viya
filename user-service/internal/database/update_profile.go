package database

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"user-service/internal/domain"

	"github.com/google/uuid"
)

func (r *Repository) UpdateProfile(ctx context.Context, userID uuid.UUID, params domain.UpdateProfileParams) error {
	// 1. Alanları ve değerlerini bir map'te toplayalım
	updates := map[string]interface{}{
		"first_name": params.FirstName,
		"last_name":  params.LastName,
		"bio":        params.Bio,
		"location":   params.Location,
		"website":    params.Website,
	}

	// Preferences özel bir işlem (JSON) gerektirdiği için onu ayrı tutabiliriz veya
	// direkt map'e ekleyip döngü içinde kontrol edebiliriz.
	if params.Preferences != nil {
		prefsJSON, _ := json.Marshal(params.Preferences)
		updates["preferences"] = prefsJSON
	}

	queryParts := []string{}
	args := []interface{}{}
	argCounter := 1

	// 2. Map üzerinde dönerek sadece nil olmayanları sorguya ekle
	for column, value := range updates {
		// Reflect kullanmak yerine basit nil kontrolü (interface{} olduğu için)
		if value != nil && !isNil(value) {
			queryParts = append(queryParts, fmt.Sprintf("%s = $%d", column, argCounter))
			args = append(args, value)
			argCounter++
		}
	}
	// if params.Preferences != nil {
	// 	prefsJSON, _ := json.Marshal(params.Preferences)
	// 	queryParts = append(queryParts, fmt.Sprintf("preferences = $%d", argCounter))
	// 	args = append(args, prefsJSON)
	// 	argCounter++
	// }

	if len(queryParts) == 0 {
		return nil // Güncellenecek bir şey yok
	}

	// 3. Sorguyu birleştir
	finalQuery := fmt.Sprintf(
		"UPDATE users SET %s, updated_at = NOW() WHERE id = $%d AND deleted_at IS NULL",
		strings.Join(queryParts, ", "),
		argCounter,
	)
	args = append(args, userID)

	_, err := r.db.ExecContext(ctx, finalQuery, args...)
	return err
}

// Helper: Interface içindeki pointer'ın nil olup olmadığını kontrol eder
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	// Sadece Pointer, Chan, Map, Slice gibi tipler IsNil() kontrolüne girebilir
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Map || v.Kind() == reflect.Slice {
		return v.IsNil()
	}
	return false
}
