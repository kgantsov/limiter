package redis_server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/kgantsov/limiter/pkg/limiter"

	log "github.com/sirupsen/logrus"
)

// ListenAndServe accepts incoming connections on the creating a new service goroutine for each.
// The service goroutines read requests and then replies to them.
// It exits program if it can not start tcp listener.
func ListenAndServe(port int, rateLimiter *limiter.RateLimiter) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Info(sig)

		log.Info("Closing rate limiter app")

		os.Exit(0)
	}()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Info("Listening on port: ", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Fatal error: ", err.Error())
			continue
		}
		go handleClient(rateLimiter, conn)
	}
}

func handleClient(rateLimiter *limiter.RateLimiter, conn net.Conn) {
	reader := bufio.NewReader(conn)
	parser := newParser(reader)
	responser := newResponser(conn)
	defer conn.Close()

	for {
		cmd, err := parser.ParseCommand()
		if err != nil {
			if err == io.EOF {
				log.Debug("Client has been disconnected")
			} else if _, ok := err.(error); ok {
				responser.sendError(err)
			} else {
				log.Debug("Errror parsing command: %s", err)
			}
			return
		}

		switch cmd.Name {
		case "REDUCE":
			var maxTokens, refillTime, refillAmount, tokens int64
			if len(cmd.Args) < 1 {
				responser.sendError(fmt.Errorf("REDUCE expects 5 argument"))
				return
			}

			key := cmd.Args[0]

			if val, err := strconv.ParseInt(cmd.Args[1], 10, 64); err == nil {
				maxTokens = val
			} else {
				responser.sendError(fmt.Errorf("REDUCE expects maxTokens to be integer"))
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[2], 10, 64); err == nil {
				refillTime = val
			} else {
				// conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
				responser.sendError(fmt.Errorf("REDUCE expects refillTime to be integer"))
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[3], 10, 64); err == nil {
				refillAmount = val
			} else {
				responser.sendError(fmt.Errorf("REDUCE expects refillAmount to be integer"))
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[4], 10, 64); err == nil {
				tokens = val
			} else {
				responser.sendError(fmt.Errorf("REDUCE expects tokens to be integer"))
				continue
			}

			tokensLeft, err := rateLimiter.Reduce(
				key, maxTokens, refillTime, refillAmount, tokens,
			)

			if err == nil {
				responser.sendInt(tokensLeft)
			} else {
				responser.sendError(err)
			}
		case "PING":
			responser.sendPong()
		default:
			responser.sendError(fmt.Errorf("unknown command '%s'\r\n", cmd.Args))
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error: ", err.Error())
	}
}
