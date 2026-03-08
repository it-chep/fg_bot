package logger

import (
	"context"
	"os"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey string

const contextLoggerKey loggerKey = "clk"

type Logger struct {
	zap *zap.Logger
}

type ErrorPayload struct {
	Message  string
	Error    error
	File     string
	Line     int
	Function string
}

type ErrorReporter func(ctx context.Context, payload ErrorPayload) error

var (
	errorReporterMu sync.RWMutex
	errorReporter   ErrorReporter
)

func New() *Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	return &Logger{
		zap: zap.New(core),
	}
}

func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, l)
}

func fromContext(ctx context.Context) *Logger {
	l, ok := ctx.Value(contextLoggerKey).(*Logger)
	if ok {
		return l
	}
	return nil
}

func Message(ctx context.Context, format string) {
	l := fromContext(ctx)
	if l != nil {
		l.zap.Info(format)
	}
}

func Error(ctx context.Context, format string, err error) {
	payload := buildErrorPayload(format, err)

	errorReporterMu.RLock()
	reporter := errorReporter
	errorReporterMu.RUnlock()

	if reporter != nil {
		if reportErr := reporter(ctx, payload); reportErr == nil {
			return
		}
	}

	l := fromContext(ctx)
	if l != nil {
		l.zap.Error(format, zap.Error(err))
	}
}

func SetErrorReporter(reporter ErrorReporter) {
	errorReporterMu.Lock()
	defer errorReporterMu.Unlock()
	errorReporter = reporter
}

func buildErrorPayload(message string, err error) ErrorPayload {
	payload := ErrorPayload{
		Message: message,
		Error:   err,
	}

	// 0 - buildErrorPayload, 1 - Error, 2 - caller of Error
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		payload.File = file
		payload.Line = line
	}
	if fn := runtime.FuncForPC(pc); fn != nil {
		payload.Function = fn.Name()
	}

	return payload
}
