package log

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"
)

type zapLogger struct {
	instance   *zap.Logger
	writer     *zapio.Writer
	onceWriter *sync.Once
	withCaller bool
}

func (z *zapLogger) WithError(err any) Logger {
	switch v := err.(type) {
	case error:
		return &zapLogger{instance: z.get().With(zap.String("error", v.Error())), onceWriter: new(sync.Once)}
	default:
		return &zapLogger{instance: z.get().With(zap.Reflect("error", err)), onceWriter: new(sync.Once)}
	}
}

func (z *zapLogger) WithField(key string, val any) Logger {
	return &zapLogger{instance: z.get().With(zap.Reflect(key, val)), onceWriter: new(sync.Once)}
}

func (z *zapLogger) WithFields(fields map[string]any) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Reflect(k, v))
	}
	return &zapLogger{instance: z.get().With(zapFields...), onceWriter: new(sync.Once)}
}

func (z *zapLogger) WithContext(ctx context.Context) Logger {
	fields := make([]zap.Field, 0)

	if requestId, ok := ctx.Value(echo.HeaderXRequestID).(string); ok {
		fields = append(fields, zap.String("request_id", requestId))
	}

	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.HasTraceID() {
		fields = append(fields, zap.String("trace_id", spanCtx.TraceID().String()))
	}
	if spanCtx.HasSpanID() {
		fields = append(fields, zap.String("span_id", spanCtx.SpanID().String()))
	}

	return &zapLogger{
		instance:   z.get().With(fields...),
		onceWriter: new(sync.Once),
		withCaller: true,
	}
}

func (z *zapLogger) Debug(message string) {
	z.get().With(z.getMandatoryFields()...).Debug(message)
}

func (z *zapLogger) Error(message string) {
	z.get().With(z.getMandatoryFields()...).Error(message)
}

func (z *zapLogger) Fatal(message string) {
	z.get().With(z.getMandatoryFields()...).Fatal(message)
}

func (z *zapLogger) Info(message string) {
	z.get().With(z.getMandatoryFields()...).Info(message)
}

func (z *zapLogger) Warn(message string) {
	z.get().With(z.getMandatoryFields()...).Warn(message)
}

func (z *zapLogger) GetWriter() io.Writer {
	z.onceWriter.Do(func() {
		z.writer = &zapio.Writer{Log: z.instance, Level: zapcore.InfoLevel}
	})
	return z.writer
}

func (z *zapLogger) Close() error {
	if z.writer != nil {
		_ = z.writer.Close()
	}
	return z.instance.Sync()
}

func (z *zapLogger) get() *zap.Logger {
	return z.instance
}

func (z *zapLogger) getMandatoryFields() []zap.Field {
	var fields []zap.Field
	if z.withCaller {
		fields = append(fields, zap.String("caller", getSource(3)))
	}
	return fields
}

func getSource(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func NewLogger(o *Option) (Logger, error) {
	var (
		logger = &zapLogger{onceWriter: new(sync.Once)}
		cfg    = zap.NewProductionConfig()
		err    error
	)

	if o == nil {
		o = NewDefaultOption()
	}

	cfg.DisableCaller = true
	cfg.OutputPaths = append(cfg.OutputPaths, o.OutputPath)
	cfg.Level, err = zap.ParseAtomicLevel(o.Level)
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true
	if err != nil {
		cfg.Level = zap.NewAtomicLevel()
	}

	logger.instance, err = cfg.Build()

	return logger, err
}
