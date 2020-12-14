package broker

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	id     string
	conn   net.Conn
	input  chan string
	output chan string
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn.RemoteAddr().String(),
		conn,
		make(chan string, 10),
		make(chan string),
	}
}

type Server struct {
	port    int
	clients map[string]*Client
	topics  map[string]map[string]bool
}

func NewServer(port int) *Server {
	clients := make(map[string]*Client)
	topics := make(map[string]map[string]bool)
	return &Server{port, clients, topics}
}

func (s *Server) Start() error {
	fmt.Println("Starting server...")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		fmt.Println("New client connected")
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}
}

func (c *Client) receive() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
		c.input <- msg
	}
}

func (c *Client) send() {
	for msg := range c.output {
		_, err := c.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("ERROR:", err)
			return
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	client := NewClient(conn)
	s.clients[client.id] = client

	go client.receive()
	go client.send()

	for {
		select {
		case msg := <-client.input:
			msg = strings.TrimSpace(msg)
			if msg == "" {
				continue
			}

			fmt.Println("Got msg:", msg)
			if strings.TrimSpace(msg) == "STOP" {
				fmt.Println("Exiting TCP server!")
				break
			}

			tokens := strings.Split(msg, " ")
			if len(tokens) < 2 {
				continue
			}

			cmd := strings.TrimSpace(tokens[0])
			topic := strings.TrimSpace(tokens[1])

			if cmd == "PUB" && len(tokens) >= 3 {
				data := strings.TrimSpace(strings.Join(tokens[2:], " "))

				subscribers, ok := s.topics[topic]

				if ok && subscribers != nil {
					for clientID, _ := range subscribers {
						fmt.Printf("Notify subscriber [%s]: %s\n", clientID, data)
						s.clients[clientID].output <- data
					}
				}
			} else if cmd == "SUB" {
				fmt.Println("Add new subscriber ", client.id)
				if subscribers, ok := s.topics[topic]; ok {
					subscribers[client.id] = true
				} else {
					subscribers := make(map[string]bool)
					subscribers[client.id] = true
					s.topics[topic] = subscribers
				}
			}

			client.output <- "OK"
		}

	}
}
