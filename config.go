package ginrus

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Option func(*Config)

// PreLogFunc is a callback that allows you to add information to a log message before it's published
type PreLogFunc func(*gin.Context, logrus.Fields)

// Config is formed through Option calls in Logger
type Config struct {
	// LogLevels defines specific log levels for types of messages
	LogLevels struct {
		// DefaultLogLevel is the default log level used for status codes below 400
		DefaultLogLevel logrus.Level

		// ErrorStatusCodeLogLevel is the log level used for status codes above 399
		ErrorStatusCodeLogLevel logrus.Level
	}

	// Messages defines all the messages reported in log messages, alongside the fields
	Messages struct {
		// DefaultMessage is the default message reported for a non-error message
		DefaultMessage string

		// DefaultErrorMessage is the default message reported for an error message
		DefaultErrorMessage string
	}

	// PreLogFunc is executed before the final log call
	PreLogFunc PreLogFunc

	// Fields allows you to toggle fields in the output
	Fields struct {
		Method       bool
		URL          bool
		StatusCode   bool
		RequestSize  bool
		ResponseSize bool
		Latency      bool
		Referrer     bool
		UserAgent    bool
		ClientIP     bool
		Path         bool
	}
}

// newDefaultConfig returns the default settings for the logger
func newDefaultConfig() *Config {
	cfg := new(Config)
	cfg.LogLevels.DefaultLogLevel = logrus.InfoLevel
	cfg.LogLevels.ErrorStatusCodeLogLevel = logrus.ErrorLevel

	cfg.Fields.Method = true
	cfg.Fields.URL = true
	cfg.Fields.StatusCode = true
	cfg.Fields.RequestSize = true
	cfg.Fields.ResponseSize = true
	cfg.Fields.Latency = true
	cfg.Fields.Referrer = true
	cfg.Fields.UserAgent = true
	cfg.Fields.ClientIP = true
	cfg.Fields.Path = true

	cfg.Messages.DefaultMessage = "request succeeded"
	cfg.Messages.DefaultErrorMessage = "request failed"

	cfg.PreLogFunc = func(*gin.Context, logrus.Fields) {}

	return cfg
}

// WithPreLog registers a callback that allows you to add information to a log message before it's published
func WithPreLog(f PreLogFunc) Option {
	return func(c *Config) {
		c.PreLogFunc = f
	}
}
