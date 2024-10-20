package store

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// We support TCP connections, and handle each of
// them in a Go routine.

type Server struct {
	host string
	port int
	data map[string][]byte
}

func NewServer(host string, port int) (*Server, error) {
	s := Server{
		host: host,
		port: port,
		data: make(map[string][]byte),
	}
	return &s, nil
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.host+":"+strconv.Itoa(s.port))
	defer listener.Close()

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer c.Close()

	// read data
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := c.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		s.handleCommand(buffer[:n], &c)
	}
}

func (s *Server) handleCommand(rawCommand []byte, c *net.Conn) {
	// Now we try and figure out the command and its arguments
	// and pass it to handlers.
	if bytes.HasPrefix(rawCommand, []byte("ping")) {
		s.handlePing(c)
	} else if bytes.HasPrefix(rawCommand, []byte("get")) {
		s.handleGet(rawCommand[4:], c)
	} else if bytes.HasPrefix(rawCommand, []byte("set")) {
		s.handleSet(rawCommand[4:], c)
	} else if bytes.HasPrefix(rawCommand, []byte("del")) {
	} else {
		conn := *c
		conn.Write([]byte("ERR\n"))
	}
}

func (s *Server) handleGet(key []byte, c *net.Conn) {
	trimmedKey := strings.TrimSpace(string(key))
	value := s.data[trimmedKey]
	conn := *c
	conn.Write(value)
}

func (s *Server) handleSet(keyVal []byte, c *net.Conn) {
	conn := *c

	splitted := bytes.SplitN(keyVal, []byte(" "), 2)
	if len(splitted) != 2 {
		conn.Write([]byte("ERR: Invalid payload\n"))
		return
	}
	s.data[string(splitted[0])] = splitted[1]
	conn.Write([]byte("OK\n"))
}

func (s *Server) handlePing(c *net.Conn) {
	pong := []byte("PONG\n")
	conn := *c
	conn.Write(pong)
}
