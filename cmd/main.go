package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	"time"

	"go-web-api/features/trip"
	"go-web-api/features/trip/domain/models"
	"go-web-api/features/trip/infra/storages"
	"go-web-api/features/user"
	user_models "go-web-api/features/user/domain/models"
	user_storage "go-web-api/features/user/infra/storages"
	"go-web-api/internal/middlewares"
	"go-web-api/internal/providers"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var (
	address          string
	jaegerUrl        string
	timeout          time.Duration
	jwtIssuer        string
	jwtSecret        string
	connectionString string
)

func main() {

	flag.StringVar(&address, "address", "127.0.0.1:4400", "Address on which server will listen")
	flag.StringVar(&jaegerUrl, "jaeger", "http://localhost:14268/api/traces", "Jaeger url")
	flag.StringVar(&jwtIssuer, "jwt-issuer", "http://127.0.0.1:4400", "JWT issuer value")
	flag.StringVar(&jwtSecret, "jwt-secret", "default", "Secret for signing JWT token")
	flag.StringVar(&connectionString, "connection-string", "host=127.0.0.1 user=postgres password=postgres dbname=go-web-app port=5432 sslmode=disable", "Connection string to postgres database")
	flag.DurationVar(&timeout, "timeout", 30, "Seconds after which request will be cancelled")

	conn, err := gorm.Open(postgres.Open(connectionString))

	if err != nil {
		panic("Failed to connect database")
	}

	tp, err := providers.CreateJaegerProvider(jaegerUrl)

	if err != nil {
		panic("Failed to connect jaeger")
	}
	otel.SetTracerProvider(tp)

	conn.AutoMigrate(&models.Trip{}, &models.TripPoint{}, &user_models.User{})

	r := mux.NewRouter()

	jwtProvider := providers.NewJwtProvider(jwtIssuer, jwtSecret)

	publicPaths := map[string]any{
		"/login":    true,
		"/register": true,
	}

	r.Use(middlewares.AuthMiddleware(jwtProvider, publicPaths))

	t := trip.NewTripController(storages.NewPostgresStorage(conn))

	u := user.NewUserController(user_storage.NewUserPostgresStorage(conn), jwtProvider)

	t.MuxRegister(r)

	u.MuxRegister(r)

	srv := http.Server{
		ReadTimeout:  time.Second * timeout,
		WriteTimeout: time.Second * timeout * 2,
		IdleTimeout:  time.Second * timeout * 4,
		Addr:         address,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancell := context.WithTimeout(context.Background(), time.Second*timeout*4)
	defer cancell()

	srv.Shutdown(ctx)

	log.Info("Shutting down...")

	os.Exit(0)

}
