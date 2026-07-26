package main

import (
	"flag"
	"log"

	"github.com/jenildhorajiya/cqrs-microservices/pkg/logger"
	"github.com/jenildhorajiya/cqrs-microservices/writer_service/config"
	"github.com/jenildhorajiya/cqrs-microservices/writer_service/internal/server"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("WriterService")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
