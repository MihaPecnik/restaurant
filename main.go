package main

import (
	"github.com/MihaPecnik/restaurant/database"
	"github.com/MihaPecnik/restaurant/handler"
	"github.com/MihaPecnik/restaurant/server"
	"log"
)
func main() {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(db)
	s := server.NewServer(h)
	log.Fatal(s.Serve(":8080"))
}
