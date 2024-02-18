package future

import (
	"fmt"
	"time"
	. "tinyquant/src/logger"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

/*
条件单的触发必须:

如果订单参数priceProtect为true:
达到触发价时，MARK_PRICE(标记价格)与CONTRACT_PRICE(合约最新价)之间的价差不能超过改symbol触发保护阈值
触发保护阈值请参考接口GET /fapi/v1/exchangeInfo 返回内容相应symbol中"triggerProtect"字段

STOP, STOP_MARKET 止损单:

买入: 最新合约价格/标记价格高于等于触发价stopPrice
卖出: 最新合约价格/标记价格低于等于触发价stopPrice

TAKE_PROFIT, TAKE_PROFIT_MARKET 止盈单:

买入: 最新合约价格/标记价格低于等于触发价stopPrice
卖出: 最新合约价格/标记价格高于等于触发价stopPrice

TRAILING_STOP_MARKET 跟踪止损单:

买入: 当合约价格/标记价格区间最低价格低于激活价格activationPrice,且最新合约价格/标记价高于等于最低价设定回调幅度。
卖出: 当合约价格/标记价格区间最高价格高于激活价格activationPrice,且最新合约价格/标记价低于等于最高价设定回调幅度。

TRAILING_STOP_MARKET 跟踪止损单如果遇到报错 {"code": -2021, "msg": "Order would immediately trigger."}
表示订单不满足以下条件:

买入: 指定的activationPrice 必须小于 latest price
卖出: 指定的activationPrice 必须大于 latest price
newOrderRespType 如果传 RESULT:

MARKET 订单将直接返回成交结果；
配合使用特殊 timeInForce 的 LIMIT 订单将直接返回成交或过期拒绝结果。
STOP_MARKET, TAKE_PROFIT_MARKET 配合 closePosition=true:

条件单触发依照上述条件单触发逻辑
条件触发后，平掉当时持有所有多头仓位(若为卖单)或当时持有所有空头仓位(若为买单)
不支持 quantity 参数
自带只平仓属性，不支持reduceOnly参数
双开模式下,LONG方向上不支持BUY; SHORT 方向上不支持SELL

*/
//下单

// 双向持仓下 做空 : side = sell , PositionSide = short
// 双向持仓下 平空 : side = buy  , PositionSide = short

// 双向持仓下 做多 : side = buy , PositionSide = long
// 双向持仓下 平多 : side = sell , PositionSide = long

func (b *Binance) NewBinanceFutureOrder(symbol string, quantity float64, price float64, stopprice float64, side binance.OrderSide, positionSide binance.PositionSide, customOrderId string) (*binance.FutureProcessedOrder, error) {

	t := binance.NewFutureOrderRequest{
		Symbol:           symbol,
		Quantity:         quantity,
		Price:            price,
		Side:             side,
		PositionSide:     positionSide,
		TimeInForce:      binance.GTC,
		Type:             binance.TypeLimit,
		Timestamp:        time.Now(),
		RecvWindow:       5 * time.Second,
		NewClientOrderID: customOrderId,
		StopPrice:        stopprice,
	}

	if t.StopPrice != 0 {
		t.Type = binance.TypeSTOP
	}

	return b.NewFutureOrder(t)

}

// 取消订单

func (b *Binance) CancelBinanceFutureOrder(symbol string, orderid int64) (*binance.CanceledFutureOrder, error) {
	t := binance.CancelFutureOrderRequest{
		Symbol:     symbol,
		OrderID:    orderid,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.CancelFutureOrder(t)

}

// 查询当前挂单
func (b *Binance) QueryBinanceOneFutureOrder(symbol string, id string) (*binance.ExecutedFutureOrder, error) {
	t := binance.QueryFutureOrderRequest{
		Symbol:            symbol,
		Timestamp:         time.Now(),
		RecvWindow:        5 * time.Second,
		OrigClientOrderId: id,
	}

	return b.QueryOneFutureOrder(t)
}

// 查询所有挂单
func (b *Binance) QueryBinanceAllFutureOrder(symbol string) ([]*binance.ExecutedFutureOrder, error) {
	t := binance.QueryFutureOrderRequest{
		Symbol:     symbol,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	return b.QueryAllFutureOrder(t)
}

//// 获取 历史 u本位 订单

func (b *Binance) GetAllFutureOrders() {

	t := binance.AllFutureOrdersRequest{
		Symbol:     "ETHUSDT",
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
		//StartTime:  s.Unix(),
		//EndTime:    time.Now().Unix(),
	}
	ts, err := b.QueryAllHistoryFutureOrders(t)
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range ts {
		fmt.Println(k, " : ", v)
	}
}

// 获取账户余额
func (b *Binance) GetFutureBalance() ([]*binance.FutureBalanceInfo, error) {

	t := binance.FutureBalanceRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}
	ts, err := b.FutureBalance(t)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (b *Binance) GetFutureAccount(symbol string) []*binance.FuturePositions {
	t := binance.FutureAccountRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	ts, err := b.FutureAccount(t)
	if err != nil {
		Logger.Error("GetFutureAccount failed", zap.Error(err))
		return nil
	}

	res := make([]*binance.FuturePositions, 0)
	for _, v := range ts.Positions {
		if v.Symbol == symbol {
			res = append(res, v)
		}
	}
	return res
}

func (b *Binance) GetDepth(symbol string, limit int) (*binance.OrderBook, error) {
	return b.OrderBook(binance.OrderBookRequest{
		Symbol: symbol,
		Limit:  limit,
	})

}

//获取手续费率

func (b *Binance) GetUserPoundageInfo() {

	now := time.Now()

	tp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()-1, 0, time.Local)

	t := binance.UserPoundageRequest{
		Symbol:     "ETHUSDT",
		Timestamp:  tp,
		RecvWindow: 5 * time.Second,
	}

	ts, err := b.UserPoundage(t)
	if err != nil {
		fmt.Println("UserPoundageInfo : ", err)
	}

	fmt.Println("xxxx", ts)
}

//调整开仓杠杆

func (b *Binance) AdjustBinanceLeverage(symbol string, leverage int) error {

	t := binance.AdjustLeverageRequest{
		Symbol:     symbol,
		Timestamp:  time.Now(),
		RecvWindow: 1 * time.Second,
		Leverage:   leverage,
	}
	_, err := b.AdjustLeverage(t)

	if err != nil {
		return err
	}
	return nil
}

//调整逐仓保证金

func (b *Binance) GetPositionMargin() {
	now := time.Now()

	tp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()-1, 0, time.Local)

	t := binance.PositionMarginRequest{
		PositionSide: "LONG",
		Symbol:       "ETHUSDT",
		Timestamp:    tp,
		RecvWindow:   5 * time.Second,
		Amount:       888.8,
		Type:         1,
	}
	ts, err := b.PositionMargin(t)
	if err != nil {
		fmt.Println("GetPositionMarginInfo: ", err)
	}
	fmt.Println("xxxx", ts)
}

//账户成交历史
func (b *Binance) GetUserTradesHistory() {
	//now := time.Now()

	//tp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()-1, 0, time.Local)

	//s := time.Date(2020, time.January, 0, 0, 0, 0, 0, time.Local)

	t := binance.UserTradesHistoryRequest{
		//StartTime:  s.UnixNano(),
		//EndTime:    time.Now().UnixNano(),
		//FromId:     1,
		//Limit:      5,
		Symbol:     "ETHUSDT",
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	ts, err := b.UserTradesHistory(t)
	if err != nil {
		fmt.Println("GetUserTradesHistoryInfo: ", err)
	}
	for k, v := range ts {
		fmt.Println(k, " : ", v)
	}
}

//最新标记价格和资金费率
func (b *Binance) GetPremiumAndFundsRate() {
	t := binance.PremiumAndFundsRateRequest{
		Symbol: "ETHUSDT",
	}
	ts, err := b.PremiumAndFundsRate(t)
	if err != nil {
		fmt.Println("GetPremiumAndFundsRate: ", err)
	}
	fmt.Println("xxxx", ts)
}

//24hr价格变动情况
func (b *Binance) GetPriceChangeSituation() {
	t := binance.PriceChangeSituationRequest{
		Symbol: "ETHUSDT",
	}
	ts, err := b.PriceChangeSituation(t)
	if err != nil {
		fmt.Println("GetPriceChangeSituation: ", err)
	}
	fmt.Println("xxxx", ts)
}

//获取未平仓合约数
func (b *Binance) GetOpenInterestNums() {
	t := binance.OpenInterestNumsRequest{
		Symbol: "ETHUSDT",
	}
	ts, err := b.OpenInterestNums(t)
	if err != nil {
		fmt.Println("GetOpenInterestNums: ", err)
	}
	fmt.Println("xxxx", ts)
}

//当前最优挂单
func (b *Binance) GetBestBookTicker() {
	t := binance.BestBookTickerRequest{
		Symbol: "ETHUSDT",
	}
	ts, err := b.BestBookTicker(t)
	if err != nil {
		fmt.Println("GetBestBookTicker: ", err)
	}
	fmt.Println("xxxx", ts)
}

//合约持仓量
func (b *Binance) GetContractPosition() {
	t := binance.ContractPositionRequest{
		Symbol: "ETHUSDT",
		Period: "4h",
		//Limit:  1,
	}
	ts, err := b.ContractPosition(t)
	if err != nil {
		fmt.Println("GetContractPosition: ", err)
	}

	for k, v := range ts {
		fmt.Println(k, "xxxx", v)
	}

}
func (b *Binance) GetTopLongShortPositionRatio() {
	t := binance.TopLongShortPositionRatioRequest{
		Symbol: "ETHUSDT",
		Period: "4h",
		//Limit:  1,
	}
	ts, err := b.TopLongShortPositionRatio(t)
	if err != nil {
		fmt.Println("GetTopLongShortPositionRatio: ", err)
	}

	for k, v := range ts {
		fmt.Println(k, "xxxx", v)
	}

}

func (b *Binance) GetGlobalLongShortAccountRatio() ([]*binance.GlobalLongShortAccountRatioInfo, error) {
	t := binance.GlobalLongShortAccountRatioRequest{
		Symbol: "ETHUSDT",
		Period: "15m",
		Limit:  10,
	}
	return b.GlobalLongShortAccountRatio(t)
}

func (b *Binance) GetTakerlongshortRatio() {
	t := binance.TakerlongshortRatioRequest{
		Symbol: "ETHUSDT",
		Period: "4h",
		//Limit:  1,
	}
	ts, err := b.TakerlongshortRatio(t)
	if err != nil {
		fmt.Println("GetTakerlongshortRatio: ", err)
	}

	for k, v := range ts {
		fmt.Println(k, "xxxx", v)
	}
}

func (b *Binance) QueryBinanceUserPositionSide() {
	t := binance.UserPositionSideRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	}

	ts, err := b.QueryUserPositionSide(t)
	if err != nil {
		fmt.Println("QueryUserPositionSide", err)
	}
	fmt.Println("xxxx", ts)
}

func (b *Binance) GetBinanceNewPrice(symbol string) (*binance.NewPriceInfo, error) {
	t := binance.NewPriceRequest{
		Symbol: symbol,
	}
	return b.GetNewPrice(t)

}

func (b *Binance) GetFutureKlines(symbol string, limit int, interval binance.Interval) ([]*binance.Kline, error) {
	t := binance.KlinesRequest{
		Symbol:   symbol,
		Interval: interval,
		Limit:    limit,
	}

	return b.FutureKlines(t)

}

func (b *Binance) ChangeBinanceMarginType(symbol string, s binance.PositionStatus) error {

	t := binance.MarginTypeRequest{
		Symbol:     util.ETHUSDT,
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
		MarginType: s,
	}

	return b.ChangeMarginType(t)

}

func (b *Binance) ChangeBinanceUserPositionSide(s binance.PosithonSideStatus) error {

	t := binance.ChangeUserPositionSideRequest{
		DualSidePosition: s,
		Timestamp:        time.Now(),
		RecvWindow:       5 * time.Second,
	}

	return b.ChangeUserPositionSide(t)
}
