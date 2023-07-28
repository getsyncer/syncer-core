package testlog

import (
	"testing"

	"go.uber.org/fx"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

type LogStorage struct {
	buf zaptest.Buffer
}

var Module = fx.Module("testlog",
	fx.Provide(NewZapSplitWriter),
	fx.Supply(new(LogStorage)),
)

func NewZapSplitWriter(t *testing.T, storage *LogStorage) *zap.Logger {
	l := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.NewMultiWriteSyncer(
			zapcore.Lock(zapcore.AddSync(&storage.buf)),
			zapcore.Lock(zapcore.AddSync(testLogAsIOWriter{t})),
		),
		zapcore.InfoLevel,
	))
	return l
}

type testLogAsIOWriter struct {
	t *testing.T
}

func (t testLogAsIOWriter) Write(p []byte) (n int, err error) {
	t.t.Log(string(p))
	return len(p), nil
}
