//go:build ignore
// +build ignore

package main

import (
	"os"

	"github.com/wzbwzt/pkg/pkg/log"

	"go.uber.org/zap/zapcore"
)

func main() {
	defer log.Sync()
	log.Debug("debug", log.String("s", "v"))
	log.Error("debug", log.String("s", "v"))

	f, _ := os.OpenFile("./log.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	logger := log.NewLogger(f, log.InfoLevel)
	defer logger.Sync()

	logger.Error("this is error",
		log.String("string", "v"),
		log.Int("int", 1),
	)

	accesslog, _ := os.OpenFile("./access.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	errorlog, _ := os.OpenFile("./error.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	tees := make([]log.Tee, 0, 2)
	tees = append(tees, log.Tee{
		W: accesslog,
		Lef: func(l zapcore.Level) bool {
			return l < log.ErrorLevel
		},
	}, log.Tee{
		W: errorlog,
		Lef: func(l zapcore.Level) bool {
			return l >= log.ErrorLevel
		},
	})
	log2 := log.NewTee(tees, log.AddCaller())
	defer log2.Sync()

	log2.Info("info:", log.String("app", "start ok"),
		log.Int("major version", 3))
	log2.Error("error:", log.String("app", "crash"),
		log.Int("reason", -1))

	logrotate := log.NewLoggerWithRotate(log.RotateConfig{
		Filename:   "rotate.log",
		MaxSize:    1,
		MaxAge:     1,
		MaxBackups: 3,
		Compress:   true,
	}, log.DebugLevel)

	for i := 0; i < 20000; i++ {
		logrotate.Info("info_for_rotate:", log.String("app", "start ok"),
			log.Int("major version", 3))
		logrotate.Error("err_for_rotate:", log.String("app", "crash"),
			log.Int("reason", -1))
	}
}
