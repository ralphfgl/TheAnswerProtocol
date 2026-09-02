package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type CLIClient struct {
	conn		net.Conn
	reader		*bufio.Reader
	username	string
}

func NewCLIClient(server string) (*CLIClient, error) {
	conn, err := net.DialTimeout("tcp", server, 5 * time.Second)
	if err != nil {
		return nil, err
	}
	return &CLIClient{
		conn:	conn,
		reader:	bufio.NewReader(conn),
	}, nil
}

func (c *CLIClient) Send(command string) error {
	if !strings.HasSuffix(command, "\n") {
		command += "\n"
	}
	_, err := c.conn.Write([]byte(command))
	return err
}

func
