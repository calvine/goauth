package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/goauth/dataaccess/memory"
	gamongo "github.com/calvine/goauth/dataaccess/mongo"
	gahttp "github.com/calvine/goauth/http"
	"github.com/calvine/goauth/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ENV_MONGO_CONNECTION_STRING = "GOAUTH_MONGO_CONNECTION_STRING"
	ENV_HTTP_ADDRESS_STRING     = "GOAUTH_HTTP_PORT_STRING"

	DEFAULT_MONGO_CONNECTION_STRING = "mongodb://root:password@localhost:27017/?authSource=admin&readPreference=primary&ssl=false"
	DEFAULT_HTTP_PORT_STRING        = ":8080"
)

var (
	//go:embed static/*
	staticFS embed.FS
	//go:embed http/templates/*
	templateFS embed.FS
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("an error occurred while starting the http server: %s", err.Error())
	}
}

func run() error {
	connectionString := utilities.GetEnv(ENV_MONGO_CONNECTION_STRING, DEFAULT_MONGO_CONNECTION_STRING)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	defer client.Disconnect(context.TODO())
	if err != nil {
		fmt.Printf("failed to connect to mongo server: %s\n", err.Error())
	}
	userRepo := gamongo.NewUserRepo(client)
	auditRepo := gamongo.NewAuditLogRepo(client)
	tokenRepo := memory.NewLocalTokenRepo()

	tokenService := service.NewTokenService(tokenRepo)
	emailService, err := service.NewEmailService(service.MockEmailService, nil)
	if err != nil {
		return err
	}
	loginService := service.NewLoginService(auditRepo, userRepo, userRepo, emailService, tokenService)

	httpStaticFS := http.FS(staticFS)
	httpServer := gahttp.NewServer(loginService, emailService, &httpStaticFS, &templateFS)
	httpServer.BuildRoutes()
	address := utilities.GetEnv(ENV_HTTP_ADDRESS_STRING, DEFAULT_HTTP_PORT_STRING)
	return http.ListenAndServe(address, &httpServer)
}
