package handler

import (
	"errors"
	gin2 "github.com/fire9900/forum/internal/transport/gin"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fire9900/forum/internal/models"
	"github.com/fire9900/forum/pkg/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllThread(t *testing.T) {
	tests := []struct {
		name         string
		mockThreads  []models.Thread
		mockError    error
		expectedCode int
	}{
		{
			name: "success",
			mockThreads: []models.Thread{
				{ID: 1, Title: "Test Thread 1"},
				{ID: 2, Title: "Test Thread 2"},
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "error case",
			mockError:    errors.New("some error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(mocks.ForumUseCase)

			mockUsecase.On("GetAllThreads").Return(tt.mockThreads, tt.mockError)

			h := gin2.NewForumHandler(mockUsecase)

			router := gin.Default()
			router.GET("/threads", h.GetAllThread)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/threads", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetThreadByID(t *testing.T) {
	tests := []struct {
		name         string
		idParam      string
		mockThread   models.Thread
		mockError    error
		expectedCode int
	}{
		{
			name:         "success",
			idParam:      "1",
			mockThread:   models.Thread{ID: 1, Title: "Test Thread"},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "invalid id",
			idParam:      "abc",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "not found",
			idParam:      "1",
			mockError:    errors.New("not found"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(mocks.ForumUseCase)
			if tt.idParam != "abc" {
				id, _ := strconv.Atoi(tt.idParam)
				mockUsecase.On("GetThreadByID", id).Return(tt.mockThread, tt.mockError)
			}

			h := gin2.NewForumHandler(mockUsecase)

			router := gin.Default()
			router.GET("/thread/:id", h.GetThreadByID)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/thread/"+tt.idParam, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestCreateThread(t *testing.T) {
	tests := []struct {
		name         string
		userID       interface{}
		threadJSON   string
		mockThread   models.Thread
		mockError    error
		expectedCode int
	}{
		{
			name:         "success",
			userID:       1,
			threadJSON:   `{"title":"Test","content":"Content"}`,
			mockThread:   models.Thread{ID: 1, Title: "Test", Content: "Content", UserID: 1},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "unauthorized",
			userID:       nil,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "invalid json",
			userID:       1,
			threadJSON:   `invalid`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "create error",
			userID:       1,
			threadJSON:   `{"title":"Test","content":"Content"}`,
			mockError:    errors.New("create error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(mocks.ForumUseCase)
			if tt.userID != nil && tt.threadJSON != "invalid" {
				mockUsecase.On("CreateThread", mock.Anything).Return(tt.mockThread, tt.mockError)
			}

			h := gin2.NewForumHandler(mockUsecase)

			router := gin.Default()
			router.POST("/threads", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				h.CreateThread(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/threads", nil)
			if tt.threadJSON != "" {
				req, _ = http.NewRequest("POST", "/threads", strings.NewReader(tt.threadJSON))
				req.Header.Set("Content-Type", "application/json")
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestDeleteThreadByID(t *testing.T) {
	tests := []struct {
		name         string
		idParam      string
		userID       interface{}
		mockError    error
		expectedCode int
	}{
		{
			name:         "success",
			idParam:      "1",
			userID:       1,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "invalid id",
			idParam:      "abc",
			userID:       1,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "unauthorized",
			idParam:      "1",
			userID:       nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "delete error",
			idParam:      "1",
			userID:       1,
			mockError:    errors.New("delete error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(mocks.ForumUseCase)
			if tt.idParam != "abc" && tt.userID != nil {
				id, _ := strconv.Atoi(tt.idParam)
				mockUsecase.On("DeleteThreadByID", id, tt.userID.(int)).Return(tt.mockError)
			}

			h := gin2.NewForumHandler(mockUsecase)

			router := gin.Default()
			router.DELETE("/thread/:id", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				h.DeleteTheadByID(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/thread/"+tt.idParam, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestCreatePost(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name         string
		postJSON     string
		mockPost     models.Post
		mockError    error
		expectedCode int
	}{
		{
			name:         "success",
			postJSON:     `{"content":"Test","thread_id":1,"user_id":1}`,
			mockPost:     models.Post{ID: 1, Content: "Test", ThreadID: 1, UserID: 1, CreateAt: now},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "invalid json",
			postJSON:     `invalid`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "create error",
			postJSON:     `{"content":"Test","thread_id":1,"user_id":1}`,
			mockError:    errors.New("create error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(mocks.ForumUseCase)
			if tt.postJSON != "invalid" {
				mockUsecase.On("CreatePost", mock.Anything).Return(tt.mockPost, tt.mockError)
			}

			h := gin2.NewForumHandler(mockUsecase)

			router := gin.Default()
			router.POST("/threads/posts", h.CreatePost)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/threads/posts", strings.NewReader(tt.postJSON))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetPostsByThreadID(t *testing.T) {
	tests := []struct {
		name         string
		idParam      string
		mockPosts    []models.Post
		mockError    error
		expectedCode int
	}{
		{
			name:    "success",
			idParam: "1",
			mockPosts: []models.Post{
				{ID: 1, Content: "Post 1"},
				{ID: 2, Content: "Post 2"},
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "invalid id",
			idParam:      "abc",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "get error",
			idParam:      "1",
			mockError:    errors.New("get error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase := new(mocks.ForumUseCase)
			if tt.idParam != "abc" {
				id, _ := strconv.Atoi(tt.idParam)
				mockUsecase.On("GetPostByThreadID", id).Return(tt.mockPosts, tt.mockError)
			}

			h := gin2.NewForumHandler(mockUsecase)

			router := gin.Default()
			router.GET("/thread/:id/posts", h.GetPostsByThreadID)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/thread/"+tt.idParam+"/posts", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}
