package domain

import (
	"time"

	"github.com/tmozzze/SasPosts/utils"
)

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      string    `json:"authorId"`
	CreatedAt     time.Time `json:"createdAt"`
	AllowComments bool      `json:"allowComments"`
}

func NewPost(title, content, authorID string, allowComments bool) *Post {
	return &Post{
		ID:            utils.GenerateID(),
		Title:         title,
		Content:       content,
		AuthorID:      authorID,
		CreatedAt:     time.Now(),
		AllowComments: allowComments,
	}
}
