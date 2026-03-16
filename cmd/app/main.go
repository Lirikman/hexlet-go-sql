package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	pkg "example.com/go-sql/internal/storage"
	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

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

	const schema = `CREATE TABLE IF NOT EXISTS courses(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slug TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		price INTEGER NOT NULL DEFAULT 0
	);`

	if _, err := db.ExecContext(ctx, schema); err != nil {
		log.Fatalf("create table: %v", err)
	}

	cmd := &cli.Command{
		Name:  "Example of working with a database on GO",
		Usage: "database SQLite; supports -a (adding a record), -l (list of courses), -f (course on id)",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "create_course",
				Aliases: []string{"a"},
				Value:   false,
				Usage:   "adding a new course to the database",
			},
			&cli.BoolFlag{
				Name:    "list_courses",
				Aliases: []string{"l"},
				Value:   false,
				Usage:   "getting a list of all courses",
			},
			&cli.BoolFlag{
				Name:    "course_on_id",
				Aliases: []string{"f"},
				Value:   false,
				Usage:   "getting a course on id",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// n := cmd.NArg()
			// if n == 0 {
			// 	return fmt.Errorf("no arguments entered, enter flag -h for help")
			// }
			if cmd.Bool("create_course") {
				res, err := pkg.CreateCourse(ctx, db, pkg.Course{Slug: "Backend", Title: "PHP", Price: 60000})
				if err != nil {
					return fmt.Errorf("Error: %w", err)
				}
				fmt.Printf("operation was completed successfully - last id: %v", res)
			}
			if cmd.Bool("list_courses") {
				res, err := pkg.ListCourses(ctx, db, 10, 0, "title ASC")
				if err != nil {
					return fmt.Errorf("error: %w", err)
				}
				fmt.Printf("List of all courses:\n")
				for _, rec := range res {
					fmt.Printf("%v. %v %v - %v RUB.\n", rec.ID, rec.Slug, rec.Title, rec.Price)
				}
			}
			if cmd.Bool("course_on_id") {
				res, err := pkg.FindCoursesByIDs(ctx, db, []int64{1, 2, 4})
				if err != nil {
					return fmt.Errorf("error: %w", err)
				}
				for _, rec := range res {
					fmt.Printf("%v %v %v %v\n", rec.ID, rec.Slug, rec.Title, rec.Price)
				}
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
