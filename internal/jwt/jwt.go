package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	config *JWTConfig

}

func NewJWTAuthenticator() *JWTAuthenticator {
	return &JWTAuthenticator{
		config: LoadAuthConfig(),
	}
}

func (j *JWTAuthenticator) GenerateToken(userId int64) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"iss": j.config.Issuer,
		"aud": j.config.Audience,
		"exp": time.Now().Add(j.config.DefaultExpiryHours).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(j.config.SecretKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWTAuthenticator) ValidateToken(encryptedToken string) (int64, error) {	
	token, err := jwt.Parse(
		encryptedToken, 
		func (token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.config.SecretKey, nil
		},
		jwt.WithAudience(j.config.Audience),
		jwt.WithIssuer(j.config.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		switch err {
			case jwt.ErrTokenExpired:
				return 0, fmt.Errorf("token expired")
			default:
				return 0, fmt.Errorf("invalid token")
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims["sub"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid subject claim")
		}
		return int64(userId), nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}