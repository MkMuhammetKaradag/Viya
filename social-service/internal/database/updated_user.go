package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (r *Repository) UpdateUserSocialInfo(ctx context.Context, id uuid.UUID, isPrivate *bool, avatarURL *string) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	parts := []string{}
	argIdx := 1

	if isPrivate != nil {
		parts = append(parts, fmt.Sprintf("is_private = $%d", argIdx))
		args = append(args, *isPrivate)
		argIdx++
	}

	if avatarURL != nil {
		parts = append(parts, fmt.Sprintf("avatar_url = $%d", argIdx))
		args = append(args, *avatarURL)
		argIdx++
	}

	if len(parts) == 0 {
		return nil
	}

	query += strings.Join(parts, ", ") + fmt.Sprintf(" WHERE id = $%d", argIdx)

	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
