package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(driverName, connStr string) (*sql.DB, error) {

	db, err := sql.Open(driverName, connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the PostgreSQL database!")

	return db, nil
}

// func CreateTables(db *sql.DB) {
// 	createBooksTable := `
// 	CREATE TABLE IF NOT EXISTS books (
// 		id SERIAL PRIMARY KEY NOT NULL,
// 		title VARCHAR(255) NOT NULL,
// 		author VARCHAR(100) NOT NULL,
// 		isbn VARCHAR(20) UNIQUE,
// 		publication_year INT,
//     	available_copies INT
// 	);
// 	`
// 	_, err := db.Exec(createBooksTable)
// 	if err != nil {
// 		log.Println(err)
// 		panic("Cannot create books table!")
// 	}
// 	log.Println("Table books was successfuly created!")

// 	createUsersTable := `
// 	CREATE TABLE IF NOT EXISTS users (
// 		id SERIAL PRIMARY KEY NOT NULL,
// 		name VARCHAR(255) NOT NULL,
// 		email VARCHAR(100) NOT NULL UNIQUE,
// 		password VARCHAR(100) NOT NULL UNIQUE
// 	);
// 	`
// 	_, err = db.Exec(createUsersTable)
// 	if err != nil {
// 		panic("Cannot create users table!")
// 	}
// 	log.Println("Table users was successfuly created!")

// 	createReservationsTable := `
// 	CREATE TABLE IF NOT EXISTS reservations (
// 		id SERIAL PRIMARY KEY NOT NULL,
// 		book_id INT NOT NULL,
// 		user_id INT NOT NULL,
// 		checkout_date TIMESTAMP NOT NULL,
// 		return_date TIMESTAMP,
// 		FOREIGN KEY (book_id) REFERENCES books(id),
//     	FOREIGN KEY (user_id) REFERENCES users(id)
// 	);
// 	`
// 	_, err = db.Exec(createReservationsTable)
// 	if err != nil {
// 		panic("Cannot create reservations table!")
// 	}
// 	log.Println("Table reservations was successfuly created!")
// }
