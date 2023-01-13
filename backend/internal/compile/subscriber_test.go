package compile

import (
	"testing"

	"github.com/eko/authz/backend/internal/event"
	"github.com/eko/authz/backend/internal/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestNewSubscriber(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	compiler := NewMockCompiler(ctrl)

	dispatcher := event.NewMockDispatcher(ctrl)

	// When
	subscriberInstance := NewSubscriber(logger, compiler, dispatcher)

	// Then
	assert := assert.New(t)

	assert.IsType(new(subscriber), subscriberInstance)

	assert.Equal(logger, subscriberInstance.logger)
	assert.Equal(compiler, subscriberInstance.compiler)
	assert.Equal(dispatcher, subscriberInstance.dispatcher)
}

func TestHandlePolicyEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	compiler := NewMockCompiler(ctrl)
	compiler.EXPECT().CompilePolicy("identifier-123").Return(nil)

	dispatcher := event.NewMockDispatcher(ctrl)

	subscriber := NewSubscriber(logger, compiler, dispatcher)

	eventChan := make(chan *event.Event, 1)

	// When - Then
	go subscriber.handlePolicyEvents(make(chan *event.Event))

	eventChan <- &event.Event{
		Data: "identifier-123",
	}

	close(eventChan)
}

func TestHandleResourceEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	compiler := NewMockCompiler(ctrl)
	compiler.EXPECT().CompileResource("identifier-123").Return(nil)

	dispatcher := event.NewMockDispatcher(ctrl)

	subscriber := NewSubscriber(logger, compiler, dispatcher)

	eventChan := make(chan *event.Event, 1)

	// When - Then
	go subscriber.handleResourceEvents(make(chan *event.Event))

	eventChan <- &event.Event{
		Data: "identifier-123",
	}

	close(eventChan)
}

func TestHandlePrincipalEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	compiler := NewMockCompiler(ctrl)
	compiler.EXPECT().CompilePrincipal("identifier-123").Return(nil)

	dispatcher := event.NewMockDispatcher(ctrl)

	subscriber := NewSubscriber(logger, compiler, dispatcher)

	eventChan := make(chan *event.Event, 1)

	// When - Then
	go subscriber.handlePrincipalEvents(make(chan *event.Event))

	eventChan <- &event.Event{
		Data: "identifier-123",
	}

	close(eventChan)
}
