package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/api"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/worker"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := queue.Connect()

	if err != nil {
		slog.Error("Error connecting to Redis", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to Redis")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker.StartWorker(ctx, i, &wg)
	}

	http.HandleFunc("/jobs", api.HandleCreateJob)
	http.HandleFunc("/metrics", api.GetMetrics)
	http.HandleFunc("/status", api.GetJobStatus)

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		slog.Info("Server running", "port", 8080)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("HTTP Server error", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	slog.Info("Shutdown signal recieved. Starting graceful exit...")

	shutDownContext, shutDownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutDownCancel()

	if err:=server.Shutdown(shutDownContext);err!=nil{
		slog.Error("API server forced to close","error",err)
	}else{
		slog.Info("API server closed gracefully.")
	}

	slog.Info("Signaling workers to stop..")

	cancel()

	wg.Wait()
	if err := queue.RDB.Close(); err != nil {
		slog.Error("Error closing Redis connection", "error", err)
	}
	slog.Info("System stopped gracefully!")

}
