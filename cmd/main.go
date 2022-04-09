package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"time"

	"go-web-api/features/trip"
	"go-web-api/features/trip/domain/models"
	"go-web-api/features/trip/infra/storages"
	"go-web-api/internal/providers"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	address   string
	jaegerUrl string
	timeout   time.Duration
)

func main() {

	flag.StringVar(&address, "address", "127.0.0.1:4400", "Address on which server will listen")
	flag.StringVar(&jaegerUrl, "jaeger", "http://localhost:14268/api/traces", "Jaeger url")
	flag.DurationVar(&timeout, "timeout", 30, "Seconds after which request will be cancelled")

	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=go-web-app port=5432 sslmode=disable"

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{SlowThreshold: time.Second, LogLevel: logger.Info, IgnoreRecordNotFoundError: true, Colorful: true})

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	tp, err := providers.CreateJaegerProvider(jaegerUrl)

	if err != nil {
		panic("Failed to connect jaeger")
	}
	otel.SetTracerProvider(tp)

	if err != nil {
		panic("Failed to connect database")
	}

	conn.AutoMigrate(models.NewTrip(), models.NewTripPoint())

	r := mux.NewRouter()

	t := trip.NewTripController(storages.NewPostgresStorage(conn), *log.Default())

	t.MuxRegister(r)

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

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancell := context.WithTimeout(context.Background(), time.Second*timeout*4)
	defer cancell()

	srv.Shutdown(ctx)

	log.Println("Shutting down...")

	os.Exit(0)

}
