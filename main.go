package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/goauth/dataaccess/memory"
	gamongo "github.com/calvine/goauth/dataaccess/mongo"
	gahttp "github.com/calvine/goauth/http"
	"github.com/calvine/goauth/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"go.uber.org/zap"
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

func setupTracing(ctx context.Context) (*sdktrace.TracerProvider, func(), error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp))

	return tp, func() { _ = tp.Shutdown(ctx) }, nil
}

func setupMetrics(ctx context.Context) (func(), error) {
	metricExporter, err := stdoutmetric.New(
		stdoutmetric.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(5*time.Second),
	)

	err = pusher.Start(ctx)
	if err != nil {
		return nil, err
	}

	return func() { _ = pusher.Stop(ctx) }, nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("an error occurred while starting the http server: %s", err.Error())
	}
}

func run() error {
	// TODO: add logger config
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}
	connectionString := utilities.GetEnv(ENV_MONGO_CONNECTION_STRING, DEFAULT_MONGO_CONNECTION_STRING)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	defer client.Disconnect(context.TODO())
	if err != nil {
		fmt.Printf("failed to connect to mongo server: %s\n", err.Error())
	}
	userRepo := gamongo.NewUserRepo(client)
	auditRepo := gamongo.NewAuditLogRepo(client)
	tokenRepo := memory.NewMemoryTokenRepo()

	tokenService := service.NewTokenService(tokenRepo)
	emailService, err := service.NewEmailService(service.MockEmailService, nil)
	if err != nil {
		return err
	}
	loginServiceOptions := service.LoginServiceOptions{
		AuditLogRepo:           auditRepo,
		UserRepo:               userRepo,
		ContactRepo:            userRepo,
		EmailService:           emailService,
		TokenService:           tokenService,
		MaxFailedLoginAttempts: 10,
		AccountLockoutDuration: time.Minute * 15,
	}
	loginService := service.NewLoginService(loginServiceOptions)

	httpStaticFS := http.FS(staticFS)
	httpServer := gahttp.NewServer(logger, loginService, emailService, tokenService, &httpStaticFS, &templateFS)
	httpServer.BuildRoutes()
	address := utilities.GetEnv(ENV_HTTP_ADDRESS_STRING, DEFAULT_HTTP_PORT_STRING)
	fmt.Printf("running http services on: %s", address)
	return http.ListenAndServe(address, &httpServer)
}
