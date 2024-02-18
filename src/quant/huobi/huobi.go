package huobi

import (
	"time"
	"tinyquant/src/mod"

	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/logging/applogger"
	"github.com/huobirdcenter/huobi_golang/pkg/client"
	"github.com/huobirdcenter/huobi_golang/pkg/client/marketwebsocketclient"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"github.com/huobirdcenter/huobi_golang/pkg/model/order"
)

var acc_cl *client.AccountClient
var order_cl *client.OrderClient
var market_cl *client.MarketClient

func init() {
	//acc_cl = new(client.AccountClient).Init(util.HUOBI_ACCESS_KEY, util.HUOBI_SECRET_KEY, util.HUOBI_HOST)
	//order_cl = new(client.OrderClient).Init(util.HUOBI_ACCESS_KEY, util.HUOBI_SECRET_KEY, util.HUOBI_HOST)
	//market_cl = new(client.MarketClient).Init(util.HUOBI_HOST)
}

/*
account: {Id:12770533 Type:spot Subtype: State:working}   // 现货账户
account: {Id:16464709 Type:otc Subtype: State:working}    // OTC账户
account: {Id:17607343 Type:margin Subtype:xrpusdt State:working}  //逐仓杠杆账户
*/

func GetAccountInfo() {

	resp, err := acc_cl.GetAccountInfo()
	if err != nil {
		applogger.Error("Get account error: %s", err)
	} else {
		applogger.Info("Get account, count=%d", len(resp))
		for _, result := range resp {
			applogger.Info("account: %+v", result)
		}
	}
}

func GetAccountBalance() {

	resp, err := acc_cl.GetAccountBalance(config.AccountId)
	if err != nil {
		applogger.Error("Get account balance error: %s", err)
	} else {
		applogger.Info("Get account balance: id=%d, type=%s, state=%s, count=%d",
			resp.Id, resp.Type, resp.State, len(resp.List))
		if resp.List != nil {
			for _, result := range resp.List {
				applogger.Info("Account balance: %+v", result)
			}
		}
	}
}

func PlaceOrder() {

	request := order.PlaceOrderRequest{
		AccountId: config.AccountId,
		Type:      "buy-limit",
		Source:    "spot-api",
		Symbol:    "btcusdt",
		Price:     "1.1",
		Amount:    "1",
	}
	resp, err := order_cl.PlaceOrder(&request)
	if err != nil {
		applogger.Error(err.Error())
	} else {
		switch resp.Status {
		case "ok":
			applogger.Info("Place order successfully, order id: %s", resp.Data)
		case "error":
			applogger.Error("Place order error: %s", resp.ErrorMessage)
		}
	}
}

// 获取 所有 交易品种的 ticker
func GetAllSymbolTickers() {
	resp, err := market_cl.GetAllSymbolsLast24hCandlesticksAskBid()
	if err != nil {
		applogger.Error(err.Error())
	} else {
		for _, tick := range resp {
			applogger.Info("Symbol: %s, High: %v, Low: %v, Ask[%v, %v], Bid[%v, %v]",
				tick.Symbol, tick.High, tick.Low, tick.Ask, tick.AskSize, tick.Bid, tick.BidSize)
		}
	}
}

// 获取 品种的 k线

func GetSymbolKline() {
	optionalRequest := market.GetCandlestickOptionalRequest{Period: market.MIN1, Size: 10}
	resp, err := market_cl.GetCandlestick("btcusdt", optionalRequest)
	if err != nil {
		applogger.Error(err.Error())
	} else {
		for _, kline := range resp {
			applogger.Info("High=%v, Low=%v", kline.High, kline.Low)
		}
	}
}

func GetKlineWs() chan *mod.Kline {
	client := new(marketwebsocketclient.CandlestickWebSocketClient).Init(config.Host)
	huobi_kline := make(chan *mod.Kline)

	client.SetHandler(
		func() {
			client.Request("filusdt", "1min", 1569361140, 1569366420, "2305")

			client.Subscribe("filusdt", "1min", "2118")
		},
		func(response interface{}) {
			resp, ok := response.(market.SubscribeCandlestickResponse)
			if ok {
				if &resp != nil {
					if resp.Tick != nil {
						t := resp.Tick

						open, _ := t.Open.Float64()
						close, _ := t.Close.Float64()
						high, _ := t.High.Float64()
						low, _ := t.Low.Float64()

						kl := &mod.Kline{
							StartTime: time.Now(),
							Open:      open,
							Close:     close,
							High:      high,
							Low:       low,
						}

						huobi_kline <- kl

					}

					if resp.Data != nil {
						// applogger.Info("WebSocket returned data, count=%d", len(resp.Data))
						for _, t := range resp.Data {

							open, _ := t.Open.Float64()
							close, _ := t.Close.Float64()
							high, _ := t.High.Float64()
							low, _ := t.Low.Float64()

							kl := &mod.Kline{
								StartTime: time.Now(),
								Open:      open,
								Close:     close,
								High:      high,
								Low:       low,
							}

							huobi_kline <- kl
						}
					}
				}
			} else {
				applogger.Warn("Unknown response: %v", resp)
			}

		})

	client.Connect(true)

	return huobi_kline
	// client.UnSubscribe("btcusdt", "1min", "2118")

	//client.Close()
	// applogger.Info("Client closed")
}
