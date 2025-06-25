package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {

	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{})
}
