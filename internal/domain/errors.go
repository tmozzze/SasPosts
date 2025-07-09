package domain

import "errors"

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotFound    = errors.New("post not found")
)
