package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	chgo "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/joho/godotenv"
)

const migrationsDir = "migrations/clickhouse"

func main() {
	_ = godotenv.Load()

	addr := envOr("CLICKHOUSE_ADDR", "clickhouse:9000")
	db := envOr("CLICKHOUSE_DB", "realtimedb")
	user := envOr("CLICKHOUSE_USER", "default")
	pass := envOr("CLICKHOUSE_PASSWORD", "")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := chgo.Open(&chgo.Options{
		Addr: []string{addr},
		Auth: chgo.Auth{Database: db, Username: user, Password: pass},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalf("clickhouse open: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("clickhouse ping: %v", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("read migrations dir: %v", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		path := filepath.Join(migrationsDir, f)
		body, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read %s: %v", path, err)
		}
		for _, stmt := range splitStatements(string(body)) {
			if strings.TrimSpace(stmt) == "" {
				continue
			}
			if err := conn.Exec(ctx, stmt); err != nil {
				log.Fatalf("exec %s: %v\nstatement:\n%s", f, err, stmt)
			}
		}
		log.Printf("applied %s", f)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// splitStatements splits SQL by ';' while ignoring comment lines so the migrate
// step can apply files containing multiple CREATE statements.
func splitStatements(s string) []string {
	var out []string
	var cur strings.Builder
	for _, line := range strings.Split(s, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "--") {
			continue
		}
		cur.WriteString(line)
		cur.WriteString("\n")
		if strings.HasSuffix(trim, ";") {
			out = append(out, cur.String())
			cur.Reset()
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
