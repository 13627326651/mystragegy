package strategy

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
	. "tinyquant/src/logger"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

type OriginOrder struct {
	Symbol       string
	Side         binance.OrderSide    // 买卖方向 SELL, BUY
	PositionSide binance.PositionSide // 持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
	Quantity     float64
	Price        float64
	ClosePrice   float64
	OrderStatus  util.ORIGIN_ORDER_STATUS
	OrderFlag    util.ORIGIN_ORDER_FLAG //仓位标志 0:手动 1:加仓单 2:减仓单
	IsTest       bool
}
type PlaceOrderManager struct {
	*sync.RWMutex
	Symbol   string
	Quantity float64 //默认下单数量

	LongLastPinPrice  float64 //最后一次插针加仓多单下单价格
	ShortLastPinPrice float64 //最后一次插针加仓空单下单价格

	LongLastDonePrice  float64 //最后一次插针加仓多单成功下单(成交)价格
	ShortLastDonePrice float64 //最后一次插针加仓空单成功下单(成交)价格

	LongPinOrderCancel  bool //多单插针单由于上一次的价格取消掉过
	ShortPinOrderCancel bool //空单插针单由于上一次的价格取消掉过

	LongLastPinPlaceOrderTime  int64 //最后一次插针加仓多单下单时间
	ShortLastPinPlaceOrderTime int64 //最后一次插针加仓空单下单时间

	LongPinCount  int32 //插针加仓多单挂单个数
	ShortPinCount int32 //插针加仓空单挂单个数

	LongContinuePlaceCount  int32 //多单连续下单的次数
	ShortContinuePlaceCount int32 //空单连续下单的次数

	OrderType     map[string]*MyFutureOrder
	Account       *BinanceFutureAsset // 账户信息
	TerracedPrice []float64           //连续开仓T度

	positionInfo PositionInfo
}

type PositionInfo interface {
	GetLongBetweenAllCloseFutureOrderAndPositionD_Value() (float64, float64, float64)
	GetShortBetweenAllCloseFutureOrderAndPositionD_Value() (float64, float64, float64)
	CancelAllCloseFutureOrder(binance.PositionSide)
	GetLongShortPinCloseFutureOrder() (bool, bool)
}

var placeLimit time.Time = time.Now()

func (p *PlaceOrderManager) MakePlaceOrder(order *OriginOrder) (*binance.FutureProcessedOrder, error) {
	p.Lock()
	defer p.Unlock()

	if order.Quantity == 0.0 {
		order.Quantity = util.Round(p.Quantity, 3)
	}

	if order.OrderFlag == util.ADDPOSITION { // 如果是加仓单
		switch order.OrderStatus {
		case util.COMMON:
			{

			}
		case util.PIN:
			{
				switch order.PositionSide {
				case binance.LONG:
					{
						if time.Now().Unix()-p.LongLastPinPlaceOrderTime > util.ContinuousOrderValidityTime*60 { //距离上一次多单时间过去5min
							p.LongContinuePlaceCount = 0 //重置连续下单的次数为0
						}
						index := int(p.LongContinuePlaceCount)
						if index >= len(p.TerracedPrice) {
							index = len(p.TerracedPrice)
						}
						if index > 0 { //说明是连续下单
							price := p.TerracedPrice[index-1] * order.Price
							if math.Abs(order.Price-p.LongLastPinPrice) < price { //和上一次下单的价格不能小于阙值
								Logger.Sugar().Infof("和上次下单价格相差小于 %v order : %+v", price, order)
								return nil, errors.New("price limit")
							} else {
								order.Quantity = util.Round(order.Quantity+util.Quantity*float64(index-1), 3)
								order.Price = util.Round(order.Price-order.Price*util.SpringPrice*float64(index), 2)
							}
						}

						positionAmt, _, entryPrice := p.positionInfo.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
						// if positionAmt >= util.Quantity*util.DoubleCreatOrderLevel { //当前仓位过大
						// 	if p.LongLastDonePrice != 0 {
						// 		limitPrice := math.Ceil(positionAmt / util.Quantity / util.DoubleCreatOrderLevel)
						// 		if p.LongLastDonePrice-order.Price < order.Price*util.Profits*limitPrice {
						// 			Logger.Sugar().Infof("取消加仓 多单总仓位  : %v 多单均价 : %v 加仓价格 : %v 加仓数量  : %v 上次加仓价格 : %v limit : %v %v", positionAmt, entryPrice, order.Price, order.Quantity, p.LongLastDonePrice, limitPrice, order.Price*util.Profits*limitPrice)
						// 			p.LongPinOrderCancel = true
						// 			return nil, errors.New("quantity limit")
						// 		}
						// 	}
						// }
						// if p.LongPinOrderCancel {
						// 	order.Quantity = util.Round(order.Quantity+util.Quantity, 3)
						// }
						if positionAmt != 0 && math.Abs(order.Price-entryPrice) > order.Price*util.IncreaseQuantityLevel {
							Logger.Sugar().Warnf("增加加仓 多单总仓位  : %v 多单均价 : %v 加仓价格 : %v 加仓数量  : %v index : %v", positionAmt, entryPrice, order.Price, order.Quantity, index)
							order.Quantity = util.Round(order.Quantity+util.Quantity, 3)
						}
						// turnPositionAmt, _, turnEntryPrice := p.positionInfo.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
						// if turnPositionAmt >= util.Quantity*6 && turnPositionAmt > positionAmt*2 && order.Quantity == util.Profits {
						// 	Logger.Sugar().Warnf("增加加仓 多单总仓位  : %v 多单均价 : %v 加仓价格 : %v 加仓数量  : %v 空单总仓位  : %v", positionAmt, entryPrice, order.Price, order.Quantity, turnPositionAmt)
						// 	order.Quantity = order.Quantity + util.Quantity
						// }

						p.LongContinuePlaceCount++
						p.LongLastPinPlaceOrderTime = time.Now().Unix()
						p.LongLastPinPrice = order.Price

					}
				case binance.SHORT:
					{
						if time.Now().Unix()-p.ShortLastPinPlaceOrderTime > util.ContinuousOrderValidityTime*60 { //距离上一次多单时间过去5min
							p.ShortContinuePlaceCount = 0 //重置连续下单的次数为0
						}
						index := int(p.ShortContinuePlaceCount)
						if index >= len(p.TerracedPrice) {
							index = len(p.TerracedPrice)
						}
						if index > 0 { //说明是连续下单
							price := p.TerracedPrice[index-1] * order.Price
							if math.Abs(order.Price-p.ShortLastPinPrice) < price { //和上一次下单的价格不能小于阙值
								Logger.Sugar().Infof("和上次下单价格相差小于 %v order : %+v", price, order)
								return nil, errors.New("price limit")
							} else {
								order.Quantity = util.Round(order.Quantity+util.Quantity*float64(index-1), 3)
								order.Price = util.Round(order.Price+order.Price*util.SpringPrice*float64(index), 2)
							}
						}

						positionAmt, _, entryPrice := p.positionInfo.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
						// if positionAmt >= util.Quantity*util.DoubleCreatOrderLevel { //当前仓位过大
						// 	if p.ShortLastDonePrice != 0 {
						// 		limitPrice := math.Ceil(positionAmt / util.Quantity / util.DoubleCreatOrderLevel)
						// 		if order.Price-p.ShortLastDonePrice < order.Price*util.Profits*limitPrice {
						// 			Logger.Sugar().Infof("取消加仓 空单总仓位  : %v 空单均价 : %v 加仓价格 : %v 加仓数量  : %v 上次加仓价格 : %v limit : %v %v", positionAmt, entryPrice, order.Price, order.Quantity, p.ShortLastDonePrice, limitPrice, order.Price*util.Profits*limitPrice)
						// 			p.ShortPinOrderCancel = true
						// 			return nil, errors.New("quantity limit")
						// 		}
						// 	}
						// }
						// if p.ShortPinOrderCancel {
						// 	order.Quantity = util.Round(order.Quantity+util.Quantity, 3)
						// }
						if positionAmt != 0 && math.Abs(order.Price-entryPrice) > order.Price*util.IncreaseQuantityLevel {
							Logger.Sugar().Warnf("增加加仓 空单总仓位  : %v 空单均价 : %v 加仓价格 : %v 加仓数量  : %v", positionAmt, entryPrice, order.Price, order.Quantity)
							order.Quantity = util.Round(order.Quantity+util.Quantity, 3)
						}
						// turnPositionAmt, _, turnEntryPrice := p.positionInfo.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
						// if turnPositionAmt >= util.Quantity*6 && turnPositionAmt > positionAmt*2 && order.Quantity == util.Profits {
						// 	Logger.Sugar().Warnf("增加加仓 空单总仓位  : %v 空单均价 : %v 加仓价格 : %v 加仓数量  : %v 多单总仓位  : %v", positionAmt, entryPrice, order.Price, order.Quantity, turnPositionAmt)
						// 	order.Quantity = order.Quantity + util.Quantity
						// }

						p.ShortContinuePlaceCount++
						p.ShortLastPinPlaceOrderTime = time.Now().Unix()
						p.ShortLastPinPrice = order.Price

					}
				}
				if (order.Price > util.PressureLevel || order.Price < util.SupportLevel) && order.OrderStatus != util.FLOW {
					Logger.Sugar().Warn("开仓价格超多压力位或者支撑位,curprice : %v PressureLevel : %v,SupportLevel : %v", order.Price, util.PressureLevel, util.SupportLevel)
					if time.Since(placeLimit) > 15*time.Minute {
						placeLimit = time.Now()
						util.SendOrderMsg(fmt.Sprintf("开仓价格超多压力位或者支撑位,请介入处理\norder price : %v \nPressureLevel : %v\n,SupportLevel : %v", order.Price, util.PressureLevel, util.SupportLevel))
					}
					return nil, nil
				}

			}
		default:

		}
	} else if order.OrderFlag == util.DELPOSITION {
		switch order.OrderStatus {

		case util.PINCLOSECOMMON:
			{
				switch order.PositionSide {
				case binance.LONG:
					ok, _ := p.positionInfo.GetLongShortPinCloseFutureOrder()
					if ok {
						return nil, errors.New("has Long pin close future order")
					}
				case binance.SHORT:
					_, ok := p.positionInfo.GetLongShortPinCloseFutureOrder()
					if ok {
						return nil, errors.New("has Short pin close future order")
					}
				}
			}
		case util.LOSSCLOSECOMMON:
			{
				Logger.Sugar().Debugf("挂止损触发单")
			}
		}
	}

	if order.IsTest {
		Logger.Info("test下单", zap.Any(order.Symbol, order))
		return nil, nil
	}

	customOrderId := strconv.FormatInt(time.Now().UnixNano(), 10)
	p.OrderType[customOrderId] = &MyFutureOrder{
		ExecutedFutureOrder: &binance.ExecutedFutureOrder{},
		OrdeType:            order.OrderStatus,
		OrderFlag:           order.OrderFlag,
	}
	// if order.Quantity >= 20 {
	// 	Logger.Error("下单拦截", zap.Any(order.Symbol, order))
	// 	return nil, nil
	// }
	resOrder, err := Binance.NewBinanceFutureOrder(order.Symbol, order.Quantity, order.Price, order.ClosePrice, order.Side, order.PositionSide, customOrderId)
	if err != nil {
		delete(p.OrderType, customOrderId)
		Logger.Error("new future order failed ", zap.Error(err), zap.Any("order", order))
		return nil, err
	}
	p.OrderType[customOrderId].ActivetePrice = resOrder.ActivatePrice
	p.OrderType[customOrderId].PriceRate = resOrder.PriceRate
	p.OrderType[customOrderId].WorkingType = resOrder.WorkingType
	p.OrderType[customOrderId].PriceProtect = resOrder.PriceProtect
	Logger.Info("开始下单", zap.Any(order.Symbol, order))
	return resOrder, nil

}

func (p *PlaceOrderManager) GetOrderInfo(customId string) *MyFutureOrder {
	p.RLock()
	defer p.RUnlock()
	return p.OrderType[customId]
}

func (p *PlaceOrderManager) AddOrderInfo(customId string, order *MyFutureOrder) {
	p.Lock()
	defer p.Unlock()
	p.OrderType[customId] = order
	if order.OrdeType == util.PIN {
		if order.PositionSide == string(binance.LONG) {
			p.LongPinCount++
		} else {
			p.ShortPinCount++
		}
	}
}

func (p *PlaceOrderManager) DelOrderInfo(customId string) {
	p.Lock()
	defer p.Unlock()
	if order, ok := p.OrderType[customId]; ok {
		if order.OrdeType == util.PIN {
			if order.PositionSide == string(binance.LONG) {
				p.LongPinCount--
				if order.Status == binance.StatusFilled {
					p.LongLastDonePrice = order.Price
					p.LongPinOrderCancel = false
				}
			} else {
				p.ShortPinCount--
				if order.Status == binance.StatusFilled {
					p.ShortLastDonePrice = order.Price
					p.ShortPinOrderCancel = false
				}
			}
		}
		delete(p.OrderType, customId)
	}
	if p.LongPinCount < 0 {
		p.LongPinCount = 0
	}
	if p.ShortPinCount < 0 {
		p.ShortPinCount = 0
	}
}
