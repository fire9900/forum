package models

import (
	"errors"
	"time"
)

var (
	ErrorNotFoundThread = errors.New("Тред не найден")
	ErrorNotFoundPost   = errors.New("Пост не найден")
	ErrorNotFoundUser   = errors.New("Пользователь не найден")
)

type User interface {
	USER_ID() int
}

type Thread struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
	UserID   int       `json:"user_ID"`
}

type Post struct {
	ID       int       `json:"id"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
	ThreadID int       `json:"thread_id"`
	UserID   int       `json:"user_id"`
}

type Chat struct {
	ThreadID int `json:"thread_id"`
	UserID   int `json:"user_id"`
	PostID   int `json:"post_id"`
}

func (t Thread) USER_ID() int {
	return t.UserID
}

func (t Post) USER_ID() int {
	return t.UserID
}

func (t Chat) USER_ID() int {
	return t.UserID
}
