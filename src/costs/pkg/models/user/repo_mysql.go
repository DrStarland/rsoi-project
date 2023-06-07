package user

// import (
// 	crypto "crypto/rand"
// 	"database/sql"
// 	"log"
// 	"math/big"
// )

// type MySQLRepository struct {
// 	DB *sql.DB
// }

// func NewMySQLRepo(db *sql.DB) *MySQLRepository {
// 	return &MySQLRepository{DB: db}
// }

// // func PasswordArgon2(plainPassword []byte) string {
// // 	return string(argon2.IDKey(plainPassword, []byte("salt"), 1, 64*1024, 4, 16))
// // }

// var GenerateID = func(crutch string) string {
// 	if crutch != "" {
// 		return crutch
// 	}

// 	safeNum, err := crypto.Int(crypto.Reader, big.NewInt(100234))
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return safeNum.String()
// }

// func (repo *MySQLRepository) GetByID(id string) (*User, error) {
// 	user := &User{}
// 	// QueryRow сам закрывает коннект
// 	err := repo.DB.
// 		QueryRow("SELECT ID, Username, Password, updated FROM Users WHERE id = ?", id).
// 		Scan(&user.ID, &user.Username, &user.password, &user.updated)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (repo *MySQLRepository) GetByUsername(username string) (*User, bool) {
// 	user := &User{}
// 	// QueryRow сам закрывает коннект
// 	err := repo.DB.
// 		QueryRow("SELECT ID, Username, Password, updated FROM Users WHERE username = ?", username).
// 		Scan(&user.ID, &user.Username, &user.password, &user.updated)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return nil, false
// 	}
// 	return user, true
// }

// func (repo *MySQLRepository) Add(user *User) (*User, error) {
// 	user.ID = GenerateID(user.ID)

// 	_, err := repo.DB.Exec(
// 		"INSERT INTO Users (`ID`, `Username`, `Password`) VALUES (?, ?, ?)",
// 		user.ID,
// 		user.Username,
// 		user.password,
// 	)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (repo *MySQLRepository) Register(login, pass string) (*User, error) {
// 	user := &User{Username: login, password: pass}
// 	return repo.Add(user)
// }

// func (repo *MySQLRepository) Authorize(login, pass string) (*User, error) {
// 	u, ok := repo.GetByUsername(login)

// 	if !ok {
// 		return nil, ErrNoUser
// 	}

// 	// dont do this un production :)
// 	// так точно, но сейчас мы не продуцируем...
// 	if u.password != pass {
// 		return nil, ErrBadPass
// 	}

// 	return u, nil
// }
