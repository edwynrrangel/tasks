package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/edwynrrangel/tasks/config"
	"github.com/edwynrrangel/tasks/logger"
)

type JWT interface {
	Generate(payload interface{}) (string, time.Time, error)
	IsValid(token string) bool
}

type jwtClaims struct {
	*jwt.StandardClaims
	Payload interface{} `json:"payload"`
}

type jwtApp struct {
	config *config.Configuration
}

func NewJWT(config *config.Configuration) JWT {
	return &jwtApp{
		config: config,
	}
}

func (j *jwtApp) Generate(payload interface{}) (string, time.Time, error) {
	// Define duration
	duration := time.Duration(j.config.JWTExpireTimeInMin) * time.Minute
	expiresAt := time.Now().Add(duration)
	// Define claims
	claims := &jwtClaims{
		Payload: payload,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return token
	tokenString, err := token.SignedString([]byte(j.config.JWTSecret))
	if err != nil {
		logger.Error(
			"error signing token",
			"func", "GenerateToken - token.SignedString",
			"error", err,
		)
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (j *jwtApp) IsValid(tokenString string) bool {
	// Validate token
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.JWTSecret), nil
	})
	if err != nil {
		logger.Error(
			"error parsing token",
			"func", "Validate - jwt.ParseWithClaims",
			"error", err,
		)
		return false
	}

	// Return true if token is valid
	return err == nil && token.Valid
}
