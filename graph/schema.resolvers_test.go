package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmozzze/SasPosts/graph/model"
	"github.com/tmozzze/SasPosts/internal/domain"
	redisMocks "github.com/tmozzze/SasPosts/internal/redis/mocks"
	"github.com/tmozzze/SasPosts/internal/repository/mocks"
)

func TestQuery_Post(t *testing.T) {
	mockPostRepo := mocks.NewPostRepository(t)
	expectedPost := &domain.Post{
		ID: "post-123", Title: "Test post", Content: "Content", Author: "Tester123",
	}

	mockPostRepo.On("GetByID", mock.Anything, "post-123").Return(expectedPost, nil)
	resolver := &Resolver{PostRepo: mockPostRepo}
	result, err := resolver.Query().Post(context.Background(), "post-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, result)
	mockPostRepo.AssertExpectations(t)
}

func TestQuery_Posts(t *testing.T) {
	mockPostRepo := mocks.NewPostRepository(t)
	expectedPosts := []*domain.Post{
		{ID: "post-1"},
		{ID: "post-2"},
	}
	mockPostRepo.On("GetAll", mock.Anything).Return(expectedPosts, nil)

	resolver := &Resolver{PostRepo: mockPostRepo}
	result, err := resolver.Query().Posts(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, result)
	mockPostRepo.AssertExpectations(t)
}

func TestMutation_CreatePost(t *testing.T) {
	mockPostRepo := mocks.NewPostRepository(t)
	input := model.NewPostInput{
		Title: "Post", Content: "Content", Author: "Author", AllowComments: true,
	}
	mockPostRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Post")).Return(nil)
	resolver := &Resolver{PostRepo: mockPostRepo}
	result, err := resolver.Mutation().CreatePost(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, input.Title, result.Title)
	assert.NotEmpty(t, result.ID)
	mockPostRepo.AssertExpectations(t)
}

func TestMutation_CreateComment(t *testing.T) {
	t.Run("valid comments create", func(t *testing.T) {
		mockPostRepo := mocks.NewPostRepository(t)
		mockCommentRepo := mocks.NewCommentRepository(t)
		mockPublisher := redisMocks.NewPubSub(t)

		input := model.NewCommentInput{
			PostID: "post-123", Author: "commenter", Content: "nice post",
		}

		mockPostRepo.On("CheckAllowedComments", mock.Anything, "post-123").Return(true, nil)
		mockCommentRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Comment")).Return(nil)

		channelName := fmt.Sprintf("comments:%s", input.PostID)

		mockPublisher.On("Publish", mock.Anything, channelName, mock.AnythingOfType("*domain.Comment")).Return(nil)

		resolver := &Resolver{
			PostRepo:    mockPostRepo,
			CommentRepo: mockCommentRepo,
			PubSub:      mockPublisher,
		}
		result, err := resolver.Mutation().CreateComment(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, input.Content, result.Content)
		mockPostRepo.AssertExpectations(t)
		mockCommentRepo.AssertExpectations(t)
		mockPublisher.AssertExpectations(t)
	})

	t.Run("error, if comments are off", func(t *testing.T) {
		mockPostRepo := mocks.NewPostRepository(t)
		input := model.NewCommentInput{PostID: "post-456"}

		mockPostRepo.On("CheckAllowedComments", mock.Anything, "post-456").Return(false, nil)
		resolver := &Resolver{PostRepo: mockPostRepo}
		_, err := resolver.Mutation().CreateComment(context.Background(), input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrCommentsOff)
		mockPostRepo.AssertExpectations(t)
	})
}

func TestMutation_ToggleComments(t *testing.T) {
	mockPostRepo := mocks.NewPostRepository(t)
	postID := "post-123"
	expectedPost := &domain.Post{ID: postID, AllowComments: false}

	mockPostRepo.On("ToggleComments", mock.Anything, postID, false).Return(nil)
	mockPostRepo.On("GetByID", mock.Anything, postID).Return(expectedPost, nil)

	resolver := &Resolver{PostRepo: mockPostRepo}
	result, err := resolver.Mutation().ToggleComments(context.Background(), postID, false)

	assert.NoError(t, err)
	assert.Equal(t, expectedPost, result)
	mockPostRepo.AssertExpectations(t)
}

func TestPostResolver_Comments(t *testing.T) {
	mockCommentRepo := mocks.NewCommentRepository(t)
	const postID = "post-with-comments"
	parentPost := &domain.Post{ID: postID}
	expectedComments := []*domain.Comment{
		{ID: "comment-1", PostID: postID, Content: "first comment"},
	}
	mockCommentRepo.On("GetByPost", mock.Anything, postID, 10, 0).Return(expectedComments, nil)
	resolver := &Resolver{CommentRepo: mockCommentRepo}
	limit, offset := 10, 0
	result, err := resolver.Post().Comments(context.Background(), parentPost, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedComments, result)
	mockCommentRepo.AssertExpectations(t)
}

func TestCommentResolver_Children(t *testing.T) {
	mockCommentRepo := mocks.NewCommentRepository(t)
	parentID := "parent-comment-1"
	parentComment := &domain.Comment{ID: parentID}
	expectedChildren := []*domain.Comment{
		{ID: "child-1", ParentID: &parentID, Content: "child comment"},
	}
	mockCommentRepo.On("GetChildren", mock.Anything, parentID, 5, 0).Return(expectedChildren, nil)
	resolver := &Resolver{CommentRepo: mockCommentRepo}
	limit, offset := 5, 0
	result, err := resolver.Comment().Children(context.Background(), parentComment, &limit, &offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedChildren, result)
	mockCommentRepo.AssertExpectations(t)
}

func TestSubscription_CommentAdded(t *testing.T) {
	mockPostRepo := mocks.NewPostRepository(t)
	mockPubSub := redisMocks.NewPubSub(t)

	postID := "post-123"
	channelName := fmt.Sprintf("comments:%s", postID)

	mockPostRepo.On("GetByID", mock.Anything, postID).Return(&domain.Post{ID: postID}, nil)

	writableChan := make(chan []byte)

	var readOnlyChan <-chan []byte = writableChan

	closeFunc := func() { close(writableChan) }

	mockPubSub.On("Subscribe", mock.Anything, channelName).Return(readOnlyChan, closeFunc)

	resolver := &Resolver{PostRepo: mockPostRepo, PubSub: mockPubSub}

	gqlChan, err := resolver.Subscription().CommentAdded(context.Background(), postID)
	assert.NoError(t, err)

	expectedComment := &domain.Comment{ID: "comment-sub-test", Content: "sas from subscription"}
	payload, _ := json.Marshal(expectedComment)

	go func() {
		writableChan <- payload
	}()

	receivedComment := <-gqlChan
	assert.Equal(t, expectedComment, receivedComment)

	mockPostRepo.AssertExpectations(t)
	mockPubSub.AssertExpectations(t)
}
