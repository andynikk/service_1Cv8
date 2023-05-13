package postgresql

import "github.com/jackc/pgx/v5/pgxpool"

type PgxpoolConn struct {
	*pgxpool.Conn
}
