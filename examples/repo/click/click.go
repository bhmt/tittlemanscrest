package click

import (
	"context"
	"log"
	"time"

	"github.com/bhmt/tittlemanscrest/repository"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type Item struct {
	id    int
	ts    int64
	state string
}

const (
	migrate = `
CREATE TABLE IF NOT EXISTS item_history (
	item_id UInt64,
	ts UInt64,
	state String
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(toDateTime(ts))
ORDER BY (item_id, ts)
`

	seed = `
INSERT INTO item_history (item_id, ts, state)
VALUES (?, ?, ?)
`

	query = `
SELECT
	item_id,
	max(ts) as latest,
	argMax(state, ts)
FROM item_history
GROUP BY item_id
ORDER BY latest DESC
`
)

var (
	now    = time.Now().Unix()
	minute = time.Now().Add(time.Minute).Unix()
	ten    = time.Now().Add(10 * time.Minute).Unix()
)

func write(session *repository.Session, ctx context.Context) {
	uow, err := repository.NewUnitOfWork(session, repository.WithTransaction(session))
	if err != nil {
		log.Fatal(err)
	}

	defer uow.Rollback()

	_, err = uow.Worker.ExecContext(ctx, migrate)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("table created")

	for _, item := range []Item{
		{1, now, "new"},
		{1, minute, "modified"},
		{2, now, "new"},
		{2, minute, "deleted"},
		{3, ten, "new"},
	} {
		if _, err := uow.Worker.ExecContext(ctx, seed, item.id, item.ts, item.state); err != nil {
			log.Fatal(err)
		}
	}

	if err := uow.Commit(); err != nil {
		log.Fatal(err)
	}

	log.Println("item inserted")
}

func read(session *repository.Session, ctx context.Context) {
	uow, err := repository.NewUnitOfWork(session)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := uow.Worker.QueryContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.id, &i.ts, &i.state); err != nil {
			log.Fatal(err)
		}
		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, i := range items {
		log.Printf("id: %d ts: %s state: %s\n", i.id, time.Unix(i.ts, 0), i.state)
	}
}

func Click() {
	session, err := repository.NewSession("clickhouse", "clickhouse://clickhouse:clickhouse@localhost:9000/clickhouse")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	write(session, ctx)
	read(session, ctx)
}
