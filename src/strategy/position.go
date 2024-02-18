package strategy

import (
	"math"
	"strings"
	"sync"
	"time"
	"tinyquant/src/util"

	. "tinyquant/src/logger"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

type Position struct {
	*binance.FuturePositions
	*sync.RWMutex
	//存储到mysql
	PlaceCount          int                       // 加仓次数
	TryCount            int                       // 尝试买的次数
	AddFutureOrder      map[string]*MyFutureOrder // 普通加仓挂单
	PinFutureOrder      map[string]*MyFutureOrder // 插针加仓挂单
	TFutureOrder        map[string]*MyFutureOrder // 做T加仓挂单
	CloseFutureOrder    map[string]*MyFutureOrder // 平仓挂单
	CloseAllFutureOrder map[string]*MyFutureOrder // 触发止损挂单
}

type MyFutureOrder struct {
	*binance.ExecutedFutureOrder
	//存储到mysql
	OrdeType  util.ORIGIN_ORDER_STATUS //挂单类型
	OrderFlag util.ORIGIN_ORDER_FLAG   //仓位标志 0:手动 1:加仓单 2:减仓单
}

func (s *Strategy) ReloadPosition() {
	timer := time.NewTimer(1 * time.Minute)
	go func() {
		for {
			select {
			case <-timer.C:
				s.LoadPosition()
				timer.Reset(1 * time.Minute)
			}
		}
	}()
}

func (s *Strategy) LoadPosition() {
	// 加载当前持仓单
	res := Binance.GetFutureAccount(s.Symbol)
	for _, v := range res {
		if v.PositionSide == string(binance.LONG) {
			s.LongPosition.Lock()
			s.LongPosition.FuturePositions = v
			s.LongPosition.FuturePositions.PositionAmt = util.Round(s.LongPosition.FuturePositions.PositionAmt, 3)
			s.LongPosition.FuturePositions.EntryPrice = util.Round(s.LongPosition.FuturePositions.EntryPrice, 2)
			s.LongPosition.Unlock()
			continue
		}
		if v.PositionSide == string(binance.SHORT) {
			s.ShortPosition.Lock()
			s.ShortPosition.FuturePositions = v
			s.ShortPosition.FuturePositions.PositionAmt = util.Round(s.ShortPosition.FuturePositions.PositionAmt, 3)
			s.ShortPosition.FuturePositions.EntryPrice = util.Round(s.ShortPosition.FuturePositions.EntryPrice, 2)
			s.ShortPosition.Unlock()
			continue
		}
	}
}

func (s *Strategy) LoadAllOpenOrder() {
	ts, err := Binance.QueryBinanceAllFutureOrder(s.Symbol)
	if err != nil {
		return
	}

	for _, order := range ts {
		Logger.Info("当前挂单 : ", zap.Any("order", order))
		s.FutureOrder[order.ClientOrderID] = &MyFutureOrder{ExecutedFutureOrder: order}
		s.FutureOrder[order.ClientOrderID].Price = util.Round(s.FutureOrder[order.ClientOrderID].Price, 2)
		s.FutureOrder[order.ClientOrderID].OrigQty = util.Round(s.FutureOrder[order.ClientOrderID].OrigQty, 3)
		s.FutureOrder[order.ClientOrderID].ExecutedQty = util.Round(s.FutureOrder[order.ClientOrderID].ExecutedQty, 3)
	}

}

func (s *Strategy) SaveFutureOrder(futureOrder *MyFutureOrder, clientOrderID string) {
	//存储挂单
	s.PlaceOrderManager.AddOrderInfo(clientOrderID, futureOrder)
	if futureOrder.OrdeType == util.COMMON {
		s.Lock()
		s.FutureOrder[clientOrderID] = futureOrder
		s.Unlock()
	} else {
		if futureOrder.PositionSide == string(binance.LONG) { //必须为双向持仓
			if futureOrder.OrdeType == util.PIN && futureOrder.OrderFlag == util.ADDPOSITION {
				s.LongPosition.Lock()
				s.LongPosition.PinFutureOrder[clientOrderID] = futureOrder
				s.LongPosition.Unlock()
			} else if (futureOrder.OrdeType == util.CLOSECOMMON || futureOrder.OrdeType == util.PINCLOSECOMMON) && futureOrder.OrderFlag == util.DELPOSITION {
				s.LongPosition.Lock()
				s.LongPosition.CloseFutureOrder[clientOrderID] = futureOrder
				s.LongPosition.Unlock()
			} else if futureOrder.OrdeType == util.LOSSCLOSECOMMON {
				s.LongPosition.Lock()
				s.LongPosition.CloseAllFutureOrder[clientOrderID] = futureOrder
				s.LongPosition.Unlock()
			} else {
				Logger.Error("", zap.Any("订单异常", futureOrder))
			}

		} else if futureOrder.PositionSide == string(binance.SHORT) {
			if futureOrder.OrdeType == util.PIN && futureOrder.OrderFlag == util.ADDPOSITION {
				s.ShortPosition.Lock()
				s.ShortPosition.PinFutureOrder[clientOrderID] = futureOrder
				s.ShortPosition.Unlock()
			} else if (futureOrder.OrdeType == util.CLOSECOMMON || futureOrder.OrdeType == util.PINCLOSECOMMON) && futureOrder.OrderFlag == util.DELPOSITION {
				s.ShortPosition.Lock()
				s.ShortPosition.CloseFutureOrder[clientOrderID] = futureOrder
				s.ShortPosition.Unlock()
			} else if futureOrder.OrdeType == util.LOSSCLOSECOMMON {
				s.ShortPosition.Lock()
				s.ShortPosition.CloseAllFutureOrder[clientOrderID] = futureOrder
				s.ShortPosition.Unlock()
			} else {
				Logger.Error("", zap.Any("订单异常", futureOrder))
			}
		} else {
			Logger.Sugar().Error("持仓方向异常")
		}
	}
}

func (s *Strategy) DelFutureOrder(futureOrder *MyFutureOrder, clientOrderID string) {
	//删除挂单
	s.PlaceOrderManager.DelOrderInfo(clientOrderID)
	if futureOrder.OrdeType == util.COMMON {
		s.Lock()
		delete(s.FutureOrder, clientOrderID)
		s.Unlock()
	} else {
		if futureOrder.PositionSide == string(binance.LONG) { //必须为双向持仓
			if futureOrder.OrdeType == util.PIN && futureOrder.OrderFlag == util.ADDPOSITION {
				s.LongPosition.Lock()
				delete(s.LongPosition.PinFutureOrder, clientOrderID)
				s.LongPosition.Unlock()
			} else if (futureOrder.OrdeType == util.CLOSECOMMON || futureOrder.OrdeType == util.PINCLOSECOMMON) && futureOrder.OrderFlag == util.DELPOSITION {
				s.LongPosition.Lock()
				delete(s.LongPosition.CloseFutureOrder, clientOrderID)
				s.LongPosition.Unlock()
			} else if futureOrder.OrdeType == util.LOSSCLOSECOMMON {
				s.LongPosition.Lock()
				delete(s.LongPosition.CloseAllFutureOrder, clientOrderID)
				s.LongPosition.Unlock()
			} else {
				Logger.Error("", zap.Any("订单异常", futureOrder))
			}
		} else if futureOrder.PositionSide == string(binance.SHORT) {
			if futureOrder.OrdeType == util.PIN && futureOrder.OrderFlag == util.ADDPOSITION {
				s.ShortPosition.Lock()
				delete(s.ShortPosition.PinFutureOrder, clientOrderID)
				s.ShortPosition.Unlock()
			} else if (futureOrder.OrdeType == util.CLOSECOMMON || futureOrder.OrdeType == util.PINCLOSECOMMON) && futureOrder.OrderFlag == util.DELPOSITION {
				s.ShortPosition.Lock()
				delete(s.ShortPosition.CloseFutureOrder, clientOrderID)
				s.ShortPosition.Unlock()
			} else if futureOrder.OrdeType == util.LOSSCLOSECOMMON {
				s.ShortPosition.Lock()
				delete(s.ShortPosition.CloseAllFutureOrder, clientOrderID)
				s.ShortPosition.Unlock()
			} else {
				Logger.Error("", zap.Any("订单异常", futureOrder))
			}
		} else {
			Logger.Sugar().Error("持仓方向异常")
		}
	}
}

//定时扫描所有开仓挂单，处理掉部分成交单
func (s *Strategy) ClearPartiallyFilledOrder() {
	timer := time.NewTimer(1 * time.Minute)
	go func() {
		for {
			select {
			case <-timer.C:
				curPrice := s.KlineManager.MinuteKlineList.GetNewPrice()
				Logger.Sugar().Debugf("curPrice : %v", curPrice)
				s.LongPosition.RLock()
				for _, order := range s.LongPosition.PinFutureOrder {
					Logger.Sugar().Debugf("LongPosition.PinFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if order.Price > util.PressureLevel || order.Price < util.SupportLevel {
						return
					}
					if (order.OrderFlag == util.ADDPOSITION && order.OrdeType == util.PIN) &&
						((order.Status == binance.StatusPartiallyFilled && math.Abs(order.Price-curPrice) > 10.0) ||
							(order.Status == binance.StatusNew && time.Since(order.UpdateTime) > 5*time.Minute)) {
						Logger.Sugar().Warnf("取消开仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
						_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
						if err != nil {
							Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
						}
					}
				}
				for _, order := range s.LongPosition.CloseFutureOrder {
					Logger.Sugar().Debugf("LongPosition.CloseFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if order.OrderFlag == util.DELPOSITION && order.OrdeType == util.PINCLOSECOMMON && time.Since(order.UpdateTime) > 15*time.Minute {
						Logger.Sugar().Warnf("取消插针平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
						_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
						if err != nil {
							Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
						}
					}
				}
				s.LongPosition.RUnlock()

				s.ShortPosition.RLock()
				for _, order := range s.ShortPosition.PinFutureOrder {
					Logger.Sugar().Debugf("ShortPosition.PinFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if order.Price > util.PressureLevel || order.Price < util.SupportLevel {
						return
					}
					if (order.OrderFlag == util.ADDPOSITION && order.OrdeType == util.PIN) &&
						((order.Status == binance.StatusPartiallyFilled && math.Abs(order.Price-curPrice) > 10.0) ||
							(order.Status == binance.StatusNew && time.Since(order.UpdateTime) > 5*time.Minute)) {
						Logger.Sugar().Warnf("取消开仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
						_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
						if err != nil {
							Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
						}
					}
				}
				for _, order := range s.ShortPosition.CloseFutureOrder {
					Logger.Sugar().Debugf("ShortPosition.CloseFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if order.OrderFlag == util.DELPOSITION && order.OrdeType == util.PINCLOSECOMMON && time.Since(order.UpdateTime) > 15*time.Minute {
						Logger.Sugar().Warnf("取消插针平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
						_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
						if err != nil {
							Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
						}
					}
				}
				s.ShortPosition.RUnlock()

				s.RLock()
				for _, order := range s.FutureOrder {
					Logger.Sugar().Debugf("s.FutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if (order.PositionSide == string(binance.LONG) && order.Side == binance.SideBuy) ||
						(order.PositionSide == string(binance.SHORT) && order.Side == binance.SideSell) {
						if order.Price > util.PressureLevel || order.Price < util.SupportLevel {
							return
						}
						if (order.Status == binance.StatusPartiallyFilled && math.Abs(order.Price-curPrice) > 10.0) ||
							(order.Status == binance.StatusNew && time.Since(order.UpdateTime) > 5*time.Minute) {
							Logger.Sugar().Warnf("取消开仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
							if err != nil {
								Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
							}
						}
					}
				}
				s.RUnlock()

				timer.Reset(1 * time.Minute)
			}
		}
	}()
}

//定时扫描所有平仓单
func (s *Strategy) ScanCloseFutureOrder() {
	timer := time.NewTimer(1 * time.Minute)
	go func() {
		for {
			select {
			case <-timer.C:
				curPrice := s.KlineManager.MinuteKlineList.GetNewPrice()
				Logger.Sugar().Debugf("curPrice : %v", curPrice)
				s.LongPosition.RLock()
				for _, order := range s.LongPosition.CloseFutureOrder {
					Logger.Sugar().Debugf("LongPosition.CloseFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if math.Abs(order.Price-curPrice) > curPrice*util.CancelCloseOrderLevel {
						Logger.Sugar().Warnf("取消平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
						_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
						if err != nil {
							Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
						}
					}
				}
				s.LongPosition.RUnlock()

				s.ShortPosition.RLock()
				for _, order := range s.ShortPosition.CloseFutureOrder {
					Logger.Sugar().Debugf("ShortPosition.CloseFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if math.Abs(order.Price-curPrice) > curPrice*util.CancelCloseOrderLevel {
						Logger.Sugar().Warnf("取消平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
						_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
						if err != nil {
							Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
						}
					}
				}
				s.ShortPosition.RUnlock()

				s.RLock()
				for _, order := range s.FutureOrder {
					Logger.Sugar().Debugf("s.FutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					if (order.Type != binance.TypeSTOP) &&
						((order.PositionSide == string(binance.LONG) && order.Side == binance.SideSell) ||
							(order.PositionSide == string(binance.SHORT) && order.Side == binance.SideBuy)) {
						if math.Abs(order.Price-curPrice) > curPrice*util.CancelCloseOrderLevel {
							Logger.Sugar().Warnf("取消平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(order.OrderID))
							if err != nil {
								Logger.Error("cancel future order failed", zap.Error(err), zap.Any("order", order))
							}
						}
					}
				}
				s.RUnlock()
				timer.Reset(1 * time.Minute)
			}
		}
	}()
}

//取消所有平仓单
func (s *Strategy) CancelAllCloseFutureOrder(positionSide binance.PositionSide) {
	sli := []int64{}
	if positionSide == binance.LONG {
		s.LongPosition.RLock()
		for _, order := range s.LongPosition.CloseFutureOrder {
			sli = append(sli, int64(order.OrderID))
		}
		s.LongPosition.RUnlock()

		s.RLock()
		for _, order := range s.FutureOrder {
			if order.PositionSide == string(binance.LONG) && order.Side == binance.SideSell {
				sli = append(sli, int64(order.OrderID))
			}
		}
		s.RUnlock()
	} else {
		s.ShortPosition.RLock()
		for _, order := range s.ShortPosition.CloseFutureOrder {
			sli = append(sli, int64(order.OrderID))
		}
		s.ShortPosition.RUnlock()

		s.RLock()
		for _, order := range s.FutureOrder {
			if order.PositionSide == string(binance.SHORT) && order.Side == binance.SideBuy {
				sli = append(sli, int64(order.OrderID))
			}
		}
		s.RUnlock()
	}

	for _, v := range sli {
		_, err := Binance.CancelBinanceFutureOrder(s.Symbol, v)
		if err != nil {
			Logger.Error("cancel future order failed", zap.Error(err), zap.Any("orderId", v))
		}
	}
}

//获取多单仓位 所有平仓挂单的仓位 当前持仓价格
func (s *Strategy) GetLongBetweenAllCloseFutureOrderAndPositionD_Value() (float64, float64, float64) {
	s.LongPosition.RLock()
	defer s.LongPosition.RUnlock()
	var quantity float64 = 0
	var quantity1 float64 = 0
	var quantity2 float64 = 0
	var quantity3 float64 = 0
	for _, order := range s.LongPosition.CloseFutureOrder {
		quantity += util.Round(order.OrigQty, 3)
		quantity1 += order.OrigQty
	}
	// Logger.Sugar().Infof("平仓挂单的仓位1 : %v", quantity)
	for _, order := range s.LongPosition.PinFutureOrder {
		if order.Status == binance.StatusPartiallyFilled {
			quantity += util.Round(order.ExecutedQty, 3)
			quantity2 += util.Round(order.ExecutedQty, 3)
		}
	}
	// Logger.Sugar().Infof("平仓挂单的仓位2 : %v", quantity)
	s.RLock()
	defer s.RUnlock()
	for _, order := range s.FutureOrder {
		if order.PositionSide == string(binance.LONG) && order.Side == binance.SideSell && order.Type != binance.TypeSTOP {
			quantity += util.Round(order.OrigQty, 3)
			quantity3 += util.Round(order.OrigQty, 3)
		}
	}
	// Logger.Sugar().Infof("平仓挂单的仓位3 : %v", quantity)
	Logger.Sugar().Infof("[%v %v %v]多单仓位 : %v 价格 : %v", util.Round(quantity1, 3), util.Round(quantity2, 3), util.Round(quantity3, 3), util.Round(s.LongPosition.PositionAmt, 3), util.Round(s.LongPosition.EntryPrice, 2))
	return util.Round(s.LongPosition.PositionAmt, 3), util.Round(quantity, 3), util.Round(s.LongPosition.EntryPrice, 2)
}

//获取空单仓位 所有平仓挂单的仓位 当前持仓价格
func (s *Strategy) GetShortBetweenAllCloseFutureOrderAndPositionD_Value() (float64, float64, float64) {
	s.ShortPosition.RLock()
	defer s.ShortPosition.RUnlock()
	var quantity float64 = 0
	var quantity1 float64 = 0
	var quantity2 float64 = 0
	var quantity3 float64 = 0
	for _, order := range s.ShortPosition.CloseFutureOrder {
		quantity += util.Round(math.Abs(order.OrigQty), 3)
		quantity1 += util.Round(math.Abs(order.OrigQty), 3)
	}
	// Logger.Sugar().Infof("平仓挂单的仓位1 : %v", quantity)
	for _, order := range s.ShortPosition.PinFutureOrder {
		if order.Status == binance.StatusPartiallyFilled {
			quantity += util.Round(math.Abs(order.ExecutedQty), 3)
			quantity2 += util.Round(math.Abs(order.ExecutedQty), 3)
		}
	}
	// Logger.Sugar().Infof("平仓挂单的仓位2 : %v", quantity)
	s.RLock()
	defer s.RUnlock()
	for _, order := range s.FutureOrder {
		if order.PositionSide == string(binance.SHORT) && order.Side == binance.SideBuy && order.Type != binance.TypeSTOP {
			quantity += util.Round(math.Abs(order.OrigQty), 3)
			quantity3 += util.Round(math.Abs(order.OrigQty), 3)
		}
	}
	// Logger.Sugar().Infof("平仓挂单的仓位3 : %v", quantity)
	Logger.Sugar().Infof("[%v %v %v]空单仓位 : %v 价格 : %v", util.Round(quantity1, 3), util.Round(quantity2, 3), util.Round(quantity3, 3), util.Round(math.Abs(s.ShortPosition.PositionAmt), 3), util.Round(s.ShortPosition.EntryPrice, 2))
	return util.Round(math.Abs(s.ShortPosition.PositionAmt), 3), util.Round(quantity, 3), util.Round(s.ShortPosition.EntryPrice, 2)
}

func (s *Strategy) GetLongShortPinCloseFutureOrder() (bool, bool) {
	long, short := false, false
	s.LongPosition.RLock()
	for _, order := range s.LongPosition.CloseFutureOrder {
		Logger.Sugar().Debugf("LongPosition.CloseFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
		if order.OrderFlag == util.DELPOSITION && order.OrdeType == util.PINCLOSECOMMON {
			Logger.Sugar().Warnf("获取到多单插针平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
			long = true
			break
		}
	}
	s.LongPosition.RUnlock()
	s.ShortPosition.RLock()
	for _, order := range s.ShortPosition.CloseFutureOrder {
		Logger.Sugar().Debugf("ShortPosition.CloseFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
		if order.OrderFlag == util.DELPOSITION && order.OrdeType == util.PINCLOSECOMMON {
			Logger.Sugar().Warnf("获取到空单插针平仓挂单 %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
			short = true
			break
		}
	}
	s.ShortPosition.RUnlock()

	return long, short
}

//定时扫描仓位,创建平仓单
func (s *Strategy) ScanPositionAndCreatCloseFutureOrder() {
	timer := time.NewTimer(5 * time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
				curPrice := s.KlineManager.MinuteKlineList.GetNewPrice()
				long_positionAmt, long_closePosition, long_entryPrice := s.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
				if long_positionAmt-long_closePosition >= util.Quantity && math.Abs(curPrice-long_entryPrice) < curPrice*util.CreatCloseOrderLevel {
					newOrder := &OriginOrder{
						Symbol:       s.Symbol,
						OrderStatus:  util.CLOSECOMMON,
						Side:         binance.SideSell,
						PositionSide: binance.LONG,
						IsTest:       util.PlaceTest,
						OrderFlag:    util.DELPOSITION,
					}
					newOrder.Quantity = util.Round(long_positionAmt-long_closePosition, 3)
					newOrder.Price = util.Round(long_entryPrice+long_entryPrice*util.Profits, 2)
					s.PlaceOrderManager.MakePlaceOrder(newOrder)
				}
				short_positionAmt, short_closePosition, short_entryPrice := s.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
				if short_positionAmt-short_closePosition >= util.Quantity && math.Abs(curPrice-short_closePosition) < curPrice*util.CreatCloseOrderLevel {
					newOrder := &OriginOrder{
						Symbol:       s.Symbol,
						OrderStatus:  util.CLOSECOMMON,
						Side:         binance.SideBuy,
						PositionSide: binance.SHORT,
						IsTest:       util.PlaceTest,
						OrderFlag:    util.DELPOSITION,
					}
					newOrder.Quantity = util.Round(short_positionAmt-short_closePosition, 3)
					newOrder.Price = util.Round(short_entryPrice-short_entryPrice*util.Profits, 2)
					s.PlaceOrderManager.MakePlaceOrder(newOrder)
				}
				//检查止损单
				var long, short float64
				long_close, short_close := []int64{}, []int64{}
				s.LongPosition.RLock()
				for _, order := range s.LongPosition.CloseAllFutureOrder {
					Logger.Sugar().Debugf("LongPosition.CloseAllFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					long += order.OrigQty
					long_close = append(long_close, int64(order.OrderID))
				}
				s.LongPosition.RUnlock()
				s.ShortPosition.RLock()
				for _, order := range s.ShortPosition.CloseAllFutureOrder {
					Logger.Sugar().Debugf("ShortPosition.CloseAllFutureOrder : %+v OrdeType : %v OrderFlag : %v", order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
					short += order.OrigQty
					short_close = append(short_close, int64(order.OrderID))
				}
				s.ShortPosition.RUnlock()
				s.RLock()
				for _, order := range s.FutureOrder {
					if order.Type == binance.TypeSTOP {
						if order.PositionSide == string(binance.LONG) && order.Side == binance.SideSell {
							long += order.OrigQty
						} else if order.PositionSide == string(binance.SHORT) && order.Side == binance.SideBuy {
							short += order.OrigQty
						}
					}
				}
				s.RUnlock()
				sli := []int64{}
				if long < long_positionAmt {
					//先取消止损单
					sli = long_close
					Logger.Sugar().Warnf("多单止损单不足")
					s.MakeCloseOrder(&MyFutureOrder{ExecutedFutureOrder: &binance.ExecutedFutureOrder{PositionSide: string(binance.LONG), Side: binance.SideBuy}})
				}
				if short < short_positionAmt {
					sli = short_close
					Logger.Sugar().Warnf("空止损单不足")
					s.MakeCloseOrder(&MyFutureOrder{ExecutedFutureOrder: &binance.ExecutedFutureOrder{PositionSide: string(binance.SHORT), Side: binance.SideSell}})
				}
				for _, v := range sli {
					_, err := Binance.CancelBinanceFutureOrder(s.Symbol, v)
					if err != nil {
						Logger.Error("cancel future order failed", zap.Error(err), zap.Any("orderId", v))
					}
				}
				timer.Reset(5 * time.Second)
			}
		}
	}()
}

//定时扫描所有挂单 防止和账户对不上
func (s *Strategy) ScanFutureOrder() {
	timer := time.NewTimer(5 * time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
				curPrice := s.KlineManager.MinuteKlineList.GetNewPrice()
				Logger.Sugar().Debugf("curPrice : %v", curPrice)

				sli := make([]*MyFutureOrder, 0)

				s.LongPosition.RLock()
				for _, order := range s.LongPosition.PinFutureOrder {
					res, err := Binance.QueryBinanceOneFutureOrder(s.Symbol, order.ClientOrderID)
					if err != nil {
						if strings.Contains(err.Error(), "Order does not exist.") {
							Logger.Sugar().Errorf("多单开仓挂单丢失 %v %+v OrdeType : %v OrderFlag : %v", err, order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							sli = append(sli, order)
						} else {
							Logger.Sugar().Errorf("%v", err)
						}
					} else {
						Logger.Sugar().Debugf("多单开仓挂单id %v 数量 : %v %v 价格 : %v 已成交数量 : %v %v", res.ClientOrderID, res.OrigQty, order.OrigQty, res.Price, res.ExecutedQty, order.ExecutedQty)
					}
				}

				for _, order := range s.LongPosition.CloseFutureOrder {
					res, err := Binance.QueryBinanceOneFutureOrder(s.Symbol, order.ClientOrderID)
					if err != nil {
						if strings.Contains(err.Error(), "Order does not exist.") {
							Logger.Sugar().Errorf("多单平仓挂单丢失 %v %+v OrdeType : %v OrderFlag : %v", err, order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							sli = append(sli, order)
						} else {
							Logger.Sugar().Errorf("%v", err)
						}
					} else {
						Logger.Sugar().Debugf("多单平仓挂单id %v 数量 : %v %v 价格 : %v 已成交数量 : %v %v", res.ClientOrderID, res.OrigQty, order.OrigQty, res.Price, res.ExecutedQty, order.ExecutedQty)
					}
				}
				s.LongPosition.RUnlock()

				s.ShortPosition.RLock()
				for _, order := range s.ShortPosition.PinFutureOrder {
					res, err := Binance.QueryBinanceOneFutureOrder(s.Symbol, order.ClientOrderID)
					if err != nil {
						if strings.Contains(err.Error(), "Order does not exist.") {
							Logger.Sugar().Errorf("空单开仓挂单丢失 %v %+v OrdeType : %v OrderFlag : %v", err, order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							sli = append(sli, order)
						} else {
							Logger.Sugar().Errorf("%v", err)
						}
					} else {
						Logger.Sugar().Debugf("空单开仓挂单id %v 数量 : %v %v 价格 : %v 已成交数量 : %v %v", res.ClientOrderID, res.OrigQty, order.OrigQty, res.Price, res.ExecutedQty, order.ExecutedQty)
					}
				}

				for _, order := range s.ShortPosition.CloseFutureOrder {
					res, err := Binance.QueryBinanceOneFutureOrder(s.Symbol, order.ClientOrderID)
					if err != nil {
						if strings.Contains(err.Error(), "Order does not exist.") {
							Logger.Sugar().Errorf("空单平仓挂单丢失 %v %+v OrdeType : %v OrderFlag : %v", err, order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							sli = append(sli, order)
						} else {
							Logger.Sugar().Errorf("%v", err)
						}
					} else {
						Logger.Sugar().Debugf("空单平仓挂单id %v 数量 : %v %v 价格 : %v 已成交数量 : %v %v", res.ClientOrderID, res.OrigQty, order.OrigQty, res.Price, res.ExecutedQty, order.ExecutedQty)
					}
				}
				s.ShortPosition.RUnlock()

				s.RLock()
				for _, order := range s.FutureOrder {
					res, err := Binance.QueryBinanceOneFutureOrder(s.Symbol, order.ClientOrderID)
					if err != nil {
						if strings.Contains(err.Error(), "Order does not exist.") {
							Logger.Sugar().Errorf("err : %v \n手动挂单丢失 %+v OrdeType : %v OrderFlag : %v", err, order.ExecutedFutureOrder, order.OrdeType, order.OrderFlag)
							sli = append(sli, order)
						} else {
							Logger.Sugar().Errorf("%v", err)
						}
					} else {
						Logger.Sugar().Debugf("手动挂单id %v 持仓方向 : %v 买卖方向 : %v 数量 : %v %v 价格 : %v 已成交数量 : %v %v", res.ClientOrderID, res.PositionSide, res.Side, res.OrigQty, order.OrigQty, res.Price, res.ExecutedQty, order.ExecutedQty)
					}
				}
				s.RUnlock()

				for _, order := range sli {
					s.DelFutureOrder(order, order.ClientOrderID)
				}

				timer.Reset(5 * time.Second)
			}
		}
	}()
}
