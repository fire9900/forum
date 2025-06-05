package models

import (
	"testing"
	"time"
)

var USER_ID int = 2

var thread = &Thread{
	ID:       1,
	Title:    "111",
	Content:  "222",
	CreateAt: time.Time{},
	UserID:   USER_ID,
}

var post = &Post{
	ID:       1,
	Content:  "222111",
	CreateAt: time.Time{},
	ThreadID: 1,
	UserID:   USER_ID,
}

var chat = &Chat{
	ThreadID: 1,
	UserID:   USER_ID,
	PostID:   1,
}

func TestThisByUser_OK(t *testing.T) {
	id1 := chat.USER_ID()
	id2 := thread.USER_ID()
	id3 := post.USER_ID()

	if id1 != USER_ID || id2 != USER_ID || id3 != USER_ID {
		t.Errorf("thread.USER_ID() = %d, post.USER_ID() = %d, chat.USER_ID() = %d", id1, id2, id3)
	}
}
