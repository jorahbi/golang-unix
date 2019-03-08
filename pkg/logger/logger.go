package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	debug *zap.Logger
	info *zap.Logger
	warn *zap.Logger
	error *zap.Logger
	dCrash *zap.Logger
	crash *zap.Logger
	fatal *zap.Logger
}

var Handle = newLogger()

func (log *Logger) Debug (msg string, fields ...zap.Field) {
	defer err(log.debug)
	log.debug.Debug(msg, fields...)
}

func (log *Logger) Info (msg string, fields ...zap.Field) {
	defer err(log.info)
	log.info.Info(msg, fields...)
}

func (log *Logger) Warn (msg string, fields ...zap.Field) {
	defer err(log.warn)
	log.warn.Warn(msg, fields...)
}

func (log *Logger) Error (msg string, fields ...zap.Field) {
	defer err(log.error)
	log.error.Error(msg, fields...)
}

func (log *Logger) DPanic (msg string, fields ...zap.Field) {
	defer err(log.dCrash)
	log.dCrash.DPanic(msg, fields...)
}

func (log *Logger) Panic (msg string, fields ...zap.Field) {
	defer err(log.crash)
	log.crash.Panic(msg, fields...)
}

func (log *Logger) Fatal (msg string, fields ...zap.Field) {
	defer err(log.fatal)
	log.fatal.Fatal(msg, fields...)
}

func err (log *zap.Logger) {
	if err := log.Sync(); err != nil {
		fmt.Println(errors.WithStack(err))
	}
}

func newLogger () *Logger{
	zap.NewAtomicLevel()
	log := &Logger{}
	basePath := "/home/bruce/Workspace/golang/src/php-go/"
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log.debug = newZapLog(fmt.Sprintf("%s%s", basePath, "debug"), zapcore.DebugLevel, encoderConfig)
	log.info = newZapLog(fmt.Sprintf("%s%s", basePath, "info"), zapcore.InfoLevel, encoderConfig)
	log.warn = newZapLog(fmt.Sprintf("%s%s", basePath, "warn"), zapcore.WarnLevel, encoderConfig)
	log.error = newZapLog(fmt.Sprintf("%s%s", basePath, "error"), zapcore.ErrorLevel, encoderConfig)
	log.dCrash = newZapLog(fmt.Sprintf("%s%s", basePath, "dPanic"), zapcore.DPanicLevel, encoderConfig)
	log.crash = newZapLog(fmt.Sprintf("%s%s", basePath, "panic"), zapcore.PanicLevel, encoderConfig)
	log.fatal = newZapLog(fmt.Sprintf("%s%s", basePath, "fatal"), zapcore.FatalLevel, encoderConfig)

	return log
}

func newZapLog (path string, level zapcore.Level, encoderConfig zapcore.EncoderConfig) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   path, // 日志文件路径
		MaxSize:    128,     // megabytes
		MaxBackups: 30,      // 最多保留300个备份
		MaxAge:     7,       // days
		Compress:   true,    // 是否压缩 disabled by default
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&hook),
		level,
	)

	return zap.New(core)
}

