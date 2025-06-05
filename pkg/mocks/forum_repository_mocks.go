package mocks

import (
	"github.com/fire9900/forum/internal/models"
	"github.com/stretchr/testify/mock"
)

type ForumRepository struct {
	mock.Mock
}

func (m *ForumRepository) GetAllThreads() ([]models.Thread, error) {
	args := m.Called()
	return args.Get(0).([]models.Thread), args.Error(1)
}

func (m *ForumRepository) GetThreadByID(id int) (models.Thread, error) {
	args := m.Called(id)
	return args.Get(0).(models.Thread), args.Error(1)
}

func (m *ForumRepository) CreateThread(thread models.Thread) (models.Thread, error) {
	args := m.Called(thread)
	return args.Get(0).(models.Thread), args.Error(1)
}

func (m *ForumRepository) DeleteThreadByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ForumRepository) GetThreadsByUserID(userId int) ([]models.Thread, error) {
	args := m.Called(userId)
	return args.Get(0).([]models.Thread), args.Error(1)
}

func (m *ForumRepository) CreatePost(post models.Post) (models.Post, error) {
	args := m.Called(post)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *ForumRepository) GetPostsByThreadID(threadID int) ([]models.Post, error) {
	args := m.Called(threadID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *ForumRepository) DeletePostByID(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ForumRepository) GetPostsByUserID(id int) ([]models.Post, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *ForumRepository) GetChatPosts(threadID int) ([]models.Post, error) {
	args := m.Called(threadID)
	return args.Get(0).([]models.Post), args.Error(1)
}

func (m *ForumRepository) LinkPostToChat(chat models.Chat) error {
	args := m.Called(chat)
	return args.Error(0)
}

func (m *ForumRepository) CheckUserByID(user models.User, id int) (bool, error) {
	args := m.Called(user, id)
	return args.Bool(0), args.Error(1)
}

func (m *ForumRepository) GetPostByID(id int) (models.Post, error) {
	args := m.Called(id)
	return args.Get(0).(models.Post), args.Error(1)
}

func (m *ForumRepository) EditThread(thread models.Thread, userID int) error {
	args := m.Called(thread, userID)
	return args.Error(0)
}
