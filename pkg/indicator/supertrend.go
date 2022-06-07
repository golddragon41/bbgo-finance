package indicator

import (
	"math"
	"time"

	"github.com/c9s/bbgo/pkg/types"
)

//go:generate callbackgen -type Supertrend
type Supertrend struct {
	types.IntervalWindow
	ATRMultiplier float64 `json:"atrMultiplier"`

	AverageTrueRange *ATR

	trendPrices types.Float64Slice

	closePrice             float64
	previousClosePrice     float64
	uptrendPrice           float64
	previousUptrendPrice   float64
	downtrendPrice         float64
	previousDowntrendPrice float64

	trend         types.Direction
	previousTrend types.Direction
	tradeSignal   types.Direction

	EndTime         time.Time
	UpdateCallbacks []func(value float64)
}

func (inc *Supertrend) Last() float64 {
	return inc.trendPrices.Last()
}

func (inc *Supertrend) Index(i int) float64 {
	length := inc.Length()
	if length == 0 || length-i-1 < 0 {
		return 0
	}
	return inc.trendPrices[length-i-1]
}

func (inc *Supertrend) Length() int {
	return len(inc.trendPrices)
}
func (inc *Supertrend) Update(highPrice, lowPrice, closePrice float64) {
	if inc.Window <= 0 {
		panic("window must be greater than 0")
	}

	// Start with DirectionUp
	if inc.trend != types.DirectionUp && inc.trend != types.DirectionDown {
		inc.trend = types.DirectionUp
	}

	// Update ATR
	inc.AverageTrueRange.Update(highPrice, lowPrice, closePrice)

	// Update last prices
	inc.previousUptrendPrice = inc.uptrendPrice
	inc.previousDowntrendPrice = inc.downtrendPrice
	inc.previousClosePrice = inc.closePrice
	inc.previousTrend = inc.trend

	inc.closePrice = closePrice

	src := (highPrice + lowPrice) / 2

	// Update uptrend
	inc.uptrendPrice = src - inc.AverageTrueRange.Last()*inc.ATRMultiplier
	if inc.previousClosePrice > inc.previousUptrendPrice {
		inc.uptrendPrice = math.Max(inc.uptrendPrice, inc.previousUptrendPrice)
	}

	// Update downtrend
	inc.downtrendPrice = src + inc.AverageTrueRange.Last()*inc.ATRMultiplier
	if inc.previousClosePrice < inc.previousDowntrendPrice {
		inc.downtrendPrice = math.Min(inc.downtrendPrice, inc.previousDowntrendPrice)
	}

	// Update trend
	if inc.previousTrend == types.DirectionUp && inc.closePrice < inc.previousUptrendPrice {
		inc.trend = types.DirectionDown
	} else if inc.previousTrend == types.DirectionDown && inc.closePrice > inc.previousDowntrendPrice {
		inc.trend = types.DirectionUp
	} else {
		inc.trend = inc.previousTrend
	}

	// Update signal
	if inc.trend == types.DirectionUp && inc.previousTrend == types.DirectionDown {
		inc.tradeSignal = types.DirectionUp
	} else if inc.trend == types.DirectionDown && inc.previousTrend == types.DirectionUp {
		inc.tradeSignal = types.DirectionDown
	} else {
		inc.tradeSignal = types.DirectionNone
	}

	// Update trend price
	if inc.trend == types.DirectionDown {
		inc.trendPrices.Push(inc.downtrendPrice)
	} else {
		inc.trendPrices.Push(inc.uptrendPrice)
	}
}

func (inc *Supertrend) GetSignal() types.Direction {
	return inc.tradeSignal
}

var _ types.Series = &Supertrend{}

func (inc *Supertrend) calculateAndUpdate(kLines []types.KLine) {
	for _, k := range kLines {
		if inc.EndTime != zeroTime && !k.EndTime.After(inc.EndTime) {
			continue
		}
		inc.Update(k.GetHigh().Float64(), k.GetLow().Float64(), k.GetClose().Float64())
	}

	inc.EmitUpdate(inc.Last())
	inc.EndTime = kLines[len(kLines)-1].EndTime.Time()
}

func (inc *Supertrend) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}

	inc.calculateAndUpdate(window)
}

func (inc *Supertrend) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
