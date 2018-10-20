package limiter

import (
	"testing"
)

func benchmarkReduce(key string, maxTokens int64, refillTime int64, refillAmount int64, tokens int64, b *testing.B) {
	b.StopTimer()
	rl := NewRateLimiter()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StartTimer()

		rl.Reduce(key, maxTokens, refillTime, refillAmount, tokens)

		b.StopTimer()
	}
}

func BenchmarkReduce_100_1000(b *testing.B) {
	benchmarkReduce("user:1001", 1000, 1, 1000, 1, b)
}
