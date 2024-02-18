package future

import (
	"context"
	"fmt"
	"time"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
)

type Binance struct {
	binance.Binance
}

func (b *Binance) InitBinance(apikey, secretkey string) {

	hmacSigner := &binance.HmacSigner{
		Key: []byte(secretkey),
	}
	ctx := context.Background()
	// use second return value for cancelling request
	binanceService := binance.NewAPIService(
		"https://dapi.binance.com",
		apikey,
		hmacSigner,
		ctx,
	)

	b.Binance = binance.NewBinance(binanceService)

}

func (b *Binance) TestBinance() {

	kl, err := b.Klines(binance.KlinesRequest{
		Symbol:   "BTCUSDT",
		Interval: binance.Minute,
	})
	if err != nil {
		panic(err)
	}

	for _, v := range kl {
		fmt.Println("xxxxxxxxxxxxxx : ", v)
	}

}

// 获取所有历史订单
func (b *Binance) GetMyTrade() {

	kl, err := b.MyTrades(binance.MyTradesRequest{
		Symbol:     "ETHUSDT",
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range kl {
		fmt.Println("xxxxxxxxxxxxxx : ", v)
	}

}

// 获取 账户信息
func (b *Binance) GetAccountInfo() {

	kl, err := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("xxxxxxxxxxx : ", kl)

	for k, v := range kl.Balances {
		fmt.Println(k, " : ", v)
	}

}

// 查询目前下单数

func (b *Binance) GetMyNowOrder() {

	kl, err := b.OpenOrders(
		binance.OpenOrdersRequest{})
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range kl {
		fmt.Println("xxxxxxxxxxxxxx : ", v)
	}

}

// 下单

func (b *Binance) NewOrder() {

	symbol := "FILUSDT"

	precision := util.Order_Precision[symbol]

	t := binance.NewOrderRequest{
		Symbol:      symbol,
		Quantity:    util.ToFloat64("1.12"),
		Price:       util.ToFloat64("56.5"),
		Side:        binance.SideSell,
		TimeInForce: binance.GTC,
		Type:        binance.TypeLimit,
		Timestamp:   time.Now(),
		Precision:   int32(precision),
	}

	err := b.NewOrderTest(t)
	if err != nil {
		fmt.Println(err)
	}

	//	fmt.Println("xxxx", kl)

}

func (b *Binance) GetWithDrawHistory() {

	wd, err := b.WithdrawHistory(binance.HistoryRequest{
		Timestamp:  time.Now(),
		RecvWindow: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("xxxx", wd)
	for k, v := range wd {
		fmt.Println(k, " : ", v)
	}
}

func (b *Binance) GetExchangeInfo() {
	str, err := b.ExchangeInfo()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("exchangeinfo : ", str)
}
