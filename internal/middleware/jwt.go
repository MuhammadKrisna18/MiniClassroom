package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"siakad-pro/internal/shared/apperrors"
)

func Protected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return apperrors.NewUnauthorized("missing jwt")
		}

		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "bearer ")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.NewUnauthorized("invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return apperrors.NewUnauthorized("invalid or expired jwt", err.Error())
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return apperrors.NewUnauthorized("invalid jwt claims")
		}

		if id, ok := claims["id"].(string); ok {
			c.Locals("userID", id)
		}
		if role, ok := claims["role"].(string); ok {
			c.Locals("userRole", role)
		}

		return c.Next()
	}
}

// GetUserID safely extracts the user ID from Fiber's Locals.
func GetUserID(c *fiber.Ctx) (string, error) {
	val := c.Locals("userID")
	if val == nil {
		return "", apperrors.NewUnauthorized("user not authenticated")
	}
	id, ok := val.(string)
	if !ok || id == "" {
		return "", apperrors.NewUnauthorized("invalid user id in context")
	}
	return id, nil
}

// GetUserRole safely extracts the user role from Fiber's Locals.
func GetUserRole(c *fiber.Ctx) (string, error) {
	val := c.Locals("userRole")
	if val == nil {
		return "", apperrors.NewUnauthorized("user not authenticated")
	}
	role, ok := val.(string)
	if !ok || role == "" {
		return "", apperrors.NewUnauthorized("invalid user role in context")
	}
	return role, nil
}

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, err := GetUserRole(c)
		if err != nil {
			return err
		}
		if userRole != role {
			return apperrors.NewForbidden("you do not have permission to perform this action")
		}
		return c.Next()
	}
}

