package helpers

import (
	log "github.com/sirupsen/logrus"
)

// GetLogLevel retuens the logrus log level corresponding to the input integer
//
// The number should be between 0 and 5 inclusive.
// Higher numbers correspond to higher log levels.
func GetLogLevel(level uint8) log.Level {
	switch level {
	case 0:
		return log.PanicLevel
	case 1:
		return log.FatalLevel
	case 2:
		return log.ErrorLevel
	case 3:
		return log.WarnLevel
	case 4:
		return log.InfoLevel
	case 5:
		return log.DebugLevel
	default:
		return log.WarnLevel
	}
}

// LogError logs err if it is non nil
func LogError(err error) {
	if err != nil {
		log.Error(err)
	}
}

// LogErrorExit logs err and exits if the input error is non nil
func LogErrorExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
