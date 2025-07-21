package localjwt

import (
	"fmt"
	"time"

	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var jwtSecret []byte

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	logger.Logger.Debug("JWT generated", zap.String("token", tokenString))
	logger.Logger.Info("userID", zap.String("userID", userID))

	return tokenString, nil
}

func CheckJWT(tokenString string) (id string, isValid bool) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	logger.Logger.Debug("token", zap.Any("token", token))

	if err != nil {
		logger.Logger.Debug("ERROR in func Parse token", zap.Error(err))
		return "", false
	}
	//Взять userid из jwt
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok = claims["user_id"].(string)
		if !ok {
			logger.Logger.Debug("ERROR in func CheckJwt(get userId from jwt.MapClaims", zap.Error(err))
			return "", false
		}
		logger.Logger.Info("userID", zap.String("userID", id))
		return id, true
	}

	logger.Logger.Error("invalid token", zap.Error(err))
	return "", false
}

func LoadJWT(secret string) {
	jwtSecret = []byte(secret)
	logger.Logger.Debug("JWT secret loaded", zap.String("secret", "********"))
}
