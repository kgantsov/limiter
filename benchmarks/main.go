package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/tidwall/redbench"
)

func main() {
	redisPort := flag.Int("redis_port", 46379, "Redis Port")
	flag.Parse()

	redbench.Bench("PING", fmt.Sprintf("127.0.0.1:%d", *redisPort), nil, nil, func(buf []byte) []byte {
		return redbench.AppendCommand(buf, "PING")
	})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	redbench.Bench("REDUCE", fmt.Sprintf("127.0.0.1:%d", *redisPort), nil, nil, func(buf []byte) []byte {
		index := r.Int31n(10000)
		key := fmt.Sprintf("key:%d", index)
		return redbench.AppendCommand(buf, "REDUCE", key, "1000", "1", "1000", "1")
	})
}
