package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *Server) authMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := s.extractToken(c)
		if err != nil {
			respondError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		claims, err := s.parseToken(tokenStr)
		if err != nil {
			respondError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		if len(requiredRoles) > 0 && !roleAllowed(claims.Role, requiredRoles) {
			respondError(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}

		c.Set(userContextKey, &AuthContext{UserID: claims.UserID, Role: claims.Role})
		c.Next()
	}
}

func (s *Server) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return "", errors.New("invalid authorization header")
		}
		return parts[1], nil
	}

	if token := c.Query("token"); token != "" {
		return token, nil
	}

	return "", errors.New("missing authorization token")
}

func (s *Server) parseToken(tokenStr string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.tokenSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func roleAllowed(role string, allowed []string) bool {
	for _, r := range allowed {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

func (s *Server) generateToken(userID uint, role string) (string, error) {
	now := time.Now()
	claims := AuthClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.tokenSecret)
}
