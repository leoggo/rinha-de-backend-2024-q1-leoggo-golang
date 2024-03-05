package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // driver do postgres
)

var (
	serverPort               = "8080"
	defaultShoutdownTimeout  = getEnv("DEFAULT_SHUTDOWN_TIME_OUT", 100)
	defaultReadTimeout       = getEnv("DEFAULT_READ_TIME_OUT", 120)
	defaultWriteTimeout      = getEnv("DEFAULT_WRITE_TIME_OUT", 120)
	defaultIdleTimeout       = getEnv("DEFAULT_IDLE_TIME_OUT", 150)
	defaultReadHeaderTimeout = getEnv("DEFAULT_READ_HEADER_TIME_OUT", 120)
	defaultMaxDBConnections  = getEnv("DEFAULT_MAX_BD_CONNECTIONS", 35)
)

var db *pgxpool.Pool

func getEnv(env string, defaultVal int) int {
	dado := os.Getenv(env)

	if dado == "" {
		return defaultVal
	}

	retval, err := strconv.Atoi(dado)

	if err != nil {
		panic(err)
	}

	return retval
}

func main() {
	psqlinfo := "postgres://user:senha@postgres:5432/clientes?sslmode=disable"

	dbConfig, err := pgxpool.ParseConfig(psqlinfo)

	if err != nil {
		panic(err)
	}

	dbConfig.MaxConns = int32(defaultMaxDBConnections)

	db, err = pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	serverMuxer := http.NewServeMux()

	buildRoutes(serverMuxer)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", serverPort),
		ReadTimeout:       time.Duration(defaultReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(defaultWriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(defaultIdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(defaultReadHeaderTimeout) * time.Second,
		Handler:           serverMuxer,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("problems with listen and serve: ", err)
		}
	}()

	fmt.Println("webserver ready to go!")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stopChan

	signal.Stop(stopChan)
	close(stopChan)

	fmt.Println("HTTP shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(defaultShoutdownTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("HTTP server Shutdown:", err)
	}

	fmt.Println("bye bye.")
}
