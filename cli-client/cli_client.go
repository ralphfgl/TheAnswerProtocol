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
	conn     net.Conn
	reader   *bufio.Reader
	username string
}

func NewCLIClient(server string) (*CLIClient, error) {
	conn, err := net.DialTimeout("tcp", server, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &CLIClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *CLIClient) Send(command string) error {
	if !strings.HasSuffix(command, "\n") {
		command += "\n"
	}
	_, err := c.conn.Write([]byte(command))
	return err
}

func (c *CLIClient) Read() (string, error) {
	return c.reader.ReadString('\n')
}

func (c *CLIClient) Close() error {
	return c.conn.Close()
}

func (c *CLIClient) displayMessage(msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}
	fmt.Println(msg)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run client.go <username>")
		return
	}
	username := os.Args[1]
	client, err := NewCLIClient("localhost:8090")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer client.Close()

	greeting, _ := client.Read()
	fmt.Print(greeting)
	client.Send("CONNECT " + username)
	response, _ := client.Read()
	client.displayMessage(response)

	go func() {
		for {
			msg, err := client.Read()
			if err != nil {
				return
			}
			client.displayMessage(msg)
			fmt.Print("> ")
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			fmt.Print("> ")
			continue
		}
		client.Send(line)
		if strings.ToUpper(line) == "QUIT" {
			break
		}
		fmt.Print("> ")
	}
}
