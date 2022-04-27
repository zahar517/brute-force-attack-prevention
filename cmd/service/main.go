package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/zahar517/brute-force-attack-prevention/internal/app"
	"github.com/zahar517/brute-force-attack-prevention/internal/config"
	"github.com/zahar517/brute-force-attack-prevention/internal/logger"
	"github.com/zahar517/brute-force-attack-prevention/internal/server"
	"github.com/zahar517/brute-force-attack-prevention/internal/storage"
)

var configFile string

func init() {
	flag.StringVarP(&configFile, "config", "c", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		log.Fatal(err)
	}

	s := storage.New(config.Database.Dsn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	if err := s.Connect(ctx); err != nil {
		cancel()
		logg.Error(err.Error())
		return
	}
	defer cancel()

	a := app.New(logg, s)

	server := grpcserver.NewServer(logg, a, config.Server.GrpcHost, config.Server.GrpcPort)

	ctxNotify, cancelNotify := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancelNotify()

	go func() {
		<-ctxNotify.Done()

		if err := server.Stop(); err != nil {
			logg.Error("failed to stop server: " + err.Error())
		}
	}()

	logg.Info("server is running...")

	if err := server.Start(); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
