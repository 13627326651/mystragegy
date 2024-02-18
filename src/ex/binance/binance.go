package binance

import (
	"fmt"
	"time"
)

// Binance is wrapper for Binance API.
//
// Read web documentation for more endpoints descriptions and list of
// mandatory and optional params. Wrapper is not responsible for client-side
// validation and only sends requests further.
//
// For each API-defined enum there's a special type and list of defined
// enum values to be used.
type Binance interface {
	// Ping tests connectivity.
	Ping() error
	// Time returns server time.
	Time() (time.Time, error)

	ExchangeInfo() (string, error)

	// 最新价格
	NewPrice(nb OrderNewPriceRequest) (*NewPrice, error)

	// OrderBook returns list of orders.
	OrderBook(obr OrderBookRequest) (*OrderBook, error)
	// AggTrades returns compressed/aggregate list of trades.
	AggTrades(atr AggTradesRequest) ([]*AggTrade, error)
	// Klines returns klines/candlestick data.
	Klines(kr KlinesRequest) ([]*Kline, error)
	// Ticker24 returns 24hr price change statistics.
	Ticker24(tr TickerRequest) (*Ticker24, error)
	// TickerAllPrices returns ticker data for symbols.
	TickerAllPrices() ([]*PriceTicker, error)
	// TickerAllBooks returns tickers for all books.
	TickerAllBooks() ([]*BookTicker, error)

	// NewOrder places new order and returns ProcessedOrder.
	NewOrder(nor NewOrderRequest) (*ProcessedOrder, error)
	// NewOrder places testing order.
	NewOrderTest(nor NewOrderRequest) error
	// QueryOrder returns data about existing order.
	QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error)
	// CancelOrder cancels order.
	CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error)
	// OpenOrders returns list of open orders.
	OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error)
	// AllOrders returns list of all previous orders.
	AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error)

	// Account returns account data.
	Account(ar AccountRequest) (*Account, error)
	// MyTrades list user's trades.
	MyTrades(mtr MyTradesRequest) ([]*Trade, error)
	// Withdraw executes withdrawal.
	Withdraw(wr WithdrawRequest) (*WithdrawResult, error)
	// DepositHistory lists deposit data.
	DepositHistory(hr HistoryRequest) ([]*Deposit, error)
	// WithdrawHistory lists withdraw data.
	WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error)

	// StartUserDataStream starts stream and returns Stream with ListenKey.
	StartUserDataStream() (*Stream, error)
	// KeepAliveUserDataStream prolongs stream livespan.
	KeepAliveUserDataStream(s *Stream) error
	// CloseUserDataStream closes opened stream.
	CloseUserDataStream(s *Stream) error

	DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)
	KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)
	TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)
	UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error)

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// 合约 接口

	// 下单
	NewFutureOrder(nfr NewFutureOrderRequest) (*FutureProcessedOrder, error)
	// 取消订单
	CancelFutureOrder(cfr CancelFutureOrderRequest) (*CanceledFutureOrder, error)

	//查询一个挂单
	QueryOneFutureOrder(qfo QueryFutureOrderRequest) (*ExecutedFutureOrder, error)

	// 查询当前所有挂单
	QueryAllFutureOrder(qfo QueryFutureOrderRequest) ([]*ExecutedFutureOrder, error)

	// 查询所有历史订单
	QueryAllHistoryFutureOrders(afo AllFutureOrdersRequest) ([]*ExecutedFutureOrder, error)

	// 获取账户余额
	FutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error)

	// 获取 账户信息
	FutureAccount(far FutureAccountRequest) (*FutureAccountInfo, error)

	//用户手续费率
	UserPoundage(upr UserPoundageRequest) (*UserPoundageInfo, error)

	//调整开仓杠杆
	AdjustLeverage(alr AdjustLeverageRequest) (*AdjustLeverageInfo, error)

	//调整逐仓保证金
	PositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error)

	//用户成交历史
	UserTradesHistory(uth UserTradesHistoryRequest) ([]*UserTradesHistoryInfo, error)

	//最新标记价格和资金费率
	PremiumAndFundsRate(pfrr PremiumAndFundsRateRequest) (*PremiumAndFundsRateInfo, error)

	//获取k线
	FutureKlines(kr KlinesRequest) ([]*Kline, error)

	// 获取最新价格
	GetNewPrice(npr NewPriceRequest) (*NewPriceInfo, error)

	//24hr价格变动情况
	PriceChangeSituation(pcsr PriceChangeSituationRequest) (*PriceChangeSituationInfo, error)

	//获取未平仓合约数
	OpenInterestNums(oir OpenInterestNumsRequest) (*OpenInterestNumsInfo, error)

	//当前最优挂单
	BestBookTicker(bbtr BestBookTickerRequest) (*BestBookTickerInfo, error)

	//合约持仓量
	ContractPosition(bpr ContractPositionRequest) ([]*ContractPositionInfo, error)

	//大户持仓量多空比
	TopLongShortPositionRatio(tspr TopLongShortPositionRatioRequest) ([]*TopLongShortPositionRatioInfo, error)

	//多空持仓人数比
	GlobalLongShortAccountRatio(glsarr GlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error)

	//合约主动买卖量
	TakerlongshortRatio(tlr TakerlongshortRatioRequest) ([]*TakerlongshortRatioInfo, error)

	// 查询持仓模式
	QueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error)

	// 变换逐全仓模式
	ChangeMarginType(mtr MarginTypeRequest) error

	// 更改持仓模式
	ChangeUserPositionSide(ups ChangeUserPositionSideRequest) error

	StartFutureUserDataStream() (*Stream, error)
	KeepAliveFutureUserDataStream(s *Stream) error
	CloseFutureUserDataStream(s *Stream) error

	FutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)
	FutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)
	FutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)
	FutureUserDataWebsocket(udwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error)

	// //币本位合约如下--------------------------------------------------------------------------------------
	//获取深度
	OrderCoinBook(obr OrderBookRequest) (*OrderBook, error)

	//K线数据
	FutureCoinKlines(kr KlinesRequest) ([]*Kline, error)

	//大户持仓量多空比
	CoinTopLongShortPositionRatio(tspr CoinTopLongShortPositionRatioRequest) ([]*CoinTopLongShortPositionRatioInfo, error)

	//合约持仓量
	CoinContractPosition(bpr CoinContractPositionRequest) ([]*CoinContractPositionInfo, error)

	// 最新价格
	CoinGetNewPrice(npr NewPriceRequest) ([]*CoinNewPriceInfo, error)

	//24hr价格变动情况
	CoinPriceChangeSituation(pcsr PriceChangeSituationRequest) ([]*CoinPriceChangeSituationInfo, error)

	//合约主动买卖量
	CoinTakerlongshortRatio(tlr CoinTakerlongshortRatioRequest) ([]*CoinTakerlongshortRatioInfo, error)

	//当前最优挂单
	CoinBestBookTicker(bbtr BestBookTickerRequest) ([]*CoinBestBookTickerInfo, error)

	// 查询当前所有挂单
	CoinQueryFutureOrder(qfo CoinQueryFutureOrderRequest) ([]*ExecutedFutureOrder, error)

	// 查询所有历史订单
	CoinAllFutureOrders(afo CoinAllFutureOrdersRequest) ([]*CoinHistoryExecutedFutureOrder, error)

	//获取未平仓合约数
	CoinOpenInterestNums(oir OpenInterestNumsRequest) (*CoinOpenInterestNumsInfo, error)

	//多空持仓人数比
	CoinGlobalLongShortAccountRatio(glsarr CoinGlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error)

	// 变换逐全仓模式
	CoinChangeMarginType(mtr MarginTypeRequest) error

	// 更改持仓模式
	CoinChangeUserPositionSide(ups ChangeUserPositionSideRequest) error

	// 查询持仓模式
	CoinQueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error)

	//调整开仓杠杆

	CoinAdjustLeverage(alr AdjustLeverageRequest) (*CoinAdjustLeverageInfo, error)

	//调整逐仓保证金
	CoinPositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error)

	//用户手续费率
	CoinUserPoundage(upr UserPoundageRequest) (*UserPoundageInfo, error)

	// 下单
	CoinNewFutureOrder(nfr NewFutureOrderRequest) (*FutureProcessedOrder, error)

	// 取消订单
	CoinCancelFutureOrder(cfr CancelFutureOrderRequest) (*CanceledFutureOrder, error)

	// 获取账户余额
	CoinFutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error)

	// 获取 账户信息
	CoinFutureAccount(far FutureAccountRequest) (*CoinFutureAccountInfo, error)

	//用户成交历史
	CoinUserTradesHistory(uth CoinUserTradesHistoryRequest) ([]*CoinUserTradesHistoryInfo, error)

	//币本位websocket  有限档深度信息
	CoinFutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)

	//websocket 归集交易
	CoinFutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)

	//K线
	CoinFutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)

	CoinFutureUserDataWebsocket(udwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error)

	//请求账户信息
	CoinAccountInfoWebsocket(udwr UserDataWebsocketRequest) (chan *CoinAccountInfo, chan struct{}, error)

	StartCoinFutureUserDataStream() (*Stream, error)
	KeepAliveCoinFutureUserDataStream(s *Stream) error
	CloseCoinFutureUserDataStream(s *Stream) error
}

type binance struct {
	Service Service
}

// Error represents Binance error structure with error code and message.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Error returns formatted error message.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// NewBinance returns Binance instance.
func NewBinance(service Service) Binance {
	return &binance{
		Service: service,
	}
}

// Ping tests connectivity.
func (b *binance) Ping() error {
	return b.Service.Ping()
}

// Time returns server time.
func (b *binance) Time() (time.Time, error) {
	return b.Service.Time()
}

// OrderBook represents Bids and Asks.
type OrderBook struct {
	LastUpdateID int
	BeforeUID    int
	UpdateID     int
	MessageTime  time.Time
	Bids         []*Order // 买方出价
	Asks         []*Order //卖方出价
}

type DepthEvent struct {
	WSEvent
	EventTime time.Time
	OrderBook
}

func (b *binance) ExchangeInfo() (string, error) {
	return b.Service.ExchangeInfo()
}

type OrderNewPriceRequest struct {
	Symbol string
}

type NewPrice struct {
	Symbol string
	Price  string
}

func (b *binance) NewPrice(nb OrderNewPriceRequest) (*NewPrice, error) {
	return b.Service.NewPrice(nb)
}

// Order represents single order information.
type Order struct {
	Price    float64
	Quantity float64
}

// OrderBookRequest represents OrderBook request data.
type OrderBookRequest struct {
	Symbol string
	Limit  int
}

// OrderBook returns list of orders.
func (b *binance) OrderBook(obr OrderBookRequest) (*OrderBook, error) {
	return b.Service.OrderBook(obr)
}

func (b *binance) OrderCoinBook(obr OrderBookRequest) (*OrderBook, error) {
	return b.Service.OrderCoinBook(obr)
}

// AggTrade represents aggregated trade.
type AggTrade struct {
	ID             int
	Price          float64
	Quantity       float64
	FirstTradeID   int
	LastTradeID    int
	Timestamp      time.Time
	BuyerMaker     bool
	BestPriceMatch bool
}

type AggTradeEvent struct {
	WSEvent
	AggTrade
}

// AggTradesRequest represents AggTrades request data.
type AggTradesRequest struct {
	Symbol    string
	FromID    int64
	StartTime int64
	EndTime   int64
	Limit     int
}

// AggTrades returns compressed/aggregate list of trades.
func (b *binance) AggTrades(atr AggTradesRequest) ([]*AggTrade, error) {
	return b.Service.AggTrades(atr)
}

// KlinesRequest represents Klines request data.
type KlinesRequest struct {
	Symbol    string
	Interval  Interval
	Limit     int
	StartTime int64
	EndTime   int64
}

// Kline represents single Kline information.
type Kline struct {
	OpenTime                 time.Time
	Open                     float64 //开盘价
	High                     float64 //最高
	Low                      float64 //最低
	Close                    float64 //收盘价
	Volume                   float64 //总成交量 ETH
	CloseTime                time.Time
	QuoteAssetVolume         float64 //总成交额 USDT
	NumberOfTrades           int     // k线成交比数
	TakerBuyBaseAssetVolume  float64 //买单成交量 ETH
	TakerBuyQuoteAssetVolume float64 //买单成交额 USDT
	Final                    bool    //这跟k线是否结束
}

type KlineEvent struct {
	WSEvent
	Interval     Interval
	FirstTradeID int64
	LastTradeID  int64
	Final        bool
	Kline
}

// Klines returns klines/candlestick data.
func (b *binance) Klines(kr KlinesRequest) ([]*Kline, error) {
	return b.Service.Klines(kr)
}

// TickerRequest represents Ticker request data.
type TickerRequest struct {
	Symbol string
}

// Ticker24 represents data for 24hr ticker.
type Ticker24 struct {
	PriceChange        float64
	PriceChangePercent float64
	WeightedAvgPrice   float64
	PrevClosePrice     float64
	LastPrice          float64
	BidPrice           float64
	AskPrice           float64
	OpenPrice          float64
	HighPrice          float64
	LowPrice           float64
	Volume             float64
	OpenTime           time.Time
	CloseTime          time.Time
	FirstID            int
	LastID             int
	Count              int
}

// Ticker24 returns 24hr price change statistics.
func (b *binance) Ticker24(tr TickerRequest) (*Ticker24, error) {
	return b.Service.Ticker24(tr)
}

// PriceTicker represents ticker data for price.
type PriceTicker struct {
	Symbol string
	Price  float64
}

// TickerAllPrices returns ticker data for symbols.
func (b *binance) TickerAllPrices() ([]*PriceTicker, error) {
	return b.Service.TickerAllPrices()
}

func (b *binance) FutureCoinKlines(kr KlinesRequest) ([]*Kline, error) {
	return b.Service.FutureCoinKlines(kr)
}

func (b *binance) CoinTopLongShortPositionRatio(tspr CoinTopLongShortPositionRatioRequest) ([]*CoinTopLongShortPositionRatioInfo, error) {
	return b.Service.CoinTopLongShortPositionRatio(tspr)
}

func (b *binance) CoinContractPosition(bpr CoinContractPositionRequest) ([]*CoinContractPositionInfo, error) {
	return b.Service.CoinContractPosition(bpr)
}

func (b *binance) CoinGetNewPrice(npr NewPriceRequest) ([]*CoinNewPriceInfo, error) {
	return b.Service.CoinGetNewPrice(npr)
}

func (b *binance) CoinOpenInterestNums(oir OpenInterestNumsRequest) (*CoinOpenInterestNumsInfo, error) {
	return b.Service.CoinOpenInterestNums(oir)
}

func (b *binance) CoinPriceChangeSituation(pcsr PriceChangeSituationRequest) ([]*CoinPriceChangeSituationInfo, error) {
	return b.Service.CoinPriceChangeSituation(pcsr)
}

func (b *binance) CoinTakerlongshortRatio(tlr CoinTakerlongshortRatioRequest) ([]*CoinTakerlongshortRatioInfo, error) {
	return b.Service.CoinTakerlongshortRatio(tlr)
}

func (b *binance) CoinBestBookTicker(bbtr BestBookTickerRequest) ([]*CoinBestBookTickerInfo, error) {
	return b.Service.CoinBestBookTicker(bbtr)
}

func (b *binance) CoinQueryFutureOrder(qfo CoinQueryFutureOrderRequest) ([]*ExecutedFutureOrder, error) {
	return b.Service.CoinQueryFutureOrder(qfo)
}

func (b *binance) CoinAllFutureOrders(afo CoinAllFutureOrdersRequest) ([]*CoinHistoryExecutedFutureOrder, error) {
	return b.Service.CoinAllFutureOrders(afo)
}

func (b *binance) CoinGlobalLongShortAccountRatio(glsarr CoinGlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error) {
	return b.Service.CoinGlobalLongShortAccountRatio(glsarr)
}

func (b *binance) CoinChangeMarginType(mtr MarginTypeRequest) error {
	return b.Service.CoinChangeMarginType(mtr)
}

func (b *binance) CoinChangeUserPositionSide(ups ChangeUserPositionSideRequest) error {
	return b.Service.CoinChangeUserPositionSide(ups)
}

func (b *binance) CoinQueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error) {
	return b.Service.CoinQueryUserPositionSide(ups)

}

func (b *binance) CoinAdjustLeverage(alr AdjustLeverageRequest) (*CoinAdjustLeverageInfo, error) {
	return b.Service.CoinAdjustLeverage(alr)
}

func (b *binance) CoinPositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error) {
	return b.Service.CoinPositionMargin(pmr)
}

func (b *binance) CoinUserPoundage(upr UserPoundageRequest) (*UserPoundageInfo, error) {
	return b.Service.CoinUserPoundage(upr)
}

func (b *binance) CoinNewFutureOrder(nfr NewFutureOrderRequest) (*FutureProcessedOrder, error) {
	return b.Service.CoinNewFutureOrder(nfr)
}

func (b *binance) CoinCancelFutureOrder(cor CancelFutureOrderRequest) (*CanceledFutureOrder, error) {
	return b.Service.CoinCancelFutureOrder(cor)
}

func (b *binance) CoinFutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error) {
	return b.Service.CoinFutureBalance(fbr)
}

func (b *binance) CoinFutureAccount(far FutureAccountRequest) (*CoinFutureAccountInfo, error) {
	return b.Service.CoinFutureAccount(far)
}

func (b *binance) CoinUserTradesHistory(uth CoinUserTradesHistoryRequest) ([]*CoinUserTradesHistoryInfo, error) {
	return b.Service.CoinUserTradesHistory(uth)
}

func (b *binance) CoinFutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	return b.Service.CoinFutureDepthWebsocket(dwr)
}

func (b *binance) CoinFutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	return b.Service.CoinFutureTradeWebsocket(twr)
}

func (b *binance) CoinFutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	return b.Service.CoinFutureKlineWebsocket(kwr)
}

func (b *binance) CoinFutureUserDataWebsocket(udwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error) {
	return b.Service.CoinFutureUserDataWebsocket(udwr)
}

func (b *binance) CoinAccountInfoWebsocket(udwr UserDataWebsocketRequest) (chan *CoinAccountInfo, chan struct{}, error) {
	return b.Service.CoinAccountInfoWebsocket(udwr)
}

// BookTicker represents book ticker data.
type BookTicker struct {
	Symbol   string
	BidPrice float64
	BidQty   float64
	AskPrice float64
	AskQty   float64
}

// TickerAllBooks returns tickers for all books.
func (b *binance) TickerAllBooks() ([]*BookTicker, error) {
	return b.Service.TickerAllBooks()
}

type NewFutureOrderRequest struct {
	Symbol           string
	Side             OrderSide    // 买卖方向 SELL, BUY
	PositionSide     PositionSide // 持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
	Type             OrderType    // 订单类型 LIMIT, MARKE
	ReduceOnly       string       // true, false; 非双开模式下默认false；双开模式下不接受此参数； 使用closePosition不支持此参数。
	Quantity         float64
	Price            float64
	NewClientOrderID string // 用户自定义的订单号，不可以重复出现在挂单中。如空缺系统会自动赋值。
	StopPrice        float64
	ClosePosition    string
	ActivationPrice  float64
	CallbackRate     float64
	TimeInForce      TimeInForce
	WorkingType      string
	PriceProject     string
	NewOrderRespType string
	RecvWindow       time.Duration
	Timestamp        time.Time
}

type FutureProcessedOrder struct {
	Symbol        string    `json:"symbol"`
	CumQuote      float64   `json:"cumQuote"`           // 成交金额
	ExecutedQty   float64   `json:"executedQty,string"` // 成交量
	ClientOrderId string    `json:"clientOrderId"`      // 用户自定义订单号
	OrderId       int64     `json:"orderId"`            // 系统订单号
	AvgPrice      float64   `json:"avgPrice,string"`    // 平均成交价
	OrigQty       float64   `json:"origQty,string"`     // 原始委托数量
	Price         float64   `json:"price,string"`       // 委托价格
	Side          string    `json:"side"`               // 买卖方向
	PositionSide  string    `json:"positionSide"`       // 持仓方向
	Status        string    `json:"status"`             // 订单状态
	StopPrice     float64   `json:"stopPrice"`          // 触发价
	ClosePosition bool      `json:"closePosition"`      // 是否条件全平仓
	TimeInForce   string    `json:"timeInForce"`        // 有效方法
	Type          string    `json:"type"`               // 订单类型
	OrigType      string    `json:"origType"`           // 触发前订单类型
	ActivatePrice float64   `json:"activatePrice"`      // 跟踪止损激活价格， 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     float64   `json:"priceRate"`          // 跟踪止损回调比例， 仅`TRAILING_STOP_MARKET` 订单返回此字段
	WorkingType   string    `json:"workingType"`        // 条件价格触发类型
	PriceProtect  bool      `json:"priceProtect"`       // 是否开启条件单触发保护
	Time          int64     `json:"time"`
	UpdateTime    time.Time `json:"updateTime"`
}

func (b *binance) NewFutureOrder(nor NewFutureOrderRequest) (*FutureProcessedOrder, error) {
	return b.Service.NewFutureOrder(nor)
}

// NewOrderRequest represents NewOrder request data.
type NewOrderRequest struct {
	Symbol           string
	Side             OrderSide
	Type             OrderType
	TimeInForce      TimeInForce
	Quantity         float64
	Price            float64
	NewClientOrderID string
	StopPrice        float64
	IcebergQty       float64
	Precision        int32 //精度
	Timestamp        time.Time
}

// ProcessedOrder represents data from processed order.
type ProcessedOrder struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string
	TransactTime  time.Time
}

// NewOrder places new order and returns ProcessedOrder.
func (b *binance) NewOrder(nor NewOrderRequest) (*ProcessedOrder, error) {
	return b.Service.NewOrder(nor)
}

// NewOrder places testing order.
func (b *binance) NewOrderTest(nor NewOrderRequest) error {
	return b.Service.NewOrderTest(nor)
}

// QueryOrderRequest represents QueryOrder request data.
type QueryOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	RecvWindow        time.Duration
	Timestamp         time.Time
}

// ExecutedOrder represents data about executed order.
type ExecutedOrder struct {
	Symbol        string
	OrderID       int
	ClientOrderID string
	Price         float64
	OrigQty       float64
	ExecutedQty   float64
	Status        OrderStatus
	TimeInForce   TimeInForce
	Type          OrderType
	Side          OrderSide
	StopPrice     float64
	IcebergQty    float64
	Time          time.Time
}

func (b *binance) QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error) {
	return b.Service.QueryOrder(qor)
}

type QueryFutureOrderRequest struct {
	Symbol            string
	RecvWindow        time.Duration
	Timestamp         time.Time
	OrderId           string
	OrigClientOrderId string
}

type ExecutedFutureOrder struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string
	CumQty        float64 //成交金额 无需关注
	Price         float64 //原始委托价格
	OrigQty       float64 //原始委托数量
	AvgPrice      string  //平均成交价
	ExecutedQty   float64 //已成交量
	Status        OrderStatus
	TimeInForce   TimeInForce
	Type          OrderType
	Side          OrderSide
	ClosePosition bool    //是否条件平仓
	StopPrice     float64 //条件订单触发价格
	ReduceOnly    bool    //是否只减仓
	OrigType      string  //触发前订单类型
	PositionSide  string
	Time          time.Time
	UpdateTime    time.Time
	ActivetePrice float64 // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     float64 // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	WorkingType   string  // 条件价格触发类型
	PriceProtect  bool    //是否开启条件触发保护
}

func (b *binance) QueryOneFutureOrder(qfo QueryFutureOrderRequest) (*ExecutedFutureOrder, error) {
	return b.Service.QueryOneFutureOrder(qfo)
}

// QueryOrder returns data about existing order.
func (b *binance) QueryAllFutureOrder(qfo QueryFutureOrderRequest) ([]*ExecutedFutureOrder, error) {
	return b.Service.QueryAllFutureOrder(qfo)
}

// 取消 合约订单
type CancelFutureOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	RecvWindow        time.Duration
	Timestamp         time.Time
}

type CanceledFutureOrder struct {
	OrderID           int64
	OrigClientOrderID string
	Symbol            string
	Price             float64
	OrigQty           float64
	ExecutedQty       float64
	Status            OrderStatus
	TimeInForce       TimeInForce
	Type              OrderType
	Side              OrderSide
	PositionSide      string    //持仓方向
	OrigType          OrderType //触发前订单类型
	ClosePosition     bool      //是否条件平仓
	StopPrice         float64
	IcebergQty        float64
	Time              time.Time
	WorkingType       string // 条件价格触发类型
	PriceProtect      bool   //是否开启条件触发保护
}

func (b *binance) CancelFutureOrder(cor CancelFutureOrderRequest) (*CanceledFutureOrder, error) {
	return b.Service.CancelFutureOrder(cor)
}

// CancelOrderRequest represents CancelOrder request data.
type CancelOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	NewClientOrderID  string

	RecvWindow time.Duration
	Timestamp  time.Time
}

// CanceledOrder represents data about canceled order.
type CanceledOrder struct {
	Symbol            string
	OrigClientOrderID string
	OrderID           int64
	ClientOrderID     string
}

// CancelOrder cancels order.
func (b *binance) CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error) {
	return b.Service.CancelOrder(cor)
}

// OpenOrdersRequest represents OpenOrders request data.
type OpenOrdersRequest struct {
	Symbol     string
	RecvWindow time.Duration
	Timestamp  time.Time
}

// OpenOrders returns list of open orders.
func (b *binance) OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error) {
	return b.Service.OpenOrders(oor)
}

// AllOrdersRequest represents AllOrders request data.
type AllOrdersRequest struct {
	Symbol     string
	OrderID    int64
	Limit      int
	RecvWindow time.Duration
	Timestamp  time.Time
}

// AllOrders returns list of all previous orders.
func (b *binance) AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error) {
	return b.Service.AllOrders(aor)
}

type AllFutureOrdersRequest struct {
	Symbol     string
	OrderID    int64
	Limit      int
	StartTime  int64
	EndTime    int64
	RecvWindow time.Duration
	Timestamp  time.Time
}

func (b *binance) QueryAllHistoryFutureOrders(afo AllFutureOrdersRequest) ([]*ExecutedFutureOrder, error) {
	return b.Service.QueryAllHistoryFutureOrders(afo)
}

// AccountRequest represents Account request data.
type AccountRequest struct {
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Account represents user's account information.
type Account struct {
	MakerCommision  int64
	TakerCommision  int64
	BuyerCommision  int64
	SellerCommision int64
	CanTrade        bool
	CanWithdraw     bool
	CanDeposit      bool
	Balances        []*Balance
}

type AccountEvent struct {
	WSEvent
	Account
}
type CoinTopLongShortPositionRatioRequest struct {
	Pair      string
	Period    string
	Limit     int
	StartTime int64
	EndTime   int64
}

type CoinTopLongShortPositionRatioInfo struct {
	Pair           string
	LongShortRatio float64
	LongAccount    float64
	ShortAccount   float64
	Timestamp      time.Time
}

type CoinContractPositionRequest struct {
	Pair         string
	ContractType string
	Period       string
	Limit        int
	StartTime    int64
	EndTime      int64
}

type CoinContractPositionInfo struct {
	Pair                 string
	ContractType         string
	SumOpenInterest      float64
	SumOpenInterestValue float64
	Timestamp            time.Time
}

type CoinTakerlongshortRatioRequest struct {
	Pair         string
	ContractType string
	Period       string
	Limit        int
	StartTime    int64
	EndTime      int64
}

type CoinTakerlongshortRatioInfo struct {
	Pair              string
	ContractType      string
	TakerBuyVol       float64
	TakerSellVol      float64
	TakerBuyVolValue  float64
	TakerSellVolValue float64
	Timestamp         time.Time
}

type CoinBestBookTickerInfo struct {
	Symbol   string
	Pair     string
	BidPrice float64
	BidQty   float64
	AskPrice float64
	AskQty   float64
	Time     time.Time
}

type CoinQueryFutureOrderRequest struct {
	Symbol     string
	RecvWindow time.Duration
	Pair       string
	Timestamp  time.Time
}

type CoinExecutedFutureOrder struct {
	Symbol        string      //1
	Pair          string      //2
	OrderID       int64       //1
	CumBase       float64     //2
	ClientOrderID string      //1
	AvgPrice      string      // 平均成交价 1
	Price         float64     // 成交金额 1
	CumQty        float64     //成交金额
	OrigQty       float64     //原始委托数量 1
	ExecutedQty   float64     // 成交量 1
	Status        OrderStatus //1
	TimeInForce   TimeInForce //1
	Type          OrderType   //1
	Side          OrderSide   //1
	ClosePosition bool        //是否条件平仓 1
	StopPrice     float64     //1
	ReduceOnly    bool        //是否只减仓 1
	OrigType      string      //触发前订单类型 1
	PositionSide  string      //1
	Time          time.Time   //1
	UpdateTime    time.Time   //1
	ActivetePrice float64     // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段 1
	PriceRate     float64     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段 1
	WorkingType   string      // 条件价格触发类型 1
	PriceProtect  bool        //是否开启条件触发保护

}

type CoinAllFutureOrdersRequest struct {
	Symbol     string
	Pair       string
	OrderID    int64
	Limit      int
	StartTime  int64
	EndTime    int64
	RecvWindow time.Duration
	Timestamp  time.Time
}

type CoinHistoryExecutedFutureOrder struct {
	Symbol        string  //1
	Pair          string  //2
	OrderID       int64   //1
	CumBase       float64 //2
	ClientOrderID string  //1
	AvgPrice      string  // 平均成交价 1
	Price         float64 // 成交金额 1
	//CumQty        float64 //成交金额
	OrigQty       float64     //原始委托数量 1
	ExecutedQty   float64     // 成交量 1
	Status        OrderStatus //1
	TimeInForce   TimeInForce //1
	Type          OrderType   //1
	Side          OrderSide   //1
	ClosePosition bool        //是否条件平仓 1
	StopPrice     float64     //1
	ReduceOnly    bool        //是否只减仓 1
	OrigType      string      //触发前订单类型 1
	PositionSide  string      //1
	Time          time.Time   //1
	UpdateTime    time.Time   //1
	ActivetePrice float64     // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段 1
	PriceRate     float64     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段 1
	WorkingType   string      // 条件价格触发类型 1
	//PriceProtect  bool    //是否开启条件触发保护

}

type CoinOpenInterestNumsInfo struct {
	OpenInterest float64
	Symbol       string
	Pair         string
	ContractType string
	Time         time.Time
}

type CoinGlobalLongShortAccountRatioRequest struct {
	Pair      string
	Period    string
	Limit     int
	StartTime int64
	EndTime   int64
}

type CoinGlobalLongShortAccountRatioInfo struct {
	Pair           string
	LongShortRatio float64
	LongAccount    float64
	ShortAccount   float64
	Timestamp      time.Time
}

type CoinAdjustLeverageInfo struct {
	Leverage int    // 杠杆倍数S
	maxQty   int    // 当前杠杆倍数下允许的最大名义价值
	Symbol   string // 交易对
}

type CoinFutureProcessedOrder struct {
	Symbol        string    `json:"symbol"`
	CumBase       float64   `json:"cumBase"`            // 成交额(标的数量)
	ExecutedQty   float64   `json:"executedQty,string"` // 成交量
	ClientOrderId string    `json:"clientOrderId"`      // 用户自定义订单号
	OrderId       int64     `json:"orderId"`            // 系统订单号
	AvgPrice      float64   `json:"avgPrice,string"`    // 平均成交价
	OrigQty       float64   `json:"origQty,string"`     // 原始委托数量
	Price         float64   `json:"price,string"`       // 委托价格
	Side          string    `json:"side"`               // 买卖方向
	PositionSide  string    `json:"positionSide"`       // 持仓方向
	Status        string    `json:"status"`             // 订单状态
	StopPrice     float64   `json:"stopPrice"`          // 触发价
	ClosePosition bool      `json:"closePosition"`      // 是否条件全平仓
	TimeInForce   string    `json:"timeInForce"`        // 有效方法
	Type          string    `json:"type"`               // 订单类型
	OrigType      string    `json:"origType"`           // 触发前订单类型
	ActivatePrice float64   `json:"activatePrice"`      // 跟踪止损激活价格， 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     float64   `json:"priceRate"`          // 跟踪止损回调比例， 仅`TRAILING_STOP_MARKET` 订单返回此字段
	WorkingType   string    `json:"workingType"`        // 条件价格触发类型
	PriceProtect  bool      `json:"priceProtect"`       // 是否开启条件单触发保护
	Time          int64     `json:"time"`
	UpdateTime    time.Time `json:"updateTime"`
}

type CoinCanceledFutureOrder struct {
	AvgPrice          float64
	OrderID           int64
	OrigClientOrderID string
	Symbol            string
	Price             float64
	OrigQty           float64
	ExecutedQty       float64
	Status            OrderStatus
	TimeInForce       TimeInForce
	Type              OrderType
	Side              OrderSide
	PositionSide      string    //持仓方向
	OrigType          OrderType //触发前订单类型
	ClosePosition     bool      //是否条件平仓
	StopPrice         float64
	IcebergQty        float64
	Time              time.Time
	WorkingType       string // 条件价格触发类型
	PriceProtect      bool   //是否开启条件触发保护
}

type CoinFutureBalanceInfo struct {
	AccountAlias       string  // 账户唯一识别码
	Asset              string  // 资产
	Balance            float64 // 总余额
	WithdrawAvailable  float64 //最大可提款金额
	CrossWalletBalance float64 // 全仓余额
	CrossUnPnl         float64 // 全仓持仓未实现盈亏
	AvailableBalance   float64 // 下单可用余额
	MaxWithdrawAmount  float64 // 最大可转出余额
	MarginAvailable    bool    // 是否可用作联合保证金
	UpdateTime         time.Time
}

type CoinFutureAccountInfo struct {
	FeeTier     int
	CanTrade    bool //是否可以交易
	CanDeposit  bool //是否可以入金
	CanWithdraw bool //是否可以出金
	UpdateTime  int
	Asset       []*CoinFutureAsset
	Positions   []*CoinFuturePositions
}

type CoinFutureAsset struct {
	Asset                  string
	WalletBalance          float64 // 余额
	UnrealizedProfit       float64 // 未实现盈亏
	MarginBalance          float64 // 保证金余额
	MaintMargin            float64 // 维持保证金
	InitialMargin          float64 // 当前所需起始保证金
	PositionInitialMargin  float64 // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin float64 // 当前挂单所需起始保证金(基于最新标记价格)
	CrossWalletBalance     float64 // 全仓账户余额2
	CrossUnPnl             float64 // 全仓持仓未实现盈亏3
	AvailableBalance       float64 // 可用余额4
	MaxWithdrawAmount      float64 // 最大可转出余额1

}

type CoinFuturePositions struct { // 头寸
	Symbol                 string
	InitialMargin          float64 // 当前所需起始保证金(基于最新标记价格)1
	MaintMargin            float64 // 维持保证金
	UnrealizedProfit       float64 // 持仓未实现盈亏
	PositionInitialMargin  float64 // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin float64 // 当前挂单所需起始保证金(基于最新标记价格)
	Leverage               float64 // 杠杆倍率
	Isolated               bool    // 是否是逐仓模式
	EntryPrice             float64 // 持仓成本价
	PositionSide           string  // 持仓方向
	PositionAmt            float64
	MaxQty                 float64 //当前杠杆下最大可开仓数(标的数量)
	UpdateTime             time.Time
}

type CoinUserTradesHistoryRequest struct {
	Symbol     string
	Pair       string
	RecvWindow time.Duration
	Timestamp  time.Time
	StartTime  int64
	EndTime    int64
	FromId     int
	Limit      int
}

type CoinUserTradesHistoryInfo struct {
	Buyer           bool
	Pair            string
	MarginAsset     string
	BaseQty         float64
	Commission      float64
	CommissionAsset string
	Id              int
	Maker           bool
	OrderId         int
	Price           float64
	Qty             float64

	RealizedPnl  float64
	Side         string
	PositionSide string
	Symbol       string
	Time         time.Time
}

type CoinDepthEvent struct {
	WSEvent
	EventTime time.Time
	CoinOrderBook
}

type CoinOrderBook struct {
	LastUpdateID int
	BeforeUID    int
	UpdateID     int
	Trading      string
	MessageTime  time.Time
	Bids         []*Order // 买方出价
	Asks         []*Order //卖方出价
}

type CoinFutureAccountEvent struct {
	EventName string
	OE        *CoinOrderEvent
	AE        *CoinAccEvent
}

type CoinOrderEvent struct {
	Type         string  // 事件类型
	EventTime    float64 // 事件时间
	Time         float64 // 撮合时间
	AccountAlias string  // 账户唯一识别码
	Order        struct {
		Symbol             string      // 交易对
		ClientOrderID      string      // 客户端自定订单ID
		Side               string      // 订单方向
		OrderType          string      // 订单类型
		TimeInForce        string      // 有效方式
		OrigQty            float64     // 订单原始数量
		Price              float64     // 订单原始价格
		AvgPrice           float64     // 订单平均价格
		StopPrice          float64     // 条件订单触发价格，对追踪止损单无效
		NewEvent           EventType   // 本次事件的具体执行类型
		OrderStatus        OrderStatus // 订单的当前状态
		ID                 int64       // 订单ID
		LastQty            float64     // 订单末次成交量
		ExecutedQty        float64     // 订单累计已成交量
		LastPrice          float64     // 订单末次成交价格
		MarginType         string      // 保证金资产类型
		RateAssetType      string      // 手续费资产类型
		RateQ              float64     // 手续费数量
		Time               time.Time   // 成交时间
		TimeID             string      // 成交ID
		BuyEquity          float64     // 买单净值
		SellEquity         float64     // 卖单净值
		IsTaker            bool        // 该成交是作为挂单成交吗？
		IsReduce           bool        // 是否是只减仓单
		NowType            OrderType   // 触发价类型
		OrigType           OrderType   // 原始订单类型
		PositionSide       string      // 持仓方向
		IsClose            bool        // 是否为触发平仓单
		Profit             float64     // 该交易实现盈亏
		TrackStopGoPrice   float64     // 追踪止损激活价格
		TrackStopBackPrice float64     // 追踪止损回调比例
		IsProtect          bool        //是否开启条件单触发保护
	}
}

type CoinAccEvent struct {
	Type         string  // 事件类型
	EventTime    float64 // 事件时间
	Time         float64 // 撮合时间
	AccountAlias string  // 账户唯一识别码
	Acc          struct {
		Event   string
		Balance []struct {
			Symbol        string
			WalletBalance float64
			CurBalance    float64
			BalanceChange float64
		}
	}
}

type CoinAccountInfo struct {
	Id int
	R  AccountInfoResult
}

type AccountInfoResult struct {
	Request string
	Respone CoinRespone
}

type CoinRespone struct {
	FeeTier      int
	CanTrade     bool
	CanDeposit   bool
	CanWithdraw  bool
	AccountAlias string
}

type CoinNewPriceInfo struct {
	Symbol     string
	Ps         string
	Price      float64
	UpdateTime time.Time
}

type CoinPriceChangeSituationInfo struct {
	Symbol             string
	Pair               string
	PriceChange        float64   //24小时价格变动
	PriceChangePercent float64   //24小时价格变动百分比
	WeightedAvgPrice   float64   //加权平均价
	LastPrice          float64   //最近一次成交价
	LastQty            float64   //最近一次成交额
	OpenPrice          float64   //24小时内第一次成交的价格
	HighPrice          float64   //24小时最高价
	LowPrice           float64   //24小时最低价
	Volume             float64   //24小时成交量
	BaseVolume         float64   //24小时成交额
	OpenTime           time.Time //24小时内，第一笔交易的发生时间
	CloseTime          time.Time //24小时内，最后一笔交易的发生时间
	FirstId            int
	LastId             int
	Count              int
}

// Balance groups balance-related information.
type Balance struct {
	Asset  string
	Free   float64
	Locked float64
}

// Account returns account data.
func (b *binance) Account(ar AccountRequest) (*Account, error) {
	return b.Service.Account(ar)
}

// MyTradesRequest represents MyTrades request data.
type MyTradesRequest struct {
	Symbol     string
	Limit      int
	FromID     int64
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Trade represents data about trade.
type Trade struct {
	ID              int64
	Price           float64
	Qty             float64
	Commission      float64
	CommissionAsset string
	Time            time.Time
	IsBuyer         bool
	IsMaker         bool
	IsBestMatch     bool
}

// MyTrades list user's trades.
func (b *binance) MyTrades(mtr MyTradesRequest) ([]*Trade, error) {
	return b.Service.MyTrades(mtr)
}

// WithdrawRequest represents Withdraw request data.
type WithdrawRequest struct {
	Asset      string
	Address    string
	Amount     float64
	Name       string
	RecvWindow time.Duration
	Timestamp  time.Time
}

// WithdrawResult represents Withdraw result.
type WithdrawResult struct {
	Success bool
	Msg     string
}

// Withdraw executes withdrawal.
func (b *binance) Withdraw(wr WithdrawRequest) (*WithdrawResult, error) {
	return b.Service.Withdraw(wr)
}

// HistoryRequest represents history-related calls request data.
type HistoryRequest struct {
	Asset      string
	Status     *int
	StartTime  time.Time
	EndTime    time.Time
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Deposit represents Deposit data.
type Deposit struct {
	InsertTime time.Time
	Amount     float64
	Asset      string
	Status     int
}

// DepositHistory lists deposit data.
func (b *binance) DepositHistory(hr HistoryRequest) ([]*Deposit, error) {
	return b.Service.DepositHistory(hr)
}

// Withdrawal represents withdrawal data.
type Withdrawal struct {
	Amount    float64
	Address   string
	TxID      string
	Asset     string
	ApplyTime time.Time
	Status    int
}

// WithdrawHistory lists withdraw data.
func (b *binance) WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error) {
	return b.Service.WithdrawHistory(hr)
}

// Stream represents stream information.
//
// Read web docs to get more information about using streams.
type Stream struct {
	ListenKey string
}

// StartUserDataStream starts stream and returns Stream with ListenKey.
func (b *binance) StartUserDataStream() (*Stream, error) {
	return b.Service.StartUserDataStream()
}

// KeepAliveUserDataStream prolongs stream livespan.
func (b *binance) KeepAliveUserDataStream(s *Stream) error {
	return b.Service.KeepAliveUserDataStream(s)
}

// CloseUserDataStream closes opened stream.
func (b *binance) CloseUserDataStream(s *Stream) error {
	return b.Service.CloseUserDataStream(s)
}

type WSEvent struct {
	Type   string
	Time   time.Time
	Symbol string
}

type DepthWebsocketRequest struct {
	Symbol string
}

func (b *binance) DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	return b.Service.DepthWebsocket(dwr)
}

type KlineWebsocketRequest struct {
	Symbol   string
	Interval Interval
}

func (b *binance) KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	return b.Service.KlineWebsocket(kwr)
}

type TradeWebsocketRequest struct {
	Symbol string
}

func (b *binance) TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	return b.Service.TradeWebsocket(twr)
}

type UserDataWebsocketRequest struct {
	ListenKey string
}

func (b *binance) UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error) {
	return b.Service.UserDataWebsocket(udwr)
}

type FutureBalanceRequest struct {
	RecvWindow time.Duration
	Timestamp  time.Time
}

type FutureBalanceInfo struct {
	AccountAlias       string  // 账户唯一识别码
	Asset              string  // 资产
	Balance            float64 // 总余额
	CrossWalletBalance float64 // 全仓余额
	CrossUnPnl         float64 // 全仓持仓未实现盈亏
	AvailableBalance   float64 // 下单可用余额
	MaxWithdrawAmount  float64 // 最大可转出余额
	MarginAvailable    bool    // 是否可用作联合保证金
	UpdateTime         time.Time
}

func (b *binance) FutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error) {
	return b.Service.FutureBalance(fbr)
}

type FutureAccountRequest struct {
	RecvWindow time.Duration
	Timestamp  time.Time
}

type FutureAccountInfo struct {
	FeeTier                     int
	CanTrade                    bool //是否可以交易
	CanDeposit                  bool //是否可以入金
	CanWithdraw                 bool //是否可以出金
	UpdateTime                  int
	TotalInitialMargin          float64 // 但前所需起始保证金总额(存在逐仓请忽略), 仅计算usdt资产
	TotalMaintMargin            float64 // 维持保证金总额, 仅计算usdt资产
	TotalWalletBalance          float64 // 账户总余额, 仅计算usdt资产
	TotalUnrealizedProfit       float64 // 持仓未实现盈亏总额, 仅计算usdt资产
	TotalMarginBalance          float64 // 保证金总余额, 仅计算usdt资产
	TotalPositionInitialMargin  float64 // 持仓所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalOpenOrderInitialMargin float64 // 当前挂单所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalCrossWalletBalance     float64 // 全仓账户余额, 仅计算usdt资产
	TotalCrossUnPnl             float64 // 全仓持仓未实现盈亏总额, 仅计算usdt资产
	AvailableBalance            float64 // 可用余额, 仅计算usdt资产
	MaxWithdrawAmount           float64 // 最大可转出余额, 仅计算usdt资产
	Asset                       []*FutureAsset
	Positions                   []*FuturePositions
}

type FutureAsset struct {
	Asset                  string
	WalletBalance          float64 // 余额
	UnrealizedProfit       float64 // 未实现盈亏
	MarginBalance          float64 // 保证金余额
	MaintMargin            float64 // 维持保证金
	InitialMargin          float64 // 当前所需起始保证金
	PositionInitialMargin  float64 // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin float64 // 当前挂单所需起始保证金(基于最新标记价格)
	CrossWalletBalance     float64 // 全仓账户余额
	CrossUnPnl             float64 // 全仓持仓未实现盈亏
	AvailableBalance       float64 // 可用余额
	MaxWithdrawAmount      float64 // 最大可转出余额
	MarginAvailable        float64 // 是否可用作联合保证金
	UpdateTime             time.Time
}

type FuturePositions struct { // 头寸
	Symbol                 string
	InitialMargin          float64 // 当前所需起始保证金(基于最新标记价格)
	MaintMargin            float64 // 维持保证金
	UnrealizedProfit       float64 // 持仓未实现盈亏
	PositionInitialMargin  float64 // 持仓所需起始保证金(基于最新标记价格)
	OpenOrderInitialMargin float64 // 当前挂单所需起始保证金(基于最新标记价格)
	Leverage               float64 // 杠杆倍率
	Isolated               bool    // 是否是逐仓模式
	EntryPrice             float64 // 持仓成本价
	MaxNotional            float64 // 当前杠杆下用户可用的最大名义价值
	PositionSide           string  // 持仓方向
	PositionAmt            float64 // 持仓数量
	UpdateTime             time.Time
}

func (b *binance) FutureAccount(far FutureAccountRequest) (*FutureAccountInfo, error) {
	return b.Service.FutureAccount(far)
}

type UserPoundageRequest struct {
	Symbol     string
	RecvWindow time.Duration
	Timestamp  time.Time
}

type UserPoundageInfo struct {
	Symbol              string
	MakerCommissionRate float64
	TakerCommissionRate float64
}

func (b *binance) UserPoundage(udr UserPoundageRequest) (*UserPoundageInfo, error) {
	return b.Service.UserPoundage(udr)
}

type AdjustLeverageRequest struct {
	Symbol     string
	Leverage   int
	RecvWindow time.Duration
	Timestamp  time.Time
}

type AdjustLeverageInfo struct {
	Leverage         int    // 杠杆倍数S
	MaxNotionalValue int    // 当前杠杆倍数下允许的最大名义价值
	Symbol           string // 交易对
}

func (b *binance) AdjustLeverage(alr AdjustLeverageRequest) (*AdjustLeverageInfo, error) {
	return b.Service.AdjustLeverage(alr)
}

type PositionMarginRequest struct {
	Symbol       string
	PositionSide string
	RecvWindow   time.Duration
	Timestamp    time.Time
	Amount       float64
	Type         int
}

type PositionMarginInfo struct {
	Amount float64
	Code   int
	Msg    string
	Type   int
}

func (b *binance) PositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error) {
	return b.Service.PositionMargin(pmr)
}

type UserTradesHistoryRequest struct {
	Symbol     string
	RecvWindow time.Duration
	Timestamp  time.Time
	StartTime  int64
	EndTime    int64
	FromId     int
	Limit      int
}

type UserTradesHistoryInfo struct {
	Buyer           bool
	Commission      float64
	CommissionAsset string
	Id              int
	Maker           bool
	OrderId         int
	Price           float64
	Qty             float64
	QuoteQty        float64
	RealizedPnl     float64
	Side            string
	PositionSide    string
	Symbol          string
	Time            time.Time
}

func (b *binance) UserTradesHistory(uth UserTradesHistoryRequest) ([]*UserTradesHistoryInfo, error) {
	return b.Service.UserTradesHistory(uth)
}

type PremiumAndFundsRateRequest struct {
	Symbol string
}

type PremiumAndFundsRateInfo struct {
	Symbol               string
	MarkPrice            float64
	IndexPrice           float64
	EstimatedSettlePrice float64
	LastFundingRate      float64
	NextFundingTime      time.Time
	InterestRate         float64
	Time                 time.Time
}

func (b *binance) PremiumAndFundsRate(pfrr PremiumAndFundsRateRequest) (*PremiumAndFundsRateInfo, error) {
	return b.Service.PremiumAndFundsRate(pfrr)
}

type PriceChangeSituationRequest struct {
	Symbol string
}

type PriceChangeSituationInfo struct {
	Symbol             string
	PriceChange        float64   //24小时价格变动
	PriceChangePercent float64   //24小时价格变动百分比
	WeightedAvgPrice   float64   //加权平均价
	LastPrice          float64   //最近一次成交价
	LastQty            float64   //最近一次成交额
	OpenPrice          float64   //24小时内第一次成交的价格
	HighPrice          float64   //24小时最高价
	LowPrice           float64   //24小时最低价
	Volume             float64   //24小时成交量
	QuoteVolume        float64   //24小时成交额
	OpenTime           time.Time //24小时内，第一笔交易的发生时间
	CloseTime          time.Time //24小时内，最后一笔交易的发生时间
	FirstId            int
	LastId             int
	Count              int
}

func (b *binance) PriceChangeSituation(pcsr PriceChangeSituationRequest) (*PriceChangeSituationInfo, error) {
	return b.Service.PriceChangeSituation(pcsr)
}

type OpenInterestNumsRequest struct {
	Symbol string
}

type OpenInterestNumsInfo struct {
	OpenInterest float64
	Symbol       string
	Time         time.Time
}

func (b *binance) OpenInterestNums(oir OpenInterestNumsRequest) (*OpenInterestNumsInfo, error) {
	return b.Service.OpenInterestNums(oir)
}

type BestBookTickerRequest struct {
	Symbol string
}

type BestBookTickerInfo struct {
	Symbol   string
	BidPrice float64
	BidQty   float64
	AskPrice float64
	AskQty   float64
	Time     time.Time
}

func (b *binance) BestBookTicker(bbtr BestBookTickerRequest) (*BestBookTickerInfo, error) {
	return b.Service.BestBookTicker(bbtr)
}

type ContractPositionRequest struct {
	Symbol    string
	Period    string
	Limit     int
	StartTime int64
	EndTime   int64
}

type ContractPositionInfo struct {
	Symbol               string
	SumOpenInterest      float64
	SumOpenInterestValue float64
	Timestamp            time.Time
}

func (b *binance) ContractPosition(bpr ContractPositionRequest) ([]*ContractPositionInfo, error) {
	return b.Service.ContractPosition(bpr)
}

type TopLongShortPositionRatioRequest struct {
	Symbol    string
	Period    string
	Limit     int
	StartTime int64
	EndTime   int64
}

type TopLongShortPositionRatioInfo struct {
	Symbol         string
	LongShortRatio float64
	LongAccount    float64
	ShortAccount   float64
	Timestamp      time.Time
}

func (b *binance) TopLongShortPositionRatio(tspr TopLongShortPositionRatioRequest) ([]*TopLongShortPositionRatioInfo, error) {
	return b.Service.TopLongShortPositionRatio(tspr)
}

type GlobalLongShortAccountRatioRequest struct {
	Symbol    string
	Period    string
	Limit     int
	StartTime int64
	EndTime   int64
}

type GlobalLongShortAccountRatioInfo struct {
	Symbol         string
	LongShortRatio float64
	LongAccount    float64
	ShortAccount   float64
	Timestamp      time.Time
}

func (b *binance) GlobalLongShortAccountRatio(glsarr GlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error) {
	return b.Service.GlobalLongShortAccountRatio(glsarr)
}

type TakerlongshortRatioRequest struct {
	Symbol    string
	Period    string
	Limit     int
	StartTime int64
	EndTime   int64
}

type TakerlongshortRatioInfo struct {
	BuySellRatio float64
	BuyVol       float64
	SellVol      float64
	Timestamp    time.Time
}

func (b *binance) TakerlongshortRatio(tlr TakerlongshortRatioRequest) ([]*TakerlongshortRatioInfo, error) {
	return b.Service.TakerlongshortRatio(tlr)
}

type UserPositionSideRequest struct {
	RecvWindow time.Duration
	Timestamp  time.Time
}

type UserPositionSideInfo struct {
	DualSidePosition bool
}

func (b *binance) QueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error) {
	return b.Service.QueryUserPositionSide(ups)
}

func (b *binance) StartFutureUserDataStream() (*Stream, error) {
	return b.Service.StartFutureUserDataStream()
}

// KeepAliveUserDataStream prolongs stream livespan.
func (b *binance) KeepAliveFutureUserDataStream(s *Stream) error {
	return b.Service.KeepAliveFutureUserDataStream(s)
}

// CloseUserDataStream closes opened stream.
func (b *binance) CloseFutureUserDataStream(s *Stream) error {
	return b.Service.CloseFutureUserDataStream(s)
}

func (b *binance) FutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	return b.Service.FutureDepthWebsocket(dwr)
}

func (b *binance) FutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	return b.Service.FutureKlineWebsocket(kwr)
}
func (b *binance) FutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	return b.Service.FutureTradeWebsocket(twr)
}
func (b *binance) FutureUserDataWebsocket(udwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error) {
	return b.Service.FutureUserDataWebsocket(udwr)
}

type FutureAccountEvent struct {
	EventName string
	OE        *OrderEvent
	AE        *AccEvent
}

type AccEvent struct {
	Type      string  // 事件类型
	EventTime float64 // 事件时间
	Time      float64 // 撮合时间
	Acc       struct {
		Event   string
		Balance []struct {
			Symbol        string
			WalletBalance float64
			CurBalance    float64
			BalanceChange float64
		}
		Property []struct {
			Symbol string  // 交易对
			Pa     float64 // 仓位
			EP     float64 // 入仓价格
			CR     float64 // (费前)累计实现损益
			UP     float64 // 持仓未实现盈亏
			MT     string  // 保证金模式
			IW     float64 // 若为逐仓，仓位保证金
			PS     string  // 持仓方向
		}
	}
}

type OrderEvent struct {
	Type      string  // 事件类型
	EventTime float64 // 事件时间
	Time      float64 // 撮合时间
	Order     struct {
		Symbol        string      // 交易对
		ClientOrderID string      // 客户端自定订单ID
		Side          string      // 订单方向
		OrderType     string      // 订单类型
		TimeInForce   string      // 有效方式
		OrigQty       float64     // 订单原始数量
		Price         float64     // 订单原始价格
		AvgPrice      float64     // 订单平均价格
		StopPrice     float64     // 条件订单触发价格，对追踪止损单无效
		NewEvent      EventType   // 本次事件的具体执行类型
		OrderStatus   OrderStatus // 订单的当前状态
		ID            int64       // 订单ID
		LastQty       float64     // 订单末次成交量
		ExecutedQty   float64     // 订单累计已成交量
		LastPrice     float64     // 订单末次成交价格
		RateAssetType string      // 手续费资产类型
		RateQ         float64     // 手续费数量
		Time          time.Time   // 成交时间
		TimeID        string      // 成交ID
		BuyEquity     float64     // 买单净值
		SellEquity    float64     // 卖单净值
		IsTaker       bool        // 该成交是作为挂单成交吗？
		IsReduce      bool        // 是否是只减仓单
		NowType       OrderType   // 触发价类型
		OrigType      OrderType   // 原始订单类型
		PositionSide  string      // 持仓方向
		IsClose       bool        // 是否为触发平仓单; 仅在条件订单情况下会推送此字段
		Profit        float64     // 该交易实现盈亏
	}
}

type NewPriceRequest struct {
	Symbol string
}

type NewPriceInfo struct {
	Symbol     string
	Price      float64
	UpdateTime time.Time
}

func (b *binance) GetNewPrice(npr NewPriceRequest) (*NewPriceInfo, error) {
	return b.Service.GetNewPrice(npr)
}

func (b *binance) FutureKlines(kr KlinesRequest) ([]*Kline, error) {
	return b.Service.FutureKlines(kr)
}

func (b *binance) ChangeMarginType(mtr MarginTypeRequest) error {
	return b.Service.ChangeMarginType(mtr)
}

type MarginTypeRequest struct {
	Symbol     string
	MarginType PositionStatus
	RecvWindow time.Duration
	Timestamp  time.Time
}

type ChangeUserPositionSideRequest struct {
	DualSidePosition PosithonSideStatus
	RecvWindow       time.Duration
	Timestamp        time.Time
}

func (b *binance) ChangeUserPositionSide(ups ChangeUserPositionSideRequest) error {
	return b.Service.ChangeUserPositionSide(ups)
}

func (b *binance) StartCoinFutureUserDataStream() (*Stream, error) {
	return b.Service.StartCoinFutureUserDataStream()
}

// KeepAliveUserDataStream prolongs stream livespan.
func (b *binance) KeepAliveCoinFutureUserDataStream(s *Stream) error {
	return b.Service.KeepAliveCoinFutureUserDataStream(s)
}

// CloseUserDataStream closes opened stream.
func (b *binance) CloseCoinFutureUserDataStream(s *Stream) error {
	return b.Service.CloseCoinFutureUserDataStream(s)
}
