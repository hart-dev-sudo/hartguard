package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/chrishartserver/hartguard/vpn-watch/internal/alerter"
	"github.com/chrishartserver/hartguard/vpn-watch/internal/checker"
)

type Config struct {
	VPNContainer      string   `yaml:"vpn_container"`
	CheckContainers   []string `yaml:"check_containers"`
	CheckURL          string   `yaml:"check_url"`
	Interval          int      `yaml:"interval"`
	LogFile           string   `yaml:"log_file"`
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

	chk := checker.New(cfg.VPNContainer, cfg.CheckContainers, cfg.CheckURL, checker.DockerExecutor{})

	run := func() {
		fmt.Printf("\n[%s] Running VPN check...\n", time.Now().Format("2006-01-02 15:04:05"))
		result := chk.Run()
		al.Report(result)
	}

	if *once || cfg.Interval == 0 {
		run()
		return
	}

	log.Printf("vpn-watch started | interval=%ds | vpn=%s | watching=%v",
		cfg.Interval, cfg.VPNContainer, cfg.CheckContainers)

	run()
	for range time.Tick(time.Duration(cfg.Interval) * time.Second) {
		run()
	}
}
