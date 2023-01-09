package fixtures

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
)

var (
	resources = map[string][]string{
		"actions":    {"list", "get"},
		"clients":    {"list", "get", "create", "delete"},
		"policies":   {"list", "get", "create", "update", "delete"},
		"principals": {"list", "get", "create", "update", "delete"},
		"resources":  {"list", "get", "create", "update", "delete"},
		"roles":      {"list", "get", "create", "update", "delete"},
		"users":      {"list", "get", "create", "delete"},
	}
)

type Initializer interface {
	Initialize() error
}

type initializer struct {
	cfg              *configs.User
	logger           *slog.Logger
	policyManager    manager.Policy
	principalManager manager.Principal
	resourceManager  manager.Resource
	roleManager      manager.Role
	userManager      manager.User
	policies         []string
}

func NewInitializer(
	cfg *configs.User,
	logger *slog.Logger,
	policyManager manager.Policy,
	principalManager manager.Principal,
	resourceManager manager.Resource,
	roleManager manager.Role,
	userManager manager.User,
) Initializer {
	return &initializer{
		cfg:              cfg,
		logger:           logger,
		policyManager:    policyManager,
		principalManager: principalManager,
		resourceManager:  resourceManager,
		roleManager:      roleManager,
		userManager:      userManager,
		policies:         []string{},
	}
}

// Initialize initializes default application resources.
func (i *initializer) Initialize() error {
	if i.checkAlreadyInitialized() {
		return nil
	}

	if err := i.initializeResources(); err != nil {
		return err
	}

	if err := i.initializePolicies(); err != nil {
		return err
	}

	if err := i.initializeRoles(); err != nil {
		return err
	}

	if err := i.initializeUser(); err != nil {
		return err
	}

	return nil
}

func (i *initializer) initializeResources() error {
	for resourceType := range resources {
		resource, err := i.resourceManager.Create(
			fmt.Sprintf("%s.%s.%s", configs.ApplicationName, resourceType, "*"),
			fmt.Sprintf("%s.%s", configs.ApplicationName, resourceType),
			"*",
			nil,
		)
		if err != nil {
			return err
		}

		resource.IsLocked = true

		if err = i.resourceManager.GetRepository().Update(resource); err != nil {
			return err
		}
	}

	return nil
}

func (i *initializer) checkAlreadyInitialized() bool {
	_, err := i.userManager.GetRepository().GetByFields(map[string]repository.FieldValue{
		"username": {Operator: "=", Value: defaultAdminUsername},
	})

	return err == nil
}

func (i *initializer) initializePolicies() error {
	for resourceType, actions := range resources {
		policy, err := i.policyManager.Create(
			fmt.Sprintf("%s-%s-admin", configs.ApplicationName, resourceType),
			[]string{
				fmt.Sprintf("%s.%s.%s", configs.ApplicationName, resourceType, "*"),
			},
			actions,
			nil,
		)
		if err != nil {
			return err
		}

		i.policies = append(i.policies, policy.ID)
	}

	return nil
}

func (i *initializer) initializeRoles() error {
	_, err := i.roleManager.Create(
		fmt.Sprintf("%s-admin", configs.ApplicationName),
		i.policies,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i *initializer) initializeUser() error {
	adminUser, err := i.userManager.GetRepository().GetByFields(map[string]repository.FieldValue{
		"username": {Operator: "=", Value: defaultAdminUsername},
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		i.logger.Error("Unable to get default admin user", err)
		return err
	}

	if adminUser != nil {
		// User already exists, nothing to do.
		return nil
	}

	// Create user "admin" and principal named "authz-user-admin"
	user, err := i.userManager.Create(defaultAdminUsername, i.cfg.AdminDefaultPassword)
	if err != nil {
		return fmt.Errorf("unable to create default admin user: %v", err)
	}

	// Retrieve principal created following the user creation
	principal, err := i.principalManager.GetRepository().Get(
		model.UserPrincipal(user.Username),
	)
	if err != nil {
		return fmt.Errorf("unable to retrieve admin principal: %v", err)
	}

	principal.IsLocked = true

	// Attach role "authz-admin" to user principal "authz-admin"
	role, err := i.roleManager.GetRepository().Get(fmt.Sprintf("%s-admin", configs.ApplicationName))
	if err != nil {
		return fmt.Errorf("unable to retrieve admin role: %v", err)
	}

	principal.Roles = []*model.Role{role}

	return i.principalManager.GetRepository().Update(principal)
}
