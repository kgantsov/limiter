package redis_server

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/kgantsov/limiter/pkg/limiter"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Port             int
	RateLimiter      *limiter.RateLimiter
	EnablePrometheus bool
	Metrics          *Metrics
	TCPListener      *net.TCPListener

	quit chan interface{}
	wg   sync.WaitGroup
}

func NewServer(port int, rateLimiter *limiter.RateLimiter, enablePrometheus bool) *Server {
	server := &Server{
		Port:             port,
		RateLimiter:      rateLimiter,
		EnablePrometheus: enablePrometheus,
		quit:             make(chan interface{}),
		Metrics:          NewMetrics("redis"),
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", server.Port))
	checkError(err)
	server.TCPListener, err = net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	return server
}

func (srv *Server) ListenAndServe() {
	log.Info().Msgf("Listening on port: %d", srv.Port)

	for {
		conn, err := srv.TCPListener.Accept()
		if err != nil {
			select {
			case <-srv.quit:
				return
			default:
				log.Error().Msgf("Fatal error: %s", err.Error())
				continue
			}
		}

		srv.wg.Add(1)

		go func() {
			srv.handleClient(conn)
			srv.wg.Done()
		}()
	}
}

func (srv *Server) Stop() {
	close(srv.quit)
	srv.TCPListener.Close()
	srv.wg.Wait()
}

func (srv *Server) handleClient(conn io.ReadWriteCloser) {
	redisProto := NewProto(srv.Metrics, srv.RateLimiter, conn, conn)
	defer conn.Close()

	redisProto.HandleRequest()
}

func checkError(err error) {
	if err != nil {
		log.Error().Msgf("Fatal error: %s", err.Error())
	}
}
