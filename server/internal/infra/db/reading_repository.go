package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/abstractions"
	"example.com/open-monastery/internal/app/apperr"
	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/enums"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingRepository struct {
	db *sql.DB
}

var _ abstractions.ReadingRepository = (*ReadingRepository)(nil)

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	return &ReadingRepository{
		db: db,
	}
}

func (r *ReadingRepository) Add(
	ctx context.Context,
	reading *models.Reading,
) error {
	query := `
		INSERT INTO reading (book_id, monk_id, customer_id, status, created_at, assigned_at, monastery_id)
			SELECT b.id, ?, ?, ?, ?, ?, b.monastery_id
			FROM book b
			WHERE b.id = ? AND b.monastery_id = ?
		`
	result, err := r.db.ExecContext(
		ctx,
		query,
		reading.MonkId,
		reading.CustomerId,
		reading.Status,
		reading.CreatedAt,
		reading.AssignedAt,
		reading.BookId,
		reading.MonasteryId,
	)

	if err != nil {
		return translateWriteError("AddReading", err)
	}
	added, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get number of readings added: %w", err)
	}
	if added == 0 {
		return apperr.NewInvalid("book does not exist or does not belong to the monastery", nil)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get inserted reading ID: %w", err)
	}
	reading.Id = id

	return nil
}

func (r *ReadingRepository) Claim(ctx context.Context, monasteryID, readingID, monkID int64) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE reading r
		JOIN users u ON u.id = ? AND u.monastery_id = r.monastery_id
		SET r.monk_id = u.id, r.status = ?, r.assigned_at = UTC_TIMESTAMP(6)
		WHERE r.id = ? AND r.monastery_id = ? AND r.monk_id IS NULL AND r.status = ?`,
		monkID, enums.Scheduled, readingID, monasteryID, enums.Pending,
	)
	if err != nil {
		return translateWriteError("ClaimReading", err)
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get number of claimed readings: %w", err)
	}
	if updated != 0 {
		return nil
	}

	var member bool
	err = r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND monastery_id = ?)`,
		monkID, monasteryID,
	).Scan(&member)
	if err != nil {
		return fmt.Errorf("check monk monastery: %w", err)
	}
	if !member {
		return apperr.NewForbidden("monastery.forbidden")
	}

	var assignee sql.NullInt64
	err = r.db.QueryRowContext(ctx,
		"SELECT monk_id FROM reading WHERE id = ? AND monastery_id = ?", readingID, monasteryID,
	).Scan(&assignee)
	if errors.Is(err, sql.ErrNoRows) {
		return apperr.NewNotFound("reading.not_found", nil)
	}
	if err != nil {
		return fmt.Errorf("get reading assignee: %w", err)
	}
	if assignee.Valid {
		return apperr.NewConflict("reading.already_claimed", nil)
	}
	return apperr.NewConflict("reading.not_claimable", nil)
}

func (r *ReadingRepository) Assign(ctx context.Context, monasteryID, readingID, monkID int64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE reading r
		JOIN users u ON u.id = ? AND u.monastery_id = r.monastery_id AND u.user_type = ?
		SET r.monk_id = u.id, r.status = ?, r.assigned_at = UTC_TIMESTAMP(6)
		WHERE r.id = ? AND r.monastery_id = ?
			AND r.status IN (?, ?)
			AND (r.status = ? OR r.monk_id IS NULL OR r.monk_id <> u.id)
	`, monkID, enums.Monk, enums.Scheduled, readingID, monasteryID,
		enums.Pending, enums.Scheduled, enums.Pending)
	if err != nil {
		return translateWriteError("assign reading", err)
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get number of assigned readings: %w", err)
	}
	if updated != 0 {
		return nil
	}

	var status enums.ReadingStatus
	var assignee sql.NullInt64
	var eligible bool
	err = r.db.QueryRowContext(ctx, `
		SELECT r.status, r.monk_id, u.id IS NOT NULL
		FROM reading r
		LEFT JOIN users u ON u.id = ? AND u.monastery_id = r.monastery_id AND u.user_type = ?
		WHERE r.id = ? AND r.monastery_id = ?
	`, monkID, enums.Monk, readingID, monasteryID).Scan(&status, &assignee, &eligible)
	if errors.Is(err, sql.ErrNoRows) {
		return apperr.NewNotFound("reading.not_found", nil)
	}
	if err != nil {
		return fmt.Errorf("check reading assignment: %w", err)
	}
	if !eligible {
		return apperr.NewNotFound("monk.not_found", nil)
	}
	if status == enums.Scheduled && assignee.Valid && assignee.Int64 == monkID {
		return nil
	}
	return apperr.NewConflict("reading.not_assignable", nil)
}

// GetReadingsByUserId implements [abstractions.ReadingRepository].
func (r *ReadingRepository) GetReadingsByUserId(
	ctx context.Context,
	userID int64,
	opts appmodels.ListOptions,
) ([]*models.Reading, error) {
	query := `
	SELECT id, book_id, monk_id, customer_id, status, created_at, assigned_at, monastery_id
	FROM reading
	WHERE customer_id = ? AND id > ?
	ORDER BY id
	LIMIT ?
	`
	args := []any{userID, opts.AfterId, opts.Limit + 1}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get readings by customer ID: %w", err)
	}
	defer rows.Close()
	readings := make([]*models.Reading, 0)
	for rows.Next() {
		var found models.Reading
		if err := rows.Scan(
			&found.Id,
			&found.BookId,
			&found.MonkId,
			&found.CustomerId,
			&found.Status,
			&found.CreatedAt,
			&found.AssignedAt,
			&found.MonasteryId,
		); err != nil {
			return nil, fmt.Errorf("scan reading: %w", err)
		}
		readings = append(readings, &found)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate readings by customer ID: %w", err)
	}

	return readings, nil
}
