package redis_server

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redis"
	"github.com/kgantsov/limiter/pkg/http_server"
	"github.com/kgantsov/limiter/pkg/limiter"
)

func reduce(redisdb *redis.Client, key string, maxTokens int64, refillTime int64, refillAmount int64, tokens int64) *redis.IntCmd {
	cmd := redis.NewIntCmd("REDUCE", key, maxTokens, refillTime, refillAmount, tokens)
	redisdb.Process(cmd)
	return cmd
}

func TestServerBasic(t *testing.T) {
	port := 56379

	app := &http_server.App{
		RateLimiter: limiter.NewRateLimiter(300 * time.Second),
		PathMap:     make(map[string]string),
	}

	go func() {
		server := NewServer(port, app.RateLimiter)
		server.ListenAndServe()
	}()

	time.Sleep(3 * time.Second)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%d", port),
		Password: "",
		DB:       0,
	})

	val, _ := reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(4), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(3), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(2), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(1), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(0), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(-1), val)
	time.Sleep(2 * time.Second)

	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(4), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(3), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(2), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(1), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(0), val)
	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(-1), val)
	time.Sleep(2 * time.Second)

	val, _ = reduce(client, "login", 5, 2, 5, 1).Result()
	assert.Equal(t, int64(4), val)

	client.Close()
}
