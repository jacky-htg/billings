package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jacky-htg/billings/pkg/config"
	"github.com/jacky-htg/billings/pkg/database"
	"github.com/jacky-htg/billings/route"
)

func main() {
	if _, ok := os.LookupEnv("APP_ENV"); !ok {
		config.Setup(".env")
	}

	log := log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	db, err := database.GetConnection()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	log.Fatal(http.ListenAndServe(":8080", route.InitRoute(db, log)))
}
