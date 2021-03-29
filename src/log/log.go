package log

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Config zap.Config

// 日志变量
var Logger *zap.Logger
var Sugar *zap.SugaredLogger

const ProgramName = "zip-tool"
const logDir = "./logs"

func init() {
	// 创建日志目录
	os.MkdirAll(logDir, os.ModePerm)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	Config = zap.Config{
		Level:         atom,          // 日志级别
		Development:   true,          // 开发模式，堆栈跟踪
		Encoding:      "console",     // 输出格式 console 或 json
		EncoderConfig: encoderConfig, // 编码器配置
		//InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths:      []string{"stdout", filepath.Join(logDir, "./zip-tool-common.log")}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	Logger, err = Config.Build()

	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	} else {
		Sugar = Logger.Sugar()
	}
}
