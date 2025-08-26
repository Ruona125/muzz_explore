package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides methods to record and query user decisions.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore constructs a new Store using the given pgx connection pool.
func NewStore(ctx context.Context, pool *pgxpool.Pool) (*Store, error) {
	s := &Store{pool: pool}
	if err := s.migrate(ctx); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return s, nil
}

// migrate creates the decisions table if it does not already exist.
func (s *Store) migrate(ctx context.Context) error {
	const createTable = `
CREATE TABLE IF NOT EXISTS decisions (
    actor_user_id     TEXT    NOT NULL,
    recipient_user_id TEXT    NOT NULL,
    liked_recipient   BOOLEAN NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (actor_user_id, recipient_user_id)
);

CREATE INDEX IF NOT EXISTS idx_decisions_recipient ON decisions (recipient_user_id);
CREATE INDEX IF NOT EXISTS idx_decisions_updated_at ON decisions (updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_decisions_actor_recipient_liked ON decisions (recipient_user_id, liked_recipient);
    `
	_, err := s.pool.Exec(ctx, createTable)
	return err
}

// PutDecision stores or updates a decision.  If liked is true the
// actor has liked the recipient; if false the actor has passed.  The
// call returns a boolean indicating whether the like is now mutual.
func (s *Store) PutDecision(ctx context.Context, actorID, recipientID string, liked bool) (bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		// If the transaction is still open, roll it back.
		_ = tx.Rollback(ctx)
	}()
	// Upsert the decision.  updated_at is set to NOW() on each write.
	const upsert = `
INSERT INTO decisions (actor_user_id, recipient_user_id, liked_recipient, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (actor_user_id, recipient_user_id)
DO UPDATE SET liked_recipient = EXCLUDED.liked_recipient, updated_at = EXCLUDED.updated_at;
    `
	if _, err := tx.Exec(ctx, upsert, actorID, recipientID, liked); err != nil {
		return false, err
	}
	var mutual bool
	if liked {
		// Check if the recipient has already liked the actor.
		const query = `
SELECT liked_recipient
FROM decisions
WHERE actor_user_id = $1 AND recipient_user_id = $2;
        `
		var likedBack bool
		err := tx.QueryRow(ctx, query, recipientID, actorID).Scan(&likedBack)
		if err == nil && likedBack {
			mutual = true
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}
	return mutual, nil
}

// Liker represents a like from an actor to a recipient.  Unix
// holds the seconds since the Unix epoch when the decision was last
// updated.
type Liker struct {
	ActorID string
	Unix    uint64
}

// ListLikedYou returns all actors who have liked the recipient.  The
// results are paginated using offset and limit.
func (s *Store) ListLikedYou(ctx context.Context, recipientID string, offset, limit int) ([]Liker, *string, error) {
	if limit <= 0 {
		return nil, nil, errors.New("limit must be positive")
	}
	const query = `
SELECT actor_user_id, extract(epoch from updated_at)::bigint
FROM decisions
WHERE recipient_user_id = $1 AND liked_recipient = TRUE
ORDER BY updated_at DESC, actor_user_id ASC
OFFSET $2 LIMIT $3;
    `
	rows, err := s.pool.Query(ctx, query, recipientID, offset, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	likers := make([]Liker, 0)
	for rows.Next() {
		var l Liker
		var ts int64
		if err := rows.Scan(&l.ActorID, &ts); err != nil {
			return nil, nil, err
		}
		l.Unix = uint64(ts)
		likers = append(likers, l)
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}
	var nextToken *string
	if len(likers) == limit {
		// Compute next offset token.
		nextOffset := offset + limit
		nt := fmt.Sprintf("%d", nextOffset)
		nextToken = &nt
	}
	return likers, nextToken, nil
}

// ListNewLikedYou returns likes where the recipient hasn't yet liked
// back.  This excludes mutual likes from the result set.  Offset
// based pagination is used for simplicity.
func (s *Store) ListNewLikedYou(ctx context.Context, recipientID string, offset, limit int) ([]Liker, *string, error) {
	if limit <= 0 {
		return nil, nil, errors.New("limit must be positive")
	}
	const query = `
SELECT d.actor_user_id, extract(epoch from d.updated_at)::bigint
FROM decisions d
LEFT JOIN decisions r ON r.actor_user_id = d.recipient_user_id AND r.recipient_user_id = d.actor_user_id AND r.liked_recipient = TRUE
WHERE d.recipient_user_id = $1 AND d.liked_recipient = TRUE AND r.actor_user_id IS NULL
ORDER BY d.updated_at DESC, d.actor_user_id ASC
OFFSET $2 LIMIT $3;
    `
	rows, err := s.pool.Query(ctx, query, recipientID, offset, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	likers := make([]Liker, 0)
	for rows.Next() {
		var l Liker
		var ts int64
		if err := rows.Scan(&l.ActorID, &ts); err != nil {
			return nil, nil, err
		}
		l.Unix = uint64(ts)
		likers = append(likers, l)
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}
	var nextToken *string
	if len(likers) == limit {
		nextOffset := offset + limit
		nt := fmt.Sprintf("%d", nextOffset)
		nextToken = &nt
	}
	return likers, nextToken, nil
}

// CountLikedYou returns the number of actors who like the recipient.
func (s *Store) CountLikedYou(ctx context.Context, recipientID string) (uint64, error) {
	const query = `
SELECT COUNT(*)::bigint
FROM decisions
WHERE recipient_user_id = $1 AND liked_recipient = TRUE;
    `
	var count uint64
	err := s.pool.QueryRow(ctx, query, recipientID).Scan(&count)
	return count, err
}
