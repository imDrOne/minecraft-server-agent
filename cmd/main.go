package main

import (
	"flag"
	"github.com/imDrOne/minecraft-server-agent/internal/config"
)

func main() {
	env := flag.String("env", "local", "Application profile")
	flag.Parse()

	cfg := config.New(*env)
	println(cfg)
}
