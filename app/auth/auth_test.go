package auth_test

import "github.com/redis/go-redis/v9"

var rds *redis.Client

func init() {
	rds = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
}
