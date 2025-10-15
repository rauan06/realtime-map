package main

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rauan06/realtime-map/go-commons/pkg/logger"
)

const (
	forceArgCount = 3
	mainArgCount  = 2
)

//nolint:cyclop // Scan method requires multiple type checks for SQL driver compatibility
func main() {
	l := logger.New("DEBUG")

	if len(os.Args) < mainArgCount {
		l.Info("Usage: migrator [up|down|drop|force <version>|version]")
		os.Exit(1)
	}

	m, err := migrate.New(
		"file://migrations",
		"postgres://postgres:example@db:5432/realtimedb?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}

		l.Info("Migrations applied successfully")

	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}

		l.Info("Rolled back one migration")

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatal(err)
		}

		l.Info("Database dropped")

	case "force":
		if len(os.Args) < forceArgCount {
			log.Fatal("force requires version number")
		}

		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("invalid version number")
		}

		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}

		l.Info("Forced to version %d\n", version)

	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}

		l.Info("Current version: %d (dirty: %t)\n", v, dirty)

	default:
		l.Info("Unknown command: %s", cmd)
		l.Info("Usage: migrator [up|down|drop|force <version>|version]")
		os.Exit(1)
	}
}
