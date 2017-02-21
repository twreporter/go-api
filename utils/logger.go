package utils

import (
	log "github.com/Sirupsen/logrus"
)

// LogError print log message and log error as error log level
func LogError(err error, message string) {
	log.WithFields(log.Fields{
		"error": err,
	}).Error(message)
}
