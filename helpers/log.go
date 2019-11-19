package helpers

import (
	log "github.com/sirupsen/logrus"
)

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

func LogError(err error) {
	if err != nil {
		log.Error(err)
	}
}

func LogErrorExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
