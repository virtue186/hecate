package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"hecate/internal/pkg/config"
	"io"
	"os"
	"path"
	"strings"
)

var log *logrus.Logger

// Init 使用提供的配置初始化全局 Logger
func Init(cfg *config.LogConfig) error {
	log = logrus.New()

	// 1. 设置日志格式
	// 检查输出中是否包含 stdout，以便决定是否启用颜色
	outputType := strings.ToLower(cfg.Output)
	isStdOut := outputType == "stdout" || outputType == "both"

	switch strings.ToLower(cfg.Format) {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		// 2. 关键改动：使用新的 Formatter
		// 只有当输出到控制台时才启用颜色
		formatter := &prefixed.TextFormatter{
			ForceColors:     true,      // 强制开启颜色
			DisableColors:   !isStdOut, // 如果不输出到控制台，则禁用颜色
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			// 你还可以自定义前缀的顺序和外观
			// EntrySectionName: "entry",
			// LevelSectionName: "level",
		}
		log.SetFormatter(formatter)
	}

	// 2. 设置日志级别 (这部分代码无需改动)
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	// 3. 设置输出 (这部分代码无需改动)
	var writers []io.Writer

	if isStdOut {
		writers = append(writers, os.Stdout)
	}
	if outputType == "file" || outputType == "both" {
		logDir := path.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writers = append(writers, file)
	}

	if len(writers) > 0 {
		log.SetOutput(io.MultiWriter(writers...))
	}

	log.WithField("service_name", "hecate-api").Info("Logger initialized successfully")
	return nil
}

// GetLogger 返回已初始化的 Logger 实例 (这部分代码无需改动)
func GetLogger() *logrus.Logger {
	if log == nil {
		defaultLog := logrus.New()
		defaultLog.SetLevel(logrus.InfoLevel)
		// 同样可以给默认 logger 也加上颜色
		defaultLog.SetFormatter(&prefixed.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		return defaultLog
	}
	return log
}
