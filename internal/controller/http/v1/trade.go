package v1

import (
    "net/http"

    "github.com/evrone/go-clean-template/internal/usecase"
    "github.com/gofiber/fiber/v2"
)

type tradeRoutes struct {
    uc *usecase.TradeStream
}

func NewTradeRoutes(handler fiber.Router, uc *usecase.TradeStream) {
    r := &tradeRoutes{uc: uc}
    handler.Post("/subscribe/:symbol", r.subscribe)
    handler.Post("/unsubscribe/:symbol", r.unsubscribe)
}

func (r *tradeRoutes) subscribe(c *fiber.Ctx) error {
    symbol := c.Params("symbol")
    _, err := r.uc.SubscribeToTrades(c.Context(), symbol)
    if err != nil {
        return err
    }
    return c.Status(http.StatusOK).SendString("Subscribed")
}

func (r *tradeRoutes) unsubscribe(c *fiber.Ctx) error {
    symbol := c.Params("symbol")
    err := r.uc.UnsubscribeFromTrades(symbol)
    if err != nil {
        return err
    }
    return c.Status(http.StatusOK).SendString("Unsubscribed")
}