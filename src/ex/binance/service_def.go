package binance

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Service represents service layer for Binance API.
//
// The main purpose for this layer is to be replaced with dummy implementation
// if necessary without need to replace Binance instance.
type Service interface {
	Ping() error
	Time() (time.Time, error)
	ExchangeInfo() (string, error)

	NewPrice(nb OrderNewPriceRequest) (*NewPrice, error)
	OrderBook(obr OrderBookRequest) (*OrderBook, error)
	AggTrades(atr AggTradesRequest) ([]*AggTrade, error)
	Klines(kr KlinesRequest) ([]*Kline, error)
	Ticker24(tr TickerRequest) (*Ticker24, error)
	TickerAllPrices() ([]*PriceTicker, error)
	TickerAllBooks() ([]*BookTicker, error)

	NewOrder(or NewOrderRequest) (*ProcessedOrder, error)
	NewOrderTest(or NewOrderRequest) error
	QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error)
	CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error)
	OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error)
	AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error)

	Account(ar AccountRequest) (*Account, error)
	MyTrades(mtr MyTradesRequest) ([]*Trade, error)
	Withdraw(wr WithdrawRequest) (*WithdrawResult, error)
	DepositHistory(hr HistoryRequest) ([]*Deposit, error)
	WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error)

	StartUserDataStream() (*Stream, error)
	KeepAliveUserDataStream(s *Stream) error
	CloseUserDataStream(s *Stream) error

	NewFutureOrder(nfr NewFutureOrderRequest) (*FutureProcessedOrder, error)

	CancelFutureOrder(cfr CancelFutureOrderRequest) (*CanceledFutureOrder, error)

	QueryOneFutureOrder(qfo QueryFutureOrderRequest) (*ExecutedFutureOrder, error)

	QueryAllFutureOrder(qfo QueryFutureOrderRequest) ([]*ExecutedFutureOrder, error)

	QueryAllHistoryFutureOrders(afo AllFutureOrdersRequest) ([]*ExecutedFutureOrder, error)

	FutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error)

	FutureAccount(far FutureAccountRequest) (*FutureAccountInfo, error)

	UserPoundage(udr UserPoundageRequest) (*UserPoundageInfo, error)

	AdjustLeverage(alr AdjustLeverageRequest) (*AdjustLeverageInfo, error)

	PositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error)

	GetNewPrice(npr NewPriceRequest) (*NewPriceInfo, error)

	UserTradesHistory(uth UserTradesHistoryRequest) ([]*UserTradesHistoryInfo, error)

	PremiumAndFundsRate(pfrr PremiumAndFundsRateRequest) (*PremiumAndFundsRateInfo, error)

	PriceChangeSituation(pcsr PriceChangeSituationRequest) (*PriceChangeSituationInfo, error)

	OpenInterestNums(oir OpenInterestNumsRequest) (*OpenInterestNumsInfo, error)

	BestBookTicker(bbtr BestBookTickerRequest) (*BestBookTickerInfo, error)

	ContractPosition(bpr ContractPositionRequest) ([]*ContractPositionInfo, error)

	TopLongShortPositionRatio(tspr TopLongShortPositionRatioRequest) ([]*TopLongShortPositionRatioInfo, error)

	GlobalLongShortAccountRatio(glsarr GlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error)

	TakerlongshortRatio(tlr TakerlongshortRatioRequest) ([]*TakerlongshortRatioInfo, error)

	QueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error)

	FutureKlines(kr KlinesRequest) ([]*Kline, error)

	ChangeMarginType(mtr MarginTypeRequest) error

	ChangeUserPositionSide(ups ChangeUserPositionSideRequest) error

	DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)
	KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)
	TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)
	UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error)

	StartFutureUserDataStream() (*Stream, error)
	KeepAliveFutureUserDataStream(s *Stream) error
	CloseFutureUserDataStream(s *Stream) error

	FutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)
	FutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)
	FutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)
	FutureUserDataWebsocket(udwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error)

	//币本位合约 如下
	OrderCoinBook(obr OrderBookRequest) (*OrderBook, error)

	FutureCoinKlines(kr KlinesRequest) ([]*Kline, error)

	CoinTopLongShortPositionRatio(tspr CoinTopLongShortPositionRatioRequest) ([]*CoinTopLongShortPositionRatioInfo, error)

	CoinContractPosition(bpr CoinContractPositionRequest) ([]*CoinContractPositionInfo, error)

	CoinGetNewPrice(npr NewPriceRequest) ([]*CoinNewPriceInfo, error)

	CoinPriceChangeSituation(pcsr PriceChangeSituationRequest) ([]*CoinPriceChangeSituationInfo, error)

	CoinTakerlongshortRatio(tlr CoinTakerlongshortRatioRequest) ([]*CoinTakerlongshortRatioInfo, error)

	CoinBestBookTicker(bbtr BestBookTickerRequest) ([]*CoinBestBookTickerInfo, error)

	CoinOpenInterestNums(oir OpenInterestNumsRequest) (*CoinOpenInterestNumsInfo, error)

	CoinGlobalLongShortAccountRatio(glsarr CoinGlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error)

	CoinChangeMarginType(mtr MarginTypeRequest) error

	CoinChangeUserPositionSide(ups ChangeUserPositionSideRequest) error

	CoinQueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error)

	CoinUserPoundage(upr UserPoundageRequest) (*UserPoundageInfo, error)

	CoinFutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error)

	CoinFutureAccount(far FutureAccountRequest) (*CoinFutureAccountInfo, error)

	CoinAllFutureOrders(afo CoinAllFutureOrdersRequest) ([]*CoinHistoryExecutedFutureOrder, error)

	//下单
	CoinNewFutureOrder(nfr NewFutureOrderRequest) (*FutureProcessedOrder, error)

	//撤销订单
	CoinCancelFutureOrder(cor CancelFutureOrderRequest) (*CanceledFutureOrder, error)

	// 查询当前所有挂单
	CoinQueryFutureOrder(qfo CoinQueryFutureOrderRequest) ([]*ExecutedFutureOrder, error)

	//websocket 币本位
	CoinFutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)

	//websocket 归集交易
	CoinFutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)

	CoinUserTradesHistory(uth CoinUserTradesHistoryRequest) ([]*CoinUserTradesHistoryInfo, error)

	CoinFutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)

	CoinFutureUserDataWebsocket(udwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error)

	CoinAccountInfoWebsocket(udwr UserDataWebsocketRequest) (chan *CoinAccountInfo, chan struct{}, error)

	CoinAdjustLeverage(alr AdjustLeverageRequest) (*CoinAdjustLeverageInfo, error)

	CoinPositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error)

	StartCoinFutureUserDataStream() (*Stream, error)
	KeepAliveCoinFutureUserDataStream(s *Stream) error
	CloseCoinFutureUserDataStream(s *Stream) error
}

type apiService struct {
	URL    string
	APIKey string
	Signer Signer
	Ctx    context.Context
}

// NewAPIService creates instance of Service.
//
// If logger or ctx are not provided, NopLogger and Background context are used as default.
// You can use context for one-time request cancel (e.g. when shutting down the app).
func NewAPIService(url, apiKey string, signer Signer, ctx context.Context) Service {

	if ctx == nil {
		ctx = context.Background()
	}
	return &apiService{
		URL:    url,
		APIKey: apiKey,
		Signer: signer,
		Ctx:    ctx,
	}
}

func (as *apiService) request(method string, endpoint string, params map[string]string,
	apiKey bool, sign bool) (*http.Response, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 5 * time.Second,
		}).DialContext,
		MaxIdleConns:        30,               //最大空闲连接数
		MaxIdleConnsPerHost: 60,               //最大与服务器的连接数  默认是2
		IdleConnTimeout:     30 * time.Second, //空闲连接保持时间
		Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:7890")
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
	}
	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("%s/%s", as.URL, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}
	req.WithContext(as.Ctx)

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	if apiKey {
		req.Header.Add("X-MBX-APIKEY", as.APIKey)
	}
	if sign {
		//level.Debug(as.Logger).Log("queryString", q.Encode())
		q.Add("signature", as.Signer.Sign([]byte(q.Encode())))
		//level.Debug(as.Logger).Log("signature", as.Signer.Sign([]byte(q.Encode())))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client do failed")
	}
	return resp, nil
}

func ReConnectWebSocket(url string) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic("connect websocket failed")
	}
	return c
}

func init() {
	websocket.DefaultDialer.Proxy = func(*http.Request) (*url.URL, error) {
		return url.Parse("http://127.0.0.1:7890")
	}
}
