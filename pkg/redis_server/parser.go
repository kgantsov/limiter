package redis_server

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type parser struct {
	reader *bufio.Reader
}

type command struct {
	Name string
	Args []string
}

func newParser(reader *bufio.Reader) *parser {
	p := &parser{}
	p.reader = reader

	return p
}

func (p *parser) ParseCommand() (*command, error) {
	for {
		line, err := p.readLine()
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
			line, err := p.readLine()
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
			if _, err := io.ReadFull(p.reader, arg); err != nil {
				return nil, err
			}

			args = append(args, string(arg[0:len(arg)-2]))
		}

		return &command{Name: args[0], Args: args[1:]}, nil
	}
}

func (p *parser) readLine() (string, error) {
	str, err := p.reader.ReadString('\n')
	if err == nil {
		return str[:len(str)-2], err
	}
	return str, err
}
