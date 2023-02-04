package compile

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/model"
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

	policy := &model.Policy{ID: "identifier-123"}

	compiler := NewMockCompiler(ctrl)
	compiler.EXPECT().CompilePolicy(policy).Return(nil)

	dispatcher := event.NewMockDispatcher(ctrl)

	subscriber := NewSubscriber(logger, compiler, dispatcher)

	eventChan := make(chan *event.Event)

	// When - Then
	go subscriber.handlePolicyEvents(eventChan)

	eventChan <- &event.Event{
		Data: &event.ItemEvent{Data: policy},
	}

	close(eventChan)
}

func TestHandleResourceEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	resource := &model.Resource{ID: "resource-123"}

	compiler := NewMockCompiler(ctrl)
	compiler.EXPECT().CompileResource(resource).Return(nil)

	dispatcher := event.NewMockDispatcher(ctrl)

	subscriber := NewSubscriber(logger, compiler, dispatcher)

	eventChan := make(chan *event.Event)

	// When - Then
	go subscriber.handleResourceEvents(eventChan)

	eventChan <- &event.Event{
		Data: &event.ItemEvent{Data: resource},
	}

	close(eventChan)
}

func TestHandlePrincipalEvents(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	logger := slog.New(log.NewNopHandler())

	principal := &model.Principal{ID: "principal-123"}

	compiler := NewMockCompiler(ctrl)
	compiler.EXPECT().CompilePrincipal(principal).Return(nil)

	dispatcher := event.NewMockDispatcher(ctrl)

	subscriber := NewSubscriber(logger, compiler, dispatcher)

	eventChan := make(chan *event.Event)

	// When - Then
	go subscriber.handlePrincipalEvents(eventChan)

	eventChan <- &event.Event{
		Data: &event.ItemEvent{Data: principal},
	}

	close(eventChan)
}
