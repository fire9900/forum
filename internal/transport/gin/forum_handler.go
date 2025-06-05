package gin

import (
	"fmt"
	"github.com/fire9900/forum/internal/models"
	"github.com/fire9900/forum/internal/usecase"
	"github.com/fire9900/forum/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// @title sigma Forum API
// @version 1.0
// @description API для постов, тредов и чатов

// @host localhost:7777
// @BasePath /api/v2
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type ForumHandler struct {
	threadCase usecase.ThreadUseCase
	postCase   usecase.PostUseCase
}

func NewForumHandler(P usecase.PostUseCase, T usecase.ThreadUseCase) *ForumHandler {
	return &ForumHandler{
		threadCase: T,
		postCase:   P,
	}
}

// @Summary Получить все треды
// @Description Получить список всех тредов
// @Tags threads
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Thread
// @Failure 400 {object} object
// @Router /threads [get]]
func (h *ForumHandler) GetAllThread(c *gin.Context) {
	threads, err := h.threadCase.GetAllThreads()
	if err != nil {
		logger.Logger.Error("Ошибка получения всех тредов",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.Info("Успешное получение всех тредов",
		zap.Int("количество", len(threads)))
	c.JSON(http.StatusOK, threads)
}

// @Summary Получить тред по ID
// @Description Получить тред по его идентификатору
// @Tags threads
// @Accept json
// @Produce json
// @Param id path int true "ID треда"
// @Success 200 {object} models.Thread
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Router /thread/{id} [get]
func (h *ForumHandler) GetThreadByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Ошибка конвертации ID треда",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	thread, err := h.threadCase.GetThreadByID(id)
	if err != nil {
		logger.Logger.Error("Ошибка получения треда по ID",
			zap.Int("id", id),
			zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Тред не найден"})
		return
	}

	logger.Logger.Info("Успешное получение треда",
		zap.Int("id", id))
	c.JSON(http.StatusOK, thread)
}

// @Summary Создать тред
// @Description Создать новый тред
// @Tags threads
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param thread body models.Thread true "Данные треда"
// @Success 200 {object} models.Thread
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /threads [post]
func (h *ForumHandler) CreateThread(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		logger.Logger.Warn("Попытка создания треда без авторизации")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Требуется авторизация",
		})
		return
	}

	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		logger.Logger.Error("Ошибка парсинга тела запроса",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		logger.Logger.Error("Неверный тип userID",
			zap.Any("userID", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка сервера",
		})
		return
	}
	thread.UserID = uid

	createdThread, err := h.threadCase.CreateThread(thread)
	if err != nil {
		logger.Logger.Error("Ошибка создания треда",
			zap.Any("thread", thread),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания треда"})
		return
	}

	logger.Logger.Info("Тред успешно создан",
		zap.Int("id", createdThread.ID),
		zap.Int("userID", uid))
	c.JSON(http.StatusOK, createdThread)
}

// @Summary Удалить тред
// @Description Удалить тред по его идентификатору
// @Tags threads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID треда"
// @Success 200
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /thread/{id} [delete]
func (f *ForumHandler) DeleteTheadByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Неверный формат ID треда для удаления",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга id"})
		return
	}

	err = f.threadCase.DeleteThreadByID(id, uid)
	if err != nil {
		logger.Logger.Error("Ошибка удаления треда",
			zap.Int("id", id),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Ошибка удаления треда: %s", err.Error())})
		return
	}

	logger.Logger.Info("Тред успешно удален",
		zap.Int("id", id))
	c.JSON(http.StatusOK, nil)
}

// @Summary Редактировать тред
// @Description Редактировать тред
// @Tags thread
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post body models.Post true "тред"
// @Success 200 {object} models.Thread
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /threads [PATCH]
func (f *ForumHandler) EditThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат данных",
			"details": err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	if err := f.threadCase.EditThread(thread, uid); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Не удалось обновить тред",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Тред успешно обновлен",
		"thread":  thread,
	})
}

// @Summary Создать пост
// @Description Создать новый пост в треде
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param post body models.Post true "Данные поста"
// @Success 200 {object} models.Post
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /threads/posts [post]
func (h *ForumHandler) CreatePost(c *gin.Context) {
	var DTOPost struct {
		Content  string `json:"content"`
		ThreadID int    `json:"thread_id"`
		UserID   int    `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&DTOPost); err != nil {
		logger.Logger.Error("Ошибка парсинга тела запроса при создании поста",
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	post := models.Post{
		Content:  DTOPost.Content,
		ThreadID: DTOPost.ThreadID,
		UserID:   DTOPost.UserID,
		CreateAt: time.Now(),
	}

	createdPost, err := h.postCase.CreatePost(post)
	if err != nil {
		logger.Logger.Error("Ошибка создания поста",
			zap.Any("post", post),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания поста"})
		return
	}

	logger.Logger.Info("Пост успешно создан",
		zap.Int("id", createdPost.ID),
		zap.Int("threadID", DTOPost.ThreadID))
	c.JSON(http.StatusOK, createdPost)
}

// @Summary Получить посты треда
// @Description Получить все посты определенного треда
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "ID треда"
// @Success 200 {array} models.Post
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /thread/{id}/posts [get]
func (h *ForumHandler) GetPostsByThreadID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Неверный формат ID треда",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	posts, err := h.postCase.GetPostByThreadID(id)
	if err != nil {
		logger.Logger.Error("Ошибка получения постов треда",
			zap.Int("threadID", id),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов"})
		return
	}

	logger.Logger.Info("Посты треда успешно получены",
		zap.Int("threadID", id),
		zap.Int("количество", len(posts)))
	c.JSON(http.StatusOK, posts)
}

// @Summary Получить посты пользователя
// @Description Получить все посты определенного пользователя
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Success 200 {array} models.Post
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /posts/user/{id} [get]
func (h *ForumHandler) GetPostsByUserID(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		logger.Logger.Warn("Попытка получения постов без авторизации")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Требуется авторизация",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		logger.Logger.Error("Неверный тип userID",
			zap.Any("userID", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка сервера",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id != uid {
		logger.Logger.Warn("Несоответствие ID пользователя",
			zap.Int("paramID", id),
			zap.Int("userID", uid))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный ID пользователя",
		})
		return
	}

	posts, err := h.postCase.GetPostsByUserID(id)
	if err != nil {
		logger.Logger.Error("Ошибка получения постов пользователя",
			zap.Int("userID", id),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка получения постов",
		})
		return
	}

	logger.Logger.Info("Посты пользователя успешно получены",
		zap.Int("userID", id),
		zap.Int("количество", len(posts)))
	c.JSON(http.StatusOK, posts)
}

// @Summary Удалить пост
// @Description Удалить пост по его идентификатору
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID поста"
// @Success 200
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /posts/{id} [delete]
func (h *ForumHandler) DeletePostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Неверный формат ID поста",
			zap.String("id", c.Param("id")),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID"})
		return
	}

	userID, exit := c.Get("userID")
	if !exit {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка парсинга id"})
		return
	}

	if err := h.postCase.DeletePostByID(id, uid); err != nil {
		logger.Logger.Error("Ошибка удаления поста",
			zap.Int("postID", id),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления поста"})
		return
	}

	logger.Logger.Info("Пост успешно удален",
		zap.Int("postID", id))
	c.JSON(http.StatusOK, nil)
}

// @Summary Получить треды пользователя
// @Description Получить все треды определенного пользователя
// @Tags threads
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID пользователя"
// @Success 200 {array} models.Thread
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 404 {object} object
// @Router /threads/user/{id} [get]
func (h *ForumHandler) GetThreadsByUserID(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		logger.Logger.Warn("Попытка получения тредов без авторизации")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Требуется авторизация",
		})
		return
	}

	uid, ok := userID.(int)
	if !ok {
		logger.Logger.Error("Неверный тип userID",
			zap.Any("userID", userID))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка сервера",
		})
		return
	}

	paramID, err := strconv.Atoi(c.Param("id"))
	if paramID != uid {
		logger.Logger.Warn("Несоответствие ID пользователя",
			zap.Int("paramID", paramID),
			zap.Int("userID", uid))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный ID пользователя",
		})
		return
	}

	threads, err := h.threadCase.GetUserThreads(paramID)
	if err != nil {
		logger.Logger.Error("Ошибка получения тредов пользователя",
			zap.Int("userID", paramID),
			zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Треды не найдены"})
		return
	}

	logger.Logger.Info("Треды пользователя успешно получены",
		zap.Int("userID", paramID),
		zap.Int("количество", len(threads)))
	c.JSON(http.StatusOK, threads)
}

// @Summary Получить сообщения чата
// @Description Получить все сообщения чата в треде
// @Tags chat
// @Accept json
// @Produce json
// @Param thread_id path int true "ID треда"
// @Success 200 {array} models.Post
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /ws/threads/{thread_id} [get]
func (h *ForumHandler) GetChatPosts(c *gin.Context) {
	threadID, err := strconv.Atoi(c.Param("thread_id"))
	if err != nil {
		logger.Logger.Error("Неверный формат ID треда",
			zap.String("thread_id", c.Param("thread_id")),
			zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID треда"})
		return
	}

	posts, err := h.postCase.GetChatPosts(threadID)
	if err != nil {
		logger.Logger.Error("Ошибка получения сообщений чата",
			zap.Int("threadID", threadID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения сообщений"})
		return
	}

	logger.Logger.Info("Сообщения чата успешно получены",
		zap.Int("threadID", threadID),
		zap.Int("количество", len(posts)))
	c.JSON(http.StatusOK, posts)
}
