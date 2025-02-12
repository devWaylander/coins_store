package auth

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"time"

	internalErrors "github.com/devWaylander/coins_store/pkg/errors"
	"github.com/devWaylander/coins_store/pkg/models"
	"github.com/devWaylander/coins_store/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
)

var unsecuredHandles = map[string]*struct{}{
	"/api/auth": {},
}

type Repository interface {
	CreateUser(ctx context.Context, username, passwordHash string) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserPassHashByUsername(ctx context.Context, username string) (string, error)
}

type middleware struct {
	repo   Repository
	jwtKey string
}

func NewMiddleware(repo Repository, jwtKey string) *middleware {
	return &middleware{
		repo:   repo,
		jwtKey: jwtKey,
	}
}

func (m *middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if unsecuredHandles[r.URL.Path] != nil {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, internalErrors.ErrAuthHeader, http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (any, error) {
			return m.jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, internalErrors.ErrInvalidToken, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*models.Claims)
		if !ok {
			http.Error(w, internalErrors.ErrInvalidClaims, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), models.UserIDkey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *middleware) LoginWithPass(ctx context.Context, username, password string) (models.AuthDTO, error) {
	user, err := m.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return models.AuthDTO{}, err
	}

	// Не зарегистрирован
	if user.ID == 0 {
		validPass := m.validatePassword(password)
		if !validPass {
			return models.AuthDTO{}, errors.New(internalErrors.ErrWrongPasswordFormat)
		}

		passHash, err := utils.PasswordHash(password)
		if err != nil {
			return models.AuthDTO{}, err
		}
		userID, err := m.repo.CreateUser(ctx, username, passHash)
		if err != nil {
			return models.AuthDTO{}, err
		}

		token, err := m.generateJWT(userID)
		if err != nil {
			return models.AuthDTO{}, err
		}

		return models.AuthDTO{Token: token}, err
	}

	passHash, err := m.repo.GetUserPassHashByUsername(ctx, username)
	if err != nil {
		return models.AuthDTO{}, err
	}
	err = utils.PasswordCompare(password, passHash)
	if err != nil {
		return models.AuthDTO{}, errors.New(internalErrors.ErrWrongPassword)
	}
	token, err := m.generateJWT(user.ID)
	if err != nil {
		return models.AuthDTO{}, err
	}

	return models.AuthDTO{Token: token}, err
}

func (m *middleware) generateJWT(userID int64) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *middleware) validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}
