package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewComment(t *testing.T) {
	t.Run("valid comment", func(t *testing.T) {
		comment, err := NewComment("post1", "author1", nil, "This comment is valid")

		require.NoError(t, err)
		require.NotNil(t, comment)

		assert.NotEmpty(t, comment.ID, "ID must not be empty")
	})

	t.Run("error, if comment is too long", func(t *testing.T) {
		longContent := strings.Repeat("a", MaxCommentLength+1)
		_, err := NewComment("post1", "author1", nil, longContent)

		require.Error(t, err)
		assert.ErrorIs(t, err, ErrCommentTooLong)
	})

	t.Run("error, if author or post id is empty", func(t *testing.T) {
		_, err := NewComment("", "author1", nil, "content")
		assert.Error(t, err)

		_, err = NewComment("post1", "", nil, "content")
		assert.Error(t, err)
	})
}
