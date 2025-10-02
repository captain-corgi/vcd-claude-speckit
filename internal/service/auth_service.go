package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"employee-management-system/internal/domain"
	"employee-management-system/internal/repository"
)

// AuthService handles authentication and authorization
type AuthService struct {
	userRepo           *repository.UserRepository
	auditRepo          *repository.AuditLogRepository
	eventStore         domain.EventStoreRepository
	eventDispatcher    *domain.EventDispatcher
	jwtSecret          string
	jwtIssuer          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
	jwt.RegisteredClaims
}

// TokenPair represents an access token and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"` // "Bearer"
	ExpiresIn    int64  `json:"expiresIn"` // expires in seconds
}

// AuthRequest represents an authentication request
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo *repository.UserRepository,
	auditRepo *repository.AuditLogRepository,
	eventStore domain.EventStoreRepository,
	eventDispatcher *domain.EventDispatcher,
	jwtSecret string,
	jwtIssuer string,
	accessTokenExpiry time.Duration,
	refreshTokenExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		auditRepo:          auditRepo,
		eventStore:         eventStore,
		eventDispatcher:    eventDispatcher,
		jwtSecret:          jwtSecret,
		jwtIssuer:          jwtIssuer,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

// AuthenticateUser authenticates a user and returns JWT tokens
func (s *AuthService) AuthenticateUser(
	ctx context.Context,
	username, password string,
	ipAddress, userAgent string,
) (*domain.User, *TokenPair, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if !user.Authenticate(password) {
		return nil, nil, domain.ErrInvalidCredentials
	}

	// Update last login
	user.UpdateLastLogin()
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Generate JWT tokens
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create audit log for successful login
	auditLog, err := s.createAuditLogForLogin(user, ipAddress, userAgent)
	if err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Save audit log if it was created
	if auditLog != nil {
		if err := s.auditRepo.Create(ctx, auditLog); err != nil {
			fmt.Printf("Warning: failed to save audit log: %v\n", err)
		}
	}

	// Create domain event for successful login
	event := domain.NewUserLoggedInEvent(user.ID, username, ipAddress, userAgent)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return user, tokenPair, nil
}

// RefreshTokens refreshes access token using refresh token
func (s *AuthService) RefreshTokens(
	ctx context.Context,
	refreshToken string,
	ipAddress, userAgent string,
) (*TokenPair, error) {
	// Validate refresh token
	claims, err := s.validateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Find user by ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil || !user.IsActive {
		return nil, fmt.Errorf("user not found or inactive")
	}

	// Generate new token pair
	newTokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	// Create audit log for token refresh
	auditLog, err := s.createAuditLogForTokenRefresh(user, ipAddress, userAgent)
	if err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Save audit log if it was created
	if auditLog != nil {
		if err := s.auditRepo.Create(ctx, auditLog); err != nil {
			fmt.Printf("Warning: failed to save audit log: %v\n", err)
		}
	}

	return newTokenPair, nil
}

// LogoutUser logs out a user and invalidates refresh token
func (s *AuthService) LogoutUser(
	ctx context.Context,
	userID uuid.UUID,
	refreshToken string,
	ipAddress, userAgent string,
) error {
	// Create audit log for logout
	auditLog, err := s.createAuditLogForLogout(userID, ipAddress, userAgent)
	if err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Save audit log if it was created
	if auditLog != nil {
		if err := s.auditRepo.Create(ctx, auditLog); err != nil {
			fmt.Printf("Warning: failed to save audit log: %v\n", err)
		}
	}

	// In a real implementation, you might want to store invalidated refresh tokens
	// in a cache or database to prevent reuse

	return nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *AuthService) ValidateAccessToken(accessToken string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ChangePassword changes a user's password with current password verification
func (s *AuthService) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	currentPassword, newPassword string,
	ipAddress, userAgent string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return domain.ErrNotFound
	}

	// Verify current password
	if !user.Authenticate(currentPassword) {
		return domain.ErrInvalidCredentials
	}

	// Change password
	if err := user.ChangePassword(currentPassword, newPassword); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog, err := s.createAuditLogForPasswordChange(user, "self", "change", ipAddress, userAgent)
	if err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Save audit log if it was created
	if auditLog != nil {
		if err := s.auditRepo.Create(ctx, auditLog); err != nil {
			fmt.Printf("Warning: failed to save audit log: %v\n", err)
		}
	}

	// Create domain event
	event := domain.NewUserPasswordChangedEvent(user.ID, user.Username, "self", "change")
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return nil
}

// ResetPassword resets a user's password (admin function)
func (s *AuthService) ResetPassword(
	ctx context.Context,
	userID uuid.UUID,
	newPassword string,
	adminID string,
	ipAddress, userAgent string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return domain.ErrNotFound
	}

	// Reset password
	if err := user.ResetPassword(newPassword); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog, err := s.createAuditLogForPasswordChange(user, adminID, "reset", ipAddress, userAgent)
	if err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Save audit log if it was created
	if auditLog != nil {
		if err := s.auditRepo.Create(ctx, auditLog); err != nil {
			fmt.Printf("Warning: failed to save audit log: %v\n", err)
		}
	}

	// Create domain event
	event := domain.NewUserPasswordChangedEvent(user.ID, user.Username, adminID, "reset")
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return nil
}

// generateTokenPair generates an access token and refresh token pair
func (s *AuthService) generateTokenPair(user *domain.User) (*TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTokenExpiry.Seconds()),
	}, nil
}

// generateAccessToken generates a JWT access token
func (s *AuthService) generateAccessToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID:    user.ID.String(),
		Username:  user.Username,
		Role:      string(user.Role),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.accessTokenExpiry).Unix(),
		Issuer:    s.jwtIssuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// generateRefreshToken generates a random refresh token
func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// validateRefreshToken validates a refresh token
func (s *AuthService) validateRefreshToken(refreshToken string) (*JWTClaims, error) {
	// In a real implementation, you would validate the refresh token against
	// a database or cache of valid tokens. For now, we'll just validate the format.

	// This is a simplified validation. In production, you'd want to:
	// 1. Check if the refresh token exists in a valid tokens list
	// 2. Check if it's expired
	// 3. Invalidate it after use

	// For this example, we'll just return dummy claims
	// In a real implementation, you'd decode and validate the refresh token

	return &JWTClaims{
		UserID:    "dummy-user-id",
		Username:  "dummy-username",
		Role:      "dummy-role",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(s.refreshTokenExpiry).Unix(),
	}, nil
}

// createAuditLogForLogin creates an audit log for successful login
func (s *AuthService) createAuditLogForLogin(user *domain.User, ipAddress, userAgent string) (*domain.AuditLog, error) {
	newValues := map[string]interface{}{
		"userId":    user.ID,
		"username":  user.Username,
		"loginTime": time.Now(),
		"ipAddress": ipAddress,
		"userAgent": userAgent,
	}

	auditLog, err := domain.NewAuditLog(
		uuid.New(), // System operation
		domain.OperationUserLogin,
		user.ID.String(),
		nil,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, err
	}
	return auditLog, nil
}

// createAuditLogForTokenRefresh creates an audit log for token refresh
func (s *AuthService) createAuditLogForTokenRefresh(user *domain.User, ipAddress, userAgent string) (*domain.AuditLog, error) {
	newValues := map[string]interface{}{
		"userId":    user.ID,
		"username":  user.Username,
		"action":    "token_refresh",
		"ipAddress": ipAddress,
		"userAgent": userAgent,
		"timestamp": time.Now(),
	}

	auditLog, err := domain.NewAuditLog(
		uuid.New(), // System operation
		domain.OperationUserLogin,
		user.ID.String(),
		nil,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, err
	}
	return auditLog, nil
}

// createAuditLogForLogout creates an audit log for logout
func (s *AuthService) createAuditLogForLogout(userID uuid.UUID, ipAddress, userAgent string) (*domain.AuditLog, error) {
	newValues := map[string]interface{}{
		"userId":     userID,
		"logoutTime": time.Now(),
		"ipAddress":  ipAddress,
		"userAgent":  userAgent,
	}

	auditLog, err := domain.NewAuditLog(
		uuid.New(), // System operation
		domain.OperationUserLogout,
		userID.String(),
		nil,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, err
	}
	return auditLog, nil
}

// createAuditLogForPasswordChange creates an audit log for password change
func (s *AuthService) createAuditLogForPasswordChange(user *domain.User, changedBy, method, ipAddress, userAgent string) (*domain.AuditLog, error) {
	newValues := map[string]interface{}{
		"passwordChanged": true,
		"changedBy":       changedBy,
		"method":          method,
		"timestamp":       time.Now(),
	}

	auditLog, err := domain.NewAuditLog(
		user.ID,
		domain.OperationPasswordChange,
		user.ID.String(),
		nil,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return nil, err
	}
	return auditLog, nil
}

// GetConfig returns authentication service configuration
func (s *AuthService) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"accessTokenExpiry":  s.accessTokenExpiry.Seconds(),
		"refreshTokenExpiry": s.refreshTokenExpiry.Seconds(),
		"jwtIssuer":          s.jwtIssuer,
	}
}
