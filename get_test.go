package gomessagestore_test

import (
	"context"
	"testing"

	. "github.com/blackhatbrigade/gomessagestore"
	"github.com/blackhatbrigade/gomessagestore/repository"
	"github.com/blackhatbrigade/gomessagestore/repository/mocks"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
)

func TestGetWithCommandStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	msg := getSampleCommand()
	ctx := context.Background()

	msgEnv := getSampleCommandAsEnvelope()

	mockRepo.
		EXPECT().
		GetAllMessagesInStream(ctx, msgEnv.Stream).
		Return([]*repository.MessageEnvelope{msgEnv}, nil)

	msgStore := NewMessageStoreFromRepository(mockRepo)
	msgs, err := msgStore.Get(ctx, CommandStream(msgEnv.StreamType))

	if err != nil {
		t.Error("An error has ocurred while getting messages from message store")
	}
	if len(msgs) != 1 {
		t.Error("Incorrect number of messages returned")
	} else {
		assertMessageMatchesCommand(t, msgs[0], msg)
	}
}

func TestGetWithoutOptionsReturnsError(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	ctx := context.Background()

	msgStore := NewMessageStoreFromRepository(mockRepo)
	_, err := msgStore.Get(ctx)

	if err != ErrMissingGetOptions {
		t.Errorf("Expected ErrMissingGetOptions and got %v", err)
	}
}

func TestGetWithEventStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	msg := getSampleEvent()
	ctx := context.Background()

	msgEnv := getSampleEventAsEnvelope()

	mockRepo.
		EXPECT().
		GetAllMessagesInStream(ctx, msgEnv.Stream).
		Return([]*repository.MessageEnvelope{msgEnv}, nil)

	msgStore := NewMessageStoreFromRepository(mockRepo)
	msgs, err := msgStore.Get(ctx, EventStream(msg.Category, msg.CategoryID))

	if err != nil {
		t.Error("An error has ocurred while getting messages from message store")
	}
	if len(msgs) != 1 {
		t.Error("Incorrect number of messages returned")
	} else {
		assertMessageMatchesEvent(t, msgs[0], msg)
	}
}

func TestGetWithCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	msg := getSampleEvent()
	ctx := context.Background()

	msgEnv := getSampleEventAsEnvelope()

	mockRepo.
		EXPECT().
		GetAllMessagesInCategory(ctx, msgEnv.StreamType).
		Return([]*repository.MessageEnvelope{msgEnv}, nil)

	msgStore := NewMessageStoreFromRepository(mockRepo)
	msgs, err := msgStore.Get(ctx, Category(msg.Category))

	if err != nil {
		t.Error("An error has ocurred while getting messages from message store")
	}
	if len(msgs) != 1 {
		t.Error("Incorrect number of messages returned")
	} else {
		assertMessageMatchesEvent(t, msgs[0], msg)
	}
}

func TestGetMessagesCannotUseBothStreamAndCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	msg := getSampleCommand()
	ctx := context.Background()

	msgStore := NewMessageStoreFromRepository(mockRepo)
	_, err := msgStore.Get(ctx, Category(msg.Category), CommandStream(msg.Category))

	if err != ErrGetMessagesCannotUseBothStreamAndCategory {
		t.Error("Expected ErrGetMessagesCannotUseBothStreamAndCategory")
	}
}

func TestGetWithEventStreamAndSince(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	msg := getSampleEvent()
	ctx := context.Background()
	var globalPosition int64

	msgStore := NewMessageStoreFromRepository(mockRepo)

	msgEnv := getSampleEventAsEnvelope()

	mockRepo.
		EXPECT().
		GetAllMessagesInStreamSince(ctx, msgEnv.Stream, globalPosition).
		Return([]*repository.MessageEnvelope{msgEnv}, nil)

	msgs, err := msgStore.Get(
		ctx,
		Since(globalPosition),
		EventStream(msg.Category, msg.CategoryID),
	)

	if err != nil {
		t.Error("An error has ocurred while getting messages from message store")
	}

	if len(msgs) != 1 {
		t.Error("Incorrect number of messages returned")
	} else {
		assertMessageMatchesEvent(t, msgs[0], msg)
	}
}

func TestGetWithCommandStreamAndSince(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)

	msg := getSampleCommand()
	ctx := context.Background()
	var globalPosition int64

	msgStore := NewMessageStoreFromRepository(mockRepo)

	msgEnv := getSampleCommandAsEnvelope()

	mockRepo.
		EXPECT().
		GetAllMessagesInStreamSince(ctx, msgEnv.Stream, globalPosition).
		Return([]*repository.MessageEnvelope{msgEnv}, nil)

	msgs, err := msgStore.Get(
		ctx,
		Since(globalPosition),
		CommandStream(msg.Category),
	)

	if err != nil {
		t.Error("An error has ocurred while getting messages from message store")
	}

	if len(msgs) != 1 {
		t.Error("Incorrect number of messages returned")
	} else {
		assertMessageMatchesCommand(t, msgs[0], msg)
	}
}

func TestGetMessagesRequiresEitherStreamOrCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var globalPosition int64

	mockRepo := mock_repository.NewMockRepository(ctrl)

	ctx := context.Background()

	msgStore := NewMessageStoreFromRepository(mockRepo)
	_, err := msgStore.Get(ctx, Since(globalPosition))

	if err != ErrGetMessagesRequiresEitherStreamOrCategory {
		t.Errorf("Expected ErrGetMessagesRequiresEitherStreamOrCategory, but got %s", err)
	}
}