package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
)

const (
	dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаём соединение с базой
	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func() {
		errCon := conn.Close(ctx)
		if errCon != nil {
			log.Fatalf("Unable to close connection: %v\n", errCon)
		}
	}()

	res, err := conn.Exec(ctx, "INSERT INTO note (title, body) VALUES ($1, $2)", gofakeit.Name(), gofakeit.Email())
	if err != nil {
		log.Fatalf("Unable to insert note: %v\n", err)
	}

	log.Printf("Inserted %d rows", res.RowsAffected())

	// Делаем запрос на выборку записей из таблицы note
	rows, err := conn.Query(ctx, "SELECT id, title, body, created_at, updated_at FROM note")
	if err != nil {
		log.Fatalf("Unable to get notes: %v\n", err)
	}

	for rows.Next() {
		var id int
		var title, body string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("Unable to scan row: %v\n", err)
		}

		log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", id, title, body, createdAt, updatedAt)
	}

}
