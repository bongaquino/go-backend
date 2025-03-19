package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func init() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer Log.Sync()
}

// Error is a wrapper for logger.Error
func Error(err error) zap.Field {
	return zap.Error(err)
}

// String is a wrapper for logger.String
func String(key string, value string) zap.Field {
	return zap.String(key, value)
}
