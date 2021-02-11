package limiter

import (
	"hash/fnv"
	"math"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

var SHARDS = uint64(128)

type Shard struct {
	Values      map[int64]int64
	LastUpdates map[int64]int64
	mu          sync.RWMutex
}

type RateLimiter struct {
	shards []*Shard
	length int64
}

func NewRateLimiter() *RateLimiter {
	rateLimiter := new(RateLimiter)

	rateLimiter.shards = make([]*Shard, SHARDS)

	for i := uint64(0); i < SHARDS; i++ {
		rateLimiter.shards[i] = &Shard{
			Values: make(map[int64]int64), LastUpdates: make(map[int64]int64),
		}
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

	value, ok := shard.Values[h]
	lastUpdate, ok1 := shard.LastUpdates[h]

	now := time.Now().Unix()

	if !ok || !ok1 {
		value = maxTokens
		lastUpdate = now

		atomic.AddInt64(&l.length, 1)
	}

	now = time.Now().Unix()
	refillCount := math.Floor(
		float64(now-lastUpdate) / float64(refillTime),
	)

	value = int64(math.Min(
		float64(maxTokens),
		float64(value)+(refillCount*float64(refillAmount)),
	))
	lastUpdate = int64(math.Min(
		float64(now),
		float64(lastUpdate)+refillCount*float64(lastUpdate),
	))
	shard.Values[h] = value
	shard.LastUpdates[h] = lastUpdate

	if tokens > value {
		return -1, nil
	}
	value = value - tokens
	shard.Values[h] = value

	return value, nil
}

func (l *RateLimiter) Len() int64 {
	return l.length
}

func (l *RateLimiter) Remove(key string) {
	h := int64(hash(key))
	shard := l.GetShard(h)

	shard.mu.Lock()
	delete(shard.Values, h)
	delete(shard.LastUpdates, h)
	shard.mu.Unlock()

	atomic.AddInt64(&l.length, -1)
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Debugf("%s took %s", name, elapsed)
}
