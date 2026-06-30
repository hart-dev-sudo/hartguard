package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/chrishartserver/linux-security-suite/port-scan-detector/internal/alerter"
	"github.com/chrishartserver/linux-security-suite/port-scan-detector/internal/detector"
	"github.com/chrishartserver/linux-security-suite/port-scan-detector/internal/sniffer"
)

type Config struct {
	Interface string   `yaml:"interface"`
	Threshold int      `yaml:"threshold"`
	Window    int      `yaml:"window"`
	Whitelist []string `yaml:"whitelist"`
	LogFile   string   `yaml:"log_file"`
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
	iface := flag.String("interface", "", "network interface (overrides config)")
	threshold := flag.Int("threshold", 0, "unique ports to trigger alert (overrides config)")
	window := flag.Int("window", 0, "time window in seconds (overrides config)")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	// CLI flags override config
	if *iface != "" {
		cfg.Interface = *iface
	}
	if *threshold > 0 {
		cfg.Threshold = *threshold
	}
	if *window > 0 {
		cfg.Window = *window
	}

	log.Printf("Starting port-scan-detector | interface=%s threshold=%d window=%ds",
		cfg.Interface, cfg.Threshold, cfg.Window)

	det := detector.New(cfg.Threshold, cfg.Window, cfg.Whitelist)

	al, err := alerter.New(cfg.LogFile)
	if err != nil {
		log.Fatalf("creating alerter: %v", err)
	}

	go al.Run(det.Alerts)

	snf := sniffer.New(cfg.Interface, det)
	if err := snf.Start(); err != nil {
		log.Fatalf("sniffer error: %v", err)
	}
}
