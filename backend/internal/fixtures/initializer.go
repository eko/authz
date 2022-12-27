package fixtures

import (
	"errors"
	"fmt"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/manager"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
)

var (
	resources = map[string][]string{
		"actions":    {"list", "get"},
		"clients":    {"list", "create"},
		"policies":   {"list", "get", "create", "update", "delete"},
		"principals": {"list", "get", "create", "update", "delete"},
		"resources":  {"list", "get", "create", "delete"},
		"roles":      {"list", "get", "create", "update", "delete"},
	}
)

type initializer struct {
	cfg      *configs.User
	logger   *slog.Logger
	manager  manager.Manager
	policies []string
}

func NewInitializer(
	cfg *configs.User,
	logger *slog.Logger,
	manager manager.Manager,
) *initializer {
	return &initializer{
		cfg:      cfg,
		logger:   logger,
		manager:  manager,
		policies: []string{},
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
	for resourceType, _ := range resources {
		resource, err := i.manager.CreateResource(
			fmt.Sprintf("%s.%s.%s", configs.ApplicationName, resourceType, "*"),
			fmt.Sprintf("%s.%s", configs.ApplicationName, resourceType),
			"*",
			nil,
		)
		if err != nil {
			return err
		}

		resource.IsLocked = true

		if err = i.manager.GetResourceRepository().Update(resource); err != nil {
			return err
		}
	}

	return nil
}

func (i *initializer) checkAlreadyInitialized() bool {
	_, err := i.manager.GetUserRepository().GetByFields(map[string]database.FieldValue{
		"username": {Operator: "=", Value: defaultAdminUsername},
	})

	return err == nil
}

func (i *initializer) initializePolicies() error {
	for resourceType, actions := range resources {
		policy, err := i.manager.CreatePolicy(
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
	_, err := i.manager.CreateRole(
		fmt.Sprintf("%s-admin", configs.ApplicationName),
		i.policies,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i *initializer) initializeUser() error {
	adminUser, err := i.manager.GetUserRepository().GetByFields(map[string]database.FieldValue{
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

	// Create user "admin" and principal named "authz-admin"
	if _, err := i.manager.CreateUser(defaultAdminUsername, i.cfg.AdminDefaultPassword); err != nil {
		return err
	}

	// Attach role "authz-admin" to user "authz-admin"
	_, err = i.manager.UpdatePrincipal("authz-admin", []string{"authz-admin"}, nil)

	return err
}
