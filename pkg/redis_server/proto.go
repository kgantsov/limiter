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
	p.Metrics.connections.Inc()
	defer func() {
		p.Metrics.connections.Dec()
	}()

	for {
		cmd, err := p.parser.ParseCommand()
		start := time.Now()
		p.Metrics.requestInFlight.Inc()
		status := "OK"

		log.Debug().Msgf("Received command: %v", cmd)

		if err != nil {
			if err == io.EOF {
				log.Debug().Msg("Client has been disconnected")
				p.Metrics.requestInFlight.Dec()
				return
			} else {
				p.responser.SendError(err)
			}
			p.Metrics.requestInFlight.Dec()
			return
		}

		switch cmd.Name {
		case "HELLO":
			p.responser.SendArr([]string{})
		case "REDUCE":
			status, err = p.handleReduceRequest(cmd)
			if err != nil {
				status = "ERROR"
			}
		case "PING":
			p.responser.SendPong()
		default:
			p.responser.SendError(fmt.Errorf("unknown command '%s'", cmd.Args))
			status = "UNKNOWN_COMMAND"
		}
		elapsed := float64(time.Since(start)) / float64(time.Second)
		p.Metrics.requestDuration.WithLabelValues(status).Observe(elapsed)
		p.Metrics.requestsTotal.WithLabelValues(status).Inc()
		p.Metrics.requestInFlight.Dec()
	}
}

func (p *Proto) handleReduceRequest(cmd *Command) (string, error) {
	var status string
	var maxTokens, refillTime, refillAmount, tokens int64
	if len(cmd.Args) < 1 {
		p.responser.SendError(fmt.Errorf("REDUCE expects 5 argument"))
		status = "ARGUMENT_ERROR"
		return status, fmt.Errorf("REDUCE expects 5 argument")
	}

	key := cmd.Args[0]

	if val, err := strconv.ParseInt(cmd.Args[1], 10, 64); err == nil {
		maxTokens = val
	} else {
		p.responser.SendError(fmt.Errorf("REDUCE expects maxTokens to be integer"))
		status = "ARGUMENT_ERROR"
		return status, fmt.Errorf("REDUCE expects maxTokens to be integer")
	}

	if val, err := strconv.ParseInt(cmd.Args[2], 10, 64); err == nil {
		refillTime = val
	} else {
		p.responser.SendError(fmt.Errorf("REDUCE expects refillTime to be integer"))
		status = "ARGUMENT_ERROR"
		return status, fmt.Errorf("REDUCE expects refillTime to be integer")
	}

	if val, err := strconv.ParseInt(cmd.Args[3], 10, 64); err == nil {
		refillAmount = val
	} else {
		p.responser.SendError(fmt.Errorf("REDUCE expects refillAmount to be integer"))
		status = "ARGUMENT_ERROR"
		return status, fmt.Errorf("REDUCE expects refillAmount to be integer")
	}

	if val, err := strconv.ParseInt(cmd.Args[4], 10, 64); err == nil {
		tokens = val
	} else {
		p.responser.SendError(fmt.Errorf("REDUCE expects tokens to be integer"))
		status = "ARGUMENT_ERROR"
		return status, fmt.Errorf("REDUCE expects tokens to be integer")
	}

	tokensLeft, err := p.rateLimiter.Reduce(
		key, maxTokens, refillTime, refillAmount, tokens,
	)

	if err != nil {
		status = "ERROR"
		p.responser.SendError(err)
		return status, err
	}

	status = "OK"
	p.responser.SendInt(tokensLeft)
	return status, nil
}
