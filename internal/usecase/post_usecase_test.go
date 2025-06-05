package usecase

import (
	"errors"
	"github.com/fire9900/forum/internal/models"
	"github.com/fire9900/forum/pkg/logger"
	"github.com/fire9900/forum/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func init() {
	logger.Logger = zap.NewNop()
}

func TestGetAllThreads(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	mockThreads := []models.Thread{
		{ID: 1, Title: "Test Thread 1"},
		{ID: 2, Title: "Test Thread 2"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetAllThreads").Return(mockThreads, nil).Once()

		u := NewPostUseCase(mockRepo)
		threads, err := u.GetAllThreads()

		assert.NoError(t, err)
		assert.Equal(t, mockThreads, threads)
		mockRepo.AssertExpectations(t)
	})

}
func TestGetThreadByID(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	mockThread := models.Thread{ID: 1, Title: "Test Thread"}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetThreadByID", 1).Return(mockThread, nil).Once()

		u := NewPostUseCase(mockRepo)
		thread, err := u.GetThreadByID(1)

		assert.NoError(t, err)
		assert.Equal(t, mockThread, thread)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("GetThreadByID", 2).Return(models.Thread{}, errors.New("error")).Once()

		u := NewPostUseCase(mockRepo)
		thread, err := u.GetThreadByID(2)

		assert.Error(t, err)
		assert.Equal(t, models.Thread{}, thread)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateThread(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	validThread := models.Thread{
		Title:   "Valid Title",
		Content: "Valid content",
		UserID:  1,
	}
	createdThread := models.Thread{
		ID:      1,
		Title:   "Valid Title",
		Content: "Valid content",
		UserID:  1,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CreateThread", validThread).Return(createdThread, nil).Once()

		u := NewPostUseCase(mockRepo)
		result, err := u.CreateThread(validThread)

		assert.NoError(t, err)
		assert.Equal(t, createdThread, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty title", func(t *testing.T) {
		invalidThread := models.Thread{
			Title:   "",
			Content: "Valid content",
			UserID:  1,
		}

		u := NewPostUseCase(mockRepo)
		_, err := u.CreateThread(invalidThread)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "CreateThread")
	})

	t.Run("empty content", func(t *testing.T) {
		invalidThread := models.Thread{
			Title:   "Valid Title",
			Content: "",
			UserID:  1,
		}

		u := NewPostUseCase(mockRepo)
		_, err := u.CreateThread(invalidThread)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "CreateThread")
	})
}

func TestDeleteThreadByID(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	thread := models.Thread{ID: 1, UserID: 1}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetThreadByID", 1).Return(thread, nil).Once()
		mockRepo.On("CheckUserByID", mock.Anything, 1).Return(true, nil).Once()
		mockRepo.On("DeleteThreadByID", 1).Return(nil).Once()

		u := NewPostUseCase(mockRepo)
		err := u.DeleteThreadByID(1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("no permissions", func(t *testing.T) {
		mockRepo.On("GetThreadByID", 1).Return(thread, nil).Once()
		mockRepo.On("CheckUserByID", mock.Anything, 2).Return(false, nil).Once()

		u := NewPostUseCase(mockRepo)
		err := u.DeleteThreadByID(1, 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "DeleteThreadByID")
	})
}

func TestCreatePost(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	validPost := models.Post{
		Content:  "Valid content",
		ThreadID: 1,
		UserID:   1,
	}
	createdPost := models.Post{
		ID:       1,
		Content:  "Valid content",
		ThreadID: 1,
		UserID:   1,
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("CreatePost", validPost).Return(createdPost, nil).Once()
		mockRepo.On("LinkPostToChat", mock.Anything).Return(nil).Once()

		u := NewPostUseCase(mockRepo)
		result, err := u.CreatePost(validPost)

		assert.NoError(t, err)
		assert.Equal(t, createdPost, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty content", func(t *testing.T) {
		invalidPost := models.Post{
			Content:  "",
			ThreadID: 1,
			UserID:   1,
		}

		u := NewPostUseCase(mockRepo)
		_, err := u.CreatePost(invalidPost)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "CreatePost")
	})
}

func TestGetChatPosts(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	mockPosts := []models.Post{
		{ID: 1, Content: "Post 1"},
		{ID: 2, Content: "Post 2"},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetChatPosts", 1).Return(mockPosts, nil).Once()

		u := NewPostUseCase(mockRepo)
		posts, err := u.GetChatPosts(1)

		assert.NoError(t, err)
		assert.Equal(t, mockPosts, posts)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeletePostByID(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	post := models.Post{ID: 1, UserID: 1}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetPostByID", 1).Return(post, nil).Once()
		mockRepo.On("CheckUserByID", mock.Anything, 1).Return(true, nil).Once()
		mockRepo.On("DeletePostByID", 1).Return(nil).Once()

		u := NewPostUseCase(mockRepo)
		err := u.DeletePostByID(1, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("no permissions", func(t *testing.T) {
		mockRepo.On("GetPostByID", 1).Return(post, nil).Once()
		mockRepo.On("CheckUserByID", mock.Anything, 2).Return(false, nil).Once()

		u := NewPostUseCase(mockRepo)
		err := u.DeletePostByID(1, 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "DeletePostByID")
	})
}

func TestGetUserThreads(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	mockThreads := []models.Thread{
		{ID: 1, Title: "User Thread 1", UserID: 1},
		{ID: 2, Title: "User Thread 2", UserID: 1},
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("GetThreadsByUserID", 1).Return(mockThreads, nil).Once()

		u := NewPostUseCase(mockRepo)
		threads, err := u.GetUserThreads(1)

		assert.NoError(t, err)
		assert.Equal(t, mockThreads, threads)
		mockRepo.AssertExpectations(t)
	})
}

func TestEditThread(t *testing.T) {
	mockRepo := new(mocks.ForumRepository)
	thread := models.Thread{ID: 1, Title: "New Title", Content: "New Content", UserID: 1}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("EditThread", thread, 1).Return(nil).Once()

		u := NewPostUseCase(mockRepo)
		err := u.EditThread(thread, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockRepo.On("EditThread", thread, 2).Return(errors.New("error")).Once()

		u := NewPostUseCase(mockRepo)
		err := u.EditThread(thread, 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
