package user

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const lovelove = "lovelove"

// go test -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

//nolint:all
func TestNew(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &MySQLRepository{
		DB: db,
	}

	funcRepo := NewMySQLRepo(db)

	if !reflect.DeepEqual(funcRepo, repo) {
		t.Errorf("results not match, want %v, have %v", repo, funcRepo)
		return
	}

}

//nolint:all
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var userID = GenerateID("")

	// good query
	rows := sqlmock.NewRows([]string{"ID", "Username", "Password", "updated"})
	expect := []*User{
		{userID, "rvasily", lovelove, sql.NullString{}},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Username, item.password, nil)
	}

	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := &MySQLRepository{
		DB: db,
	}
	item, err := repo.GetByID(userID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(userID).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.GetByID(userID)
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// row scan error
	rows = sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "title")

	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(userID).
		WillReturnRows(rows)

	_, err = repo.GetByID(userID)
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

//nolint:all
func TestGetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	var username = "rvasily"

	// good query
	rows := sqlmock.NewRows([]string{"ID", "Username", "Password", "updated"})
	expect := []*User{
		{"1", username, lovelove, sql.NullString{}},
	}
	for _, item := range expect {
		rows = rows.AddRow(item.ID, item.Username, item.password, nil)
	}

	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(username).
		WillReturnRows(rows)

	repo := &MySQLRepository{
		DB: db,
	}
	item, flag := repo.GetByUsername(username)
	if flag != true {
		t.Errorf("unexpected err: ")
		return
	}
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(username).
		WillReturnError(fmt.Errorf("db_error"))

	_, flag = repo.GetByUsername(username)
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if flag == true {
		t.Errorf("expected error (false), got true")
		return
	}

	rows = sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "title")

	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(username).
		WillReturnRows(rows)

	_, flag = repo.GetByUsername(username)
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if flag == true {
		t.Errorf("expected error (false), got true")
		return
	}
}

//nolint:all
func TestAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &MySQLRepository{
		DB: db,
	}

	ID := GenerateID("1")
	Username := "rvasily"
	password := lovelove
	testItem := &User{
		ID:       ID,
		Username: Username,
		password: password,
	}

	// ok query
	mock.
		ExpectExec("INSERT INTO Users").
		WithArgs(ID, Username, password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	newUser, err := repo.Add(testItem)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if newUser.ID != ID {
		t.Errorf("bad id: want %v, have %v", "1", newUser.ID)
		return
	}

	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
	}

	// query error
	mock.
		ExpectExec("INSERT INTO Users").
		WithArgs(ID, Username, password).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.Add(testItem)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
	}
}

//nolint:all
func TestMySQLRepository_Register(t *testing.T) {
	GenerateID = func(crutch string) string {
		if crutch != "" {
			return crutch
		}

		return "1"
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &MySQLRepository{
		DB: db,
	}

	Username := "rvasily"
	password := lovelove

	// ok query
	mock.
		ExpectExec("INSERT INTO Users").
		WithArgs("1", Username, password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	newUser, err := repo.Register(Username, password)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if newUser.Username != Username {
		t.Errorf("bad id: want %v, have %v", Username, newUser.Username)
		return
	}
}

//nolint:all
func TestAuthorize(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	// good query
	correctRow := sqlmock.NewRows([]string{"ID", "Username", "Password", "updated"})
	badRow := sqlmock.NewRows([]string{"ID", "Username", "Password", "updated"})

	var username = "rvasily"
	password := lovelove
	expect := []*User{
		{"1", username, lovelove, sql.NullString{}},
	}
	unexpect := []*User{
		{"2", username, "hatehate", sql.NullString{}},
	}
	for _, item := range expect {
		correctRow = correctRow.AddRow(item.ID, item.Username, item.password, nil)
	}
	for _, item := range unexpect {
		badRow = badRow.AddRow(item.ID, item.Username, item.password, nil)
	}

	repo := &MySQLRepository{
		DB: db,
	}

	// correct authorization
	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(username).
		WillReturnRows(correctRow)

	item, err := repo.Authorize(username, password)
	if err != nil {
		t.Errorf("unexpected err: %v", err)
		return
	}
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if !reflect.DeepEqual(item, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], item)
		return
	}

	// query error
	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(username).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = repo.Authorize(username, password)
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// bad password
	mock.
		ExpectQuery("SELECT ID, Username, Password, updated FROM Users WHERE").
		WithArgs(username).
		WillReturnRows(badRow)

	_, err = repo.Authorize(username, password)
	if shadowErr := mock.ExpectationsWereMet(); shadowErr != nil {
		t.Errorf("there were unfulfilled expectations: %s", shadowErr)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
