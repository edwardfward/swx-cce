package main

import (
	"../pkg/cce"
	"log"
)

func main() {

	s, err := cce.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	s.StartServer()

}
