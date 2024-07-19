package redis_server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/kgantsov/limiter/pkg/limiter"
	"github.com/rs/zerolog/log"
)

type Proto struct {
	parser      *Parser
	responser   *Responser
	rateLimiter *limiter.RateLimiter
	Metrics     *Metrics
}

func NewProto(Metrics *Metrics, rateLimiter *limiter.RateLimiter, reader io.Reader, writer io.Writer) *Proto {
	r := bufio.NewReader(reader)
	parser := NewParser(r)
	responser := NewResponser(writer)

	p := &Proto{
		parser:      parser,
		responser:   responser,
		rateLimiter: rateLimiter,
		Metrics:     Metrics,
	}

	return p
}

func (p *Proto) HandleRequest() {
	for {
		cmd, err := p.parser.ParseCommand()
		start := time.Now()
		status := "OK"
		// log.Debug().Msgf("Received command: %v", cmd)

		if err != nil {
			if err == io.EOF {
				log.Debug().Msg("Client has been disconnected")
				return
			} else {
				p.responser.SendError(err)
			}

			return
		}

		switch cmd.Name {
		case "HELLO":
			p.responser.SendArr([]string{})
		case "REDUCE":
			var maxTokens, refillTime, refillAmount, tokens int64
			if len(cmd.Args) < 1 {
				p.responser.SendError(fmt.Errorf("REDUCE expects 5 argument"))
				status = "ARGUMENT_ERROR"
				continue
			}

			key := cmd.Args[0]

			if val, err := strconv.ParseInt(cmd.Args[1], 10, 64); err == nil {
				maxTokens = val
			} else {
				p.responser.SendError(fmt.Errorf("REDUCE expects maxTokens to be integer"))
				status = "ARGUMENT_ERROR"
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[2], 10, 64); err == nil {
				refillTime = val
			} else {
				p.responser.SendError(fmt.Errorf("REDUCE expects refillTime to be integer"))
				status = "ARGUMENT_ERROR"
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[3], 10, 64); err == nil {
				refillAmount = val
			} else {
				p.responser.SendError(fmt.Errorf("REDUCE expects refillAmount to be integer"))
				status = "ARGUMENT_ERROR"
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[4], 10, 64); err == nil {
				tokens = val
			} else {
				p.responser.SendError(fmt.Errorf("REDUCE expects tokens to be integer"))
				status = "ARGUMENT_ERROR"
				continue
			}

			tokensLeft, err := p.rateLimiter.Reduce(
				key, maxTokens, refillTime, refillAmount, tokens,
			)

			if err == nil {
				p.responser.SendInt(tokensLeft)
			} else {
				p.responser.SendError(err)
			}
		case "PING":
			p.responser.SendPong()
		default:
			p.responser.SendError(fmt.Errorf("unknown command '%s'", cmd.Args))
			status = "UNKNOWN_COMMAND"
		}
		if p.Metrics != nil {
			elapsed := float64(time.Since(start)) / float64(time.Second)
			p.Metrics.reqDurations.Observe(elapsed)
			p.Metrics.reqCount.WithLabelValues(status).Inc()
		}
	}
}
