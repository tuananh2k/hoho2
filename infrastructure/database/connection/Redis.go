package connection

import (
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

// A PostgresCon represents a mysql connection
type RedisCon struct{}

// NewPostgresCon creates a PostgresCon struct.
func NewRedisCon() RedisInterface {
	return &RedisCon{}
}

// Conn get a mysql connection.
func (mc *RedisCon) Conn() (*redis.Client, error) {
	// dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", , , , , )
	// conn, err := sql.Open(os.Getenv("DB_DRIVER"), dataSourceName)
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return conn, nil
}
