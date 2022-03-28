package logging

import (
	"fmt"
	"log"

	"go.uber.org/zap"
)

var myLogger *zap.Logger
var logConfig zap.Config

func init() {
	instantiateLogger()
}

func instantiateLogger() {
	_, err := getLogger()
	if err != nil {
		log.Printf("Error bootstrapping logger: %s", err)
		panic(err)
	}

}

// GetLogger give the instance of logger
func GetLogger() (logger *zap.Logger) {
	return myLogger
}

func getLogger() (*zap.Logger, error) {
	var err error
	logConfig = zap.NewDevelopmentConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	myLogger, err = logConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("Error while attempting to acquire logger: %w", err)
	}
	return myLogger, nil
}
