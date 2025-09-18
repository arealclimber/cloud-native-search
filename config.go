package main

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 預設值
	}
	log.Printf("Loaded config: port=%s", port)
	return &Config{Port: ":" + port}
}
