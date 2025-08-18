// pkg/redis/client.go
package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient membuat koneksi baru ke server Redis.
func NewRedisClient(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr, // e.g., "localhost:6379"
		Password: "",   // no password set
		DB:       0,    // use default DB
	})

	// Lakukan ping untuk memastikan koneksi berhasil
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		// Dalam aplikasi sungguhan, ini harus di-handle dengan lebih baik,
		// mungkin dengan panic atau logger.fatal
		panic(err)
	}

	return rdb
}
