package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	// defaultAddr is the default address the HTTP proxy and admin interface listen on.
	// Changed from :8080 to 127.0.0.1:8080 to avoid binding on all interfaces by default.
	defaultAddr = "127.0.0.1:8080"
	// defaultAdminPath is the default path prefix for the admin interface.
	defaultAdminPath = "/hetty/"
	// defaultCertsSubdir is the subdirectory under the user config dir for certs.
	defaultCertsSubdir = "/hetty/certs"
	// defaultDBSubdir is the subdirectory under the user cache dir for the database.
	defaultDBSubdir = "/hetty/db"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	// Parse command-line flags.
	addr := flag.String("addr", defaultAddr, "TCP address to listen on (e.g. \"127.0.0.1:8080\")")
	adminPath := flag.String("adminPath", defaultAdminPath, "Path prefix for admin interface")
	certsDir := flag.String("certsDir", "", "Directory for storing CA certificate and key (defaults to system config dir)")
	dbPath := flag.String("db", "", "Path to database file (defaults to system data dir)")
	upstreamProxy := flag.String("upstreamProxy", "", "Optional upstream proxy URL (e.g. http://proxy:8080)")
	printVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *printVersion {
		fmt.Printf("hetty %v\n", version)
		os.Exit(0)
	}

	log.Printf("[INFO] Starting hetty %v", version)

	// Resolve default certs directory if not provided.
	if *certsDir == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalf("[FATAL] Could not determine user config directory: %v", err)
		}
		*certsDir = configDir + defaultCertsSubdir
	}

	// Resolve default database path if not provided.
	if *dbPath == "" {
		dataDir, err := os.UserCacheDir()
		if err != nil {
			log.Fatalf("[FATAL] Could not determine user cache directory: %v", err)
		}
		*dbPath = dataDir + defaultDBSubdir
	}

	log.Printf("[INFO] Using certs directory: %v", *certsDir)
	log.Printf("[INFO] Using database path: %v", *dbPath)

	if *upstreamProxy != "" {
		log.Printf("[INFO] Using upstream proxy: %v", *upstreamProxy)
	}

	log.Printf("[INFO] Listening on %v", *addr)
	log.Printf("[INFO] Admin interface available at http://%v%v", *addr, *adminPath)

	// Reminder: open the admin UI in the browser after startup for easier access.
	// e.g. http://127.0.0.1:8080/hetty/

	// TODO: Initialize proxy, database, and HTTP server.
	// This will be wired up as the project progresses.
	_ = adminPath
}
