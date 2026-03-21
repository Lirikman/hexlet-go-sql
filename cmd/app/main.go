package main

import (
	"context"
	"database/sql"
	"encoding/json"
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

	const schema = `CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50),
		email VARCHAR(50) UNIQUE,
		age INTEGER,
		created_at TIMESTAMP DEFAULT(datetime('now', 'localtime'))
	);`

	if _, err := db.ExecContext(ctx, schema); err != nil {
		log.Fatalf("create table: %v", err)
	}

	cmd := &cli.Command{

		Name:  "Example of working with a database on GO",
		Usage: "database SQLite; supports -a (adding user), -l (list of users), -f (get user on id)",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a user to the database",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
					&cli.StringFlag{Name: "email", Aliases: []string{"e"}},
					&cli.Int64Flag{Name: "age", Aliases: []string{"a"}},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := sql.NullString{String: cmd.String("name"), Valid: cmd.String("name") != ""}
					email := sql.NullString{String: cmd.String("email"), Valid: cmd.String("email") != ""}
					age := sql.NullInt64{Int64: cmd.Int64("age"), Valid: cmd.Int64("age") > 0}

					dto := pkg.CreateUserDTO{
						Name:  name,
						Email: email,
						Age:   age,
					}

					res, err := pkg.CreateUser(ctx, db, dto)
					if err != nil {
						return err
					}

					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					fmt.Println(encoder.Encode(res))
					return nil
				},
			},
			{
				Name:  "update",
				Usage: "updating user data",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "id"},
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
					&cli.StringFlag{Name: "email", Aliases: []string{"e"}},
					&cli.Int64Flag{Name: "age", Aliases: []string{"a"}},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Int("id")
					name := sql.NullString{String: cmd.String("name"), Valid: cmd.String("name") != ""}
					email := sql.NullString{String: cmd.String("email"), Valid: cmd.String("email") != ""}
					age := sql.NullInt64{Int64: cmd.Int64("age"), Valid: cmd.Int64("age") > 0}

					dto := pkg.UpdateUserDTO{
						ID:    id,
						Name:  name,
						Email: email,
						Age:   age,
					}

					res, err := pkg.UpdateUser(ctx, db, dto)
					if err != nil {
						return err
					}

					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					fmt.Println(encoder.Encode(res))
					return nil
				},
			},
			{
				Name:  "get",
				Usage: "getting user by id",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "id"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Int("id")
					res, err := pkg.GetUser(ctx, db, id)

					if err != nil {
						return err
					}

					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					fmt.Println(encoder.Encode(res))
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "get a list of all users",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					res, err := pkg.ListUsers(ctx, db)

					if err != nil {
						return err
					}

					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					fmt.Println(encoder.Encode(res))
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "deleting a user by id",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "id"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					id := cmd.Int("id")

					res, err := pkg.DeleteUser(ctx, db, id)

					if err != nil {
						return err
					}

					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					fmt.Println(encoder.Encode(res))
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
