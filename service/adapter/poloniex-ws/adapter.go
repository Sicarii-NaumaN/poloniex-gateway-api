package poloniex_ws

import (
	"context"
	"encoding/json"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"time"

	"github.com/gorilla/websocket"
)

// НЕ УСПЕЛ СДЕЛАТЬ КРАСИВЕЕ

var pingMsg = []byte(`{"event": "ping"}`)

type adapter struct {
	conn *websocket.Conn

	reconnectURL string
}

func NewAdapter(conn *websocket.Conn, reconnectURL string) (*adapter, error) {
	return &adapter{conn: conn, reconnectURL: reconnectURL}, nil
}

func (a *adapter) Subscribe(ctx context.Context, channel string, symbols []poloniex.Pair) error {
	err := a.conn.WriteJSON(subscriptionMessage{
		Event:   "subscribe",
		Channel: []string{channel},
		Symbols: symbols,
	})
	if err != nil {
		return err
	}

	go a.keepAlive(ctx)

	return nil
}

// Listen слушает входящие сообщения и обрабатывает данные о сделках
func (a *adapter) Listen(ctx context.Context, tradeChan chan<- poloniex.RecentTrade) error {
	defer func() {
		close(tradeChan) // Закрываем канал после выхода из функции
	}()
	for {
		_, message, err := a.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				a.reconnect(ctx)
				continue
			} else {
				return err
			}
		}

		var incoming incomingMessage
		err = json.Unmarshal(message, &incoming)
		if err != nil {
			logger.Errorf("error in Listen json.Unmarshal: %v", err)
			continue
		}

		// Записываем в канал данные бачами
		if incoming.Channel == "trades" {
			for _, trade := range incoming.Data {
				tradeChan <- repackRT(trade)
			}
		}
	}
}

func (a *adapter) reconnect(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.conn, _, err = websocket.DefaultDialer.Dial(a.reconnectURL, nil)
			if err == nil {
				return
			}
		}
	}
}

func (a *adapter) keepAlive(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.conn.WriteMessage(websocket.TextMessage, pingMsg); err != nil {
				logger.Errorf("Error sending ping: %v", err)
				return
			}
			logger.Info("Sending ping to poloniex...")
		}
	}
}

func (a *adapter) Close() {
	a.conn.Close()
}
