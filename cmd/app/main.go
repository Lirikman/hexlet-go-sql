package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	FisrttName string `json:"first_name"`
	LastName   string `json:"last_name"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := sql.Open("sqlite", "file:data.db?_foreign_keys=on&_busy_timeout=5000")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	const schema = `CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		first_name TEXT,
		last_name TEXT
	);`
	if _, err := db.ExecContext(ctx, schema); err != nil {
		log.Fatalf("create table: %v", err)
	}

	const insert = `INSERT INTO users(email, first_name, last_name) VALUES (?, ?, ?) ON CONFLICT(email) DO NOTHING;`
	if _, err := db.ExecContext(ctx, insert, "ivan_iv@mail.ru", "Ivan", "Petrov"); err != nil {
		log.Fatalf("insert user: %v", err)
	}
	if _, err := db.ExecContext(ctx, insert, "max20@ya.ru", "Maksim", "Rozov"); err != nil {
		log.Fatalf("insert user: %v", err)
	}

	var u User
	err = db.QueryRowContext(ctx,
		`SELECT id, email, first_name, last_name FROM users WHERE email = ?`,
		"ivan_iv@mail.ru",
	).Scan(&u.ID, &u.Email, &u.FisrttName, &u.LastName)
	if err != nil {
		log.Fatalf("select user: %v", err)
	}

	payload, _ := json.MarshalIndent(u, "", "  ")
	log.Printf("loaded user: %s", payload)
}
