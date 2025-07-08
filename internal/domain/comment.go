package domain

import (
	"errors"
	"time"
	"unicode/utf8"
)

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	AuthorID  string    `json:"authorId"`
	ParentID  *string   `json:"parentId,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Path      string    `json:"path"`
	Depth     int       `json:"depth"`
}

const MaxCommentLength = 2000

func NewComment(postID, authorID, parentID *string, content string) (*Comment, error) {
	if utf8.RuneCountInString(content) > MaxCommentLength {
		return nil, ErrCommentTooLong
	}

	comment := &Comment{
		ID:        generateID(),
		PostID:    postID,
		AuthorID:  authorID,
		ParentID:  parentID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	return comment, nil
}

var ErrCommentTooLong = errors.New("comment is too long")
var ErrCommentNotFound = errors.New("comment not found")

//TODO: generateID
