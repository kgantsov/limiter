package limiter

import (
	"hash/fnv"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

var SHARDS = uint64(1024)

type Bucket struct {
	Value        int64
	WillBeFullAt int64
	UpdatedAt    int64
}

type Shard struct {
	Buckets map[string]Bucket
	mu      sync.RWMutex
}

type RateLimiter struct {
	shards []*Shard
	length int64
}

func NewRateLimiter() *RateLimiter {
	rateLimiter := new(RateLimiter)

	rateLimiter.shards = make([]*Shard, SHARDS)

	for i := uint64(0); i < SHARDS; i++ {
		rateLimiter.shards[i] = &Shard{Buckets: make(map[string]Bucket)}
	}

	return rateLimiter
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (l *RateLimiter) GetShard(key int64) *Shard {
	return l.shards[uint64(key)%SHARDS]
}

func (l *RateLimiter) Reduce(key string, maxTokens int64, refillTime int64, refillAmount int64, tokens int64) (int64, error) {
	// if log.GetLevel() == log.DebugLevel {
	// 	defer TimeTrack(time.Now(), "RateLimiter.Reduce")
	// }
	h := int64(hash(key))

	shard := l.GetShard(h)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	bucket, ok := shard.Buckets[key]

	now := time.Now().Unix()

	log.Debug().Msgf("Bucket for key: %s found: %t now %d: %+v", key, ok, now, bucket)

	if !ok {
		value := maxTokens - tokens
		tokensNeeded := maxTokens - value
		refillCyclesNeeded := int64(math.Ceil(float64(tokensNeeded) / float64(refillAmount)))
		WillBeFullAt := now + (refillCyclesNeeded * refillTime)

		bucket = Bucket{
			Value:        value,
			WillBeFullAt: WillBeFullAt,
			UpdatedAt:    now,
		}
		shard.Buckets[key] = bucket

		atomic.AddInt64(&l.length, 1)
		return bucket.Value, nil
	}

	refillCount := math.Floor(
		float64(now-bucket.UpdatedAt) / float64(refillTime),
	)

	value := int64(math.Min(
		float64(maxTokens),
		float64(bucket.Value)+(refillCount*float64(refillAmount)),
	))
	lastUpdate := int64(math.Min(
		float64(now),
		float64(bucket.UpdatedAt)+refillCount*float64(bucket.UpdatedAt),
	))

	if tokens > value {
		return -1, nil
	}

	value = value - tokens
	tokensNeeded := maxTokens - value
	refillCyclesNeeded := int64(math.Ceil(float64(tokensNeeded) / float64(refillAmount)))
	WillBeFullAt := now + (refillCyclesNeeded * refillTime)

	bucket = Bucket{Value: value, WillBeFullAt: WillBeFullAt, UpdatedAt: lastUpdate}
	shard.Buckets[key] = bucket

	return bucket.Value, nil
}

func (l *RateLimiter) Len() int64 {
	return atomic.LoadInt64(&l.length)
}

func (l *RateLimiter) CleanUpFullBuckets() {

	for _, shard := range l.shards {
		shard.mu.Lock()
		for key, bucket := range shard.Buckets {
			if bucket.WillBeFullAt < time.Now().Unix() {
				delete(shard.Buckets, key)
				atomic.AddInt64(&l.length, -1)
			}
		}
		shard.mu.Unlock()
	}
}

func (l *RateLimiter) StartCleanUpFullBuckets() {
	for {
		time.Sleep(60 * time.Second)
		started := time.Now()

		l.CleanUpFullBuckets()

		TimeTrack(started, "StartCleanUpFullBuckets took")
	}
}

func (l *RateLimiter) Remove(key string) {
	h := int64(hash(key))
	shard := l.GetShard(h)

	shard.mu.Lock()
	delete(shard.Buckets, key)
	shard.mu.Unlock()

	atomic.AddInt64(&l.length, -1)
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debug().Msgf("%s took %s", name, elapsed)
}
