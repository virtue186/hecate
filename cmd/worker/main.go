package main

import "hecate/internal/pkg/config"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	print(cfg.Server.Port)

}
