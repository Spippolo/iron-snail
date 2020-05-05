package main

import (
	"github.com/Spippolo/iron-snail/game"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
