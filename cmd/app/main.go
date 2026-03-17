package main

import (
	"flag"
	"log"

	"github.com/nurtikaga/tms_proj/internal/app"
	"github.com/nurtikaga/tms_proj/internal/config"
)

func main() {
	cfgPath := flag.String("config", "internal/config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err = app.Run(cfg); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
