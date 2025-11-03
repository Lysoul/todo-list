package main

import (
	"github.com/Lysoul/gocommon/monitoring"
)

func main() {
	log := monitoring.Logger()
	log.Info("Application started")
}
