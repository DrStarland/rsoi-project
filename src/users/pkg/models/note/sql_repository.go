package note

import (
	"context"

	"users/pkg/dbcontext"

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

// func (repo *PostgresRepository) GetPrivilegeByUsername(username string) (*Privilege, error) {
// 	var privilege Privilege

// 	log.Printf(">>>> username: %s", username)
// 	row := repo.DB.QueryRow("SELECT * FROM privilege where username = $1;", username)
// 	err := row.Scan(&privilege.ID, &privilege.Username, &privilege.Status, &privilege.Balance)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return &privilege, err
// 		}
// 	}

// 	return &privilege, nil
// }

// func (repo *PostgresRepository) GetHistoryById(privilegeID string) ([]*PrivilegeHistory, error) {
// 	var history []*PrivilegeHistory
// 	rows, err := repo.DB.Query("SELECT * FROM privilege_history where privilege_id = $1;", privilegeID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to execute the query: %w", err)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("failed to execute the query: %s", err)
// 	}

// 	for rows.Next() {
// 		row := new(PrivilegeHistory)
// 		rows.Scan(
// 			&row.ID,
// 			&row.PrivilegeID,
// 			&row.TicketUID,
// 			&row.Date,
// 			&row.BalanceDiff,
// 			&row.OperationType,
// 		)

// 		if err != nil {
// 			return nil, fmt.Errorf("failed to execute the query: %s", err)
// 		}

// 		history = append(history, row)
// 	}

// 	return history, nil
// }

// ID           SERIAL PRIMARY KEY,
// scope_id    int REFERENCES scope(ID),
// -- связь с тегами через доп таблицу
// author_id   int,
// title           VARCHAR(50),
// content       VARCHAR(2000),
// CreatedAt     TIMESTAMP default current_timestamp,
// UpdatedAt   TIMESTAMP default current_timestamp

// Get reads the album with the specified ID from the database.
func (r SqlRepository) Get(ctx context.Context, id int) (note Note, err error) {
	err = r.db.With(ctx).Select().Model(id, &note)

	// row := r.db.DB().DB().QueryRow("SELECT scope_id, author_id, CreatedAt, UpdatedAt FROM $1 where ID = $2;", note.TableName(), id)
	// err := row.Scan(&note.Scope, &privilege.Username, &privilege.Status, &privilege.Balance)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return &privilege, err
	// 	}
	// }

	return note, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r SqlRepository) Create(ctx context.Context, note *Note) error {
	return r.db.With(ctx).Model(note).Insert()
}

// Update saves the changes to an album in the database.
func (r SqlRepository) Update(ctx context.Context, note Note) error {
	return r.db.With(ctx).Model(&note).Update()
}

// Delete deletes an album with the specified ID from the database.
func (r SqlRepository) Delete(ctx context.Context, id int) error {
	note, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&note).Delete()
}

// Count returns the number of the album records in the database.
func (r SqlRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("notes").Row(&count)
	return count, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r SqlRepository) Query(ctx context.Context, offset, limit int) (notes []Note, err error) {
	err = r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&notes)
	return notes, err
}
