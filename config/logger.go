package config

import "github.com/sirupsen/logrus"

// / logger log levels.
const (
	warningLoggerLevel = "warning"
	errorLoggerLevel   = "error"
	fatalLoggerLevel   = "fatal"
	infoLoggerLevel    = "info"
)

// initLogger init logrus preferences.
func InitLogger(cfg LoggerCfg, version string) *logrus.Entry {
	logger := logrus.New()

	logger.SetReportCaller(cfg.Caller)

	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	var level logrus.Level

	switch cfg.Level {
	case warningLoggerLevel:
		level = logrus.WarnLevel
	case errorLoggerLevel:
		level = logrus.ErrorLevel
	case fatalLoggerLevel:
		level = logrus.FatalLevel
	case infoLoggerLevel:
		level = logrus.InfoLevel
	default:
		level = logrus.DebugLevel
	}

	logger.SetLevel(level)

	return logger.WithField("version", version)
}
