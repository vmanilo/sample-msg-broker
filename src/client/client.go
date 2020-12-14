package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	address string
	income chan string
	output chan string
	subscribers []chan string
}

func NewClient(address string) *Client {
	return &Client{address, make(chan string, 10), make(chan string), []chan string{}}
}

func (c *Client) Publish(topic string, message string) {
	c.output <- fmt.Sprintf("PUB %s %s\n", topic, message)
}

func (c *Client) Subscribe(topic string) {
	c.output <- fmt.Sprintf("SUB %s\n", topic)
}

func (c *Client) GetSubscription() chan string {
	subscriber := make(chan string, 10)
	c.subscribers = append(c.subscribers, subscriber)
	return subscriber
}


func (c *Client) receive(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if message == "" {
			continue
		}
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		c.income <- message
	}
}

func (c *Client) readConsoleInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		c.output <- text
	}
}


func (c *Client) Connect() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", c.address)
	if err != nil {
		fmt.Println(err)
		return
	}


	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	go c.receive(conn)
	go c.readConsoleInput()

	go func() {
		for {
			select {
			case msg := <-c.income:
				fmt.Print("->: " + msg)
				if msg != "OK" {
					for _, subscriber := range c.subscribers {
						subscriber <- msg
					}
				}

			case msg := <-c.output:
				conn.Write([]byte(msg+"\n"))

				if strings.TrimSpace(msg) == "STOP" {
					fmt.Println("TCP client exiting...")
					return
				}
			}
		}
	}()

}