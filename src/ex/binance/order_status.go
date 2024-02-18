package binance

// OrderStatus represents order status enum.
type OrderStatus string

// OrderType represents order type enum.
type OrderType string

// OrderSide represents order side enum.
type OrderSide string

type PositionSide string

type PositionStatus string

type PosithonSideStatus string

type EventType string

var (
	StatusNew             = OrderStatus("NEW")
	StatusPartiallyFilled = OrderStatus("PARTIALLY_FILLED")
	StatusFilled          = OrderStatus("FILLED")
	StatusCancelled       = OrderStatus("CANCELED")
	StatusExpired         = OrderStatus("EXPIRED")
	StatusInsurance       = OrderStatus("NEW_INSURANCE")
	StatusADL             = OrderStatus("NEW_ADL")

	TypeLimit      = OrderType("LIMIT")
	TypeMarket     = OrderType("MARKET")
	TypeSTOP       = OrderType("STOP")
	TypeTakeProfit = OrderType("TAKE_PROFIT")

	SideBuy  = OrderSide("BUY")
	SideSell = OrderSide("SELL")

	BOTH  = PositionSide("BOTH")
	LONG  = PositionSide("LONG")
	SHORT = PositionSide("SHORT")

	EventNew        = EventType("NEW")
	EventCanceled   = EventType("CANCELED")
	EventCalCulated = EventType("CALCULATED")
	EventExpired    = EventType("EXPIRED")
	EventTrade      = EventType("TRADE")

	POSITION_ISOLATED = PositionStatus("ISOLATED") // 逐仓
	POSITION_CROSSED  = PositionStatus("CROSSED")  // 全仓

	PosithonSingleSide = PosithonSideStatus("false") //单向持仓
	PosithonBothSide   = PosithonSideStatus("true")  //双向持仓
)
