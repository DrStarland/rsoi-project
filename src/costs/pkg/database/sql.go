package database

import (
	"fmt"
	"time"

	// _ "github.com/lib/pq"

	// _ "github.com/go-sql-driver/postgresql"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

// func main() {
// 	db, _ := dbx.Open("mysql", "user:pass@/example")

// 	// create a new query
// 	// q := db.NewQuery("SELECT id, name FROM users LIMIT 10")

// 	// // fetch all rows into a struct array
// 	// var users []struct {
// 	// 	ID, Name string
// 	// }
// 	// err := q.All(&users)

// 	// // fetch a single row into a struct
// 	// var user struct {
// 	// 	ID, Name string
// 	// }
// 	// err = q.One(&user)

// 	// // fetch a single row into a string map
// 	// data := dbx.NullStringMap{}
// 	// err = q.One(data)

// 	// // fetch row by row
// 	// rows2, _ := q.Rows()
// 	// for rows2.Next() {
// 	// 	_ = rows2.ScanStruct(&user)
// 	// 	// rows.ScanMap(data)
// 	// 	// rows.Scan(&id, &name)
// 	// }
// }

func CreateConnection() (*dbx.DB, error) {
	port := 5432
	host := "0.0.0.0" //"postgres"

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, "postgres", "notes", "postgres")

	db, err := dbx.MustOpen("postgres", dsn)
	if err != nil {
		time.Sleep(10 * time.Second)
		db, err = dbx.MustOpen("postgres", dsn)
		if err != nil {
			return nil, err
		}
	}

	err = db.DB().Ping()
	// db.DB().Query("SELECT * ")

	// db.Select()

	return db, err
}
