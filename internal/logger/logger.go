package logger

import (
	"learn-gin/internal/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局Logger实例
var Log *zap.Logger

// Init 初始化日志
func Init() {
	cfg := config.Conf.Log

	// 配置Lumberjack
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename: cfg.Filename,
		MaxSize: cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge: cfg.MaxAge,
		Compress: true, // 旧文件压缩
	})

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式：2023-01-01T00:00:00.000Z
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 级别格式：INFO， ERROR

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
        level = zap.InfoLevel
    case "error":
        level = zap.ErrorLevel
    default:
        level = zap.InfoLevel
	}

	// 创建core
	// 可以让日志同时输出到控制台和文件
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(os.Stdout, writeSyncer),
		level,
	)

	// 初始化全局变量
	Log = zap.New(core, zap.AddCaller()) // AddCaller会显示日志是那行代码带出来的
}