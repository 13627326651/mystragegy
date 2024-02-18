package strategy

import "github.com/rootpd/binance"

func CalProfit(start_price, stop_price, quote float64, rate float64) float64 {

	return 0

}

// 开单方向预测
/*
	1. 今天的恐慌贪婪指数  20%
	2. 订单簿 多空双方的博弈  20%
	3. k线走势 20%
	4. 最大回撤 20%
	5. 日线或者月线的影线 20%
*/

func CalSide() binance.OrderSide {
	return binance.SideBuy
}
