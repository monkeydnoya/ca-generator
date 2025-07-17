package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/joho/godotenv"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.MemProfileRate = 0 // Disable profiling overhead

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("API_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8000"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	var server *http.Server

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		signal := <-sigchan
		slog.Info("Received", "signal", signal)

		server.Close()
		os.Exit(0)
	}()

	gofakeit.Seed(0)

	server = &http.Server{Addr: addr, Handler: RegisterHandlers()}
	slog.Info("Listening on", "addr", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
