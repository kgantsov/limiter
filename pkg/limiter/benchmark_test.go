package limiter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func benchmarkReduce(numberOfKeys int, maxTokens int64, refillTime int64, refillAmount int64, tokens int64, b *testing.B) {
	b.StopTimer()
	rl := NewRateLimiter()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		index := r.Int31n(int32(numberOfKeys))
		key := fmt.Sprintf("user:%d", index)

		b.StartTimer()

		rl.Reduce(key, maxTokens, refillTime, refillAmount, tokens)

		b.StopTimer()
	}
}

func BenchmarkReduce_100_1000(b *testing.B) {
	benchmarkReduce(100, 1000, 1, 1000, 1, b)
}
