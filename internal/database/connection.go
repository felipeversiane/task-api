package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Connection *pgxpool.Pool

func Connect(ctx context.Context) error {
	dsn := getConnectionString()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}

	Connection, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	if Connection == nil {
		return
	}
	Connection.Close()
}

func getConnectionString() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	dbport := os.Getenv("POSTGRES_PORT")
	dbhost := os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", user, password, dbname, dbport, dbhost)

	return dsn
}
