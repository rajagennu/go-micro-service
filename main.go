package main

import (
	"context"
	"fmt"
	"go-micro-service/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// create new object for hello handler
	_ = handlers.NewHello(l)

	// new handler for goodbye
	gb := handlers.NewGoodBye(l)

	// product handler
	ph := handlers.NewProducts(l)

	// create a serveMux
	mux := http.NewServeMux()
	mux.Handle("/", ph)
	mux.Handle("/goodbye", gb)

	// create the server object

	s := http.Server{
		Addr:         ":9000",
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// create a go routine to start the server in async
	go func() {
		s.ListenAndServe()
	}()

	// implement graceful shutdown

	notifyChannel := make(chan os.Signal, 1)
	signal.Notify(notifyChannel, os.Interrupt)
	signal.Notify(notifyChannel, os.Kill)

	// now block the channel, so that it would listen
	sig := <-notifyChannel
	fmt.Println("Received shutdown/interrupt signal", sig)

	// lets create the context for graceful shtudown
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
