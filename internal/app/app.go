package app

import (
	"fmt"
	"id-maker/config"
	"id-maker/pkg/logger"
	"id-maker/pkg/mysql"
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
}
