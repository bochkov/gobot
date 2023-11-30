package db

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool interface {
	getPool() *pgxpool.Pool
	setPool(pool *pgxpool.Pool)
}

type connectionPool struct {
	pool *pgxpool.Pool
	sync.RWMutex
}

func (cp *connectionPool) getPool() *pgxpool.Pool {
	cp.RLock()
	defer cp.RUnlock()
	return cp.pool
}

func (cp *connectionPool) setPool(pool *pgxpool.Pool) {
	cp.Lock()
	defer cp.Unlock()
	cp.pool = pool
}

var dbPool *connectionPool
var once sync.Once

func NewPool(ctx context.Context, url string) *pgxpool.Pool {
	once.Do(func() {
		dbPool = &connectionPool{}
		pool, err := pgxpool.New(ctx, url)
		if err != nil {
			slog.Warn(err.Error())
			os.Exit(1)
		}
		dbPool.setPool(pool)
	})
	return dbPool.getPool()
}
