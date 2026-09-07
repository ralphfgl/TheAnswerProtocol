package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"the_answer_protocol/common"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

type CLIClient struct {
	conn        net.Conn
	reader      *bufio.Reader
	username    string
	currentRoom string
}

// NewCLIClient constructor
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
	if strings.HasPrefix(msg, "OK") && strings.Contains(msg, `"id"`) {
		c.displayRoom(msg)
		return
	}
	fmt.Println(msg)
}

func (c *CLIClient) displayRoom(response string) {
	if strings.HasPrefix(response, "OK") {
		response = strings.TrimPrefix(response, "OK")
		response = strings.TrimSpace(response)
	}

	var lookResp common.LookResponse
	if err := json.Unmarshal([]byte(response), &lookResp); err != nil {
		fmt.Printf("%s[Error parsing room data: %v]%s\n", ColorRed, err, ColorReset)
		fmt.Printf("%sRaw: %s%s\n", ColorRed, response, ColorReset)
		return
	}
	c.currentRoom = lookResp.Room.Id
	//fmt.Print("\033[2J\033[H")
	fmt.Printf("=== %s ===\n", lookResp.Room.Name)
	fmt.Printf("%s\n\n", lookResp.Room.Description)
	// Exits
	if len(lookResp.Room.Exits) > 0 {
		exits := make([]string, 0, len(lookResp.Room.Exits))
		for dir := range lookResp.Room.Exits {
			exits = append(exits, dir)
		}
		fmt.Printf("Exits: %s\n", strings.Join(exits, ", "))
	}
	// Players
	if len(lookResp.Players) > 0 {
		fmt.Printf("Players: %s\n", strings.Join(lookResp.Players, ", "))
	}
	// Items
	if len(lookResp.Items) > 0 {
		fmt.Printf("Items: %s\n", strings.Join(lookResp.Items, ", "))
	}
	// NPCs
	if len(lookResp.NPCs) > 0 {
		fmt.Printf("NPCs: %s\n", strings.Join(lookResp.NPCs, ", "))
	}
	fmt.Println()
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
			if !strings.Contains(msg, `"id"`) {
				fmt.Print("> ")
			}
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
