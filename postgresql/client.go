package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	_ "github.com/jackc/pgx/v4/stdlib" // is required by std sql lib

	"github.com/triabokon/gotagv/flags"
)

var StatementBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DSNFromConfig(c Config) string { //nolint:gocritic
	var parts []string

	for key, value := range map[string]string{
		"host":     c.Host,
		"port":     c.Port,
		"user":     c.User,
		"password": c.Password,
		"dbname":   c.Database,
		"sslmode":  c.SSLMode,
	} {
		if value != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(parts, " ")
}

func (c *Config) Flags(name, prefix string) *pflag.FlagSet {
	f := pflag.NewFlagSet(name, pflag.PanicOnError)
	f.StringVar(&c.Host, "host", "127.0.0.1", "")
	f.StringVar(&c.Port, "port", "5432", "")
	f.StringVar(&c.User, "user", "postgres", "")
	f.StringVar(&c.Password, "password", "", "[secret]")
	f.StringVar(&c.Database, "database", "postgres", "database which connect to")
	f.StringVar(&c.SSLMode, "ssl_mode", "disable", "ssl mode")

	return flags.MapWithPrefix(f, name, pflag.PanicOnError, prefix)
}

type Client struct {
	DB *pgxpool.Pool
}

func (c *Client) Stat() *pgxpool.Stat {
	return c.DB.Stat()
}

func New(ctx context.Context, conf Config) (*Client, func() error, error) { //nolint:gocritic
	dsn := DSNFromConfig(conf)
	parsedConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse config")
	}

	db, err := pgxpool.ConnectConfig(ctx, parsedConf)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create pgx pool connection")
	}
	closer := func() error {
		if db != nil {
			db.Close()
		}
		return nil
	}

	return &Client{DB: db}, closer, nil
}

func NewStdSQL(config Config) (*sql.DB, func() error, error) { //nolint:gocritic
	dsn := DSNFromConfig(config)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open connection")
	}
	return db, func() error { return errors.WithStack(db.Close()) }, nil
}
