package app

import (
	"github.com/fire9900/forum/internal/repository"
	"github.com/fire9900/forum/internal/transport/gin"
	"github.com/fire9900/forum/internal/usecase"
	"github.com/fire9900/forum/pkg/database"
	"github.com/fire9900/forum/pkg/logger"
	"github.com/fire9900/forum/pkg/wsserver"
	"go.uber.org/zap"
)

func RunMain() {
	db, err := database.NewSQLiteConnection()
	if err != nil {
		logger.Logger.Fatal("Ошибка подключения к базе данных",
			zap.Error(err),
			zap.String("component", "database"))
	}
	logger.Logger.Info("Подключение к базе данных прошло успешно")

	forumRepo := repository.NewForumRepository(db, logger.Logger)
	p := usecase.NewPostUseCase(forumRepo)
	t := usecase.NewThreadUseCase(forumRepo)
	hub := wsserver.NewHub(p, logger.Logger)

	router := gin.SetupRouter(p, t, ClientStart(), hub)
	logger.Logger.Info("Сервер стартует на порту :7777")
	if err := router.Run(":7777"); err != nil {
		logger.Logger.Fatal("Ошибка запуска сервера",
			zap.Error(err),
			zap.String("component", "http-server"))
	}
}
