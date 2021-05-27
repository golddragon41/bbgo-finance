package okex

import (
	"fmt"
	"strings"

	"github.com/c9s/bbgo/pkg/exchange/okex/okexapi"
	"github.com/c9s/bbgo/pkg/types"
)

func toGlobalSymbol(symbol string) string {
	return strings.ReplaceAll(symbol, "-", "")
}

////go:generate sh -c "echo \"package okex\nvar spotSymbolMap = map[string]string{\n\" $(curl -s -L 'https://okex.com/api/v5/public/instruments?instType=SPOT' | jq -r '.data[] | \"\\(.instId | sub(\"-\" ; \"\") | tojson ): \\( .instId | tojson),\n\"') \"\n}\" > symbols.go"
//go:generate go run gensymbols.go
func toLocalSymbol(symbol string) string {
	if s, ok := spotSymbolMap[symbol]; ok {
		return s
	}

	log.Errorf("failed to look up local symbol from %s", symbol)
	return symbol
}

func toGlobalTicker(marketTicker okexapi.MarketTicker) *types.Ticker {
	return &types.Ticker{
		Time:   marketTicker.Timestamp.Time(),
		Volume: marketTicker.Volume24H.Float64(),
		Last:   marketTicker.Last.Float64(),
		Open:   marketTicker.Open24H.Float64(),
		High:   marketTicker.High24H.Float64(),
		Low:    marketTicker.Low24H.Float64(),
		Buy:    marketTicker.BidPrice.Float64(),
		Sell:   marketTicker.AskPrice.Float64(),
	}
}

func toGlobalBalance(balanceSummaries []okexapi.BalanceSummary) types.BalanceMap {
	var balanceMap = types.BalanceMap{}
	for _, balanceSummary := range balanceSummaries {
		for _, balanceDetail := range balanceSummary.Details {
			balanceMap[balanceDetail.Currency] = types.Balance{
				Currency:  balanceDetail.Currency,
				Available: balanceDetail.CashBalance,
				Locked:    balanceDetail.Frozen,
			}
		}
	}
	return balanceMap
}

type WebsocketSubscription struct {
	Channel      string `json:"channel"`
	InstrumentID string `json:"instId,omitempty"`
	InstrumentType string `json:"instType,omitempty"`
}

var CandleChannels = []string{
	"candle1Y",
	"candle6M", "candle3M", "candle1M",
	"candle1W",
	"candle1D", "candle2D", "candle3D", "candle5D",
	"candle12H", "candle6H", "candle4H", "candle2H", "candle1H",
	"candle30m", "candle15m", "candle5m", "candle3m", "candle1m",
}

func convertIntervalToCandle(interval string) string {
	switch interval {

	case "1h", "2h", "4h", "6h", "12h", "1d", "3d":
		return "candle" + strings.ToUpper(interval)

	case "1m", "5m", "15m", "30m":
		return "candle" + interval

	}

	return "candle" + interval
}

func convertSubscription(s types.Subscription) (WebsocketSubscription, error) {
	// binance uses lower case symbol name,
	// for kline, it's "<symbol>@kline_<interval>"
	// for depth, it's "<symbol>@depth OR <symbol>@depth@100ms"
	switch s.Channel {
	case types.KLineChannel:
		// Channel names are:
		return WebsocketSubscription{
			Channel:      convertIntervalToCandle(s.Options.Interval),
			InstrumentID: toLocalSymbol(s.Symbol),
		}, nil

	case types.BookChannel:
		return WebsocketSubscription{
			Channel:      "books",
			InstrumentID: toLocalSymbol(s.Symbol),
		}, nil
	}

	return WebsocketSubscription{}, fmt.Errorf("unsupported public stream channel %s", s.Channel)
}
