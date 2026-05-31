package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"notebook/internal/handler"
	"notebook/internal/repository"
)

func main() {
	addr := env("APP_ADDR", ":8080")
	dsn := env("DATABASE_URL", "postgres://notebook:notebook@db:5432/notebook?sslmode=disable")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	if err := waitDB(ctx, pool); err != nil {
		log.Fatalf("db not ready: %v", err)
	}

	repo := repository.New(pool)
	h := handler.New(repo)

	mux := http.NewServeMux()
	h.Register(mux)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func waitDB(ctx context.Context, pool *pgxpool.Pool) error {
	for {
		if err := pool.Ping(ctx); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
