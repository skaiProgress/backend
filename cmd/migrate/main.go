package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"aiqadam-backend/internal/config"
	appmigrate "aiqadam-backend/internal/migrate"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	command := flag.String("command", "", "migration command: up, down, status")
	flag.Parse()

	if *command == "" {
		log.Fatal("missing -command flag (up, down, status)")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	migrationsPath := appmigrate.ResolveMigrationsPath()
	runner, err := appmigrate.NewRunner(cfg.DatabaseURL, migrationsPath)
	if err != nil {
		log.Fatalf("init migrator: %v", err)
	}
	defer func() {
		if err := runner.Close(); err != nil {
			log.Printf("close migrator: %v", err)
		}
	}()

	switch *command {
	case "up":
		if err := runner.Up(); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrations applied successfully")
	case "down":
		if err := runner.Down(); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("migration rolled back successfully")
	case "status":
		version, dirty, err := runner.Status()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				fmt.Println("status: no migrations applied")
				return
			}
			log.Fatalf("migrate status: %v", err)
		}
		fmt.Printf("version: %d\n", version)
		fmt.Printf("dirty: %v\n", dirty)
	default:
		log.Fatalf("unknown command: %s (use up, down, status)", *command)
	}
}
