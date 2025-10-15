package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr := c.Get("Authorization")
        if tokenStr == "" {
            return fiber.ErrUnauthorized
        }
        
        tokenParts := strings.Split(tokenStr, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return fiber.ErrUnauthorized
        }

        tokenStr = tokenParts[1]

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return fiber.ErrUnauthorized
        }
        return c.Next()
    }
}