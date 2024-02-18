package future_test

import (
	"fmt"
	"testing"
	fb "tinyquant/src/quant/coin_future_binance"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
)

var Binance fb.Binance

func init() {

	//	config.InitConfig()
	//	util.InitParam()
	//	logger.InitLogger()

	Binance = fb.Binance{}
	Binance.InitBinance(util.BINANCE_API_KEY, util.BINANCE_SECRET_KEY)
}

func Test_FutureCoinKlines(t *testing.T) {
	//volume := make([]float64, 0)
	ts, err := Binance.GetFutureKlines("ETHUSD_PERP", 120, binance.FifteenMinutes)

	if err != nil {
		t.Error(err)
	}
	for _, v := range ts {
		fmt.Println(" : ", v)
		//volume = append(volume, v.Volume)
	}

}

func Test_GetDepth(t *testing.T) {
	orderbook, err := Binance.GetDepth(util.ETHUSDT, 10)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(orderbook.LastUpdateID)
	fmt.Println(orderbook.MessageTime)
	for k, v := range orderbook.Asks {
		fmt.Println(k, " : ", v)
	}
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	for k, v := range orderbook.Bids {
		fmt.Println(k, " : ", v)
	}

}

func Test_GetCoinTopLongShortPositionRatio(t *testing.T) {
	Binance.GetTopLongShortPositionRatio("", "")
}

func Test_GetCoinContractPosition(t *testing.T) {
	Binance.GetContractPosition("", "", "")
}

func Test_CoinGetNewPrice(t *testing.T) {
	gp, err := Binance.GetBinanceNewPrice("ETHUSD_PERP")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("xxx111", gp)
}

func Test_CoinGetPriceChangeSituation(t *testing.T) {
	Binance.GetPriceChangeSituation("")
}

func Test_GetCoinTakerlongshortRatio(t *testing.T) {
	Binance.GetTakerlongshortRatio("", "", "")
}

func Test_GetCoinBestBookTicker(t *testing.T) {
	Binance.GetBestBookTicker("")
}

func Test_GetCoinAllFutureOrders(t *testing.T) {
	Binance.GetAllFutureOrders("")
}

func Test_GetCoinOpenInterestNums(t *testing.T) {
	Binance.GetOpenInterestNums("")
}

func Test_CoinQueryUserPositionSide(t *testing.T) {
	Binance.QueryUserPositionSide()
}

func Test_GetCoinAdjustLeverage(t *testing.T) {
	// Binance.AdjustLeverage("BTCUSD_PERP", 1)
}

func Test_GetCoinPositionMargin(t *testing.T) {
	Binance.GetPositionMargin("", "", 2.0, 1)
}

func Test_CoinFutureBalance(t *testing.T) {
	Binance.GetFutureBalance()
}

func Test_CoinFutureAccount(t *testing.T) {
	Binance.GetFutureAccount()
}

func Test_NewFutureOrder(t *testing.T) {

	symbol := "ETHUSD_PERP"
	quantity := 5.0
	price := 3160.0
	side := binance.SideBuy
	position_side := binance.LONG

	/*	err := future.ChangeUserPositionSide(binance.PosithonBothSide)
		if err != nil {
			//t.Error(err)
		}

		// 改变全仓模式

		err = future.ChangeMarginType(symbol, binance.POSITION_ISOLATED)
		if err != nil {

			//t.Error(err)
		}
	*/
	order, err := Binance.NewBinanceFutureOrder(symbol, quantity, price, side, position_side)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("xxxx", order)
}

func Test_GetGlobalLongShortAccountRatio(t *testing.T) {

	xx, err := Binance.GetGlobalLongShortAccountRatio()
	if err != nil {
		t.Error(err)
	}

	for _, v := range xx {
		fmt.Println(v)
	}
}
