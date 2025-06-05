package usecase

import (
	"fmt"
	entity "github.com/fire9900/forum/internal/models"
	"github.com/fire9900/forum/internal/repository"
	"github.com/fire9900/forum/pkg/logger"
	"go.uber.org/zap"
)

type PostUseCase interface {
	CreatePost(post entity.Post) (entity.Post, error)
	GetChatPosts(threadID int) ([]entity.Post, error)
	GetPostByThreadID(threadID int) ([]entity.Post, error)
	DeletePostByID(id int, userID int) error
	CheckUserByID(any entity.User, id int) (bool, error)
	GetPostsByUserID(id int) ([]entity.Post, error)
}

type PUseCase struct {
	repo repository.ForumRepository
}

func NewPostUseCase(repo repository.ForumRepository) *PUseCase {
	return &PUseCase{repo: repo}
}

func (f *PUseCase) CreatePost(post entity.Post) (entity.Post, error) {
	if post.Content == "" || len(post.Content) > 5000 {
		err := fmt.Errorf("Недопустимый размер описания! Описание == 0 || > 5000")
		logger.Logger.Error("Невалидное содержание поста",
			zap.Error(err),
			zap.Int("contentLength", len(post.Content)))
		return entity.Post{}, err
	}

	createdPost, err := f.repo.CreatePost(post)
	if err != nil {
		return entity.Post{}, err
	}

	if err := f.repo.LinkPostToChat(entity.Chat{
		post.ThreadID,
		post.UserID,
		createdPost.ID,
	}); err != nil {
		return entity.Post{}, fmt.Errorf("Ошибка создания поста в чат: %w", err)
	}
	return createdPost, nil
}

func (f *PUseCase) GetChatPosts(threadID int) ([]entity.Post, error) {
	return f.repo.GetChatPosts(threadID)
}

func (f *PUseCase) GetPostByThreadID(threadID int) ([]entity.Post, error) {
	logger.Logger.Debug("Получение постов по ID треда", zap.Int("threadID", threadID))
	posts, err := f.repo.GetPostsByThreadID(threadID)
	if err != nil {
		logger.Logger.Error("Ошибка при получении постов треда",
			zap.Int("threadID", threadID),
			zap.Error(err))
		return nil, err
	}
	logger.Logger.Debug("Посты треда успешно получены",
		zap.Int("threadID", threadID),
		zap.Int("count", len(posts)))
	return posts, nil
}

func (f *PUseCase) DeletePostByID(id int, userID int) error {
	logger.Logger.Info("Удаление поста", zap.Int("id", id))

	post, err := f.repo.GetPostByID(id)
	if err != nil {
		return err
	}

	valid, err := f.repo.CheckUserByID(post, userID)
	if !valid || err != nil {
		return fmt.Errorf("Нет прав: %w", err)
	}

	err = f.repo.DeletePostByID(id)
	if err != nil {
		logger.Logger.Error("Ошибка при удалении поста",
			zap.Int("id", id),
			zap.Error(err))
		return err
	}
	logger.Logger.Info("Пост успешно удален", zap.Int("id", id))
	return nil
}

func (f *PUseCase) CheckUserByID(any entity.User, id int) (bool, error) {
	return f.repo.CheckUserByID(any, id)
}

func (f *PUseCase) GetPostsByUserID(id int) ([]entity.Post, error) {
	return f.repo.GetPostsByUserID(id)
}
