package client

import (
	"bufio"
	"net"
)

// Reader ..
type Reader interface {
	Read() (string, error)
}

// NewlineReader ..
type NewlineReader struct {
	conn *bufio.Reader
}

// BufferedReader ..
type BufferedReader struct {
	conn net.Conn
}

func (nr BufferedReader) Read() (string, error) {
	recvBuf := make([]byte, 1024)
	n, err := nr.conn.Read(recvBuf[:])
	if err != nil {
		return "", err
	}
	return string(recvBuf[:n]), nil
}

func (nr NewlineReader) Read() (string, error) {
	s, err := nr.conn.ReadString('\n')
	if err != nil {
		return "", err
	}
	return s, nil
}
