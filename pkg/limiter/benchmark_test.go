package limiter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkReduce_100_1000_1_1000_1(b *testing.B) {
	numberOfKeys := 1000
	rl := NewRateLimiter(300 * time.Second)

	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			index := r.Int31n(int32(numberOfKeys))
			key := fmt.Sprintf("user:%d", index)
			rl.Reduce(key, 1000, 1, 1000, 1)
		}
	})
}
func BenchmarkReduce_10000000_1000_10_1000_10(b *testing.B) {
	numberOfKeys := 10000000
	rl := NewRateLimiter(300 * time.Second)

	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			index := r.Int31n(int32(numberOfKeys))
			key := fmt.Sprintf("user:%d", index)
			rl.Reduce(key, 1000, 10, 1000, 10)
		}
	})
}
