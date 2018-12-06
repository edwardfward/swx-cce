package main

import (
	"../pkg/app"
	"log"
)

func main() {

	s, err := app.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	s.StartServer()

}
