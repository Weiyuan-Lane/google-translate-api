package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	wrappedLogger *zap.Logger
}

func New(appName string, isDevEnv bool) *Logger {
	var cfg zap.Config

	if isDevEnv {
		cfg = newDevelopmentLoggerConfig(appName)
	} else {
		cfg = newProductionLoggerConfig(appName)
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		wrappedLogger: zapLogger,
	}
}

func (l *Logger) Debug(msg string, inputFields ...map[string]string) {
	fields := []zap.Field{}

	if len(inputFields) > 0 {
		fields = transformStrMapToFields(inputFields[0])
	}

	l.wrappedLogger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, inputFields ...map[string]string) {
	fields := []zap.Field{}

	if len(inputFields) > 0 {
		fields = transformStrMapToFields(inputFields[0])
	}

	l.wrappedLogger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, inputFields ...map[string]string) {
	fields := []zap.Field{}

	if len(inputFields) > 0 {
		fields = transformStrMapToFields(inputFields[0])
	}

	l.wrappedLogger.Fatal(msg, fields...)
}

func (l *Logger) Info(msg string, inputFields ...map[string]string) {
	fields := []zap.Field{}

	if len(inputFields) > 0 {
		fields = transformStrMapToFields(inputFields[0])
	}

	l.wrappedLogger.Info(msg, fields...)
}

func transformStrMapToFields(strMap map[string]string) []zap.Field {
	fields := []zap.Field{}
	for k, v := range strMap {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}
