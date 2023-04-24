package note

import (
	"context"

	"notes/pkg/dbcontext"

	"go.uber.org/zap"
)

// repository persists albums in database
type SqlRepository struct {
	db     *dbcontext.DB
	logger *zap.SugaredLogger
}

// NewRepository creates a new album repository
func NewRepository(db *dbcontext.DB, logger *zap.SugaredLogger) SqlRepository {
	return SqlRepository{db, logger}
}

// Get reads the album with the specified ID from the database.
func (r SqlRepository) Get(ctx context.Context, id string) (Note, error) {
	var note Note
	err := r.db.With(ctx).Select().Model(id, &note)
	return note, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r SqlRepository) Create(ctx context.Context, note Note) error {
	return r.db.With(ctx).Model(&note).Insert()
}

// Update saves the changes to an album in the database.
func (r SqlRepository) Update(ctx context.Context, note Note) error {
	return r.db.With(ctx).Model(&note).Update()
}

// Delete deletes an album with the specified ID from the database.
func (r SqlRepository) Delete(ctx context.Context, id string) error {
	album, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&album).Delete()
}

// Count returns the number of the album records in the database.
func (r SqlRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("album").Row(&count)
	return count, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r SqlRepository) Query(ctx context.Context, offset, limit int) ([]Note, error) {
	var albums []Note
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&albums)
	return albums, err
}
