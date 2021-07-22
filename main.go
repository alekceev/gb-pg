package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := context.Background()

	// Строка для подключения к базе данных
	url := "postgres://www:pass@127.0.0.1:15434/todo"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}

	// Pool соединений обязательно ограничивать сверху,
	// так как иначе есть потенциальная опасность превысить лимит соединений с базой.
	cfg.MaxConns = 16
	cfg.MinConns = 16

	// HealthCheckPeriod - частота проверки работоспособности
	// соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute

	// MaxConnLifetime - сколько времени будет жить соединение.
	// Так как большого смысла удалять живые соединения нет,
	// можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour

	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	// если запросов не поступало, то соединение закроется.
	cfg.MaxConnIdleTime = 30 * time.Minute
	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second

	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	// pgx предоставляет набор адаптеров для популярных логеров
	// это позволяет организовать сбор ошибок при работе с базой
	// @see: https://github.com/jackc/pgx/tree/master/log
	// cfg.ConnConfig = zerologadapter.NewLogger(logger)

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	duration := time.Duration(10 * time.Second)
	threads := 1000
	fmt.Printf("start attack: \n\tMinConns=%d\n\tMaxConns=%d\n", cfg.MinConns, cfg.MaxConns)
	res := attack(ctx, duration, threads, dbpool)

	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)

}

type (
	Name  string
	Email string
)

type EmailSearchHint struct {
	Name     Name
	Email    Email
	ListsCnt int
}

// search ищет всех сотрудников, email которых начинается с prefix.
// Из функции возвращается список EmailSearchHint, отсортированный по Email.
// Размер возвращаемого списка ограничен значением limit.
func search(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int) ([]EmailSearchHint, error) {
	const sql = `
select
	u.email,
	u.name,
	count(ul.list_id) as cnt
from lists l
join users u on l.user_uuid = u.uuid
where u.email like $1
group by u.email, u.name
order by cnt desc
limit $2;
`

	pattern := prefix + "%"
	rows, err := dbpool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []EmailSearchHint

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint EmailSearchHint

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Email, &hint.Name, &hint.ListsCnt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		hints = append(hints, hint)
	}

	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}

	return hints, nil
}

type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

func attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			_, err := search(ctx, dbpool, "user", 5)
			if err != nil {
				log.Fatal(err)
			}

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}

	wg.Wait()

	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}
