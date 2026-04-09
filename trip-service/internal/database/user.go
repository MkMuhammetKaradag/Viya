package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (r *Repository) UpdateUserInterest(ctx context.Context, userID uuid.UUID, tripID uuid.UUID, weight float32) error {
	var oldInterestStr, tripVecStr sql.NullString

	// 1. Verileri string (metin) olarak çekiyoruz (Hata almamak için)
	err := r.db.QueryRowContext(ctx, "SELECT interest_vector::text FROM users WHERE id = $1", userID).Scan(&oldInterestStr)
	err = r.db.QueryRowContext(ctx, "SELECT content_vector::text FROM trips WHERE id = $1", tripID).Scan(&tripVecStr)

	if err != nil || !tripVecStr.Valid {
		return fmt.Errorf("trip vector not found or null: %w", err)
	}

	// 2. Bu stringleri (ör: "[0.1,0.2]") Go'nun anlayacağı []float32'ye çeviriyoruz
	oldInterestVec := parseVectorString(oldInterestStr.String)
	tripVec := parseVectorString(tripVecStr.String)

	// 3. Go tarafında ağırlıklı ortalamayı hesapla
	newInterestVec := calculateWeightedAverage(oldInterestVec, tripVec, weight)

	// 4. Tekrar stringe çevir ve DB'ye yaz
	strForDB := formatVector(newInterestVec)
	_, err = r.db.ExecContext(ctx, "UPDATE users SET interest_vector = $1 WHERE id = $2", strForDB, userID)

	return err
}
func parseVectorString(str string) []float32 {
	if str == "" || str == "NULL" {
		return nil
	}
	// Köşeli parantezleri temizle: [0.1,0.2] -> 0.1,0.2
	str = strings.Trim(str, "[]")
	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	vec := make([]float32, len(parts))
	for i, p := range parts {
		val, _ := strconv.ParseFloat(strings.TrimSpace(p), 32)
		vec[i] = float32(val)
	}
	return vec
}
func calculateWeightedAverage(oldVec, newVec []float32, weight float32) []float32 {
	if len(oldVec) == 0 {
		return newVec
	}

	result := make([]float32, len(oldVec))
	for i := range oldVec {
		// Formül: (Eski * 0.95) + (Yeni * 0.05)
		result[i] = (oldVec[i] * (1.0 - weight)) + (newVec[i] * weight)
	}
	return result
}
func formatVector(v []float32) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range v {
		// Sayıları metne çeviriyoruz (f formatı: ondalıklı sayı)
		sb.WriteString(fmt.Sprintf("%f", val))
		if i < len(v)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")
	return sb.String()
}
