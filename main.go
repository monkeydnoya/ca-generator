package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/trace"
	"syscall"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("GOMAXPROCS=%d\n", runtime.NumCPU())

	addr := "0.0.0.0:7000"

	slog.SetLogLoggerLevel(slog.LevelDebug)

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
