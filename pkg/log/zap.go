package log

import (
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  lumberjack "gopkg.in/natefinch/lumberjack.v2"
  "sync"
)

var (
  once sync.Once
  logger *zap.Logger
)

func InitLogger(logpath string, loglevel string) *zap.Logger {
  hook := lumberjack.Logger{
    Filename:   logpath, // ⽇志⽂件路径
    MaxSize:    1024,    // megabytes
    MaxBackups: 10,       // 最多保留3个备份
    MaxAge:     7,       //days
    Compress:   true,    // 是否压缩 disabled by default
  }
  w := zapcore.AddSync(&hook)
  var level zapcore.Level
  switch loglevel {
  case "debug":
    level = zap.DebugLevel
  case "info":
    level = zap.InfoLevel
  case "error":
    level = zap.ErrorLevel
  default:
    level = zap.InfoLevel
  }
  encoderConfig := zap.NewProductionEncoderConfig()
  encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
  core := zapcore.NewCore(
    zapcore.NewConsoleEncoder(encoderConfig),
    w,
    level,
  )
  logger := zap.New(core)
  logger.Info("DefaultLogger init success")
  return logger
}

func GetLogger() *zap.Logger {
  once.Do(func() {
    logger = InitLogger("./logs/webspider.log", "warn")
  })
  return logger
}

//func main() {
//  logger := initLogger("all.log", "info")
//  logger.Info("test log", zap.Int("line", 47))
//  logger.Warn("testlog", zap.Int("line", 47))
//}
