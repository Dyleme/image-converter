package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dyleme/image-coverter/internal/model"
	"golang.org/x/crypto/bcrypt"
)

// HashGenerator interface providing you the ability to generate password hash
// and compare it with pure text passoword.
type HashGenerator interface {
	GeneratePasswordHash(password string) string
	IsValidPassword(password string, hash []byte) bool
}

// HashGen struct is realiszation of the HashGenerator interface with the bcrypt package.
type HashGen struct{}

// GeneratePasswordHash generates hash from the password.
func (h *HashGen) GeneratePasswordHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hash)
}

// IsValidPasword compare the password and password hash,
// returns true if they corresponds, flase in the other situations.
func (h *HashGen) IsValidPassword(password string, hash []byte) bool {
	errNotEqual := bcrypt.CompareHashAndPassword(hash, []byte(password))

	return errNotEqual == nil
}

// AuthRepo is an interface which provides methods to implement with repository.
type AuthRepo interface {
	// CreateUserMethod creates user in the repository.
	CreateUser(ctx context.Context, user model.User) (int, error)

	// GetPasswordHashAndID returns user password hash and id.
	GetPasswordHashAndID(ctx context.Context, nickname string) (hash []byte, userID int, err error)
}

// JwtGenerator is an interface that provides method to create jwt tokens.
type JwtGenerator interface {
	// Create token is method that can create jwt tokens.
	CreateToken(ctx context.Context, id int) (string, error)
}

// Auth struct provides the ability to create user and validate it.
type Auth struct {
	repo    AuthRepo
	hashGen HashGenerator
	jwtGen  JwtGenerator
}

// NewAuthService is the constructor to the AuthService.
func NewAuth(repo AuthRepo, hashGen HashGenerator, jwtGen JwtGenerator) *Auth {
	return &Auth{repo: repo, hashGen: hashGen, jwtGen: jwtGen}
}

// CreateUser function returns the id of the created user or error if any occures.
// Function get password hash of the user and creates user and calls CreateUser method of repository.
func (s *Auth) CreateUser(ctx context.Context, user model.User) (int, error) {
	user.Password = s.hashGen.GeneratePasswordHash(user.Password)
	return s.repo.CreateUser(ctx, user)
}

var ErrWrongPassword = errors.New("wrong password")

// ValidateUser returns the jwt token of the user, if the provided user exists  in repo and password is correct.
// In any other sitationds function returns ("", err).
// Method get password and if calling repo.GetPasswordHashAndID then validates it with the hashGen.IsValidPassword,
// and create token with the help jwtGen.CreateToken.
func (s *Auth) ValidateUser(ctx context.Context, user model.User) (string, error) {
	hash, id, err := s.repo.GetPasswordHashAndID(ctx, user.Nickname)
	if err != nil {
		return "", fmt.Errorf("validate user %w", err)
	}

	if !s.hashGen.IsValidPassword(user.Password, hash) {
		return "", ErrWrongPassword
	}

	return s.jwtGen.CreateToken(ctx, id)
}
