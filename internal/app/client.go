package app

import (
	"github.com/fire9900/auth/pkg/client"
	"github.com/fire9900/forum/pkg/logger"
	"go.uber.org/zap"
)

func ClientStart() *client.AuthClient {
	authClient, err := client.NewAuthClient("localhost:50051")
	if err != nil {
		logger.Logger.Fatal("не удалось инициализировать auth-client",
			zap.Error(err),
			zap.String("component", "auth-client"))
	}

	logger.Logger.Info("Успешное подключение к auth-client")
	return authClient
}
