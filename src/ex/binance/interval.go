package binance

// Interval represents interval enum.
type Interval string

var (
	Minute         = Interval("1m")
	ThreeMinutes   = Interval("3m")
	FiveMinutes    = Interval("5m")
	FifteenMinutes = Interval("15m")
	ThirtyMinutes  = Interval("30m")
	Hour           = Interval("1h")
	TwoHours       = Interval("2h")
	FourHours      = Interval("4h")
	EightHours     = Interval("8h")
	SixHours       = Interval("6h")
	TwelveHours    = Interval("12h")
	Day            = Interval("1d")
	ThreeDays      = Interval("3d")
	Week           = Interval("1w")
	Month          = Interval("1M")
)

// TimeInForce represents timeInForce enum.
type TimeInForce string

var (
	GTC = TimeInForce("GTC") // 订单一直有效直到取消
	IOC = TimeInForce("IOC") // 订单必须立即以限价或更佳的价格成交，否则自动取消
)
