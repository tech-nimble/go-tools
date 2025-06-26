package pg

import (
	"context"
	"fmt"

	"github.com/gobuffalo/envy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Initialize() (*pgxpool.Pool, error) {
	DBHost := envy.Get("DB_HOST", "localhost")
	DBPort := envy.Get("DB_PORT", "5432")
	DBName := envy.Get("DB_NAME", "users")
	DBPass := envy.Get("DB_PASS", "")
	DBUser := envy.Get("DB_USER", "root")

	DBUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPass, DBHost, DBPort, DBName)

	return pgxpool.New(context.Background(), DBUrl)
}
