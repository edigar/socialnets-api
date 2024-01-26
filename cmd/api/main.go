package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/edigar/socialnets-api/internal/config"
	"github.com/edigar/socialnets-api/internal/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.Load()
	r := router.Generate()

	server := &http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: r}
	fmt.Printf("SocialNets API is running on port %d...\n", config.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Stopping...")

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Server stopped")
}
