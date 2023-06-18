package income

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"
)

// repository persists albums in database
type MemoryRepository struct {
	logger     *zap.SugaredLogger
	storage    map[int]Income
	current_id int
	mu         *sync.Mutex
}

var global MemoryRepository

// NewRepository creates a new album repository
func NewMemoryRepository(logger *zap.SugaredLogger) MemoryRepository {
	global = MemoryRepository{
		storage:    make(map[int]Income),
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
func (r MemoryRepository) Get(ctx context.Context, id int) (income Income, err error) {
	global.mu.Lock()
	income, ok := global.storage[id] //global.db.With(ctx).Select().Model(id, &Income)
	global.mu.Unlock()

	if !ok {
		return Income{}, errors.New("no such Income")
	}
	// row := global.db.DB().DB().QueryRow("SELECT scope_id, author_id, CreatedAt, UpdatedAt FROM $1 where ID = $2;", Income.TableName(), id)
	// err := row.Scan(&Income.Scope, &privilege.Username, &privilege.Status, &privilege.Balance)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return &privilege, err
	// 	}
	// }
	return income, err
}

// Create saves a new album record in the database.
// It returns the ID of the newly inserted album record.
func (r MemoryRepository) Create(ctx context.Context, income *Income) error {
	global.mu.Lock()
	income.ID = global.current_id
	global.storage[global.current_id] = *income
	global.current_id++
	global.mu.Unlock()
	return nil
}

// Update saves the changes to an album in the database.
func (r MemoryRepository) Update(ctx context.Context, income Income) error {
	global.mu.Lock()
	global.storage[income.ID] = income
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
	// err := global.db.With(ctx).Select("COUNT(*)").From("incomes").Row(&count)
	return count, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r MemoryRepository) Query(ctx context.Context, offset, limit int) (incomes []Income, err error) {
	// err = global.db.With(ctx).
	// 	Select().
	// 	OrderBy("id").
	// 	Offset(int64(offset)).
	// 	Limit(int64(limit)).
	// 	All(&incomes)
	incomes = make([]Income, 0)
	for _, n := range global.storage {
		incomes = append(incomes, n)
	}

	return incomes, nil
}
