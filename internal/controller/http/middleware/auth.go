package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr := c.Get("Authorization")
        if tokenStr == "" {
            return fiber.ErrUnauthorized
        }
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return fiber.ErrUnauthorized
        }
        return c.Next()
    }
}