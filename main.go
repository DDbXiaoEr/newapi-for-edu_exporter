package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/quanxiaoer/newapi_exporter/config"
	"github.com/quanxiaoer/newapi_exporter/db"
	"github.com/quanxiaoer/newapi_exporter/exporter"
	"github.com/quanxiaoer/newapi_exporter/metrics"
)

func main() {
	configPath := flag.String("config", "", "path to YAML configuration file")
	dbType := flag.String("db-type", "", "database type (sqlite/mysql)")
	dbPath := flag.String("db", "", "path to SQLite database file")
	mysqlHost := flag.String("mysql-host", "", "MySQL host")
	mysqlPort := flag.String("mysql-port", "", "MySQL port")
	mysqlUser := flag.String("mysql-user", "", "MySQL user")
	mysqlPassword := flag.String("mysql-password", "", "MySQL password")
	mysqlDatabase := flag.String("mysql-database", "", "MySQL database name")
	outputDir := flag.String("out", "output", "output directory")
	outputFormat := flag.String("format", "", "output format (json/csv)")
	servePort := flag.String("serve", "", "start prometheus HTTP server on :port")
	flag.Parse()

	cfg := config.LoadWithConfigFile(*configPath)

	if *dbType != "" {
		cfg.DBType = *dbType
	}
	if *dbPath != "" {
		cfg.DatabasePath = *dbPath
	}
	if *mysqlHost != "" {
		cfg.MySQLHost = *mysqlHost
	}
	if *mysqlPort != "" {
		cfg.MySQLPort = *mysqlPort
	}
	if *mysqlUser != "" {
		cfg.MySQLUser = *mysqlUser
	}
	if *mysqlPassword != "" {
		cfg.MySQLPassword = *mysqlPassword
	}
	if *mysqlDatabase != "" {
		cfg.MySQLDatabase = *mysqlDatabase
	}
	if *outputFormat != "" {
		cfg.OutputFormat = *outputFormat
	}

	database, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	if *servePort != "" {
		runServer(database, *servePort)
		return
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	exp := exporter.New(database, cfg)

	tables := map[string]func() (interface{}, error){
		"statistics": func() (interface{}, error) { return exp.ExportStatistics() },
		"quota_data": func() (interface{}, error) { return exp.ExportQuotaData() },
		"tokens":     func() (interface{}, error) { return exp.ExportTokens() },
		"users":      func() (interface{}, error) { return exp.ExportUsers() },
	}

	for name, exportFn := range tables {
		data, err := exportFn()
		if err != nil {
			log.Printf("warning: failed to export %s: %v", name, err)
			continue
		}

		path := fmt.Sprintf("%s/%s.%s", *outputDir, name, cfg.OutputFormat)
		if err := exp.WriteOutput(data, path, cfg.OutputFormat); err != nil {
			log.Printf("warning: failed to write %s: %v", name, err)
			continue
		}
		fmt.Printf("exported: %s\n", path)
	}
}

func runServer(database *db.Database, port string) {
	collector := metrics.NewCollector(database.Conn())
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())

	addr := ":" + port
	fmt.Printf("serving metrics on %s/metrics\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}