package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/chrishartserver/hartguard/host-watch/internal/alerter"
	"github.com/chrishartserver/hartguard/host-watch/internal/checker"
)

type ServiceURL struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Config struct {
	DiskPaths    []string     `yaml:"disk_paths"`
	DiskWarnPct  int          `yaml:"disk_warn_percent"`
	Containers   []string     `yaml:"containers"`
	ServiceURLs  []ServiceURL `yaml:"service_urls"`
	Interval     int          `yaml:"interval"`
	LogFile      string       `yaml:"log_file"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	once := flag.Bool("once", false, "run a single check and exit")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	al, err := alerter.New(cfg.LogFile)
	if err != nil {
		log.Fatalf("creating alerter: %v", err)
	}

	// Convert config service URLs to checker type
	svcURLs := make([]checker.ServiceURL, len(cfg.ServiceURLs))
	for i, s := range cfg.ServiceURLs {
		svcURLs[i] = checker.ServiceURL{Name: s.Name, URL: s.URL}
	}

	chkCfg := checker.Config{
		DiskPaths:   cfg.DiskPaths,
		DiskWarnPct: cfg.DiskWarnPct,
		Containers:  cfg.Containers,
		ServiceURLs: svcURLs,
	}

	run := func() {
		fmt.Printf("\n[%s] Host check\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("========================================")
		result := checker.Run(chkCfg)
		al.Report(result, cfg.DiskWarnPct)
		fmt.Println("========================================")
	}

	if *once || cfg.Interval == 0 {
		run()
		return
	}

	log.Printf("host-watch started | interval=%ds", cfg.Interval)
	run()
	for range time.Tick(time.Duration(cfg.Interval) * time.Second) {
		run()
	}
}
