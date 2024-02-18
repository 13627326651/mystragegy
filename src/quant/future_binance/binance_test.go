package future_test

import (
	"fmt"
	"sort"
	"testing"

	. "tinyquant/src/logger"
	fb "tinyquant/src/quant/future_binance"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
)

var Binance fb.Binance

func init() {

	//util.InitParam()
	InitLogger()
	Binance = fb.Binance{}
	Binance.InitBinance("QNmPFa4ztYzfuvurAORGONxzvYDuipGcZDTaCWDv5HT5cQbMzXchewQKIOjs2G4B", "HTXrk9v65XMYRA5TH8lDnE05SYcy1aZRS7opTJVJuUtwtpGhGO0TqhF7WcrTJO3r")
}

func Test_TestBinance(t *testing.T) {
	Binance.TestBinance()
}

func Test_GetMyTrade(t *testing.T) {

	Binance.GetMyTrade()
}

func Test_GetAccountInfo(t *testing.T) {
	Binance.GetAccountInfo()
}

func Test_NewOrder(t *testing.T) {
	Binance.NewOrder()
}
func Test_GetWithDrawHistory(t *testing.T) {
	Binance.GetWithDrawHistory()
}

func Test_GetTradeWs(t *testing.T) {

	Binance.GetTradeWs()
}

func Test_GetKlineWs(t *testing.T) {
	//future.GetKlineWs()
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

func Test_ExchangeInfo(t *testing.T) {
	s, _ := Binance.ExchangeInfo()
	fmt.Println(s)
}

func Test_GetAllFutureOrders(t *testing.T) {
	Binance.GetAllFutureOrders()
}

func Test_GetFutureBalance(t *testing.T) {
	ts, err := Binance.GetFutureBalance()
	if err != nil {
		t.Error(err)
	}

	for _, v := range ts {
		if v.Asset == "USDT" {
			fmt.Printf("%+v\n", v)
		}
	}
}

func Test_GetFutureAccount(t *testing.T) {
	Binance.GetFutureAccount("ETHUSDT")
}

func Test_GetUserPoundageInfo(t *testing.T) {
	Binance.GetUserPoundageInfo()
}

func Test_GetAdjustLeverage(t *testing.T) {
	//	future.GetAdjustLeverage()
}

func Test_GetPositionMargin(t *testing.T) {
	Binance.GetPositionMargin()
}

func TestUserTradesHistory(t *testing.T) {
	Binance.GetUserTradesHistory()
}

func Test_GetPremiumAndFundsRate(t *testing.T) {
	Binance.GetPremiumAndFundsRate()
}

func Test_GetPriceChangeSituation(t *testing.T) {
	Binance.GetPriceChangeSituation()
}

func Test_GetOpenInterestNums(t *testing.T) {
	Binance.GetOpenInterestNums()
}

func Test_GetBestBookTicker(t *testing.T) {
	Binance.GetBestBookTicker()
}

func Test_GetContractPosition(t *testing.T) {
	Binance.GetContractPosition()
}

func Test_GetTopLongShortPositionRatio(t *testing.T) {
	Binance.GetTopLongShortPositionRatio()
}

func Test_GetGlobalLongShortAccountRatio(t *testing.T) {
	ts, err := Binance.GetGlobalLongShortAccountRatio()
	fmt.Println(err)
	for _, v := range ts {
		fmt.Println("xxxxxxxxxxxxxxxxxxx", v)
	}
}

func Test_GetTakerlongshortRatio(t *testing.T) {
	Binance.GetTakerlongshortRatio()
}
func Test_QueryUserPositionSide(t *testing.T) {
	Binance.QueryBinanceUserPositionSide()
}

func Test_NewFutureOrder(t *testing.T) {

	symbol := "ETHUSDT"
	quantity := 0.1
	price := 3105.0
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
	order, err := Binance.NewBinanceFutureOrder(symbol, quantity, price, 0, side, position_side, "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("xxxx", order)
}

func Test_CancelFutureOrder(t *testing.T) {
	//future.CancelFutureOrder()
}

func Test_GetAccountWs(t *testing.T) {
	Binance.GetAccountWs()
}

func Test_QueryFutureOrder(t *testing.T) {
	ts, err := Binance.QueryBinanceAllFutureOrder(util.ETHUSDT)
	if err != nil {
		t.Error(err)
	}
	for _, v := range ts {
		fmt.Printf("%+v\n", v)
	}
}

func Test_GetNewPrice(t *testing.T) {
	price, err := Binance.GetBinanceNewPrice(util.ETHUSDT)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("xxxx", price)
}

func Test_GetFutureKlines(t *testing.T) {

	volume := make([]float64, 0)

	buy_volume := make([]float64, 0)
	sell_volume := make([]float64, 0)
	ts, err := Binance.GetFutureKlines(util.ETHUSDT, 24*60, binance.Minute)
	if err != nil {
		t.Error(err)
	}
	for _, v := range ts {
		fmt.Println("  ", v.OpenTime, v.TakerBuyBaseAssetVolume, "          ", v.Volume-v.TakerBuyBaseAssetVolume, "          ", v.Volume)
		volume = append(volume, v.Volume)
		buy_volume = append(buy_volume, v.TakerBuyBaseAssetVolume)
		sell_volume = append(sell_volume, v.Volume-v.TakerBuyBaseAssetVolume)
	}

	sort.Float64s(volume)
	sort.Float64s(buy_volume)
	sort.Float64s(sell_volume)
	vs := len(volume)
	avgVolume := volume[vs/2]
	mid_buy_volume := buy_volume[len(buy_volume)/2]
	mid_sell_volume := sell_volume[len(sell_volume)/2]

	volumeSum := 0.0
	for _, v := range volume {
		volumeSum += v
	}

	fmt.Println("xxxxxx", avgVolume, mid_buy_volume, mid_sell_volume, volumeSum/float64(vs))

}

func Test_ChangeMarginType(t *testing.T) {

	err := Binance.ChangeBinanceMarginType(util.ETHUSDT, binance.POSITION_CROSSED)
	if err != nil {
		t.Error(err)
	}
}

func Test_ChangeUserPositionSide(t *testing.T) {
	err := Binance.ChangeBinanceUserPositionSide(binance.PosithonBothSide)
	if err != nil {
		t.Error(err)
	}
}
