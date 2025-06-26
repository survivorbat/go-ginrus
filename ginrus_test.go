package ginrus

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type logrusOutput struct {
	ClientIP     string `json:"client_ip"`
	RequestSize  int64  `json:"request_size"`
	ResponseSize int64  `json:"response_size"`
	Level        string `json:"level"`
	Method       string `json:"method"`
	Msg          string `json:"msg"`
	Path         string `json:"path"`
	Referer      string `json:"referer"`
	StatusCode   int    `json:"status_code"`
	URL          string `json:"url"`
	UserAgent    string `json:"user_agent"`

	// These are too difficult to test because of timing
	// Latency    float64 `json:"latency"`
	// Time string `json:"time"`
}

func TestNew_LogsIncomingRequests(t *testing.T) {
	t.Parallel()
	// Arrange
	out := new(bytes.Buffer)

	logger := logrus.New()
	logger.Out = out

	// Easier to parse
	logger.SetFormatter(&logrus.JSONFormatter{})

	subject := New(logger)

	httpWriter := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(httpWriter)

	// This sets content-length to 5
	body := strings.NewReader("abcde")

	req := httptest.NewRequest(http.MethodPut, "https://localhost:8080/foo/bar", body)
	req.Header.Set("User-Agent", "Test Agent")
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	req.Header.Set("Referer", "localhost.net")

	ginContext.Request = req

	// Write some data to the writer
	ginContext.JSON(http.StatusCreated, map[string]any{"foo": "bar"})

	// Act
	subject(ginContext)

	// Assert
	var actual logrusOutput
	err := json.Unmarshal(out.Bytes(), &actual)
	require.NoError(t, err)

	expected := logrusOutput{
		ClientIP:     "1.2.3.4",
		RequestSize:  5,
		ResponseSize: 13,
		Level:        logrus.InfoLevel.String(),
		Method:       http.MethodPut,
		Msg:          "request succeeded",
		Path:         "/foo/bar",
		Referer:      "localhost.net",
		StatusCode:   http.StatusCreated,
		URL:          "https://localhost:8080/foo/bar",
		UserAgent:    "Test Agent",
	}

	assert.Equal(t, expected, actual)
}

func TestNew_CallsPreLogCallback(t *testing.T) {
	t.Parallel()
	// Arrange
	logger := logrus.New()
	logger.Out = io.Discard

	var callbackCalled bool
	callback := func(*gin.Context, logrus.Fields) {
		callbackCalled = true
	}

	subject := New(logger, WithPreLog(callback))

	httpWriter := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(httpWriter)
	ginContext.Request = httptest.NewRequest(http.MethodGet, "https://localhost", http.NoBody)

	// Act
	subject(ginContext)

	// Assert
	assert.True(t, callbackCalled, "expected callback to have been called")
}

func TestNew_IgnoresPathsIfRequested(t *testing.T) {
	t.Parallel()
	// Arrange
	out := new(bytes.Buffer)

	logger := logrus.New()
	logger.Out = out

	subject := New(logger, WithIgnore("/health"))

	httpWriter := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(httpWriter)
	ginContext.Request = httptest.NewRequest(http.MethodGet, "https://localhost/health", http.NoBody)

	// Act
	subject(ginContext)

	// Assert
	assert.Empty(t, out.Bytes())
}
