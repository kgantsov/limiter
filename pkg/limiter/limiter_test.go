package limiter

import (
	"fmt"
	"testing"
	"time"
)

func assetEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		fmt.Printf("Expected `%t`. Got `%t`\n", expected, actual)
		t.Errorf("Expected `%#v`. Got `%#v`\n", expected, actual)
	}
}

func TestSlowRateLimiter(t *testing.T) {

	rl := NewRateLimiter()
	rl.AddLimiter("login", 5, 2, 5)

	val, _ := rl.Reduce("login", 1)
	assetEqual(t, int64(4), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(3), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(2), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(1), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(0), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(-1), val)
	time.Sleep(2 * time.Second)

	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(4), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(3), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(2), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(1), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(0), val)
	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(-1), val)
	time.Sleep(2 * time.Second)

	val, _ = rl.Reduce("login", 1)
	assetEqual(t, int64(4), val)
}

func TestFastRateLimiter(t *testing.T) {

	rl := NewRateLimiter()
	rl.AddLimiter("api_call", 1000, 1, 1000)

	for i := 999; i >= 0; i-- {
		val, _ := rl.Reduce("api_call", 1)
		assetEqual(t, int64(i), val)
	}
	val, _ := rl.Reduce("api_call", 1)
	assetEqual(t, int64(-1), val)

	time.Sleep(1 * time.Second)

	for i := 999; i >= 0; i-- {
		val, _ := rl.Reduce("api_call", 1)
		assetEqual(t, int64(i), val)
	}
	val, _ = rl.Reduce("api_call", 1)
	assetEqual(t, int64(-1), val)
}
