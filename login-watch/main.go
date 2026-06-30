package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/hart-dev-sudo/hartguard/login-watch/internal/alerter"
	"github.com/hart-dev-sudo/hartguard/login-watch/internal/detector"
	"github.com/hart-dev-sudo/hartguard/login-watch/internal/parser"
)

type Config struct {
	LogFile    string   `yaml:"log_file"`
	AuthLog    string   `yaml:"auth_log"`
	Threshold  int      `yaml:"threshold"`
	Window     int      `yaml:"window"`
	Whitelist  []string `yaml:"whitelist"`
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
	scan := flag.Bool("scan", false, "scan existing log and exit (no tail)")
	flag.Parse()

	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	al, err := alerter.New(cfg.LogFile)
	if err != nil {
		log.Fatalf("creating alerter: %v", err)
	}

	det := detector.New(cfg.Threshold, cfg.Window, cfg.Whitelist)

	var wg sync.WaitGroup
	wg.Add(1)
	go al.Run(det.Alerts, wg.Done)

	f, err := os.Open(cfg.AuthLog)
	if err != nil {
		log.Fatalf("opening auth log %s: %v", cfg.AuthLog, err)
	}
	defer f.Close()

	if *scan {
		fmt.Println("========================================")
		fmt.Printf("  login-watch — scan mode\n")
		fmt.Printf("  log:       %s\n", cfg.AuthLog)
		fmt.Printf("  threshold: %d failures in %ds\n", cfg.Threshold, cfg.Window)
		fmt.Println("========================================")
		scanExisting(f, det)
		close(det.Alerts)
		wg.Wait()
		fmt.Println("========================================")
		fmt.Println("  scan complete")
		fmt.Println("========================================")
		return
	}

	// Seek to end for live tailing
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		log.Fatalf("seeking log: %v", err)
	}

	fmt.Println("========================================")
	fmt.Printf("  login-watch — live mode\n")
	fmt.Printf("  log:       %s\n", cfg.AuthLog)
	fmt.Printf("  threshold: %d failures in %ds\n", cfg.Threshold, cfg.Window)
	fmt.Println("========================================")

	tail(f, det)
}

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func scanExisting(f *os.File, det *detector.Detector) {
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if e := parser.Parse(scanner.Text()); e != nil {
			det.Process(e.SrcIP, e.EventType, e.Username)
			count++
		}
	}
	if count == 0 {
		fmt.Printf("%s[OK]%s    No failure events found\n", colorGreen, colorReset)
	} else {
		fmt.Printf("[*]     %d failure events processed\n", count)
	}
}

func tail(f *os.File, det *detector.Detector) {
	reader := bufio.NewReader(f)
	heartbeat := time.NewTicker(60 * time.Second)
	defer heartbeat.Stop()

	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if e := parser.Parse(line); e != nil {
				det.Process(e.SrcIP, e.EventType, e.Username)
			}
		}
		if err == io.EOF {
			select {
			case t := <-heartbeat.C:
				fmt.Printf("%s[*]%s     still watching — %s\n",
					colorGreen, colorReset, t.Format("15:04:05"))
			default:
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err != nil {
			log.Printf("read error: %v", err)
			return
		}
	}
}
