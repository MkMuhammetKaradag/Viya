package database

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) UpdateBanner(ctx context.Context, userID uuid.UUID, banerURL string) error {
	query := `UPDATE users SET banner_url = $1 WHERE id = $2` //updated_at = NOW()
	_, err := r.db.ExecContext(ctx, query, banerURL, userID)
	return err
}
