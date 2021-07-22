// +build integration

package storage_test

import (
	"context"
	"os"
	"testing"

	"gb-pg/pkg/todo/storage"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestIntegrationSearch(t *testing.T) {
	ctx := context.Background()
	dbpool := connect(ctx)
	defer dbpool.Close()

	tests := []struct {
		name    string
		store   *storage.PG
		ctx     context.Context
		prefix  string
		limit   int
		prepare func(*pgxpool.Pool)
		check   func(*testing.T, []storage.EmailSearchHint, error)
	}{
		{
			name:   "success",
			store:  storage.NewPG(dbpool),
			ctx:    context.Background(),
			prefix: "user",
			limit:  5,
			prepare: func(dbpool *pgxpool.Pool) {
				// Подготовка тестовых данных
				dbpool.Exec(context.Background(), `insert into users (name, email, pass, salt) values
				('User1', 'user1@mail.ru', 'd36f9d30acb4e2857a2818aa8420f7b7', '111'),
				('Admin', 'admin@mail.ru', 'd36f9d30acb4e2857a2818aa8420f7b7', '111');
				insert into lists (user_uuid, title, description) values ((select uuid from users where email = 'user1@mail.ru'), 'Купить в магазине', '');
				`)
			},
			check: func(t *testing.T, hints []storage.EmailSearchHint, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, hints)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(dbpool)
			hints, err := tt.store.Search(tt.ctx, tt.prefix, tt.limit)
			tt.check(t, hints, err)
		})
	}
}

// Соединение с экземпляром Postgres
func connect(ctx context.Context) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return dbpool
}
