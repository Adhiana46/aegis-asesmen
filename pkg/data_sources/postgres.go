package data_sources

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	*sqlx.DB
}

type logEntry struct {
	Query           string `json:"query"`
	Args            any    `json:"args"`
	ExecutionTimeMs int64  `json:"execution_time_ms"`
	Error           error  `json:"error,omitempty"`
}

var (
	postgresDBInstance PostgresDB
	postgresDBOnce     sync.Once
)

func NewPostgresDb(host, port, user, pass, dbname string) *PostgresDB {
	postgresDBOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host,
			user,
			pass,
			dbname,
			port,
		)

		conn, err := sqlx.Connect("pgx", dsn)
		if err != nil {
			slog.Error(err.Error(), slog.String("dsn", dsn))
			panic(0)
		}

		conn.SetMaxOpenConns(60)
		conn.SetConnMaxLifetime(120 * time.Second)
		conn.SetMaxIdleConns(30)
		conn.SetConnMaxIdleTime(20 * time.Second)

		// try to ping
		if err := conn.Ping(); err != nil {
			slog.Error(err.Error())
			panic(0)
		}

		slog.Info(fmt.Sprintf(
			"User '%s' successfully connected to PostgreSQL database '%s'@'%s'",
			user,
			dbname,
			host,
		))

		postgresDBInstance = PostgresDB{
			DB: conn,
		}
	})

	return &postgresDBInstance
}

func (r *PostgresDB) Close() error {
	return r.DB.Close()
}

func (r *PostgresDB) logQuery(entry *logEntry) {
	slog.Debug("Executing query", slog.Any("details", entry))
}

func (r *PostgresDB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	startTime := time.Now()

	// Execute query
	rows, err := r.DB.QueryxContext(ctx, query, args...)

	// Get Execution time
	executionTime := time.Since(startTime)

	// log
	r.logQuery(&logEntry{
		Query:           query,
		Args:            args,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Error:           err,
	})

	// RETURN RESULT
	return rows, err
}

func (r *PostgresDB) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	startTime := time.Now()

	// Execute query
	row := r.DB.QueryRowxContext(ctx, query, args...)

	// Get Execution time
	executionTime := time.Since(startTime)

	// log
	r.logQuery(&logEntry{
		Query:           query,
		Args:            args,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Error:           nil,
	})

	// RETURN RESULT
	return row
}

func (r *PostgresDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	startTime := time.Now()

	// Execute query
	result, err := r.DB.ExecContext(ctx, query, args...)

	// Get Execution time
	executionTime := time.Since(startTime)

	// log
	r.logQuery(&logEntry{
		Query:           query,
		Args:            args,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Error:           err,
	})

	// RETURN RESULT
	return result, err
}

func (r *PostgresDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	startTime := time.Now()

	// Execute query
	err := r.DB.SelectContext(ctx, dest, query, args...)

	// Get Execution time
	executionTime := time.Since(startTime)

	// log
	r.logQuery(&logEntry{
		Query:           query,
		Args:            args,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Error:           err,
	})

	// RETURN RESULT
	return err
}

func (r *PostgresDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	startTime := time.Now()

	// Execute query
	err := r.DB.GetContext(ctx, dest, query, args...)

	// Get Execution time
	executionTime := time.Since(startTime)

	// log
	r.logQuery(&logEntry{
		Query:           query,
		Args:            args,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Error:           err,
	})

	// RETURN RESULT
	return err
}

func (r *PostgresDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	startTime := time.Now()

	// Execute query
	result, err := r.DB.NamedExecContext(ctx, query, arg)

	// Get Execution time
	executionTime := time.Since(startTime)

	// log
	r.logQuery(&logEntry{
		Query:           query,
		Args:            arg,
		ExecutionTimeMs: executionTime.Milliseconds(),
		Error:           err,
	})

	// RETURN RESULT
	return result, err
}
