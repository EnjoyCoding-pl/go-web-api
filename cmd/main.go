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
	"go-web-api/features/trip/domain"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	address string
	timeout time.Duration
)

func main() {

	flag.StringVar(&address, "address", "127.0.0.1:4400", "Address on which server will listen")
	flag.DurationVar(&timeout, "timeout", 30, "Seconds after which request will be cancelled")

	dsn := "host=127.0.0.1 user=postgres password=postgres dbname=go-web-app port=5432 sslmode=disable"

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	conn.AutoMigrate(domain.NewTrip())

	r := mux.NewRouter()

	tsrv := &trip.TripService{
		Db: conn,
	}

	tsrv.MuxRegister(r)

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

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("pong"))
}
