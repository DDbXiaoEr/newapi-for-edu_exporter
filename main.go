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

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func printHelp() {
	fmt.Printf(`newapi-exporter %s (commit: %s, built: %s)

Usage:
  newapi-exporter [flags]

Flags:
  -config string    Path to YAML configuration file (default: auto-detect config.yaml/config.yml)
  -out string       Output directory for file export (default: "output")
  -format string    Output format: json or csv (default: json)
  -serve string     Start Prometheus HTTP server on :port
  -version          Print version and exit
  -help             Print help and exit

Database configuration is specified via config file or environment variables.
See AGENTS.md for full documentation.
`, Version, Commit, BuildTime)
}

func main() {
	configPath := flag.String("config", "", "path to YAML configuration file")
	outputDir := flag.String("out", "output", "output directory")
	outputFormat := flag.String("format", "", "output format (json/csv)")
	servePort := flag.String("serve", "", "start prometheus HTTP server on :port")
	version := flag.Bool("version", false, "print version and exit")
	help := flag.Bool("help", false, "print help and exit")

	flag.Usage = printHelp
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *version {
		fmt.Printf("newapi-exporter %s (commit: %s, built: %s)\n", Version, Commit, BuildTime)
		return
	}

	cfg := config.LoadWithConfigFile(*configPath)

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
	if database.HasLogs() {
		tables["logs"] = func() (interface{}, error) { return exp.ExportLogs() }
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
	collector := metrics.NewCollector(database)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())

	addr := ":" + port
	fmt.Printf("serving metrics on %s/metrics\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}