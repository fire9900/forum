package usecase

import (
	"fmt"
	"github.com/fire9900/forum/internal/models"
	"github.com/fire9900/forum/internal/repository"
	"github.com/fire9900/forum/pkg/logger"
	"go.uber.org/zap"
)

type ThreadUseCase interface {
	GetUserThreads(userId int) ([]models.Thread, error)
	GetAllThreads() ([]models.Thread, error)
	GetThreadByID(id int) (models.Thread, error)
	CreateThread(thread models.Thread) (models.Thread, error)
	DeleteThreadByID(id int, userID int) error
	EditThread(thread models.Thread, userID int) error
	CheckUserByID(any models.User, id int) (bool, error)
}

type TUseCase struct {
	repo repository.ForumRepository
}

func NewThreadUseCase(repo repository.ForumRepository) ThreadUseCase {
	return &TUseCase{repo: repo}
}

func (f *TUseCase) CheckUserByID(any models.User, id int) (bool, error) {
	return f.repo.CheckUserByID(any, id)
}

func (f *TUseCase) GetUserThreads(userId int) ([]models.Thread, error) {
	return f.repo.GetThreadsByUserID(userId)
}

func (f *TUseCase) EditThread(thread models.Thread, userID int) error {
	return f.repo.EditThread(thread, userID)
}

func (f *TUseCase) GetAllThreads() ([]models.Thread, error) {
	logger.Logger.Debug("Получение всех тредов")
	threads, err := f.repo.GetAllThreads()
	if err != nil {
		logger.Logger.Error("Ошибка при получении всех тредов",
			zap.Error(err))
		return nil, err
	}
	logger.Logger.Debug("Успешно получены все треды",
		zap.Int("count", len(threads)))
	return threads, nil
}

func (f *TUseCase) GetThreadByID(id int) (models.Thread, error) {
	logger.Logger.Debug("Получение треда по ID", zap.Int("id", id))
	thread, err := f.repo.GetThreadByID(id)
	if err != nil {
		logger.Logger.Error("Ошибка при получении треда",
			zap.Int("id", id),
			zap.Error(err))
		return models.Thread{}, err
	}
	logger.Logger.Debug("Тред успешно получен",
		zap.Int("id", id),
		zap.String("title", thread.Title))
	return thread, nil
}

func (f *TUseCase) CreateThread(thread models.Thread) (models.Thread, error) {
	logger.Logger.Debug("Проверка валидности данных треда",
		zap.Int("userID", thread.UserID),
		zap.String("title", thread.Title))

	if thread.Content == "" || len(thread.Content) > 5000 {
		err := fmt.Errorf("Недопустимый размер описания! Описание == 0 || > 5000")
		logger.Logger.Error("Невалидное содержание треда",
			zap.Error(err),
			zap.Int("contentLength", len(thread.Content)))
		return models.Thread{}, err
	}
	if thread.Title == "" || len(thread.Title) > 500 {
		err := fmt.Errorf("Недопустимый размер заголовка! Заголовк == 0 || > 1000")
		logger.Logger.Error("Невалидный заголовок треда",
			zap.Error(err),
			zap.Int("titleLength", len(thread.Title)))
		return models.Thread{}, err
	}

	logger.Logger.Info("Создание нового треда",
		zap.Int("userID", thread.UserID),
		zap.String("title", thread.Title))

	createdThread, err := f.repo.CreateThread(thread)
	if err != nil {
		logger.Logger.Error("Ошибка при создании треда",
			zap.Any("thread", thread),
			zap.Error(err))
		return models.Thread{}, err
	}

	logger.Logger.Info("Тред успешно создан",
		zap.Int("id", createdThread.ID),
		zap.String("title", createdThread.Title))
	return createdThread, nil
}

func (f *TUseCase) DeleteThreadByID(id int, userID int) error {
	logger.Logger.Info("Удаление треда", zap.Int("id", id))

	thread, err := f.repo.GetThreadByID(id)
	if err != nil {
		return err
	}

	valid, err := f.CheckUserByID(thread, userID)
	if !valid || err != nil {
		if !valid {
			return fmt.Errorf("Нет прав доступа")
		}
		return fmt.Errorf("Нет прав доступа, тк ошибка")
	}

	if err := f.repo.DeleteThreadByID(id); err != nil {
		logger.Logger.Error("Ошибка при удалении треда",
			zap.Int("id", id),
			zap.Error(err))
		return err
	}
	logger.Logger.Info("Тред успешно удален", zap.Int("id", id))
	return nil
}
