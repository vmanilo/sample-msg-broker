package main

import (
	"fmt"
	"sample-msg-broker/src/broker"
)

func main() {

	srv := broker.NewServer(9123)
	err := srv.Start()
	if err != nil {
		fmt.Println("ERROR:", err)
	}

}
