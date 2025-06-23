package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hecate/internal/pkg/logger"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		log := logger.GetLogger()

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		entry := log.WithFields(logrus.Fields{
			"statusCode": statusCode,
			"latency":    latency,
			"clientIP":   clientIP,
			"method":     method,
			"path":       path,
		})

		if raw != "" {
			entry = entry.WithField("query", raw)
		}

		if len(c.Errors) > 0 {
			// 如果有错误，记录错误信息
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			// 否则，记录请求信息
			msg := "Request handled"
			if statusCode >= 500 {
				entry.Error(msg)
			} else if statusCode >= 400 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}
