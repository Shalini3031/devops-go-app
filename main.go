package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// --------------------
// Prometheus Metrics
// --------------------
var httpRequestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
)

func main() {

	// --------------------
	// Read Environment Variables
	// --------------------
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("One or more required environment variables are missing")
	}

	// --------------------
	// Database Connection (with retry)
	// --------------------
	dsn := "host=" + dbHost +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" sslmode=disable"

	var db *sql.DB
	var err error

	for i := 1; i <= 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			log.Println("Successfully connected to DB")
			break
		}

		log.Printf("DB not ready (attempt %d/10): %v", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to DB after retries:", err)
	}

	// --------------------
	// Register Prometheus Metrics
	// --------------------
	prometheus.MustRegister(httpRequestsTotal)

	// --------------------
	// HTTP Handlers
	// --------------------

	// Health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		httpRequestsTotal.Inc()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// --------------------
	// Start HTTP Server
	// --------------------
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

