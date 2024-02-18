package db

import (
	"time"
	. "tinyquant/src/logger"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

type Kline struct {
	OpenTime    time.Time `xorm:"open_time"`
	CloseTime   time.Time `xorm:"close_time"`
	Open        float64   `xorm:"open"`
	Close       float64   `xorm:"close"`
	High        float64   `xorm:"high"`
	Low         float64   `xorm:"low"`
	Volume      float64   `xorm:"volume"`
	BuyVolume   float64   `xorm:"buy_volume"`
	SellVolume  float64   `xorm:"sell_volume"`
	Quote       float64   `xorm:"quote"`
	BuyQuote    float64   `xorm:"buy_quote"`
	SellQuote   float64   `xorm:"sell_quote"`
	TradeNumber int       `xorm:"trade_number"`
}

func PutKline(kline *Kline) error {
	_, err := GetSession().Table("kline").Insert(kline)
	if err != nil {
		Logger.Error("insert user income  failed", zap.Error(err))
		return err
	}

	return nil
}

type Order struct {
	OrderID         int64                    `xorm:"order_id"`
	Symbol          string                   `xorm:"symbol"`
	Status          binance.OrderStatus      `xorm:"status"` //订单状态 NEW  PARTIALLY_FILLED  FILLED CANCELED EXPIRED NEW_INSURANCE 风险保障基金(强平)  NEW_ADL 自动减仓序列(强平)
	ClientOrderID   string                   `xorm:"client_order_id"`
	Price           float64                  `xorm:"price"`     //委托价格
	AvgPrice        float64                  `xorm:"avg_price"` // 平均成交价
	OrigQty         float64                  `xorm:"origqty"`   // 原始委托数量
	OrderType       binance.OrderType        `xorm:"order_type"`
	Side            binance.OrderSide        `xorm:"side"`
	PositionSide    binance.PositionSide     `xorm:"position_side"`
	StopPrice       float64                  `xorm:"stop_price"`
	Fee             float64                  `xorm:"fee"` //手续费
	LastVolume      float64                  `xorm:"last_volume"`
	Volume          float64                  `xorm:"volume"` //总成交量
	UpdateTime      time.Time                `xorm:"update_time"`
	Profit          float64                  `xorm:"profit"`
	PositionStatus  int                      `xorm:"position_status"` // 仓位状态 0: 开仓单 ， 1 : 平仓单
	IsOrderFinish   int                      `xorm:"is_order_finish"` // 订单流程结束 0 : 以结束 , 1 : 持仓中
	IsReduce        bool                     `xorm:"is_reduce"`
	Leverage        int                      `xorm:"leverage"` //杠杆倍数
	OrigOrderStatus util.ORIGIN_ORDER_STATUS `xorm:"orig_order_status"`
	ProfitLossClose bool                     `xorm:"-"`
}

func PutOrder(order *Order) error {
	_, err := GetSession().Table("order").Insert(order)
	if err != nil {
		Logger.Error("insert order  failed", zap.Error(err))
		return err
	}

	return nil
}

func UpdateOrder(order *Order) error {

	od := new(Order)
	od.IsOrderFinish = order.IsOrderFinish
	od.UpdateTime = order.UpdateTime

	_, err := GetSession().Table("order").Where("symbol = ? and order_id = ?", order.Symbol, order.OrderID).Cols("is_order_finish").Update(od)
	if err != nil {
		//Logger.Error("Update order  failed", zap.Error(err))
		return err
	}
	return nil
}

func UpdateOrderByOrderType(order *Order) error {

	od := new(Order)
	od.IsOrderFinish = order.IsOrderFinish
	od.UpdateTime = order.UpdateTime
	od.OrigOrderStatus = order.OrigOrderStatus

	_, err := GetSession().Table("order").Where("symbol = ? and order_id = ? and orig_order_status = ?", order.Symbol, order.OrderID, od.OrigOrderStatus).Cols("is_order_finish").Update(od)
	if err != nil {
		//Logger.Error("Update order  failed", zap.Error(err))
		return err
	}
	return nil
}

func GetNewOrder() ([]*Order, error) {

	var OrderList []*Order

	err := GetSession().Table("order").Where("is_order_finish = ?", true).Find(&OrderList)
	if err != nil {
		Logger.Error("get order list failed", zap.Error(err))
		return nil, err
	}
	return OrderList, nil

}
