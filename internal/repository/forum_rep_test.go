package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fire9900/forum/internal/models"
	"go.uber.org/zap"
	"testing"
	"time"
)

func setupLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func Test_forumRepository_GetAllThreads(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	rows := sqlmock.NewRows([]string{"id", "title", "content", "create_at", "user_id"}).
		AddRow(1, "Thread 1", "Content 1", time.Now(), 1).
		AddRow(2, "Thread 2", "Content 2", time.Now(), 2)

	mock.ExpectQuery("SELECT id, title, content, create_at, user_id FROM threads ORDER BY create_at DESC").
		WillReturnRows(rows)

	threads, err := repo.GetAllThreads()
	if err != nil {
		t.Errorf("ошибка не ожидалась при получении тем: %s", err)
	}

	if len(threads) != 2 {
		t.Errorf("ожидалось 2 темы, получено %d", len(threads))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_GetThreadByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testID := 1
	rows := sqlmock.NewRows([]string{"id", "title", "content", "create_at", "user_id"}).
		AddRow(testID, "Test Thread", "Test Content", time.Now(), 1)

	mock.ExpectQuery("SELECT id, title, content, create_at, user_id FROM threads WHERE id = \\$1 ORDER BY create_at DESC").
		WithArgs(testID).
		WillReturnRows(rows)

	thread, err := repo.GetThreadByID(testID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при получении темы: %s", err)
	}

	if thread.ID != testID {
		t.Errorf("ожидался ID темы %d, получено %d", testID, thread.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_CreateThread(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	newThread := models.Thread{
		Title:   "New Thread",
		Content: "New Content",
		UserID:  1,
	}

	mock.ExpectQuery("INSERT INTO threads").
		WithArgs(newThread.Title, newThread.Content, sqlmock.AnyArg(), newThread.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "create_at", "user_id"}).AddRow(1, newThread.Title, newThread.Content, time.Now(), newThread.UserID))

	createdThread, err := repo.CreateThread(newThread)
	if err != nil {
		t.Errorf("ошибка не ожидалась при создании темы: %s", err)
	}

	if createdThread.Title != newThread.Title {
		t.Errorf("ожидалось название темы %s, получено %s", newThread.Title, createdThread.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_DeleteThreadByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testID := 1
	mock.ExpectExec("DELETE FROM threads WHERE id = \\$1").
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteThreadByID(testID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при удалении темы: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_CreatePost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	newPost := models.Post{
		Content:  "New Post",
		ThreadID: 1,
		UserID:   1,
		CreateAt: time.Now(),
	}

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(newPost.Content, newPost.CreateAt, newPost.ThreadID, newPost.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "content", "create_at", "thread_id", "user_id"}).
			AddRow(1, newPost.Content, newPost.CreateAt, newPost.ThreadID, newPost.UserID))

	createdPost, err := repo.CreatePost(newPost)
	if err != nil {
		t.Errorf("ошибка не ожидалась при создании поста: %s", err)
	}

	if createdPost.Content != newPost.Content {
		t.Errorf("ожидалось содержимое поста %s, получено %s", newPost.Content, createdPost.Content)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_GetPostsByThreadID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testThreadID := 1
	rows := sqlmock.NewRows([]string{"id", "content", "create_at", "thread_id", "user_id"}).
		AddRow(1, "Post 1", time.Now(), testThreadID, 1).
		AddRow(2, "Post 2", time.Now(), testThreadID, 2)

	mock.ExpectQuery("SELECT id, content, create_at, thread_id, user_id FROM posts WHERE thread_id = \\$1").
		WithArgs(testThreadID).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByThreadID(testThreadID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при получении постов: %s", err)
	}

	if len(posts) != 2 {
		t.Errorf("ожидалось 2 поста, получено %d", len(posts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_GetPostsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testUserID := 1
	rows := sqlmock.NewRows([]string{"id", "content", "create_at", "thread_id", "user_id"}).
		AddRow(1, "Post 1", time.Now(), 1, testUserID).
		AddRow(2, "Post 2", time.Now(), 2, testUserID)

	mock.ExpectQuery("SELECT id, content, create_at, thread_id, user_id FROM posts WHERE user_id = \\$1 ORDER BY create_at DESC").
		WithArgs(testUserID).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByUserID(testUserID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при получении постов: %s", err)
	}

	if len(posts) != 2 {
		t.Errorf("ожидалось 2 поста, получено %d", len(posts))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_GetPostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testPostID := 1
	rows := sqlmock.NewRows([]string{"id", "content", "create_at", "thread_id", "user_id"}).
		AddRow(testPostID, "Test Post", time.Now(), 1, 1)

	mock.ExpectQuery("SELECT \\* FROM posts WHERE id = \\$1").
		WithArgs(testPostID).
		WillReturnRows(rows)

	post, err := repo.GetPostByID(testPostID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при получении поста: %s", err)
	}

	if post.ID != testPostID {
		t.Errorf("ожидался ID поста %d, получено %d", testPostID, post.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_DeletePostByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testPostID := 1
	mock.ExpectExec("DELETE FROM posts WHERE id = \\$1").
		WithArgs(testPostID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeletePostByID(testPostID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при удалении поста: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_GetThreadsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	testUserID := 1
	rows := sqlmock.NewRows([]string{"id", "title", "content", "create_at", "user_id"}).
		AddRow(1, "Thread 1", "Content 1", time.Now(), testUserID).
		AddRow(2, "Thread 2", "Content 2", time.Now(), testUserID)

	mock.ExpectQuery("SELECT \\* FROM threads WHERE user_ID = \\$1 ORDER BY create_at DESC").
		WithArgs(testUserID).
		WillReturnRows(rows)

	threads, err := repo.GetThreadsByUserID(testUserID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при получении тем: %s", err)
	}

	if len(threads) != 2 {
		t.Errorf("ожидалось 2 темы, получено %d", len(threads))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}

func Test_forumRepository_CheckUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка '%s' не ожидалась при открытии заглушки подключения к базе данных", err)
	}
	defer db.Close()

	logger := setupLogger()
	repo := NewForumRepository(db, logger)

	adminID := 1
	rows := sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(adminID, "admin", "admin@test.com", "admin")

	mock.ExpectQuery("SELECT id, name, email, role FROM users WHERE id = \\$1").
		WithArgs(adminID).
		WillReturnRows(rows)

	valid, err := repo.CheckUserByID(models.Thread{UserID: adminID}, adminID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при проверке пользователя-администратора: %s", err)
	}
	if !valid {
		t.Error("ожидалось, что администратор будет валидным")
	}

	userID := 2
	rows = sqlmock.NewRows([]string{"id", "name", "email", "role"}).
		AddRow(userID, "user", "user@test.com", "user")

	mock.ExpectQuery("SELECT id, name, email, role FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(rows)

	valid, err = repo.CheckUserByID(models.Thread{UserID: userID}, userID)
	if err != nil {
		t.Errorf("ошибка не ожидалась при проверке обычного пользователя: %s", err)
	}
	if !valid {
		t.Error("ожидалось, что обычный пользователь будет валидным, когда ID совпадают")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("не все ожидания были выполнены: %s", err)
	}
}
