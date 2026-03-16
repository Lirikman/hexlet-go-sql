package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Course struct {
	ID    int
	Slug  string
	Title string
	Price int
}

func CreateCourse(ctx context.Context, db *sql.DB, c Course) (int64, error) {
	const query = `
        INSERT INTO courses(slug, title, price)
        VALUES(?, ?, ?)
        RETURNING id`

	var id int64
	if err := db.QueryRowContext(ctx, query, c.Slug, c.Title, c.Price).Scan(&id); err != nil {
		return 0, fmt.Errorf("create course: %w", err)
	}
	return id, nil
}

var allowedOrder = map[string]string{
	"price_asc":  "price ASC",
	"price_desc": "price DESC",
	"title_asc":  "title ASC",
}

func ListCourses(ctx context.Context, db *sql.DB, limit, offset int, order string) ([]Course, error) {
	ord, ok := allowedOrder[order]
	if !ok {
		ord = "id ASC"
	}

	query := `
        SELECT id, slug, title, price 
        FROM courses
        ORDER BY ` + ord + ` LIMIT ? OFFSET ?`

	rows, err := db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list courses: %w", err)
	}
	defer rows.Close()

	var out []Course

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Price); err != nil {
			return nil, fmt.Errorf("scan course: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func FindCoursesByIDs(ctx context.Context, db *sql.DB, listIds []int64) ([]Course, error) {
	if len(listIds) == 0 {
		return []Course{}, nil
	}

	placeholders := strings.Repeat("?,", len(listIds))
	placeholders = placeholders[:len(placeholders)-1]

	args := make([]interface{}, len(listIds))
	for i, v := range listIds {
		args[i] = v
	}

	query := `SELECT id, slug, title, price FROM courses WHERE id IN (` + placeholders + `)`

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find courses: %w", err)
	}
	defer rows.Close()

	var out []Course

	for rows.Next() {
		var c Course
		if err := rows.Scan(&c.ID, &c.Slug, &c.Title, &c.Price); err != nil {
			return nil, fmt.Errorf("scan course: %w", err)
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
