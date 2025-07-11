package domain

import "errors"

var (
	ErrCommentNotFound       = errors.New("comment not found")
	ErrParentCommentNotFound = errors.New("parent comment not found")
	ErrPostNotFound          = errors.New("post not found")
	ErrCommentsOff           = errors.New("comments off for this post")
)
