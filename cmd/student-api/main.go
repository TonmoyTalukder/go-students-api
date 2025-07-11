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
	"github.com/TonmoyTalukder/go-students-api/internal/http/handlers/student"
	"github.com/TonmoyTalukder/go-students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// logger if any
	slog.Info("Student API: Go REST")

	// database set up
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /api/", func(w http.ResponseWriter,  r *http.Request){
		w.Write([]byte("Welcome to the server..."))
	})

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	router.HandleFunc("PUT /api/students/{id}", student.UpdateById(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))


	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started at", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server.")
		}
	}()

	<-done

	slog.Info("Shutting down the server.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully.")
}
