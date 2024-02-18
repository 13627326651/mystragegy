package strategy

import (
	"math"
	"tinyquant/src/util"

	. "tinyquant/src/logger"

	"github.com/rootpd/binance"
)

//创建平仓单
func (s *Strategy) MakePlaceOrder(futureOrder *MyFutureOrder) {
	//创建平仓单
	if futureOrder.OrderFlag == util.ADDPOSITION && futureOrder.OrdeType == util.PIN {
		s.LoadPosition() //实时更新下仓位
		var positionAmt, closePosition, entryPrice float64 = 0.0, 0.0, 0.0
		if futureOrder.PositionSide == string(binance.LONG) {
			positionAmt, closePosition, entryPrice = s.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
		} else {
			positionAmt, closePosition, entryPrice = s.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
		}
		newOrder := &OriginOrder{
			Symbol:       s.Symbol,
			OrderStatus:  util.CLOSECOMMON,
			Side:         futureOrder.Side,
			PositionSide: binance.PositionSide(futureOrder.PositionSide),
			IsTest:       util.PlaceTest,
			OrderFlag:    util.DELPOSITION,
		}
		if futureOrder.Side == binance.SideBuy {
			newOrder.Side = binance.SideSell
		} else {
			newOrder.Side = binance.SideBuy
		}

		if futureOrder.Status == binance.StatusCancelled || futureOrder.Status == binance.StatusExpired {
			newOrder.Quantity = util.Round(futureOrder.ExecutedQty, 3)
			if futureOrder.Side == binance.SideBuy {
				newOrder.Price = util.Round(futureOrder.Price+futureOrder.Price*util.Profits, 2)
			} else {
				newOrder.Price = util.Round(futureOrder.Price-futureOrder.Price*util.Profits, 2)
			}
			s.PlaceOrderManager.MakePlaceOrder(newOrder)
			return
		}
		Logger.Sugar().Infof("positionAmt : %v closePosition : %v futureOrder.OrigQty : %v", positionAmt, closePosition, futureOrder.OrigQty)
		if math.Abs(positionAmt-closePosition-futureOrder.OrigQty) >= 0.001 {
			Logger.Sugar().Infof("positionAmt : %v closePosition : %v futureOrder.OrigQty : %v", positionAmt, closePosition, futureOrder.OrigQty)
			//说明有仓位没有对应挂单
			//case1 手动取消了平仓挂单
			//case2 由于价格相差过大自动取消了平仓挂单
			newOrder.Quantity = util.Round(positionAmt-closePosition, 3)
			if futureOrder.Side == binance.SideBuy {
				newOrder.Price = util.Round(entryPrice+entryPrice*util.Profits/5.0, 2)
			} else {
				newOrder.Price = util.Round(entryPrice-entryPrice*util.Profits/5.0, 2)
			}
			curPrice := s.KlineManager.MinuteKlineList.GetNewPrice()
			//这里就是做T逻辑
			if math.Abs(newOrder.Price-curPrice) > curPrice*util.CancelCloseOrderLevel {
				newOrder.Quantity = util.Round(futureOrder.OrigQty, 3)
				if futureOrder.Side == binance.SideBuy {
					newOrder.Price = util.Round(futureOrder.Price+entryPrice*util.Profits, 2)
				} else {
					newOrder.Price = util.Round(futureOrder.Price-entryPrice*util.Profits, 2)
				}
			}
		} else {
			newOrder.Quantity = util.Round(futureOrder.OrigQty, 3)
			if futureOrder.Side == binance.SideBuy {
				newOrder.Price = util.Round(futureOrder.Price+entryPrice*util.Profits, 2)
			} else {
				newOrder.Price = util.Round(futureOrder.Price-entryPrice*util.Profits, 2)
			}
		}

		s.PlaceOrderManager.MakePlaceOrder(newOrder)
	}
}

//挂止损单
func (s *Strategy) MakeCloseOrder(futureOrder *MyFutureOrder) {
	//更新
	s.LoadPosition() //实时更新下仓位
	var positionAmt, _, _ float64 = 0.0, 0.0, 0.0
	if futureOrder.PositionSide == string(binance.LONG) {
		positionAmt, _, _ = s.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
	} else {
		positionAmt, _, _ = s.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
	}
	newOrder := &OriginOrder{
		Symbol:       s.Symbol,
		OrderStatus:  util.LOSSCLOSECOMMON,
		Side:         futureOrder.Side,
		PositionSide: binance.PositionSide(futureOrder.PositionSide),
		IsTest:       util.PlaceTest,
		OrderFlag:    util.DELPOSITION,
		Quantity:     positionAmt,
	}
	if futureOrder.Side == binance.SideBuy {
		newOrder.Side = binance.SideSell
		newOrder.Price = util.SupportLevel - 10
		newOrder.ClosePrice = util.SupportLevel
	} else {
		newOrder.Side = binance.SideBuy
		newOrder.Price = util.PressureLevel + 10
		newOrder.ClosePrice = util.PressureLevel
	}
	s.PlaceOrderManager.MakePlaceOrder(newOrder)
}
