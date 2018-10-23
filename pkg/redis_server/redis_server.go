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
	defer conn.Close()

	for {
		cmd, err := parseCommand(reader)
		if err != nil {
			if err == io.EOF {
				log.Debug("Client has been disconnected")
			} else if _, ok := err.(error); ok {
				conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
				log.Debug(fmt.Sprintf("-ERR %s\r\n", err))
			} else {
				log.Debug("Errror parsing command: %s", err)
			}
			return
		}

		switch cmd.Name {
		case "REDUCE":
			var maxTokens, refillTime, refillAmount, tokens int64
			if len(cmd.Args) < 1 {
				conn.Write([]byte("REDUCE expects 5 argument"))
				return
			}

			key := cmd.Args[0]

			if val, err := strconv.ParseInt(cmd.Args[1], 10, 64); err == nil {
				maxTokens = val
			} else {
				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[2], 10, 64); err == nil {
				refillTime = val
			} else {
				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[3], 10, 64); err == nil {
				refillAmount = val
			} else {
				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
				continue
			}

			if val, err := strconv.ParseInt(cmd.Args[4], 10, 64); err == nil {
				tokens = val
			} else {
				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
				continue
			}

			tokensLeft, err := rateLimiter.Reduce(
				key, maxTokens, refillTime, refillAmount, tokens,
			)

			if err == nil {
				conn.Write([]byte(fmt.Sprintf(":%d\r\n", tokensLeft)))
			} else {
				conn.Write([]byte(fmt.Sprintf("$-1\r\n")))
			}
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		default:
			conn.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd.Args)))
		}
	}
}

type command struct {
	Name string
	Args []string
}

func parseCommand(reader *bufio.Reader) (*command, error) {
	for {
		line, err := readLine(reader)
		if err != nil {
			return nil, err
		}

		if line[0] != '*' {
			return &command{Name: line}, nil
		}

		argcStr := line[1:]
		argc, err := strconv.ParseUint(argcStr, 10, 64)

		if err != nil || argc < 1 {
			return nil, fmt.Errorf("Error parsing number of arguments %s", err)
		}

		args := make([]string, 0, argc)
		for i := 0; i < int(argc); i++ {
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}

			if line[0] != '$' {
				return nil, fmt.Errorf("Error parsing argument %s", line)
			}

			argLenStr := line[1:]
			argLen, err := strconv.ParseUint(argLenStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Error parsing argument length %s", argLenStr)
			}

			arg := make([]byte, argLen+2)
			if _, err := io.ReadFull(reader, arg); err != nil {
				return nil, err
			}

			args = append(args, string(arg[0:len(arg)-2]))
		}

		return &command{Name: args[0], Args: args[1:]}, nil
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\n')
	if err == nil {
		return str[:len(str)-2], err
	}
	return str, err
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error: ", err.Error())
	}
}
