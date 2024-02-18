package strategy

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	. "tinyquant/src/logger"
	"tinyquant/src/mod"
	"tinyquant/src/util"

	quant "tinyquant/src/quant"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

var Binance quant.Binance

type Strategy struct {
	*sync.RWMutex
	Symbol            string
	LongPosition      Position                         //多单持仓信息
	ShortPosition     Position                         //空单持仓信息
	FutureOrder       map[string]*MyFutureOrder        //所有手动的挂单
	KlineWs           chan *mod.Kline                  //K线事件
	Ch15Kline         chan *mod.Kline                  //K线事件
	Ch4hKline         chan *mod.Kline                  //K线事件
	AccWs             chan *binance.FutureAccountEvent //账户变动事件
	KlineManager      *Market                          //K线
	OBM               *OrderBookMap                    //深度
	PlaceOrderManager *PlaceOrderManager               //开单管理
}

func (s *Strategy) placeAssert(ke *mod.Kline, kqueue *MyKlineQueue) {
	//更新本地K线
	kqueue.EnQqueu(&Kline{
		Open:      ke.Open,
		Close:     ke.Close,
		High:      ke.High,
		Low:       ke.Low,
		Volume:    ke.Volume,
		CloseTime: ke.CloseTime,
		BuyVolume: ke.BuyVolume,
	})
	if ke.Final {
		//k线结束更新均值
		kqueue.UpdateUpDownLink(true)
	}
	go func() {
		//插针下单判断
		upl := kqueue.GetUpDownLink()
		if time.Now().Unix()%5 == 0 {
			Logger.Sugar().Debugf("平均成交量 * %v : %v half 采样点 : %v K线当前成交量  : %v k线当前价格 : %v 均价 : %v",
				util.VolumeIncrease, upl.AvgVolume*util.VolumeIncrease, upl.HalfSampleAvgPrice*util.VolumeIncrease, ke.Volume, ke.Close, upl.AvgPrice)
		}
		if ke.Volume > upl.AvgVolume*util.VolumeIncreaseForClose && ke.Volume > upl.HalfSampleAvgPrice*util.VolumeIncreaseForClose {
			go kqueue.UpdateUpDownLink(true) //先更新

			if ke.Open > ke.Close && ke.Close < upl.AvgPrice-ke.Close*util.SpringPrice { //向下插针
				turnPositionAmt, _, turnEntryPrice := s.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
				if turnPositionAmt != 0 && turnEntryPrice > ke.Close {
					newOrder := &OriginOrder{
						Symbol:       s.Symbol,
						OrderStatus:  util.PINCLOSECOMMON,
						Side:         binance.SideBuy,
						PositionSide: binance.PositionSide(binance.SHORT),
						IsTest:       util.PlaceTest,
						OrderFlag:    util.DELPOSITION,
						Quantity:     util.Round(turnPositionAmt, 3),
						Price:        util.Round(ke.Close+ke.Close*util.SpringPrice/2, 2),
					}
					Logger.Sugar().Debugf("插针取消平仓单,创建新的平仓单 %+v", newOrder)
					// p.positionInfo.CancelAllCloseFutureOrder(newOrder.PositionSide)
					s.PlaceOrderManager.MakePlaceOrder(newOrder)
				} else {
					Logger.Sugar().Debugf("turnEntryPrice : %v", turnEntryPrice)
				}
			} else if ke.Open < ke.Close && ke.Close > upl.AvgPrice+ke.Close*util.SpringPrice { //向上插针
				turnPositionAmt, _, turnEntryPrice := s.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
				if turnPositionAmt != 0 && turnEntryPrice < ke.Close {
					newOrder := &OriginOrder{
						Symbol:       s.Symbol,
						OrderStatus:  util.PINCLOSECOMMON,
						Side:         binance.SideSell,
						PositionSide: binance.PositionSide(binance.LONG),
						IsTest:       util.PlaceTest,
						OrderFlag:    util.DELPOSITION,
						Quantity:     util.Round(turnPositionAmt, 3),
						Price:        util.Round(ke.Close-ke.Close*util.SpringPrice/2, 2),
					}
					Logger.Sugar().Debugf("插针取消平仓单,创建新的平仓单 %+v", newOrder)
					// p.positionInfo.CancelAllCloseFutureOrder(newOrder.PositionSide)
					s.PlaceOrderManager.MakePlaceOrder(newOrder)
				} else {
					Logger.Sugar().Debugf("turnEntryPrice : %v", turnEntryPrice)
				}
			}

			if ke.Volume > upl.AvgVolume*util.VolumeIncrease && ke.Volume > upl.HalfSampleAvgPrice*util.VolumeIncrease {
				if ke.Open > ke.Close && ke.Close < upl.AvgPrice-ke.Close*util.SpringPrice { //向下插针
					order := &OriginOrder{
						Symbol:       s.Symbol,
						OrderStatus:  util.PIN,
						Side:         binance.SideBuy,
						PositionSide: binance.LONG,
						OrderFlag:    util.ADDPOSITION,
						IsTest:       util.PlaceTest,
					}
					order.Price = util.Round(ke.Close-ke.Close*util.SpringPrice, 2)
					order.Quantity = util.Round(util.Quantity, 3)
					Logger.Sugar().Infof("向下插针 分钟平均成交量 * %v : %v K线当前成交量 : %v k线当前价格 : %v 创建开仓单价格 : %v",
						util.VolumeIncrease, upl.AvgVolume*util.VolumeIncrease, ke.Volume, ke.Close, order.Price)
					s.PlaceOrderManager.MakePlaceOrder(order)
				} else if ke.Open < ke.Close && ke.Close > upl.AvgPrice+ke.Close*util.SpringPrice { //向上插针
					order := &OriginOrder{
						Symbol:       s.Symbol,
						OrderStatus:  util.PIN,
						Side:         binance.SideSell,
						PositionSide: binance.SHORT,
						OrderFlag:    util.ADDPOSITION,
						IsTest:       util.PlaceTest,
					}
					order.Price = util.Round(ke.Close+ke.Close*util.SpringPrice, 2)
					order.Quantity = util.Round(util.Quantity, 3)
					Logger.Sugar().Infof("向上插针 分钟平均成交量 * %v : %v K线当前成交量 : %v k线当前价格 : %v 创建开仓单价格 : %v",
						util.VolumeIncrease, upl.AvgVolume*util.VolumeIncrease, ke.Volume, ke.Close, order.Price)
					s.PlaceOrderManager.MakePlaceOrder(order)
				}
			}
		}
	}()
}

func (s *Strategy) StrategyLoop(ct bool) error {
	Logger.Info("开启策略")

	for {
		select {
		case ke := <-s.KlineWs:
			s.placeAssert(ke, s.KlineManager.MinuteKlineList)
		case ke := <-s.Ch15Kline:
			s.placeAssert(ke, s.KlineManager.FifteenMinuteKlineList)
		case ke := <-s.Ch4hKline:
			s.placeAssert(ke, s.KlineManager.FourHourKlineList)
		case acc := <-s.AccWs:
			switch acc.EventName {
			case util.ACCOUNT_UPDATE: //TODO 需要定时去更新最新可下单余额
				Logger.Debug("ACCOUNT_UPDATE")
				for _, v := range acc.AE.Acc.Balance {
					if v.Symbol != util.ACCOUNTASSET[s.Symbol] {
						continue
					}
					Logger.Sugar().Debugf("%+v", v)
					s.PlaceOrderManager.Account.Lock()
					s.PlaceOrderManager.Account.AvailableBalance = v.CurBalance
					s.PlaceOrderManager.Account.Unlock()
				}

				for _, v := range acc.AE.Acc.Property {
					if v.Symbol != s.Symbol {
						continue
					}
					Logger.Sugar().Debugf("%+v", v)
					switch v.PS {
					case string(binance.LONG):
						//仓位没了
						s.LongPosition.Lock()
						Logger.Sugar().Infof("做多方向仓位变动,原始持仓数量 : %v 价格 : %v 未实现盈亏 : %v 变动后 持仓数量 : %v 价格 : %v 未实现盈亏 : %v",
							s.LongPosition.PositionAmt, s.LongPosition.EntryPrice, s.LongPosition.UnrealizedProfit, v.Pa, v.EP, v.UP)
						if v.Pa == 0 {
							//删除本地
							for _, v := range s.LongPosition.CloseFutureOrder {
								Logger.Sugar().Errorf("取消平仓单 %+v", v.ExecutedFutureOrder)
								_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(v.OrderID))
								if err != nil {
									Logger.Error("cancel future order failed", zap.Error(err))
								}
								delete(s.LongPosition.CloseFutureOrder, v.ClientOrderID)
							}
						}
						s.LongPosition.UnrealizedProfit = v.UP
						s.LongPosition.EntryPrice = util.Round(v.EP, 2)
						s.LongPosition.PositionAmt = util.Round(v.Pa, 3)
						s.LongPosition.UpdateTime = time.Now()
						s.LongPosition.Unlock()

					case string(binance.SHORT):
						//仓位没了
						//有时候是自动平的有时候是手动平的
						s.ShortPosition.Lock()
						Logger.Sugar().Infof("做空方向仓位变动,原始持仓数量 : %v 价格 : %v 未实现盈亏 : %v 变动后 持仓数量 : %v 价格 : %v 未实现盈亏 : %v",
							s.ShortPosition.PositionAmt, s.ShortPosition.EntryPrice, s.ShortPosition.UnrealizedProfit, v.Pa, v.EP, v.UP)
						if v.Pa == 0 {
							//删除本地
							for _, v := range s.ShortPosition.CloseFutureOrder {
								Logger.Sugar().Errorf("取消平仓单 %+v", v.ExecutedFutureOrder)
								_, err := Binance.CancelBinanceFutureOrder(s.Symbol, int64(v.OrderID))
								if err != nil {
									Logger.Error("cancel future order failed", zap.Error(err))
								}
								delete(s.ShortPosition.CloseFutureOrder, v.ClientOrderID)
							}
						}
						s.ShortPosition.UnrealizedProfit = v.UP
						s.ShortPosition.EntryPrice = util.Round(v.EP, 2)
						s.ShortPosition.PositionAmt = util.Round(v.Pa, 3)
						s.ShortPosition.UpdateTime = time.Now()
						s.ShortPosition.Unlock()
					}
				}
			case util.ORDER_TRADE_UPDATE:
				order := acc.OE.Order
				if order.Symbol != s.Symbol {
					continue
				}
				Logger.Info("ORDER_TRADE_UPDATE")
				Logger.Sugar().Debugf("%+v", order)

				positionSide := "多单"
				if order.PositionSide == string(binance.SHORT) {
					positionSide = "空单"
				}
				side := "买"
				if order.Side == string(binance.SideSell) {
					side = "卖"
				}

				var futureOrder *MyFutureOrder = nil
				if futureOrder = s.PlaceOrderManager.GetOrderInfo(order.ClientOrderID); futureOrder == nil {
					futureOrder = &MyFutureOrder{ExecutedFutureOrder: &binance.ExecutedFutureOrder{}}
				}

				orderFlag := ""
				if futureOrder.OrderFlag == util.ADDPOSITION && futureOrder.OrdeType == util.PIN {
					orderFlag = "自动加仓单"
				} else if futureOrder.OrderFlag == util.DELPOSITION && (futureOrder.OrdeType == util.CLOSECOMMON || futureOrder.OrdeType == util.PINCLOSECOMMON) {
					orderFlag = "自动减仓单"
				} else if futureOrder.OrdeType == util.COMMON && futureOrder.OrderFlag == util.UNKNNOW {
					if (order.PositionSide == string(binance.LONG) && order.Side == string(binance.SideBuy)) ||
						(order.PositionSide == string(binance.SHORT) && order.Side == string(binance.SideSell)) {
						Logger.Info("手动加仓")
					} else if (order.PositionSide == string(binance.LONG) && order.Side == string(binance.SideSell)) ||
						(order.PositionSide == string(binance.SHORT) && order.Side == string(binance.SideBuy)) {
						Logger.Info("手动减仓")
					}
				} else {
					Logger.Error("异常")
				}

				futureOrder.Symbol = order.Symbol
				futureOrder.OrderID = order.ID
				futureOrder.ClientOrderID = order.ClientOrderID
				futureOrder.Price = util.Round(order.Price, 2)
				futureOrder.OrigQty = util.Round(order.OrigQty, 3)
				futureOrder.AvgPrice = strconv.FormatFloat(order.AvgPrice, 'f', 10, 64)
				futureOrder.ExecutedQty = util.Round(order.ExecutedQty, 3)
				futureOrder.Status = order.OrderStatus
				futureOrder.TimeInForce = binance.TimeInForce(order.TimeInForce)
				futureOrder.Type = binance.OrderType(order.OrderType)
				futureOrder.OrigType = string(order.OrigType)
				futureOrder.Side = binance.OrderSide(order.Side)
				futureOrder.ClosePosition = order.IsClose
				futureOrder.StopPrice = order.StopPrice
				futureOrder.ReduceOnly = order.IsReduce
				futureOrder.PositionSide = order.PositionSide
				futureOrder.Time = order.Time       //新订单是否有值？
				futureOrder.UpdateTime = order.Time //新订单是否有值？
				Logger.Debug("", zap.Any(s.Symbol, futureOrder))
				Logger.Sugar().Infof("价格 : %v 数量 : %v 买卖方向 : %v 持仓方向 : %v 类型 : %v", order.Price, order.OrigQty, side, positionSide, orderFlag)
				switch order.NewEvent {
				case binance.EventNew: //新挂单
					{
						Logger.Info("新挂单", zap.Any(s.Symbol, order))
						s.SaveFutureOrder(futureOrder, order.ClientOrderID)
					}

				case binance.EventCanceled: //挂单取消
					{
						Logger.Info("挂单取消", zap.Any(s.Symbol, order))
						s.DelFutureOrder(futureOrder, order.ClientOrderID)
						if futureOrder.ExecutedQty != 0 &&
							((futureOrder.PositionSide == string(binance.LONG) && futureOrder.Side == binance.SideBuy) ||
								(futureOrder.PositionSide == string(binance.SHORT) && futureOrder.Side == binance.SideSell)) {
							//部分成交的加仓挂单创建对应平仓单
							s.MakePlaceOrder(futureOrder)
							s.MakeCloseOrder(futureOrder)
						}
					}
				case binance.EventCalCulated: //挂单计算？
					{
						Logger.Info("挂单计算？", zap.Any(s.Symbol, order))
					}
				case binance.EventTrade: //挂单成交
					{
						Logger.Info("挂单成交", zap.Any(s.Symbol, order))
						switch order.OrderStatus {
						case binance.StatusNew:
							{

							}
						case binance.StatusPartiallyFilled:
							{
								s.SaveFutureOrder(futureOrder, order.ClientOrderID)
								// util.SendOrderMsg("")
							}
						case binance.StatusFilled:
							{

								s.DelFutureOrder(futureOrder, order.ClientOrderID)
								//创建平仓单
								s.MakePlaceOrder(futureOrder)
								s.MakeCloseOrder(futureOrder)
								fx := "开仓"
								var f1, f2, f3 float64
								if futureOrder.PositionSide == string(binance.LONG) {
									if futureOrder.Side == binance.SideSell {
										fx = "平仓"
									}
									f1, f2, f3 = s.GetLongBetweenAllCloseFutureOrderAndPositionD_Value()
								}
								if futureOrder.PositionSide == string(binance.SHORT) {
									if futureOrder.Side == binance.SideBuy {
										fx = "平仓"
									}
									f1, f2, f3 = s.GetShortBetweenAllCloseFutureOrderAndPositionD_Value()
								}
								msg := fmt.Sprintf("订单类型 : %s \n订单品种 : %s  \n订单方向 : %s  \n成交价格 : %f  \n成交数量 :  %f \n盈利 : %f \n仓位 : %f \n所有平仓挂单的仓位 : %f \n当前持仓价格 : %f \n",
									fx, "ETHUSDT", futureOrder.PositionSide, futureOrder.Price, futureOrder.OrigQty, order.Profit, f1, f2, f3)
								util.SendOrderMsg(msg)
							}
						case binance.StatusCancelled:
							{
								s.DelFutureOrder(futureOrder, order.ClientOrderID)
							}
						case binance.StatusExpired:
							{
								s.DelFutureOrder(futureOrder, order.ClientOrderID)
							}
						case binance.StatusInsurance:
							{
								s.DelFutureOrder(futureOrder, order.ClientOrderID)
							}
						case binance.StatusADL:
							{
								s.DelFutureOrder(futureOrder, order.ClientOrderID)
							}
						default:
							{
								Logger.Sugar().Errorf("未知订单状态 : %v", order.OrderStatus)
							}
						}
					}
				case binance.EventExpired: //挂单过期
					{
						Logger.Info("挂单过期", zap.Any(s.Symbol, order))
						s.DelFutureOrder(futureOrder, order.ClientOrderID)
						if futureOrder.ExecutedQty != 0 &&
							((futureOrder.PositionSide == string(binance.LONG) && futureOrder.Side == binance.SideBuy) ||
								(futureOrder.PositionSide == string(binance.SHORT) && futureOrder.Side == binance.SideSell)) {
							//部分成交的加仓挂单创建对应平仓单
							s.MakePlaceOrder(futureOrder)
							s.MakeCloseOrder(futureOrder)
						}
					}
				default:
					{
						Logger.Info("挂单未知事件类型", zap.Any(s.Symbol, order))
						Logger.Sugar().Errorf("未知事件类型 : %v", order.NewEvent)
					}
				}
			}
		}
	}
}

func (s *Strategy) InitStrategy(symbol string) error {
	//初始化结构
	{
		s.Symbol = symbol
		s.FutureOrder = make(map[string]*MyFutureOrder)
		s.LongPosition.RWMutex = &sync.RWMutex{}
		s.ShortPosition.RWMutex = &sync.RWMutex{}
		s.LongPosition.PinFutureOrder = make(map[string]*MyFutureOrder)
		s.ShortPosition.PinFutureOrder = make(map[string]*MyFutureOrder)
		s.LongPosition.CloseAllFutureOrder = make(map[string]*MyFutureOrder)
		s.ShortPosition.CloseAllFutureOrder = make(map[string]*MyFutureOrder)
		s.LongPosition.CloseFutureOrder = make(map[string]*MyFutureOrder)
		s.ShortPosition.CloseFutureOrder = make(map[string]*MyFutureOrder)
		s.PlaceOrderManager = &PlaceOrderManager{
			RWMutex:       &sync.RWMutex{},
			Symbol:        symbol,
			Quantity:      util.Quantity,
			OrderType:     make(map[string]*MyFutureOrder),
			TerracedPrice: util.TerracedPrice,
			positionInfo:  s,
		}
		s.KlineManager = &Market{
			MinuteKlineList:        NewQueue(60),     //一小时
			FifteenMinuteKlineList: NewQueue(16),     //4小时
			OneHourKlineList:       NewQueue(24),     //一天
			FourHourKlineList:      NewQueue(6 * 10), //10天
			DayKlineList:           NewQueue(30),     //30天
		}
	}
	//初始化账户
	acc := &BinanceFutureAsset{RWMutex: &sync.RWMutex{}}
	acc.InitAccount(symbol)
	s.PlaceOrderManager.Account = acc

	//加载当前持仓
	s.LoadPosition()
	// s.ReloadPosition()

	//加载当前挂单
	s.LoadAllOpenOrder()

	//启动所有定时任务
	s.ClearPartiallyFilledOrder()
	s.ScanCloseFutureOrder()
	s.ScanPositionAndCreatCloseFutureOrder()
	s.ScanFutureOrder()

	//初始化K线事件
	s.KlineWs = Binance.GetKlineWs(util.ETHUSDT, binance.Minute)
	// s.Ch15Kline = Binance.GetKlineWs(util.ETHUSDT, binance.FiveMinutes)

	//初始化账户事件
	s.AccWs, _ = Binance.GetAccountWs()

	//初始化K线
	s.KlineManager.InitMarket(symbol)

	return nil
}
