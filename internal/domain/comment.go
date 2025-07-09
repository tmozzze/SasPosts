package domain

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/tmozzze/SasPosts/utils"
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

func NewComment(postID, authorID string, parentID *string, content string) (*Comment, error) {
	if utf8.RuneCountInString(content) > MaxCommentLength {
		return nil, ErrCommentTooLong
	}
	if postID == "" || authorID == "" {
		return nil, errors.New("postID and authorID cannot be empty")
	}

	comment := &Comment{
		ID:        utils.GenerateID(),
		PostID:    postID,
		AuthorID:  authorID,
		ParentID:  parentID, // nil если коммент - корень
		Content:   content,
		CreatedAt: time.Now(),
	}

	return comment, nil
}

var ErrCommentTooLong = errors.New("comment is too long")
