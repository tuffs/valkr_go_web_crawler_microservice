package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Init(dbURL string) error {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	return pool.Ping(context.Background())
}

func Pool() *pgxpool.Pool {
	return pool
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

// Upsert webpage (ON CONFLICT DO NOTHING)
func UpsertWebpage(ctx context.Context, p Page) error {
	query := `
		INSERT INTO webpages (title, url, description, thumbnail_url, verified, ai_result)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (url) DO NOTHING
	`

	_, err := Pool().Exec(ctx, query,
		p.Title,
		p.URL,
		p.Description,
		nullString(p.ThumbnailURL),
		false,
		false,
	)
	return err
}

func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// Insert into crawl_queues
func InsertCrawlQueue(ctx context.Context, topic string, urls []string) error {
	query := `
		INSERT INTO crawl_queues (topic, urls, status)
		VALUES ($1, $2, 'completed')
	`

	_, err := Pool().Exec(ctx, query, topic, urls)
	return err
}