package webapi

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"

    "github.com/evrone/go-clean-template/internal/entity"
    "github.com/gorilla/websocket"
)

type finnhubMessage struct {
    Type string          `json:"type"`
    Data []entity.TradeDetails `json:"data"`
}

type FinnhubRepoImpl struct {
    conn      *websocket.Conn
    apiKey    string
    mu        sync.Mutex
    channels  map[string]chan entity.TradeDetails
    done      chan struct{}
    errChan   chan error
}

func NewFinnhubRepo(apiKey string) (*FinnhubRepoImpl, error) {
    url := fmt.Sprintf("wss://ws.finnhub.io?token=%s", apiKey)
    conn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return nil, err
    }

    repo := &FinnhubRepoImpl{
        conn:     conn,
        apiKey:   apiKey,
        channels: make(map[string]chan entity.TradeDetails),
        done:     make(chan struct{}),
        errChan:  make(chan error, 1),
    }

    go repo.readLoop()  // Background goroutine untuk baca messages
    return repo, nil
}

func (r *FinnhubRepoImpl) readLoop() {
    defer close(r.done)
    for {
        _, message, err := r.conn.ReadMessage()
        if err != nil {
            r.errChan <- err
            return
        }
        var msg finnhubMessage
        if err := json.Unmarshal(message, &msg); err != nil {
            continue  // Log error jika perlu
        }
        if msg.Type == "trade" {
            r.mu.Lock()
            for _, trade := range msg.Data {
                if ch, ok := r.channels[trade.Symbol]; ok {
                    ch <- trade
                }
            }
            r.mu.Unlock()
        }
    }
}

func (r *FinnhubRepoImpl) Subscribe(ctx context.Context, symbol string) (<-chan entity.TradeDetails, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.channels[symbol]; ok {
        return nil, fmt.Errorf("already subscribed to %s", symbol)
    }

    ch := make(chan entity.TradeDetails, 100)  // Buffered channel
    r.channels[symbol] = ch

    subMsg := map[string]string{"type": "subscribe", "symbol": symbol}
    if err := r.conn.WriteJSON(subMsg); err != nil {
        delete(r.channels, symbol)
        close(ch)
        return nil, err
    }

    return ch, nil
}

func (r *FinnhubRepoImpl) Unsubscribe(symbol string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if ch, ok := r.channels[symbol]; ok {
        close(ch)
        delete(r.channels, symbol)
    }

    unsubMsg := map[string]string{"type": "unsubscribe", "symbol": symbol}
    return r.conn.WriteJSON(unsubMsg)
}

func (r *FinnhubRepoImpl) Close() error {
    r.mu.Lock()
    for symbol := range r.channels {
        r.Unsubscribe(symbol)
    }
    r.mu.Unlock()
    return r.conn.Close()
}