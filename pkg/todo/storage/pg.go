package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PG struct {
	dbpool *pgxpool.Pool
}

func NewPG(dbpool *pgxpool.Pool) *PG {
	return &PG{dbpool}
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
func (s *PG) Search(ctx context.Context, prefix string, limit int) ([]EmailSearchHint, error) {
	const sql = `
select
	u.email,
	u.name,
	count(ul.list_id) as cnt
from users_lists ul
join users u on ul.user_uuid = u.uuid
where u.email like $1
group by u.email, u.name
order by cnt desc
limit $2;
`

	pattern := prefix + "%"
	rows, err := s.dbpool.Query(ctx, sql, pattern, limit)
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
