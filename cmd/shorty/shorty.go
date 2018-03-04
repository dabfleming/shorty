package main

import (
	"log"

	"github.com/dabfleming/shorty/cmd/shorty/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	s.Go()
}
