package main

import (
	"log"

	"github.com/dabfleming/shorty/cmd/shorty/server"
	"github.com/dabfleming/shorty/internal/datastore"
	"github.com/dabfleming/shorty/internal/platform/mysql"
	"github.com/ua-parser/uap-go/uaparser"
)

func main() {
	// Connect to DB
	db, err := mysql.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Instantiate datastore
	ds, err := datastore.New(db)
	if err != nil {
		log.Fatalf("Error creating datastore: %v", err)
	}

	// User-Agent Parser
	parser := uaparser.NewFromSaved()

	// Start Server
	s, err := server.New(ds, parser)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	s.Go()
}
