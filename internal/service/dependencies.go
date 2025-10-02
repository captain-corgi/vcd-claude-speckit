package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
	"employee-management-system/internal/repository"
)

// ServiceDependencies contains all dependencies needed for service layer
type ServiceDependencies struct {
	DB *gorm.DB
}

// NewServiceDependencies creates a new service dependencies container
func NewServiceDependencies(db *gorm.DB) *ServiceDependencies {
	return &ServiceDependencies{DB: db}
}

// EmployeeServiceFactory creates a new EmployeeService with all dependencies
func (sd *ServiceDependencies) EmployeeService() *EmployeeService {
	return &EmployeeService{
		employeeRepo:    repository.NewEmployeeRepository(sd.DB),
		auditRepo:       repository.NewAuditLogRepository(sd.DB),
		eventStore:      repository.NewEventStoreRepository(sd.DB),
		eventDispatcher: domain.NewEventDispatcher(),
	}
}

// UserServiceFactory creates a new UserService with all dependencies
func (sd *ServiceDependencies) UserService() *UserService {
	return &UserService{
		userRepo:        repository.NewUserRepository(sd.DB),
		auditRepo:       repository.NewAuditLogRepository(sd.DB),
		eventStore:      repository.NewEventStoreRepository(sd.DB),
		eventDispatcher: domain.NewEventDispatcher(),
	}
}

// AuthServiceFactory creates a new AuthService with all dependencies and configuration
func (sd *ServiceDependencies) AuthService(
	jwtSecret string,
	jwtIssuer string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:           repository.NewUserRepository(sd.DB),
		auditRepo:          repository.NewAuditLogRepository(sd.DB),
		eventStore:         repository.NewEventStoreRepository(sd.DB),
		eventDispatcher:    domain.NewEventDispatcher(),
		jwtSecret:          jwtSecret,
		jwtIssuer:          jwtIssuer,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

// ServiceContainer provides a centralized container for all services
type ServiceContainer struct {
	EmployeeService *EmployeeService
	UserService     *UserService
	AuthService     *AuthService
}

// NewServiceContainer creates a new service container with all services
func NewServiceContainer(
	db *gorm.DB,
	jwtSecret string,
	jwtIssuer string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
) *ServiceContainer {
	deps := NewServiceDependencies(db)

	return &ServiceContainer{
		EmployeeService: deps.EmployeeService(),
		UserService:     deps.UserService(),
		AuthService:     deps.AuthService(jwtSecret, jwtIssuer, accessTokenExpiry, refreshTokenExpiry),
	}
}

// SystemUser provides a system user ID for internal operations
var SystemUser = uuid.Nil.String()

// DefaultContext returns a context with system user information
func DefaultContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "system_user", SystemUser)
}
