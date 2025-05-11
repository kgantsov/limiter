package limiter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlowRateLimiter(t *testing.T) {

	rl := NewRateLimiter(300 * time.Second)

	val, _ := rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(4), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(3), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(2), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(1), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(0), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(-1), val)
	time.Sleep(2 * time.Second)

	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(4), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(3), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(2), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(1), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(0), val)
	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(-1), val)
	time.Sleep(2 * time.Second)

	val, _ = rl.Reduce("login", 5, 2, 5, 1)
	assert.Equal(t, int64(4), val)
}

func TestFastRateLimiter(t *testing.T) {

	rl := NewRateLimiter(300 * time.Second)

	for i := 999; i >= 0; i-- {
		val, _ := rl.Reduce("api_call", 1000, 1, 1000, 1)
		assert.Equal(t, int64(i), val)
	}
	val, _ := rl.Reduce("api_call", 1000, 1, 1000, 1)
	assert.Equal(t, int64(-1), val)

	time.Sleep(1 * time.Second)

	for i := 999; i >= 0; i-- {
		val, _ := rl.Reduce("api_call", 1000, 1, 1000, 1)
		assert.Equal(t, int64(i), val)
	}
	val, _ = rl.Reduce("api_call", 1000, 1, 1000, 1)
	assert.Equal(t, int64(-1), val)
}

func TestRateLimiterWithManyKeys(t *testing.T) {
	rl := NewRateLimiter(300 * time.Second)

	for i := 1000_000; i >= 0; i-- {
		val, _ := rl.Reduce(fmt.Sprintf("api_call:%d", i), 1000, 1, 1000, 1)
		assert.Equal(t, int64(999), val)
	}
}

func TestRateLimiterReuseTheSameKey(t *testing.T) {

	rl := NewRateLimiter(300 * time.Second)

	val, _ := rl.Reduce("user:123", 1000, 1, 50, 1)
	assert.Equal(t, int64(999), val)
	val, _ = rl.Reduce("user:123", 1000, 1, 50, 1)
	assert.Equal(t, int64(998), val)
	val, _ = rl.Reduce("user:123", 1000, 1, 50, 1)
	assert.Equal(t, int64(997), val)

	val, _ = rl.Reduce("user:123", 50, 1, 50, 1)
	assert.Equal(t, int64(49), val)

	val, _ = rl.Reduce("user:123", 1000, 1, 50, 1)
	assert.Equal(t, int64(996), val)
}

func TestRateLimiterD(t *testing.T) {
	rl := NewRateLimiter(300 * time.Second)

	val, _ := rl.Reduce("api_call", 100, 1, 10, 10)
	assert.Equal(t, int64(90), val)
	val, _ = rl.Reduce("api_call", 100, 1, 10, 10)
	assert.Equal(t, int64(80), val)
	val, _ = rl.Reduce("api_call", 100, 1, 10, 10)
	assert.Equal(t, int64(70), val)
	time.Sleep(1 * time.Second)
	val, _ = rl.Reduce("api_call", 100, 1, 10, 50)
	assert.Equal(t, int64(30), val)
}

func TestCleanUpFullBuckets(t *testing.T) {
	rl := NewRateLimiter(300 * time.Second)

	assert.Equal(t, int64(0), rl.Len())

	rl.Reduce("api_call", 100, 1, 1, 1)
	assert.Equal(t, int64(1), rl.Len())

	time.Sleep(2 * time.Second)

	rl.CleanUpFullBuckets()

	assert.Equal(t, int64(0), rl.Len())
}
