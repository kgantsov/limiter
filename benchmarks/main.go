package main

import (
	"flag"
	"fmt"

	"github.com/tidwall/redbench"
)

func main() {
	redisPort := flag.Int("redis_port", 46379, "Redis Port")
	flag.Parse()

	redbench.Bench("PING", fmt.Sprintf("127.0.0.1:%d", *redisPort), nil, nil, func(buf []byte) []byte {
		return redbench.AppendCommand(buf, "PING")
	})
	redbench.Bench("REDUCE", fmt.Sprintf("127.0.0.1:%d", *redisPort), nil, nil, func(buf []byte) []byte {
		return redbench.AppendCommand(buf, "REDUCE", "key", "1000", "1", "1000", "1")
	})
}
