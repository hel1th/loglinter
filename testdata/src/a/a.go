package testdata

import (
	"log"
	"log/slog"

	"go.uber.org/zap"
)


func testAllRulesValid() {
	// корректные логи
	log.Print("server started on port 8080")
	log.Println("database connection established")
	log.Printf("processed %d requests", 100)
	slog.Info("user authenticated successfully")
	slog.Error("failed to read configuration file")
	slog.Warn("cache miss on key user_session")
	slog.Debug("background worker initialized")

	logger, _ := zap.NewProduction()
	logger.Info("job queue worker started")
	logger.Error("connection pool exhausted")
	logger.Warn("disk usage above 80 percent")
	zap.L().Debug("health check passed")
}

func testAllRulesLowercase() {
	// заглавная буква в начале
	log.Print("Server started")              // want "the log message must start with lowercase letter"
	slog.Error("Database connection failed") // want "the log message must start with lowercase letter"

	logger, _ := zap.NewProduction()
	logger.Info("Worker started")     // want "the log message must start with lowercase letter"
	zap.L().Warn("High memory usage") // want "the log message must start with lowercase letter"
}

func testAllRulesEnglish() {
	// не английский язык
	log.Print("сервер запущен")           // want "the log message must be in english"
	slog.Error("ошибка подключения к бд") // want "the log message must be in english"

	logger, _ := zap.NewProduction()
	logger.Info("воркер остановлен") // want "the log message must be in english"
	zap.L().Warn("память на исходе") // want "the log message must be in english"
}

func testAllRulesSpecialSymbols() {
	// спецсимволы и эмодзи
	log.Print("server started!")       // want "the log message must not contain any special symbols"
	slog.Error("connection failed...") // want "the log message must not contain any special symbols"
	log.Println("deploy done 🚀")       // want "the log message must not contain any special symbols"

	logger, _ := zap.NewProduction()
	logger.Info("all systems go ✅")     // want "the log message must not contain any special symbols"
	logger.Error("critical failure!!!") // want "the log message must not contain any special symbols"
}

func testAllRulesSensitive() {
	password := "hunter2"
	apiKey := "sk-live-abc"
	token := "eyJhbGci"

	// чувствительные данные
	log.Println("password : " + password) // want "the log message must not contain any sensitive data: password"
	slog.Info("api_key=" + apiKey)        // want "the log message must not contain any sensitive data: api_key"
	slog.Error("token: " + token)         // want "the log message must not contain any sensitive data: token"

	logger, _ := zap.NewProduction()
	logger.Error("user password: " + password) // want "the log message must not contain any sensitive data: password"
	zap.L().Info("bearer is" + token)          // want "the log message must not contain any sensitive data: bearer"
}

func testAllRulesEdgeCases() {
	// граничные случаи
	log.Print("")
	log.Println("   ")
	log.Printf("attempt %d of %d", 1, 3)
	slog.Info("retry-job triggered")
	slog.Debug("task_queue length is 0")

	logger, _ := zap.NewProduction()
	logger.Info("graceful shutdown initiated")
	logger.Warn("rate limit approaching threshold")
}
