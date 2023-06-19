package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"id-maker/config"
	"id-maker/internal/app"
	"log"
)

func main() {
	var cfg config.Config
	err := cleanenv.ReadConfig("./config/config.yml", &cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(&cfg)

}
