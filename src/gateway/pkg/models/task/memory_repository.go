package task

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
)

// repository persists albums in database
type MemoryRepository struct {
	logger     *zap.SugaredLogger
	storage    map[int]Task
	current_id int
	mu         *sync.Mutex
}

var global MemoryRepository

// NewRepository creates a new album repository
func NewMemoryRepository(logger *zap.SugaredLogger) MemoryRepository {
	global = MemoryRepository{
		storage:    make(map[int]Task),
		logger:     logger,
		current_id: 1,
		mu:         &sync.Mutex{},
	}
	return global
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
func (r MemoryRepository) Get(ctx context.Context, id int) (task Task, err error) {
	global.mu.Lock()
	task, ok := global.storage[id] //global.db.With(ctx).Select().Model(id, &task)
	global.mu.Unlock()

	if !ok {
		return Task{}, errors.New("no such Task")
	}
	// row := global.db.DB().DB().QueryRow("SELECT scope_id, author_id, CreatedAt, UpdatedAt FROM $1 where ID = $2;", task.TableName(), id)
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
func (r MemoryRepository) Create(ctx context.Context, task *Task) error {
	global.mu.Lock()
	task.ID = global.current_id
	task.mu = &sync.Mutex{}
	task.Comments = make([]Comment, 0)
	global.storage[global.current_id] = *task
	global.current_id++
	global.mu.Unlock()
	return nil
}

// Update saves the changes to an album in the database.
func (r MemoryRepository) Update(ctx context.Context, task Task) error {
	global.mu.Lock()
	global.storage[task.ID] = task
	global.mu.Unlock()
	return nil
}

// Delete deletes an album with the specified ID from the database.
func (r MemoryRepository) Delete(ctx context.Context, id int) error {
	global.mu.Lock()
	delete(global.storage, id)
	global.mu.Unlock()
	return nil
}

// Count returns the number of the album records in the database.
func (r MemoryRepository) Count(ctx context.Context) (int, error) {
	var count int
	// err := global.db.With(ctx).Select("COUNT(*)").From("tasks").Row(&count)
	return count, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r MemoryRepository) Query(ctx context.Context, offset, limit int) (tasks []Task, err error) {
	// err = global.db.With(ctx).
	// 	Select().
	// 	OrderBy("id").
	// 	Offset(int64(offset)).
	// 	Limit(int64(limit)).
	// 	All(&tasks)
	tasks = make([]Task, 0)
	for _, n := range global.storage {
		tasks = append(tasks, n)
	}

	return tasks, nil
}

func (r MemoryRepository) AddComment(ctx context.Context, id int, comment Comment) error {
	global.mu.Lock()
	task := global.storage[id]
	// r.logger.Infoln("!!!!!!")
	task.Comments = append(task.Comments, comment)
	global.storage[id] = task
	global.mu.Unlock()
	return nil
}
