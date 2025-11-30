package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitZap() {
	// 1. 定义编码器配置（强制开启颜色）
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey, // 省略函数名（可选）
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 颜色级别（关键）
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短调用路径（如 pkg/file.go:123）
	}

	// 2. 选择输出目标（控制台 + 文件，可选）
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	level := zap.DebugLevel
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

	logger := zap.New(core, zap.AddCaller())

	// 替换全局 logger，后续可通过 zap.L() 调用
	zap.ReplaceGlobals(logger)

}
