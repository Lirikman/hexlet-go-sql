package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      *string   `json:"name"`
	Email     *string   `json:"email"`
	Age       *int64    `json:"age"`
	CreatedAt time.Time `json:"created at"`
}

type CreateUserDTO struct {
	ID        int
	Name      sql.NullString
	Email     sql.NullString
	Age       sql.NullInt64
	CreatedAt time.Time
}

type UpdateUserDTO struct {
	ID        int
	Name      sql.NullString
	Email     sql.NullString
	Age       sql.NullInt64
	CreatedAt time.Time
}

func CreateUser(ctx context.Context, db *sql.DB, dto CreateUserDTO) (User, error) {
	const query = `
        INSERT INTO users(name, email, age)
        VALUES(?, ?, ?)
        RETURNING id, name, email, age, created_at
		`
	var out User
	if err := db.QueryRowContext(ctx, query, dto.Name, dto.Email, dto.Age).Scan(&out.ID, &out.Name, &out.Email, &out.Age, &out.CreatedAt); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return out, nil
}

func UpdateUser(ctx context.Context, db *sql.DB, dto UpdateUserDTO) (User, error) {
	const query = `
		UPDATE users SET name = COALESCE(?, name), age = COALESCE(?, age), email = COALESCE(?, email)
		WHERE id = ?
		RETURNING id, name, email, age, created_at
		`
	var out User
	if err := db.QueryRowContext(ctx, query, dto.Name, dto.Age, dto.Email, dto.ID).Scan(&out.ID, &out.Name, &out.Email, &out.Age, &out.CreatedAt); err != nil {
		return User{}, fmt.Errorf("update user: %w", err)
	}
	return out, nil
}

func GetUser(ctx context.Context, db *sql.DB, id int) (User, error) {
	const query = `
		SELECT id, name, email, age, created_at FROM users 
		WHERE id = ?
		`
	var u User
	if err := db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.CreatedAt); err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func ListUsers(ctx context.Context, db *sql.DB) ([]User, error) {
	const query = `
		SELECT id, name, email, age, created_at FROM users
		ORDER BY id
		`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(ctx context.Context, db *sql.DB, id int) (string, error) {
	const query = `
		DELETE FROM users 
		WHERE id = ?
		`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return "", fmt.Errorf("delete user: %w", err)
	}
	msg := fmt.Sprintf("user with the email - %d has been successfully deleted", id)
	return msg, nil
}
