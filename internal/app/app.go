package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"id-maker/config"
	v1 "id-maker/internal/controller/http/v1"
	"id-maker/internal/controller/rpc"
	"id-maker/internal/usecase"
	"id-maker/internal/usecase/repo"
	"id-maker/pkg/grpcserver"
	"id-maker/pkg/httpserver"
	"id-maker/pkg/logger"
	"id-maker/pkg/mysql"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	mysql, err := mysql.New(
		cfg.MySQL.URL,
		mysql.MaxIdleConns(cfg.MaxIdleConns),
		mysql.MaxOpenConns(cfg.MaxOpenConns),
	)

	if err != nil {
		l.Fatal(fmt.Errorf("app - run - mysql.New: %w", err))
	}

	defer mysql.Close()

	segmentUseCase := usecase.New(repo.New(mysql))

	handler := gin.New()

	v1.NewRouter(handler, l, segmentUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	grpcServer := grpcserver.New(grpcserver.Port(cfg.GRPC.Port))
	rpc.NewRouter(segmentUseCase, l)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	case err = <-grpcServer.Notify():
		l.Error(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	grpcServer.Shutdown()
}
