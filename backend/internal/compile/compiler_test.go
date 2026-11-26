package compile

import (
	"testing"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/helper/time"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewCompiler(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)

	clock := time.NewMockClock(ctrl)
	compiledManager := manager.NewMockCompiledPolicy(ctrl)
	policyManager := manager.NewMockPolicy(ctrl)
	principalManager := manager.NewMockPrincipal(ctrl)
	resourceManager := manager.NewMockResource(ctrl)

	// When
	compilerInstance := NewCompiler(
		clock,
		compiledManager,
		policyManager,
		principalManager,
		resourceManager,
	)

	// Then
	assert := assert.New(t)

	assert.IsType(new(compiler), compilerInstance)

	assert.Equal(clock, compilerInstance.clock)
	assert.Equal(compiledManager, compilerInstance.compiledManager)
	assert.Equal(policyManager, compilerInstance.policyManager)
	assert.Equal(principalManager, compilerInstance.principalManager)
	assert.Equal(resourceManager, compilerInstance.resourceManager)
}
