package main

import (
	"flag"
	"log"

	"github.com/jenildhorajiya/cqrs-microservices/api_gateway_service/config"
	"github.com/jenildhorajiya/cqrs-microservices/api_gateway_service/internal/server"
	"github.com/jenildhorajiya/cqrs-microservices/pkg/logger"
)

// @contact.name Jenil Dhorajiya
// @contact.url https://github.com/jenildhorajiya
// @contact.email jenild2002@gmail.com
func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ApiGateway")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
