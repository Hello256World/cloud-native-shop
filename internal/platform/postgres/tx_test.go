//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Hello256World/cloud-native-shop/internal/platform/postgres"
)

func TestTxManager(t *testing.T) {
	ctx := context.Background()
	pool, err := postgres.Connect(ctx, postgres.Options{
		URI: os.Getenv("TEST_DATABASE_URL"), MinPoolSize: 1, MaxPoolSize: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	tm := postgres.NewTxManager(pool)
	count := func() int {
		var n int
		_ = pool.QueryRow(ctx, `SELECT count(*) FROM platform.outbox_events WHERE aggregate_type = 'tx_test'`).Scan(&n)
		return n
	}
	insert := func(ctx context.Context) error {
		_, err := postgres.Executor(ctx, pool).Exec(ctx,
			`INSERT INTO platform.outbox_events (aggregate_type, aggregate_id, event_type, payload)
			 VALUES ('tx_test', '1', 'test', '{}')`)
		return err
	}
	t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM platform.outbox_events WHERE aggregate_type = 'tx_test'`) })

	boom := errors.New("boom")
	err = tm.WithinTx(ctx, func(ctx context.Context) error {
		if err := insert(ctx); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) || count() != 0 {
		t.Fatalf("rollback failed: err=%v count=%d", err, count())
	}

	if err := tm.WithinTx(ctx, insert); err != nil || count() != 1 {
		t.Fatalf("commit failed: err=%v count=%d", err, count())
	}
}