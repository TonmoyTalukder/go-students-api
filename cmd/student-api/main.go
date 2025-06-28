package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TonmoyTalukder/go-students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// logger if any

	// database set up

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter,  r *http.Request){
		w.Write([]byte("Welcome to student api..."))
	})

	// setup server
	server := http.Server {
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started at", slog.String("address", cfg.Addr))
	// fmt.Printf("Server started at %s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server.")
		}
	}()

	<-done

	slog.Info("Shutting down the server.")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully.")
}