package poloniex_api

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/models/poloniex"
	"io"
	"net/http"
	"time"
)

// НЕ УСПЕЛ СДЕЛАТЬ КРАСИВЕЕ

const baseURL = "https://api.poloniex.com/markets/"

// Adapter poloniex adapter
type adapter struct {
	client *http.Client
}

// NewAdapter impl
func NewAdapter() *adapter {
	return &adapter{client: &http.Client{
		Transport: &http.Transport{ // Не критично не закрыть конекшн
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS13,
			},
		},
		Timeout: 15 * time.Second,
	}}
}

func (a *adapter) GetCandleSticks(ctx context.Context, req poloniex.GetCandleStickReq) (poloniex.GetCandleStickResp, error) {
	url := fmt.Sprintf("%s%s/candles?interval=%s&startTime=%d&endTime=%d&limit=%d",
		baseURL,
		req.Pair,
		req.Interval,
		req.StartTime,
		req.EndTime,
		req.Limit,
	)
	resp, err := a.client.Get(url)
	if err != nil {
		return poloniex.GetCandleStickResp{}, fmt.Errorf("error in client.Get: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return poloniex.GetCandleStickResp{}, fmt.Errorf("error in io.ReadAll: %w", err)
	}
	return repackGetCandleSticks(req.Pair, respBody)
}

func (a *adapter) GetTrades(ctx context.Context, req poloniex.GetTradesReq) (poloniex.GetTradesResp, error) {
	url := fmt.Sprintf("%s%s/trades?&limit=%d",
		baseURL, req.Pair, req.Limit,
	)
	resp, err := a.client.Get(url)
	if err != nil {
		return poloniex.GetTradesResp{}, fmt.Errorf("error in client.Get: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return poloniex.GetTradesResp{}, fmt.Errorf("error in io.ReadAll: %w", err)
	}

	var trades = make([]tradeData, 0, 0)
	err = json.Unmarshal(respBody, &trades)
	if err != nil {
		return poloniex.GetTradesResp{}, fmt.Errorf("error in json.Unmarshal: %w", err)
	}
	return repackRT(req.Pair, trades)
}
