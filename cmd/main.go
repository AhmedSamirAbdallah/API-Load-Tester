package main

import (
	"api-load-tester/config"
	"api-load-tester/pkg/services"
)

func main() {
	cfg := config.LoadConfig()

	loadTestservice := services.NewLoadTesterService(cfg)

	loadTestservice.RunLoadTest()
}
