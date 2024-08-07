package main

import (
	"github.com/venusforest2013/config/application"
	"log"
)

func main() {

	if err := application.New().Start(); err != nil {
		log.Fatal(err)
		return
	}
}
