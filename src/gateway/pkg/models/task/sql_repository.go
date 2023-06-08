package task

import (
	"context"

	"tasks/pkg/dbcontext"

	"go.uber.org/zap"
)

// repository persists albums in database
type SqlRepository struct {
	db     *dbcontext.DB
	logger *zap.SugaredLogger
}

// NewRepository creates a new album repository
func NewSQLRepository(db *dbcontext.DB, logger *zap.SugaredLogger) SqlRepository {
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
func (r SqlRepository) Get(ctx context.Context, id int) (task Task, err error) {
	err = r.db.With(ctx).Select().Model(id, &task)

	// row := r.db.DB().DB().QueryRow("SELECT scope_id, author_id, CreatedAt, UpdatedAt FROM $1 where ID = $2;", task.TableName(), id)
	// err := row.Scan(&task.Scope, &privilege.Username, &privilege.Status, &privilege.Balance)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return &privilege, err
	// 	}
	// }

	return task, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r SqlRepository) Create(ctx context.Context, task *Task) error {
	return r.db.With(ctx).Model(task).Insert()
}

// Update saves the changes to an album in the database.
func (r SqlRepository) Update(ctx context.Context, task Task) error {
	return r.db.With(ctx).Model(&task).Update()
}

// Delete deletes an album with the specified ID from the database.
func (r SqlRepository) Delete(ctx context.Context, id int) error {
	task, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&task).Delete()
}

// Count returns the number of the album records in the database.
func (r SqlRepository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("tasks").Row(&count)
	return count, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r SqlRepository) Query(ctx context.Context, offset, limit int) (tasks []Task, err error) {
	err = r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&tasks)
	return tasks, err
}
