package utils

import (
	"os"

	"github.com/uber-go/zap"
)

type LoggerConfig struct {
	Level string `toml:"level"`
	Path  string `toml:"path"`
}

func NewLog(conf *LoggerConfig) zap.Logger {
	f, err := os.OpenFile(conf.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}

	zlevel := new(zap.Level)
	err = zlevel.UnmarshalText([]byte(conf.Level))
	if err != nil {
		panic(err)
	}

	logger := zap.New(
		zap.NewJSONEncoder(
			zap.RFC3339Formatter("timestamp"), // human-readable timestamps
			zap.MessageKey("message"),         // customize the message key
			zap.LevelString("level"),          // stringify the log level
		),
		zap.Output(f),
		*zlevel,
	)

	return logger
}
