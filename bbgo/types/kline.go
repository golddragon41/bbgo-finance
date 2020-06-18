package types

type KLine interface {
	GetTrend() int
	GetChange() float64
	GetMaxChange() float64
	GetThickness() float64

	GetOpen() float64
	GetClose() float64
	GetHigh() float64
	GetLow() float64

	BounceUp() bool
	BounceDown() bool
	GetUpperShadowRatio() float64
	GetLowerShadowRatio() float64
}

