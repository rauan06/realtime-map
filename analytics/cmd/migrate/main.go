package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrator [up|down|drop|force <version>|version]")
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
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Rolled back one migration")

	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Database dropped")

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("force requires version number")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("invalid version number")
		}
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Forced to version %d\n", version)

	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Current version: %d (dirty: %t)\n", v, dirty)

	default:
		fmt.Println("Unknown command:", cmd)
		fmt.Println("Usage: migrator [up|down|drop|force <version>|version]")
		os.Exit(1)
	}
}
