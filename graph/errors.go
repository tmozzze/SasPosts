package graph

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/tmozzze/SasPosts/internal/domain"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {

	if errors.Is(err, domain.ErrCommentTooLong) {
		return &gqlerror.Error{
			Message: err.Error(),
			Extensions: map[string]interface{}{
				"code":      "COMMENT_TOO_LONG",
				"maxLength": domain.MaxCommentLength,
			},
		}
	}
	if errors.Is(err, domain.ErrCommentsOff) {
		return &gqlerror.Error{
			Message: err.Error(),
			Extensions: map[string]interface{}{
				"code": "COMMENT_OFF",
			},
		}
	}

	if errors.Is(err, domain.ErrPostNotFound) {
		return &gqlerror.Error{
			Message: err.Error(),
			Extensions: map[string]interface{}{
				"code": "POST_NOT_FOUND",
			},
		}
	}
	if errors.Is(err, domain.ErrParentCommentNotFound) {
		return &gqlerror.Error{
			Message: err.Error(),
			Extensions: map[string]interface{}{
				"code": "COMMENT_PARENT_NOT_FOUND",
			},
		}
	}

	return graphql.DefaultErrorPresenter(ctx, err)
}
