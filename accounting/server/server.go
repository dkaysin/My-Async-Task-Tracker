package main

import (
	"async_course/accounting"
	database "async_course/accounting/internal/database"
	reader "async_course/accounting/internal/event_reader"
	writer "async_course/accounting/internal/event_writer"
	httpAPI "async_course/accounting/internal/http_api"
	service "async_course/accounting/internal/service"
	schema "async_course/schema_registry"

	"log/slog"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// env vars
const (
	listenAddressEnvVar  = "LISTEN_ADDRESS"
	defaultListenAddress = ":4080"

	kafkaBrokersEnvVar        = "KAFKA_BROKERS"
	defaultKafkaBrokersEnvVar = "localhost:9092"

	pgConnStringEnvVar        = "PG_CONN_STRING"
	defaultPgConnStringEnvVar = "postgres://postgres:postgres@127.0.0.1:5432/postgres"

	signingKeyEnvVar = "SIGNING_KEY"
)

func main() {

	// set config
	config := viper.New()
	config.SetEnvPrefix("ACCOUNTING")
	config.AutomaticEnv()
	config.SetDefault(listenAddressEnvVar, defaultListenAddress)
	config.SetDefault(kafkaBrokersEnvVar, defaultKafkaBrokersEnvVar)
	config.SetDefault(pgConnStringEnvVar, defaultPgConnStringEnvVar)

	// set signing key
	signingKey := config.GetString(signingKeyEnvVar)
	if signingKey == "" {
		slog.Error("signing key not provided")
		os.Exit(1)
	}

	// set database
	db, err := database.NewDatabase(config.GetString(pgConnStringEnvVar))
	if err != nil {
		slog.Error("fatal error while initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// set schema registry
	sr := schema.NewSchemaRegistry(accounting.ProducerName)

	// set event writer
	brokers := strings.Split(config.GetString(kafkaBrokersEnvVar), ",")
	ew := writer.NewEventWriter(brokers, sr)
	defer ew.Close()

	// set service
	s := service.NewService(config, db, ew, signingKey)

	// set event reader
	er := reader.NewEventReader(s, sr)
	er.StartReaders(brokers, accounting.KafkaConsumerGroupID)

	// set http handler
	h := httpAPI.NewHttpAPI(config, s)

	// set server and API
	e := echo.New()

	public := e.Group("")
	h.RegisterPublic(public)

	api := e.Group("/api")
	// parse jwt token into "user" context key
	api.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(accounting.JwtCustomClaims)
		},
		ErrorHandler: httpAPI.JwtMiddlewareErrorHandler,
		SigningKey:   []byte(signingKey),
	}))
	h.RegisterAPI(api)

	// set echo logger
	e.Logger.SetPrefix("main")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
			`"error_code":"${header:x-hoop-error-code}"}` +
			"\n",
	}))

	// start server
	e.Logger.Fatal(e.Start(config.GetString(listenAddressEnvVar)))
}
