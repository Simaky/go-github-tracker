package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof" // side-effect: registers /debug/pprof on the default mux
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // side-effect: registers the "pgx" DB driver

	"github.com/Simaky/go-github-tracker/backend/app"
	"github.com/Simaky/go-github-tracker/backend/app/config"
	"github.com/Simaky/go-github-tracker/backend/app/storage"
	"github.com/Simaky/go-github-tracker/backend/server"
)

// appVersion is injected at link time via -ldflags.
var appVersion = "dev"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) == 2 && os.Args[1] == "-version" {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	cfg := config.LoadConfig()

	store, err := storage.New(context.Background(), cfg.Storage.DSN)
	if err != nil {
		log.Fatalf("opening storage: %s", err)
	}

	appInst := app.New(store)

	if err := server.New(appInst, cfg).Run(appVersion); err != nil {
		log.Fatalf("running server: %s", err)
	}
}
