package bun

import (
	"context"
	"database/sql"
	"time"

	"mattermost/pkg/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// Database config that contains dsn string.
type Config struct {
	URL                   string        `yaml:"url" mapstructure:"url"`
	DialTimeout           time.Duration `yaml:"dial_timeout"  mapstructure:"dial_timeout"`
	IdleTimeout           time.Duration `yaml:"idle_timeout"  mapstructure:"idle_timeout"`
	ReadTimeout           time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout          time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	RetryStatementTimeout bool          `yaml:"retry_statement_timeout"  mapstructure:"retry_statement_timeout"`
	MaxRetries            int           `yaml:"max_retries"  mapstructure:"max_retries"`
	MaxRetryBackoff       time.Duration `yaml:"max_retry_backoff" mapstructure:"max_retry_backoff"`
	PoolSize              int           `yaml:"pool_size"  mapstructure:"pool_size"`
	PoolTimeout           time.Duration `yaml:"pool_timeout" mapstructure:"pool_timeout"`
}

// NewDatabase creates new DatabaseConfig instance
// with no default configurations.
func NewConfig() (db *Config) {
	db = new(Config)
	return
}

// The DB represents Database.
type DB struct {
	// underlying bun-db instance, embed
	*bun.DB
	// wrapped zap.Logger instance.
	log *logger.Logger
}

// Connect to database using by DSN string.
func New(c *Config, log *logger.Logger) (db *DB, err error) {

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(c.URL)))

	db = new(DB)
	db.log = log
	db.DB = bun.NewDB(sqldb, pgdialect.New())

	// check out db server and fail if not responding.
	if err = db.Ping(context.Background()); err != nil {
		db.Close()
		return
	}

	return
}

// Ping DB Server.
func (db *DB) Ping(ctx context.Context) (err error) {
	return db.DB.PingContext(ctx)
}

// Close DB.
func (db *DB) Close() (err error) {
	return db.DB.Close()
}
