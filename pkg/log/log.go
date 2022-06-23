package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	l *zap.Logger
}

type level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

var (
	String     = zap.String
	Stringp    = zap.Stringp
	Skip       = zap.Skip
	Binary     = zap.Binary
	Bool       = zap.Bool
	Boolp      = zap.Boolp
	ByteString = zap.ByteString
	Float64    = zap.Float64
	Float64p   = zap.Float64p
	Float32    = zap.Float32
	Float32p   = zap.Float32p
	Durationp  = zap.Durationp
	Any        = zap.Any
	Int        = zap.Int
	Intp       = zap.Intp
)

var (
	Debug  = std.Debug
	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal

	Sync = std.Sync
)

type Option = zap.Option

var (
	//AddCaller configures the Logger to annotate each message with the filename,
	//line number, and function name of zap's caller.
	AddCaller = zap.AddCaller
	//AddStacktrace configures the Logger to record a stack trace for all
	//messages at or above a given level.
	AddStacktrace = zap.AddStacktrace
)

var std = NewLogger(os.Stderr, DebugLevel, AddCaller(), AddStacktrace(ErrorLevel))

func NewLogger(w io.Writer, level level, ops ...Option) *Logger {
	cf := zap.NewProductionConfig()
	//日期格式
	cf.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cf.EncoderConfig),
		zapcore.AddSync(w),
		level,
	)
	log := zap.New(core, ops...)

	return &Logger{l: log}
}

type Field = zapcore.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}
func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}
func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

//一个log实例写入多个log文件，类似nginx的access.log和error.log
type Tee struct {
	W   io.Writer
	Lef zap.LevelEnablerFunc
}

func NewTee(tees []Tee, ops ...Option) *Logger {
	var cores = make([]zapcore.Core, 0, len(tees))

	cf := zap.NewProductionConfig()
	//日期格式
	cf.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}

	for _, tee := range tees {
		tee := tee
		lv := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return tee.Lef(l)
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cf.EncoderConfig),
			zapcore.AddSync(tee.W),
			lv,
		)

		cores = append(cores, core)
	}

	return &Logger{l: zap.New(zapcore.NewTee(cores...), ops...)}
}

type RotateConfig = lumberjack.Logger

func NewLoggerWithRotate(rotateconfig RotateConfig, level level, ops ...Option) *Logger {
	cf := zap.NewProductionConfig()
	//日期格式
	cf.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format(time.RFC3339))
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cf.EncoderConfig),
		zapcore.AddSync(&rotateconfig),
		level,
	)
	log := zap.New(core, ops...)

	return &Logger{l: log}
}
