package connection

import "hoho-framework-v2/adapters/repository"

// A MySQLConn represents a database connection.
type PostgresConInterface interface {
	Conn() (*repository.SymperOrm, error)
}
