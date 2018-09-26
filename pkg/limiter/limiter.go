package limiter

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	Value        int64
	MaxTokens    int64
	RefillTime   int64
	RefillAmount int64
	LastUpdate   int64
}

type RateLimiter struct {
	Buckets map[string]*Bucket
	mu      sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
	rateLimiter := new(RateLimiter)
	rateLimiter.Buckets = make(map[string]*Bucket)

	return rateLimiter
}

func (l *RateLimiter) Reduce(key string, max_tokens int64, refill_time int64, refill_amount int64, tokens int64) (int64, error) {
	l.mu.RLock()
	bucket, ok := l.Buckets[key]
	l.mu.RUnlock()

	if !ok {
		bucket = &Bucket{
			Value:        max_tokens,
			MaxTokens:    max_tokens,
			RefillTime:   refill_time,
			RefillAmount: refill_amount,
			LastUpdate:   time.Now().Unix(),
		}
		l.mu.Lock()
		l.Buckets[key] = bucket
		l.mu.Unlock()
	}

	now := time.Now().Unix()
	refillCount := math.Floor(
		float64(now-bucket.LastUpdate) / float64(bucket.RefillTime),
	)

	atomic.StoreInt64(
		&bucket.Value,
		int64(math.Min(
			float64(bucket.MaxTokens),
			float64(bucket.Value)+(refillCount*float64(bucket.RefillAmount)),
		)),
	)

	atomic.StoreInt64(
		&bucket.LastUpdate,
		int64(math.Min(
			float64(now),
			float64(bucket.LastUpdate)+refillCount*float64(bucket.LastUpdate),
		)),
	)

	if tokens > bucket.Value {
		return -1, nil
	}

	atomic.AddInt64(
		&bucket.Value,
		-tokens,
	)

	return bucket.Value, nil
}
