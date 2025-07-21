package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware создает middleware для логирования запросов
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Старт таймера
		start := time.Now()

		// Запоминаем некоторые данные запроса ДО обработки
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		userAgent := c.Request.UserAgent()

		// Для логирования тела запроса (с ограничением)
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, 1024))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Перехватываем ответ для логирования
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Обрабатываем запрос
		c.Next()

		// Данные ПОСЛЕ обработки
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Логируем основные параметры
		logger.Logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.Int("response-size", responseSize),
			zap.String("error", errorMessage),
			zap.String("request-id", c.GetHeader("X-Request-ID")),
		)

		// Дополнительное логирование для важных событий
		if statusCode >= 400 {
			// Логируем тело запроса для ошибок
			logger.Logger.Error("Request Error Details",
				zap.ByteString("request-body", requestBody),
				zap.String("response-body", blw.body.String()),
			)
		}

		// Логирование критических ошибок
		if statusCode >= 500 {
			logger.Logger.Error("Server Error",
				zap.String("path", path),
				zap.Int("status", statusCode),
			)
		}
	}
}

// Вспомогательная структура для перехвата тела ответа
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
