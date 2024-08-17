package main

import (
	"context"
	"log"
	"main/config"
	_ "main/docs"
	ticketServices "main/internal/v1/ticket/application/services"
	ticketHandler "main/internal/v1/ticket/handler/http"
	ticketRepos "main/internal/v1/ticket/infrastructure/repository"
	"main/pkg/databases/postgresql"
	"main/pkg/logger"
	"main/pkg/server"
	"main/pkg/utils/common"
	"main/pkg/utils/graceful_exit"
	"os"
)

// @title Ticket API
// @version 1.0
// @description Ticket service broker with REST endpoints
// @contact.email semerci394@gmail.com
// @BasePath /v1
func main() {
	common.InitializeI18N()

	os.Setenv("APP_ENV", "dev")

	cfg, errConfig := config.ParseConfig()
	if errConfig != nil {
		log.Fatal(errConfig)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()

	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.APP_VERSION, cfg.Logger.LEVEL, cfg.Server.APP_ENV)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init Clients
	postgresqlDB, err := postgresql.NewPostgresqlDB(cfg)
	if err != nil {
		appLogger.Fatal("Error when tyring to connect to Postgresql")
	} else {
		appLogger.Info("Postgresql connected")
	}

	// Init repositories
	pgRepoTicket := ticketRepos.NewPostgresqlRepository(ctx, postgresqlDB, appLogger)

	// Init services
	svcTicket := ticketServices.NewTicketService(cfg, pgRepoTicket, appLogger)

	// Interceptors
	//

	servers := server.NewServer(cfg, &ctx, appLogger)

	httpServer, errHttpServer := servers.NewHttpServer()
	if errHttpServer != nil {
		println(errHttpServer.Error())
	}
	versioning := httpServer.Group(cfg.Server.API_VER)

	// Init handlers for HTTP Server
	hndlTicket := ticketHandler.NewHttpHandler(ctx, cfg, svcTicket, appLogger)

	// Init routes for HTTP Server
	ticketHandler.MapRoutes(hndlTicket, versioning)

	servers.Listen(httpServer)

	// Exit from application gracefully
	graceful_exit.TerminateApp(ctx)

	appLogger.Info("Server Exited Properly")
}
