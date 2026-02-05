package jwt

import (
	"errors"
	"minichat/internal/model"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // access | refresh
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidTokenType  = errors.New("invalid token type")
	ErrTokenExpired      = errors.New("token expired")
	ErrUnexpectedSigning = errors.New("unexpected signing method")
)

// Configuration via env:
// - MINICHAT_JWT_SECRET (recommended)
// - MINICHAT_JWT_ACCESS_TTL_MINUTES (default 120)
// - MINICHAT_JWT_REFRESH_TTL_HOURS (default 168)

func secretKey() []byte {
	sec := os.Getenv("MINICHAT_JWT_SECRET")
	if sec == "" {
		// dev fallback; set env in real deployments.
		sec = "minichat-dev-secret-change-me"
	}
	return []byte(sec)
}

func accessTTL() time.Duration {
	v := os.Getenv("MINICHAT_JWT_ACCESS_TTL_MINUTES")
	if v == "" {
		return 2 * time.Hour
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return 2 * time.Hour
	}
	return time.Duration(n) * time.Minute
}

func refreshTTL() time.Duration {
	v := os.Getenv("MINICHAT_JWT_REFRESH_TTL_HOURS")
	if v == "" {
		return 7 * 24 * time.Hour
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return 7 * 24 * time.Hour
	}
	return time.Duration(n) * time.Hour
}

func CreateClaims(user model.User, tokenType string, ttl time.Duration) MyClaims {
	now := time.Now()
	return MyClaims{
		Id:        user.ID,
		Username:  user.Username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "minichat",
		},
	}
}

// GenerateToken returns (accessToken, refreshToken, error)
func GenerateToken(user model.User) (string, string, error) {
	key := secretKey()

	accessClaims := CreateClaims(user, "access", accessTTL())
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(key)
	if err != nil {
		return "", "", err
	}

	refreshClaims := CreateClaims(user, "refresh", refreshTTL())
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(key)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

// ParseToken parses and validates signature/exp. It returns claims when ok.
func ParseToken(tokenString string) (*MyClaims, error) {
	key := secretKey()
	claims := &MyClaims{}

	tok, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigning
		}
		return key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}
	if !tok.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func ValidateAccessToken(tokenString string) (*MyClaims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*MyClaims, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, ErrInvalidTokenType
	}
	return claims, nil
}

// RefreshTokens validates the refresh token and issues new (access, refresh).
// If you want refresh-token rotation / revoke, persist a token ID (jti) and check it in DB.
func RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	user := model.User{ID: claims.Id, Username: claims.Username}
	return GenerateToken(user)
}
