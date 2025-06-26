# 🪵 Ginrus

Ginrus forwards `gin` logs to logrus.

## ⬇️ Installation

`go get github.com/survivorbat/go-ginrus`

## 📋 Usage

```go
func main() {
  logger := logrus.New()

  engine := gin.New()

  // This will configure ginrus with its default configuration
  engine.Use(ginrus.New(logger))
}
```

## ⚙️ Configuration

In [config.go](./config.go) you'll find the complete configuration struct.
For customisation, apply the following pattern:

```go
func configureLogger(cfg *ginrus.Config) {
  // Change the config values here
  cfg.Fields.ResponseSize = false
}

func main() {
  logger := logrus.New()

  engine := gin.New()

  engine.Use(ginrus.New(logger), configureLogger)
}
```

It is also possible to configure a callback that is fired right before the logger is called using `WithPreLog` as an option.

## 🔭 Plans

Not much yet.
