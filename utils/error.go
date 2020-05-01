package utils

import log "github.com/sirupsen/logrus"

func CheckErr(err error, label string) {
	if err != nil {
		log.Fatalf("%s: %v", label, err)
	}
}
