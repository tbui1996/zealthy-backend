package logging

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

func Must(logger *zap.Logger, err error) *zap.Logger {
	if err != nil {
		log.Panicln(err)
	}
	return logger
}

func NewLoggerFromEvent(event interface{}) (*zap.Logger, error) {
	logger := LoggerFieldsFromEvent(event)

	if logger.Error != nil {
		return nil, logger.Error
	}

	if logger.SkipValidation {
		return logger.Fields.NewLogger()
	}

	err := logger.Fields.Validate()
	if err != nil {
		return nil, err
	}

	return logger.Fields.NewLogger()
}

func NewBasicLogger() (*zap.Logger, error) {
	var lf LoggerFields
	return lf.NewLogger()
}

func (fields *LoggerFields) NewLogger() (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	if os.Getenv("environment") == "prod" {
		// production logger: Does not print debug statements, but
		// will print everything in JSON format so that it is easily
		// searchable by Cloudwatch by default

		// example
		// command: logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
		// output: {"level":"info","ts":1632446521.735787,"caller":"cmd/main.go:45","msg":"This is an INFO message with fields","region":"us-west","id":2}
		loggerConfig := zap.NewProductionConfig()

		// Cloudwatch doesn't differentiate between STDOUT and STDERR, all end up as lines in the logs
		// Configure everything to go through STDOUT so that we can ignore an invalid argument error when calling Sync
		// https://github.com/uber-go/zap/issues/328  https://github.com/influxdata/influxdb/pull/20448
		loggerConfig.OutputPaths = []string{"stdout"}
		logger, err = loggerConfig.Build()

		if err != nil {
			return nil, fmt.Errorf("error getting production logger")
		}
	} else {
		// development logger: Logs that are a little frendlier on
		// developers.Prints only the additional fields in JSON

		// example
		// command: logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
		// output: 2021-09-23T21:22:01.735-0400    INFO    cmd/main.go:31  This is an INFO message with fields     {"region": "us-west", "id": 2}
		loggerConfig := zap.NewDevelopmentConfig()

		loggerConfig.OutputPaths = []string{"stdout"}
		logger, err = loggerConfig.Build()

		if err != nil {
			return nil, fmt.Errorf("error getting development logger")
		}
	}

	if fields.SourceIP != "" {
		logger = logger.With(zap.String("sourceIP", fields.SourceIP))
	}
	if fields.RouteKey != "" {
		logger = logger.With(zap.String("routeKey", fields.RouteKey))
	}
	if fields.RequestID != "" {
		logger = logger.With(zap.String("requestID", fields.RequestID))
	}
	if fields.UserAgent != "" {
		logger = logger.With(zap.String("userAgent", fields.UserAgent))
	}
	if fields.Email != "" {
		logger = logger.With(zap.String("email", fields.Email))
	}
	if fields.UserID != "" {
		logger = logger.With(zap.String("userID", fields.UserID))
	}

	return logger, nil
}

func SyncLogger(logger *zap.Logger) {
	if err := logger.Sync(); err != nil && err.Error() != "sync /dev/stdout: invalid argument" {
		log.Println(err)
	}
}
