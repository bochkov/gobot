package db

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
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

func NewPool(ctx context.Context, url string) ConnectionPool {
	once.Do(func() {
		dbPool = &connectionPool{}
		pool, err := pgxpool.New(ctx, url)
		if err != nil {
			slog.Warn(err.Error())
			os.Exit(1)
		}
		dbPool.setPool(pool)
	})
	return dbPool
}

func GetPool() *pgxpool.Pool {
	return dbPool.getPool()
}
