package mocks


import (
	"github.com/fire9900/forum/internal/models"
	"github.com/stretchr/testify/mock"
)

type ForumUseCase struct {
	mock.Mock
}

func (m *ForumUseCase) GetAllThreads() ([]models.Thread, error) {
	args := m.Called()
	return args.Get(0).([]models.Thread), args.Error(1)
}

func (m *ForumUseCase) GetThreadByID(id int) (models.Thread, error) {
	args := m.Called(id)
	return args.Get(0).(models.Thread), args.Error(1)
}

func (m *ForumUseCase) CreateThread(thread models.Thread) (models.Thread, error) {
	args := m.Called(thread)
	return args.Get(0).(models.Thread), args.Error(1)
}

func (m *ForumUseCase) DeleteThreadByID(id int, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *ForumUseCase) CreatePost(post models.Post) (models.Post, error) {
	args := m.Called(post)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *ForumUseCase) GetChatPosts(threadID int) ([]models.Post, error) {
	args := m.Called(threadID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *ForumUseCase) GetPostByThreadID(threadID int) ([]models.Post, error) {
	args := m.Called(threadID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *ForumUseCase) DeletePostByID(id int, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *ForumUseCase) GetPostsByUserID(id int) ([]models.Post, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *ForumUseCase) CheckUserByID(any models.User, id int) (bool, error) {
	args := m.Called(any, id)
	return args.Bool(0), args.Error(1)
}

func (m *ForumUseCase) GetUserThreads(userId int) ([]models.Thread, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.Thread), args.Error(1)
}

func (m *ForumUseCase) EditThread(thread models.Thread, userID int) error {
	args := m.Called(thread, userID)
	return args.Error(0)
}
