package metadata

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MensahPrince/fileserver/types"
)

type MySQLMetadata struct { // ← declared here, in metadata
	db *sql.DB
}

var ErrNotFound = errors.New("metadata: not found")

func NewMySQLMetadata(db *sql.DB) *MySQLMetadata {
	return &MySQLMetadata{db: db}
}

func (m *MySQLMetadata) Save(ctx context.Context, meta types.FileMeta) error {
	_, err := m.db.ExecContext(ctx, " INSERT INTO metadata (id, name, hash, created_at, size,  owner) VALUES (?,?,?,?,?,?)",
		meta.ID,
		meta.Name,
		meta.Hash,
		meta.CreatedAt,
		meta.Size,
		meta.Owner,
	)
	if err != nil {
		return fmt.Errorf("metadata: save failed: %w", err)
	}
	return nil
}

func (m *MySQLMetadata) Get(ctx context.Context, id string) (types.FileMeta, error) {
	var meta types.FileMeta
	err := m.db.QueryRowContext(ctx,
		"SELECT id, name, hash, created_at, size, owner FROM metadata WHERE id = ?",
		id,
	).Scan(&meta.ID, &meta.Name, &meta.Hash, &meta.CreatedAt, &meta.Size, &meta.Owner)

	if errors.Is(err, sql.ErrNoRows) {
		return types.FileMeta{}, ErrNotFound
	}
	if err != nil {
		return types.FileMeta{}, fmt.Errorf("metadata: get failed: %w", err)
	}
	return meta, nil
}

func (m *MySQLMetadata) Delete(ctx context.Context, id string) error {
	result, err := m.db.ExecContext(ctx, "DELETE FROM metadata WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("metadata: delete failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("metadata: could not determine rows affected: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (m *MySQLMetadata) List(ctx context.Context) ([]types.FileMeta, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT id, name, hash, created_at, size, owner
		FROM metadata
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []types.FileMeta

	for rows.Next() {
		var meta types.FileMeta

		if err := rows.Scan(
			&meta.ID,
			&meta.Name,
			&meta.Hash,
			&meta.CreatedAt,
			&meta.Size,
			&meta.Owner,
		); err != nil {
			return nil, err
		}

		results = append(results, meta)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
