package ginrus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// New instantiates a new ginrus with the provided configuration options
//
//nolint:cyclop // Acceptable in this context, many if-statrements for fields
func New(log *logrus.Logger, opts ...Option) gin.HandlerFunc {
	cfg := newDefaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	ignorePaths := make(map[string]struct{}, len(cfg.IgnorePaths))
	for _, ignoredPath := range cfg.IgnorePaths {
		ignorePaths[ignoredPath] = struct{}{}
	}

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		requestPath := c.Request.URL.Path

		if _, ok := ignorePaths[requestPath]; ok {
			return
		}

		fields := logrus.Fields{}

		if cfg.Fields.Path {
			fields["path"] = requestPath
		}

		if cfg.Fields.Latency {
			fields["latency"] = time.Since(start).String()
		}

		if cfg.Fields.Method {
			fields["method"] = c.Request.Method
		}

		if cfg.Fields.URL {
			fields["url"] = c.Request.URL.String()
		}

		if cfg.Fields.StatusCode {
			fields["status_code"] = c.Writer.Status()
		}

		if cfg.Fields.Latency {
			fields["latency"] = time.Since(start).String()
		}

		if cfg.Fields.RequestSize {
			fields["request_size"] = c.Request.ContentLength
		}

		if cfg.Fields.ResponseSize {
			fields["response_size"] = c.Writer.Size()
		}

		if cfg.Fields.ClientIP {
			fields["client_ip"] = c.ClientIP()
		}

		if cfg.Fields.UserAgent {
			fields["user_agent"] = c.Request.UserAgent()
		}

		if cfg.Fields.Referrer {
			fields["referer"] = c.Request.Referer()
		}

		cfg.PreLogFunc(c, fields)

		if c.Writer.Status() > 399 || len(c.Errors) > 0 {
			log.WithContext(c.Request.Context()).WithFields(fields).Log(cfg.LogLevels.ErrorStatusCodeLogLevel, cfg.Messages.DefaultErrorMessage)
			return
		}

		log.WithContext(c.Request.Context()).WithFields(fields).Log(cfg.LogLevels.DefaultLogLevel, cfg.Messages.DefaultMessage)
	}
}
