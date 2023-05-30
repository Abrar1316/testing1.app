package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/helpers/util"
	"github.com/go-pg/pg/v10"
)

var (
	lock sync.Mutex
	pgDB = make(map[string]*pg.DB)
)

// `getDB` returns, or initialize DB handler for a particular class
func getDB(opts ...ConfigFunc) *pg.DB {
	const section = "pg"

	lock.Lock()
	defer lock.Unlock()
	db, ok := pgDB[section]
	if ok {
		return db
	}

	cfg := util.GetBillAppConfigInstance()

	// cfg.LoadConfig(filepath.Join("config", "config.ini"))

	postgresIni := cfg.GetPostgresIni(false, section)

	pgDB[section] = NewDB(
		postgresIni.GetHost(),
		postgresIni.GetUsername(),
		postgresIni.GetPassword(),
		postgresIni.GetDatabase(),
		opts...,
	)
	pgDB[section].AddQueryHook(dbLogger{})

	return pgDB[section]
}

// `GetDB` returns the db handle for a particular class of DB.
// Defaults to emulating the old behavior - no parms means `rw`.
func GetDB(opts ...ConfigFunc) *pg.DB {
	return getDB(opts...)
}

// `getDBForTests` returns a handle that must point to a test DB
// (validated in `GetPostgresIni`). Not a fan of essentially duplicating
// `getDB` here but the control flow makes it difficult otherwise
func getDBForTests() *pg.DB {
	// Test DB handle must be rw because it's used for destructive ops.
	section := "pg"

	lock.Lock()
	defer lock.Unlock()
	db, ok := pgDB[section]
	if ok {
		return db
	}
	cfg := util.GetBillAppConfigInstance()
	// `true` here means "must be a test host"
	postgresIni := cfg.GetPostgresIni(true, section)

	connection := NewDB(
		postgresIni.GetHost(),
		postgresIni.GetUsername(),
		postgresIni.GetPassword(),
		postgresIni.GetDatabase(),
		WithPoolSize(100),
	)
	// Uncomment the line below to print raw queries
	// connection.AddQueryHook(dbLogger{})
	pgDB[section] = connection

	return pgDB[section]
}

// dblogger is used for debugging of Postgres queries, by printing the raw queries to stdout.
//
//lint:ignore U1000 Used by uncommenting a line above.
type dbLogger struct{}

type ConfigFunc func(opts *pg.Options)

func defaultConfig(addr string, user string, password string, database string) *pg.Options {
	return &pg.Options{
		Addr:                  addr,
		User:                  user,
		Password:              password,
		Database:              database,
		PoolSize:              400,
		IdleTimeout:           time.Minute * 2,
		DialTimeout:           time.Second * 2,
		PoolTimeout:           time.Second * 2,
		RetryStatementTimeout: false,
		MaxConnAge:            time.Minute * 2,
	}
}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, _ := q.FormattedQuery()
	fmt.Println(string(query))
	return nil
}

// Insert Exported pg Insert wrapper
func Insert(object interface{}) error {
	_, err := getDB().Model(object).Insert()
	return err
}

// Update Exported pg Update wrapper
func Update(object interface{}) error {
	result, err := getDB().Model(object).WherePK().Update()
	if err != nil {
		return err
	}
	return assertOneRow(result.RowsAffected())
}

// Delete Exported pg Delete wrapper
func Delete(object interface{}) error {
	result, err := getDB().Model(object).WherePK().Delete()
	if err != nil {
		return err
	}
	return assertOneRow(result.RowsAffected())
}

// Select Exported pg Select wrapper
func Select(object interface{}) error {
	return getDB().Model(object).Select()
}

func assertOneRow(l int) error {
	switch {
	case l == 0:
		return pg.ErrNoRows
	case l > 1:
		return pg.ErrMultiRows
	default:
		return nil
	}
}

// NewDB sets up a connection to a Postgres DB instance ,
// using the provided mandatory values.
//
// If additional configuration needs to be tweaked, pass
// the postgres.WithX functions as opts arguments.
func NewDB(
	address string,
	user string,
	password string,
	database string,
	opts ...ConfigFunc,
) *pg.DB {
	config := defaultConfig(address, user, password, database)
	for _, f := range opts {
		f(config)
	}
	connection := pg.Connect(config)
	// Uncomment the line below to print raw queries.
	// connection.AddQueryHook(dbLogger{})
	return connection
}

// WithPoolSize sets the connection pool size for a Postgres config.
// This connection will not make more connections at a time and
// will instead wait for connections to be released to the queue first.
func WithPoolSize(size int) func(opts *pg.Options) {
	return func(opts *pg.Options) {
		opts.PoolSize = size
	}
}

// WithPoolTimeout sets the amount of time to wait for an available
// connection from the pool, for a Postgres config.
func WithPoolTimeout(duration time.Duration) func(opts *pg.Options) {
	return func(opts *pg.Options) {
		opts.PoolTimeout = duration
	}
}

// WithIdleTimeout sets the idle timeout for a Postgres config to the given duration
func WithIdleTimeout(duration time.Duration) func(opts *pg.Options) {
	return func(opts *pg.Options) {
		opts.IdleTimeout = duration
	}
}

// WithIdleTimeout sets the dial timeout for a Postgres config to the given duration
func WithDialTimeout(duration time.Duration) func(opts *pg.Options) {
	return func(opts *pg.Options) {
		opts.DialTimeout = duration
	}
}

// WithRetryStatementTimeout instructs a Postgres config that
// queries which timed out should be retried.
func WithRetryStatementTimeout() func(opts *pg.Options) {
	return func(opts *pg.Options) {
		opts.RetryStatementTimeout = true
	}
}

// WithMaxConnAge sets the maximum connection age for a Postgres
// config to the given duration. Connection older than this will
// be closed, which is useful with proxies like PgBouncer.
func WithMaxConnAge(duration time.Duration) func(opts *pg.Options) {
	return func(opts *pg.Options) {
		opts.MaxConnAge = duration
	}
}
