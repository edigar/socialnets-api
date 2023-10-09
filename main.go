package main

import (
	"fmt"
	"github.com/edigar/socialnets-api/src/config"
	"github.com/edigar/socialnets-api/src/router"
	"log"
	"net/http"
)

func main() {
	config.Load()
	r := router.Generate()

	fmt.Printf("SocialNets API is running on port %d...\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
