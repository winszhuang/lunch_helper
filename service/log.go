package service

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var sugarLogger *zap.SugaredLogger

type LogService struct {
	logger *zap.SugaredLogger
}

func NewLogService(logFile, logErrorFile string) *LogService {
	encoder := getEncoder()

	// handle debug level
	debugCore := zapcore.NewCore(encoder, getLogWriter(logFile), zapcore.DebugLevel)
	// handle error level
	coreError := zapcore.NewCore(encoder, getLogWriter(logErrorFile), zapcore.ErrorLevel)

	core := zapcore.NewTee(debugCore, coreError)
	// AddCallerSkip確保不會每次打印出來都是service/log.go
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger = logger.Sugar()

	return &LogService{
		logger: sugarLogger,
	}
}

func (ls *LogService) Sync() {
	ls.logger.Sync()
}

func (ls *LogService) Infof(template string, args ...interface{}) {
	ls.logger.Infof(template, args...)
}

func (ls *LogService) Error(args ...interface{}) {
	ls.logger.Error(args...)
}

func (ls *LogService) Errorf(template string, args ...interface{}) {
	ls.logger.Errorf(template, args...)
}

func (ls *LogService) Debugf(template string, args ...interface{}) {
	ls.logger.Debugf(template, args...)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(fileName string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		// Filename:   "log/test.log",
		Filename:   fileName,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	// 寫入檔案且打印至terminal
	mw := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(mw)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 將時間轉換為台北時區
	loc, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		return
	}
	t = t.In(loc)

	enc.AppendString(t.Format("2006-01-02T15:04:05.000-0700"))
}
