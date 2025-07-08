package domain

import (
	"time"
)

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorId      string    `json:"authorId"`
	CreatedAt     time.Time `json:"createdAt"`
	AllowComments bool      `json:"allowComments"`
}

func NewPost(title, content, authorId string, allowComments bool) *Post {
	return &Post{
		ID:            generateID(),
		Title:         title,
		Content:       content,
		AuthorId:      authorId,
		CreatedAt:     time.Now(),
		AllowComments: allowComments,
	}
}

//TODO: generateID
