package connection

import (
	"github.com/go-redis/redis/v8"
)

// A MySQLConn represents a database connection.
type RedisInterface interface {
	Conn() (*redis.Client, error)
}
