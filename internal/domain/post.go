package domain

import (
	"time"

	"github.com/tmozzze/SasPosts/utils"
)

type Post struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Author        string    `json:"author"`
	CreatedAt     time.Time `json:"createdAt"`
	AllowComments bool      `json:"allowComments"`
}

func NewPost(title, content, author string, allowComments bool) *Post {
	return &Post{
		ID:            utils.GenerateID(),
		Title:         title,
		Content:       content,
		Author:        author,
		CreatedAt:     time.Now(),
		AllowComments: allowComments,
	}
}
