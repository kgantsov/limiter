package redis_server

import (
	"fmt"
	"net"
)

type responser struct {
	conn net.Conn
}

func newResponser(conn net.Conn) *responser {
	r := &responser{conn}

	return r
}

func (r *responser) sendError(err error) {
	r.conn.Write([]byte(fmt.Sprintf("-ERR %s\r\n", err)))
	// log.Debug(fmt.Sprintf("-ERR %s\r\n", err))
}

func (r *responser) sendPong() {
	r.conn.Write([]byte("+PONG\r\n"))
}

func (r *responser) sendInt(value int64) {
	r.conn.Write([]byte(fmt.Sprintf(":%d\r\n", value)))
}
