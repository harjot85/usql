package pgx

import (
	"database/sql"

	// DRIVER: pgx
	"github.com/jackc/pgx/stdlib"

	"github.com/jackc/pgx"
	"github.com/knq/dburl"

	"github.com/knq/usql/drivers"
)

const (
	pgxMaxConnections = 3
)

func init() {
	drivers.Register("pgx", drivers.Driver{
		V: func(db *sql.DB) (string, error) {
			var ver string
			err := db.QueryRow(`show server_version`).Scan(&ver)
			if err != nil {
				return "", err
			}
			return "PostgreSQL " + ver, nil
		},
		O: func(u *dburl.URL) (func(string, string) (*sql.DB, error), error) {
			var err error

			u.DSN, err = dburl.GenPostgres(u)
			if err != nil {
				return nil, err
			}

			cfg, err := pgx.ParseDSN(u.DSN)
			if err != nil {
				return nil, err
			}

			pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
				ConnConfig:     cfg,
				MaxConnections: pgxMaxConnections,
			})
			if err != nil {
				return nil, err
			}

			return func(string, string) (*sql.DB, error) {
				return stdlib.OpenFromConnPool(pool)
			}, nil
		},
		E: func(err error) (string, string) {
			if e, ok := err.(pgx.PgError); ok {
				return e.Code, e.Message
			}
			return "", err.Error()
		},
		PwErr: func(err error) bool {
			if e, ok := err.(pgx.PgError); ok {
				return e.Code == "28P01"
			}
			return false
		},
	})
}
