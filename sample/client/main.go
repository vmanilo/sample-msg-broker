package main

import (
	"fmt"
	"os"
	"os/signal"
	"sample-msg-broker/src/client"
	"syscall"
	"time"
)

func main() {
	client := client.NewClient("localhost:9123")
	client.Connect()

	topic := "qwerty"

	client.Subscribe(topic)

	go func() {
		for msg := range client.GetSubscription() {
			fmt.Println("Received from subscription:", msg)
		}
	}()

	for _, msg := range []string{"one", "two", "three"} {
		time.Sleep(time.Second)
		client.Publish(topic, msg)
	}

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan
}
