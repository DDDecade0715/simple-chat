package zap

import (
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"time"
)

type Zap struct {
	SugarLogger *zap.SugaredLogger
}

var SugarLoggerPool = make(map[string]*Zap)

func GetZap(name string, p string) *zap.SugaredLogger {
	if sugarLogger, ok := SugarLoggerPool[name]; ok {
		return sugarLogger.SugarLogger
	}
	encoder := getEncoder()
	writeSyncer := getLogWriter(p)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	SugarLoggerPool[name] = &Zap{
		SugarLogger: logger.Sugar(),
	}
	return logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(p string) zapcore.WriteSyncer {
	fileWriter, _ := zaprotatelogs.New(
		path.Join(p, "%Y-%m-%d.log"),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	return zapcore.AddSync(fileWriter)
}

func (s *Zap) Debug(args ...interface{}) {
	s.SugarLogger.Debug(args...)
}

func (s *Zap) Debugf(template string, args ...interface{}) {
	s.SugarLogger.Debugf(template, args...)
}

func (s *Zap) Info(args ...interface{}) {
	s.SugarLogger.Info(args...)
}

func (s *Zap) Infof(template string, args ...interface{}) {
	s.SugarLogger.Infof(template, args...)
}

func (s *Zap) Infow(template string, args ...interface{}) {
	s.SugarLogger.Infow(template, args...)
}

func (s *Zap) Warn(args ...interface{}) {
	s.SugarLogger.Warn(args...)
}

func (s *Zap) Warnf(template string, args ...interface{}) {
	s.SugarLogger.Warnf(template, args...)
}

func (s *Zap) Error(args ...interface{}) {
	s.SugarLogger.Error(args...)
}

func (s *Zap) Errorf(template string, args ...interface{}) {
	s.SugarLogger.Errorf(template, args...)
}

func (s *Zap) DPanic(args ...interface{}) {
	s.SugarLogger.DPanic(args...)
}

func (s *Zap) DPanicf(template string, args ...interface{}) {
	s.SugarLogger.DPanicf(template, args...)
}

func (s *Zap) Panic(args ...interface{}) {
	s.SugarLogger.Panic(args...)
}

func (s *Zap) Panicf(template string, args ...interface{}) {
	s.SugarLogger.Panicf(template, args...)
}

func (s *Zap) Fatal(args ...interface{}) {
	s.SugarLogger.Fatal(args...)
}

func (s *Zap) Fatalf(template string, args ...interface{}) {
	s.SugarLogger.Fatalf(template, args...)
}
