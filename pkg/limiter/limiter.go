package limiter

import (
	"hash/fnv"
	"math"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

var SHARDS = uint64(1024)

type Bucket struct {
	Value     int64
	UpdatedAt int64
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
	if log.GetLevel() == log.DebugLevel {
		defer TimeTrack(time.Now(), "RateLimiter.Reduce")
	}
	h := int64(hash(key))

	shard := l.GetShard(h)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	bucket, ok := shard.Buckets[key]

	now := time.Now().Unix()

	if !ok {
		value := maxTokens - tokens

		bucket = Bucket{
			Value:     value,
			UpdatedAt: now,
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

	bucket = Bucket{Value: value, UpdatedAt: lastUpdate}
	shard.Buckets[key] = bucket

	return bucket.Value, nil
}

func (l *RateLimiter) Len() int64 {
	return l.length
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
	log.Debugf("%s took %s", name, elapsed)
}
