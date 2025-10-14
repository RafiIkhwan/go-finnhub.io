package v1

import (
    "github.com/evrone/go-clean-template/internal/usecase"
    "github.com/gofiber/fiber/v2"
)

type authRoutes struct {
    a *usecase.Auth
}

func NewAuthRoutes(handler fiber.Router, a usecase.Auth) {
    r := &authRoutes{a: &a}
    handler.Post("/register", r.register)
    handler.Post("/login", r.login)
}

func (r *authRoutes) register(c *fiber.Ctx) error {
    var req struct { Email, Password string }
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    return r.a.Register(c.Context(), req.Email, req.Password)
}

func (r *authRoutes) login(c *fiber.Ctx) error {
    var req struct { Email, Password string }
    if err := c.BodyParser(&req); err != nil {
        return err
    }
    token, err := r.a.Login(c.Context(), req.Email, req.Password)
    if err != nil {
        return err
    }
    return c.JSON(map[string]string{"token": token})
}