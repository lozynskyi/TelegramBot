package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/lozynskyi/TelegramBot/internal/app/telegrambot"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/main.toml", "Path to config file.")
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flag.Parse()

	config := telegrambot.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	t := telegrambot.New(config)
	if err := t.Start(); err != nil {
		log.Fatal(err)
	}
}
