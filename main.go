package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/calvine/goauth/core/jwt"
	"github.com/calvine/goauth/core/models"
	"github.com/calvine/goauth/core/nullable"
	"github.com/calvine/goauth/core/utilities"
	"github.com/calvine/goauth/dataaccess/memory"
	gamongo "github.com/calvine/goauth/dataaccess/mongo"
	gahttp "github.com/calvine/goauth/http"
	"github.com/calvine/goauth/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"go.uber.org/zap"
)

const (
	ENV_MONGO_CONNECTION_STRING = "GOAUTH_MONGO_CONNECTION_STRING"
	ENV_HTTP_ADDRESS_STRING     = "GOAUTH_HTTP_PORT_STRING"

	DEFAULT_MONGO_CONNECTION_STRING = "mongodb://root:password@localhost:27017/?authSource=admin&ssl=false&replicaSet=goauth_test&connect=direct"
	DEFAULT_HTTP_PORT_STRING        = "0.0.0.0:8080"
)

var (
	//go:embed static/*
	staticFS embed.FS
	//go:embed http/templates/*
	templateFS embed.FS
)

// https://opentelemetry.io/docs/go/getting-started/

func setupTelemetry(ctx context.Context) (func(), error) {
	tracePusher, traceCleanup, err := setupTracing(ctx)
	if err != nil {
		return nil, err
	}
	metricsPusher, metricsCleanup, err := setupMetrics(ctx)
	if err != nil {
		return nil, err
	}

	otel.SetTracerProvider(tracePusher)
	global.SetMeterProvider(metricsPusher.MeterProvider())

	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	return func() {
		traceCleanup()
		metricsCleanup()
	}, nil

}

func setupTracing(ctx context.Context) (*sdktrace.TracerProvider, func(), error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(
			resource.NewWithAttributes("com.calvinechols.goauth/resources", semconv.ServiceNameKey.String("goauth")),
		),
	)

	return tp, func() { _ = tp.Shutdown(ctx) }, nil
}

func setupMetrics(ctx context.Context) (*controller.Controller, func(), error) {
	metricExporter, err := stdoutmetric.New(
		stdoutmetric.WithPrettyPrint(),
	)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	return pusher, func() { _ = pusher.Stop(ctx) }, nil
}

func main() {
	ctx := context.Background()
	telemetryCleanup, err := setupTelemetry(ctx)
	if err != nil {
		log.Fatalf("failed to start telemetry: %s", err.Error())
	}
	defer telemetryCleanup()
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
	logger = logger.With(zap.String("app_name", "goauth"))
	defer logger.Sync()
	connectionString := utilities.GetEnv(ENV_MONGO_CONNECTION_STRING, DEFAULT_MONGO_CONNECTION_STRING)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	defer client.Disconnect(context.TODO())
	if err != nil {
		fmt.Printf("failed to connect to mongo server: %s\n", err.Error())
	}
	userRepo := gamongo.NewUserRepo(client)
	auditRepo := gamongo.NewAuditLogRepo(client)
	jsmRepo := memory.NewMemoryJWTSigningMaterialRepo()
	appRepo := memory.NewMemoryAppRepo()
	tokenRepo := memory.NewMemoryTokenRepo()

	tokenService := service.NewTokenService(tokenRepo)

	// TODO: set this up from configuration
	fsEmailServiceOptions := service.FSEmailServiceOptions{
		MessageDir: path.Join(".", "test_emails"),
	}
	emailService, err := service.NewEmailService(service.FSEmailService, fsEmailServiceOptions)
	if err != nil {
		return err
	}
	appService := service.NewAppService(appRepo, auditRepo)
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
	userService := service.NewUserService(userRepo, userRepo, tokenService, emailService)
	jsmService := service.NewJWTSigningMaterialService(jsmRepo)
	cachedJSMService := service.NewCachedJWTSigningMaterialService(jsmService, time.Minute*15)
	httpStaticFS := http.FS(staticFS)

	// add a jwt signing material to the repo so we can test...
	// TODO: later we will not add this per run, perhaps another executable to generate jwt signing material?
	testJSM := models.NewHMACJWTSigningMaterial("secret123", nullable.NullableTime{HasValue: false})
	err = cachedJSMService.AddJWTSigningMaterial(context.TODO(), logger, &testJSM, "startup")
	if err != nil {
		return err
	}

	httpServerOptions := gahttp.HTTPServerOptions{
		Logger:                     logger,
		LoginService:               loginService,
		UserService:                userService,
		EmailService:               emailService,
		TokenService:               tokenService,
		AppService:                 appService,
		JsmService:                 cachedJSMService,
		StaticFS:                   &httpStaticFS,
		TemplateFS:                 &templateFS,
		TokenSigningAlgorithmTypes: []jwt.JWTSingingAlgorithmFamily{jwt.HMAC}, // TODO: implment other signers of other types
	}
	httpServer, err := gahttp.NewServer(context.Background(), httpServerOptions)
	if err != nil {
		logger.Error("failed to create the http server", zap.Any("error", err))
		return err
	}
	httpServer.BuildRoutes()
	address := utilities.GetEnv(ENV_HTTP_ADDRESS_STRING, DEFAULT_HTTP_PORT_STRING)
	fmt.Printf("running http services on: %s", address)
	return http.ListenAndServe(address, &httpServer)
}
