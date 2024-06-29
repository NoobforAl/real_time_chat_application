package logging

import (
	"os"
	"sync"

	"github.com/NoobforAl/real_time_chat_application/src/config"
	"github.com/NoobforAl/real_time_chat_application/src/contract"
	"github.com/sirupsen/logrus"
)

var once sync.Once
var logger *logrus.Logger

func New() contract.Logger {
	once.Do(
		func() {
			level := logrus.InfoLevel

			if config.Debug {
				level = logrus.DebugLevel
			}

			logger = &logrus.Logger{
				Out:          os.Stderr,
				Formatter:    new(logrus.TextFormatter),
				Hooks:        make(logrus.LevelHooks),
				Level:        level,
				ExitFunc:     os.Exit,
				ReportCaller: false,
			}
		})

	return logger
}
