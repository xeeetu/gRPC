package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=54321 dbname=note user=note-user password=note-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("Error connect to db: %v", err)
	}
	defer pool.Close()

	// Вставляем заметку в базу и возвращаем id
	builderInsert := sq.Insert("note").
		PlaceholderFormat(sq.Dollar).
		Columns("title", "body").
		Values(gofakeit.Name(), gofakeit.Email()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("Error building insert query: %v", err)
	}

	row := pool.QueryRow(ctx, query, args...)

	var idInsert int64
	err = row.Scan(&idInsert)

	if err != nil {
		log.Fatalf("Error scan id: %v", err)
	}

	fmt.Printf("Note inserted with id: %d ", idInsert)

	// Выбираем заметки из базы (максимум 2)
	builderSelect := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		OrderBy("id ASC").
		Limit(2)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("Error building select query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("Error executing select query: %v", err)
	}
	defer rows.Close()

	var id int64
	var title, body string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {

		err = rows.Scan(&id, &title, &body, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("Error scan row: %v", err)
		}

		fmt.Printf("Note { id: %d, title: %s, body: %s, createdAt: %v, updatedAt: %v }", id, title, body, createdAt, updatedAt)
	}

	// Обновляем заметку по id
	updateBuilder := sq.Update("note").
		PlaceholderFormat(sq.Dollar).
		Set("title", gofakeit.Name()).
		Set("body", gofakeit.Email()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": idInsert})

	query, args, err = updateBuilder.ToSql()
	if err != nil {
		log.Fatalf("Error building update query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("Error executing update query: %v", err)
	}

	fmt.Printf("Rows updated: %d", res.RowsAffected())

	// Выводим заметку по id
	builderSelectOne := sq.Select("id", "title", "body", "created_at", "updated_at").
		From("note").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": idInsert}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("Error building select one query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &title, &body, &createdAt, &updatedAt)

	fmt.Printf("Note { id: %d, title: %s, body: %s, createdAt: %v, updatedAt: %v }", id, title, body, createdAt, updatedAt)
}
