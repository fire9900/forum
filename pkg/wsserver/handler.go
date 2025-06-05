package wsserver

import (
	"encoding/json"
	"github.com/fire9900/forum/internal/models"
	"github.com/fire9900/forum/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (hub *Hub) ThreadChat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Logger.Error("Ошибка при переходе на WebSocket соединение",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Некорректный ID треда в WebSocket запросе",
			zap.Error(err),
			zap.String("параметр", c.Param("id")))
		conn.WriteMessage(websocket.CloseMessage, []byte("Invalid thread ID"))
		conn.Close()
		return
	}

	logger.Logger.Info("Новое WebSocket соединение",
		zap.Int("threadID", id))

	client := &Client{
		conn:     conn,
		send:     make(chan models.Post, 256),
		threadID: id,
	}

	hub.register <- client

	go func() {
		posts, err := hub.UseCase.GetChatPosts(id)
		if err != nil {
			logger.Logger.Error("Ошибка при получении сообщений чата",
				zap.Int("threadID", id),
				zap.Error(err))
			return
		}

		logger.Logger.Debug("Отправка истории сообщений новому клиенту",
			zap.Int("threadID", id),
			zap.Int("количество сообщений", len(posts)))

		for _, post := range posts {
			select {
			case client.send <- post:
			default:
				logger.Logger.Warn("Канал отправки переполнен, отключаем клиента",
					zap.Int("threadID", id))
				hub.unregister <- client
				conn.Close()
				return
			}
		}
	}()

	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
			logger.Logger.Info("WebSocket соединение закрыто",
				zap.Int("threadID", id))
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Logger.Debug("Ошибка чтения сообщения из WebSocket",
					zap.Int("threadID", id),
					zap.Error(err))
				break
			}

			var post models.Post
			if err := json.Unmarshal(message, &post); err != nil {
				logger.Logger.Warn("Некорректный формат сообщения",
					zap.Int("threadID", id),
					zap.Error(err))
				conn.WriteJSON(map[string]string{"error": "invalid message format"})
				continue
			}

			post.ThreadID = id
			createdPost, err := hub.UseCase.CreatePost(post)
			if err != nil {
				logger.Logger.Error("Ошибка при создании сообщения",
					zap.Int("threadID", id),
					zap.Error(err))
				conn.WriteJSON(map[string]string{"error": "failed to create post"})
				continue
			}

			logger.Logger.Debug("Новое сообщение создано",
				zap.Int("threadID", id),
				zap.Int("userID", post.UserID),
				zap.String("content", post.Content))

			hub.chat <- createdPost
		}
	}()

	go func() {
		defer conn.Close()

		for message := range client.send {
			postBytes, err := json.Marshal(message)
			if err != nil {
				logger.Logger.Error("Ошибка при сериализации сообщения",
					zap.Int("threadID", message.ThreadID),
					zap.Error(err))
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, postBytes); err != nil {
				logger.Logger.Debug("Ошибка отправки сообщения через WebSocket",
					zap.Int("threadID", message.ThreadID),
					zap.Error(err))
				break
			}
		}
	}()
}
