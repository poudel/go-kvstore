package main

import (
	"fmt"
	"kvstore/store"
)

func main() {
	fmt.Println("KV Store")
	s, err := store.NewServer("127.0.0.1", 1116)
	if err != nil {
		panic(err)
	}

	// Listen to the Port
	fmt.Println("Listening on 127.0.0.1:1116")
	panic(s.Listen())
}
