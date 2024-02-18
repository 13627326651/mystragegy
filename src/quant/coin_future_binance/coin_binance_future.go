package future

import (
	"fmt"
	"time"

	"github.com/rootpd/binance"
)

func (b *Binance) GetFutureKlines(symbol string, limit int, interval binance.Interval) ([]*binance.Kline, error) {
	t := binance.KlinesRequest{
		Symbol:   symbol,
		Interval: interval,
		Limit:    limit,
	}

	return b.FutureCoinKlines(t)
}

func (b *Binance) GetTopLongShortPositionRatio(Pair string, Period string) ([]*binance.CoinTopLongShortPositionRatioInfo, error) {
	t := binance.CoinTopLongShortPositionRatioRequest{
		Pair:   Pair,
		Period: Period,
		//Limit:  1,
	}
	return b.CoinTopLongShortPositionRatio(t)

}

func (b *Binance) GetDepth(symbol string, limit int) (*binance.OrderBook, error) {
	return b.OrderCoinBook(binance.OrderBookRequest{
		Symbol: symbol,
		Limit:  limit,
	})

}

//合约持仓量  COIN
func (b *Binance) GetContractPosition(Pair string, ContractType string, Period string) ([]*binance.CoinContractPositionInfo, error) {
	t := binance.CoinContractPositionRequest{
		Pair:         Pair,
		ContractType: ContractType,
		Period:       Period,
		//Limit:  1,
	}
	return b.CoinContractPosition(t)

}

func (b *Binance) GetBinanceNewPrice(symbol string) (*binance.NewPriceInfo, error) {
	t := binance.NewPriceRequest{
		Symbol: symbol,
	}
	ts, err := b.CoinGetNewPrice(t)
	if err != nil {
		return nil, err
	}

	for _, pr := range ts {

		if pr.Symbol == symbol {
			return &binance.NewPriceInfo{
				Symbol:     pr.Symbol,
				Price:      pr.Price,
				UpdateTime: pr.UpdateTime,
			}, nil
		}
	}
	return nil, fmt.Errorf("can not get new price")

}

func (b *Binance) GetPriceChangeSituation(symbol string) ([]*binance.CoinPriceChangeSituationInfo, error) {
	t := binance.PriceChangeSituationRequest{
		Symbol: symbol,
	}
	return b.CoinPriceChangeSituation(t)

}

func (b *Binance) GetTakerlongshortRatio(Pair string, ContractType string, Period string) ([]*binance.CoinTakerlongshortRatioInfo, error) {
	t := binance.CoinTakerlongshortRatioRequest{
		Pair:         Pair,
		ContractType: ContractType,
		Period:       Period,
		//Limit:  1,
	}
	return b.CoinTakerlongshortRatio(t)

}

//当前最优挂单
func (b *Binance) GetBestBookTicker(symbol string) ([]*binance.CoinBestBookTickerInfo, error) {
	t := binance.BestBookTickerRequest{
		Symbol: symbol,
	}
	return b.CoinBestBookTicker(t)

}

// 查询当前所有挂单 币本位
func (b *Binance) QueryAllFutureOrder(symbol string) ([]*binance.ExecutedFutureOrder, error) {

	t := binance.CoinQueryFutureOrderRequest{
		Symbol:     symbol,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.CoinQueryFutureOrder(t)
}

func (b *Binance) GetAllFutureOrders(Symbol string) ([]*binance.CoinHistoryExecutedFutureOrder, error) {

	t := binance.CoinAllFutureOrdersRequest{
		Symbol:     Symbol,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}
	return b.CoinAllFutureOrders(t)

}

//获取未平仓合约数 币本位
func (b *Binance) GetOpenInterestNums(Symbol string) (*binance.CoinOpenInterestNumsInfo, error) {
	t := binance.OpenInterestNumsRequest{
		Symbol: Symbol,
	}
	return b.CoinOpenInterestNums(t)

}

func (b *Binance) GetGlobalLongShortAccountRatio() ([]*binance.GlobalLongShortAccountRatioInfo, error) {
	t := binance.CoinGlobalLongShortAccountRatioRequest{
		Pair:   "ETHUSD",
		Period: "15m",
		Limit:  10,
	}
	return b.CoinGlobalLongShortAccountRatio(t)
}

func (b *Binance) ChangeBinanceMarginType(symbol string, s binance.PositionStatus) error {

	t := binance.MarginTypeRequest{
		Symbol:     symbol,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
		MarginType: s,
	}

	return b.CoinChangeMarginType(t)

}

func (b *Binance) ChangeBinanceUserPositionSide(s binance.PosithonSideStatus) error {

	t := binance.ChangeUserPositionSideRequest{
		DualSidePosition: s,
		Timestamp:        time.Now(),
		RecvWindow:       5 * time.Second,
	}

	return b.CoinChangeUserPositionSide(t)
}

func (b *Binance) QueryUserPositionSide() (*binance.UserPositionSideInfo, error) {

	t := binance.UserPositionSideRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.CoinQueryUserPositionSide(t)

}

func (b *Binance) AdjustBinanceLeverage(symbol string, leverage int) error {

	t := binance.AdjustLeverageRequest{
		Symbol:     symbol,
		Timestamp:  time.Now(),
		RecvWindow: 1 * time.Second,
		Leverage:   leverage,
	}
	_, err := b.CoinAdjustLeverage(t)
	if err != nil {
		return err
	}
	return nil

}

//调整逐仓保证金

func (b *Binance) GetPositionMargin(Symbol string, PositionSide string, Amount float64, Type int) (*binance.PositionMarginInfo, error) {

	t := binance.PositionMarginRequest{
		PositionSide: PositionSide,
		Symbol:       Symbol,
		Timestamp:    time.Now(),
		RecvWindow:   5 * time.Second,
		Amount:       Amount,
		Type:         Type,
	}
	return b.CoinPositionMargin(t)

}

//获取手续费率 币本位

func (b *Binance) GetUserPoundageInfo(Symbol string) (*binance.UserPoundageInfo, error) {

	t := binance.UserPoundageRequest{
		Symbol:     Symbol,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.CoinUserPoundage(t)

}

func (b *Binance) NewBinanceFutureOrder(symbol string, quantity float64, price float64, side binance.OrderSide, positionSide binance.PositionSide) (*binance.FutureProcessedOrder, error) {

	t := binance.NewFutureOrderRequest{
		Symbol:       symbol,
		Quantity:     quantity,
		Price:        price,
		Side:         side,
		PositionSide: positionSide,
		TimeInForce:  binance.GTC,
		Type:         binance.TypeLimit,
		Timestamp:    time.Now(),
		//StopPrice:    3102.0,
		RecvWindow: 5 * time.Second,
	}

	return b.CoinNewFutureOrder(t)

}

func (b *Binance) CancelBinanceFutureOrder(symbol string, orderid int64) (*binance.CanceledFutureOrder, error) {

	t := binance.CancelFutureOrderRequest{
		Symbol:     symbol,
		OrderID:    orderid,
		Timestamp:  time.Now(),
		RecvWindow: 2 * time.Second,
	}

	return b.CoinCancelFutureOrder(t)

}

// 获取账户余额 币本位
func (b *Binance) GetFutureBalance() ([]*binance.FutureBalanceInfo, error) {

	t := binance.FutureBalanceRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}
	return b.CoinFutureBalance(t)

}

func (b *Binance) GetFutureAccount() (*binance.CoinFutureAccountInfo, error) {

	t := binance.FutureAccountRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.CoinFutureAccount(t)

}

//账户成交历史
func (b *Binance) GetUserTradesHistory(Symbol string) ([]*binance.CoinUserTradesHistoryInfo, error) {

	t := binance.CoinUserTradesHistoryRequest{

		Symbol:     Symbol,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.CoinUserTradesHistory(t)

}
