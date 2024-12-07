package main

import (
	"os"

	"github.com/KutsDenis/logzap"
	"go.uber.org/zap"

	"fincraft-finance/internal/config"
	"fincraft-finance/internal/infrastructure"
	"fincraft-finance/internal/interfaces"
	"fincraft-finance/internal/metrics"
	"fincraft-finance/internal/server"
	"fincraft-finance/internal/usecases"
)

func main() {
	// Инициализация логгера
	appEnv := os.Getenv("APP_ENV")
	logzap.Init(appEnv)
	defer logzap.Sync()
	log := logzap.Logger

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
		os.Exit(1)
	}
	log.Info("Configuration loaded successfully")

	// Метрики
	go metrics.StartMetricsServer(cfg.MetricsPort)

	// Подключение к базе данных
	db, err := infrastructure.NewDBConnection(cfg.DBDSN)
	if err != nil {
		log.Fatal("Failed to connect to the database", zap.Error(err))
		os.Exit(1)
	}
	//noinspection GoUnhandledErrorResult
	defer db.Close()
	log.Info("Database connection established")

	// Создание зависимостей
	incomeRepo := infrastructure.NewIncomeRepository(db)
	incomeUsecase := usecases.NewIncomeUseCase(incomeRepo)
	financeHandler := interfaces.NewFinanceHandler(incomeUsecase)

	// Запуск сервера
	if err := server.RunGRPCServer(cfg.GRPCPort, financeHandler); err != nil {
		log.Fatal("Failed to start gRPC server", zap.Error(err))
		os.Exit(1)
	}
}
