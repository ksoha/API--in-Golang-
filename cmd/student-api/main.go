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

	"github.com/ksoha/API-in-Golang/internal/config"
	"github.com/ksoha/API-in-Golang/internal/http/handlers/student"
	"github.com/ksoha/API-in-Golang/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to connect to the datbase", err)
	}

	slog.Info("datbase connected successfully, storage initialized", slog.String("env", cfg.Env))

	//setup routes
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	//server setup

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))

	//creating a chanel
	//the chanel will store the signal used by the operating system
	done := make(chan os.Signal, 1)

	//to get the signal inside the chanel we use the signal package
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) //this will listen to the interrupt signal from the OS

	//to gracefully shutdown the sever while in prod we listen the server in a go routine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start the server")
		}
	}()

	<-done //this will block the main tread until we recieve a signal from the OS

	//logic to stop the server gracefully

	slog.Info("Server shutting down")

	//context is used to set a timeout for the shutdown process so that it does not take too long to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //this will create a context with a timeout of 5 seconds
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
