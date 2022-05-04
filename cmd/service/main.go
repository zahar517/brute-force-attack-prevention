package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	flag "github.com/spf13/pflag"
	"github.com/zahar517/brute-force-attack-prevention/internal/app"
	"github.com/zahar517/brute-force-attack-prevention/internal/limiter"
	"github.com/zahar517/brute-force-attack-prevention/internal/logger"
	"github.com/zahar517/brute-force-attack-prevention/internal/server"
	"github.com/zahar517/brute-force-attack-prevention/internal/storage"
)

const interval = 60

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logLevel := os.Getenv("LOGG_LEVEL")
	logFile := os.Getenv("LOG_FILE")

	logg, err := logger.New(logLevel, logFile)
	if err != nil {
		log.Fatal(err)
	}

	loginLimit, err := strconv.ParseInt(os.Getenv("LOGIN_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	passLimit, err := strconv.ParseInt(os.Getenv("PASSWORD_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	ipLimit, err := strconv.ParseInt(os.Getenv("IP_LIMIT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	l := limiter.New(loginLimit, passLimit, ipLimit, interval)

	dsn := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	s := storage.New(dsn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	if err := s.Connect(ctx); err != nil {
		cancel()
		logg.Error(err.Error())
		return
	}
	defer cancel()

	a := app.New(logg, s, l)

	grpcHost := os.Getenv("GRPC_HOST")
	grpcPort := os.Getenv("GRPC_PORT")

	server := grpcserver.NewServer(logg, a, grpcHost, grpcPort)

	ctxNotify, cancelNotify := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancelNotify()

	go func() {
		<-ctxNotify.Done()

		l.Stop()
		if err := server.Stop(); err != nil {
			logg.Error("failed to stop server: " + err.Error())
		}
	}()

	l.Start()

	logg.Info("server is running...")

	if err := server.Start(); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
