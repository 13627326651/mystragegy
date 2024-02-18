package quant

import (
	"tinyquant/src/mod"

	"github.com/rootpd/binance"
)

type Binance interface {
	InitBinance(apikey, secretkey string)
	NewBinanceFutureOrder(symbol string, quantity float64, price float64, stopprice float64, side binance.OrderSide, positionSide binance.PositionSide, id string) (*binance.FutureProcessedOrder, error)

	CancelBinanceFutureOrder(symbol string, orderid int64) (*binance.CanceledFutureOrder, error)
	QueryBinanceOneFutureOrder(symbol string, id string) (*binance.ExecutedFutureOrder, error)
	QueryBinanceAllFutureOrder(symbol string) ([]*binance.ExecutedFutureOrder, error)

	GetFutureBalance() ([]*binance.FutureBalanceInfo, error)
	GetFutureAccount(symbol string) []*binance.FuturePositions

	GetDepth(symbol string, limit int) (*binance.OrderBook, error)

	AdjustBinanceLeverage(symbol string, leverage int) error

	GetGlobalLongShortAccountRatio() ([]*binance.GlobalLongShortAccountRatioInfo, error)
	GetBinanceNewPrice(symbol string) (*binance.NewPriceInfo, error)
	ChangeBinanceMarginType(symbol string, s binance.PositionStatus) error
	ChangeBinanceUserPositionSide(s binance.PosithonSideStatus) error
	GetFutureDepthWs(symbol string) (chan *binance.DepthEvent, chan struct{})
	GetAccountWs() (chan *binance.FutureAccountEvent, chan struct{})

	GetFutureKlines(symbol string, limit int, interval binance.Interval) ([]*binance.Kline, error)

	GetKlineWs(symbol string, interval binance.Interval) chan *mod.Kline
}
