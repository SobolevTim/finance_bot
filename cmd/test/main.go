package main

import (
	"fmt"
	"log"

	"github.com/SobolevTim/finance_bot/config"
)

func main() {
	// Подключаем конфигурацию
	config, err := config.LoadConfig("config")
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}
	fmt.Println("Config loaded:", config)
}
