package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	. "tinyquant/src/logger"
	"tinyquant/src/util"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (as *apiService) FutureDepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@depth@100ms", strings.ToLower(dwr.Symbol))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	dech := make(chan *DepthEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				Logger.Error("websocket recived depth failed")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					c = ReConnectWebSocket(url)
					Logger.Error("trade websocket read failed ", zap.Error(err))
					continue
				}
				rawDepth := struct {
					Type          string          `json:"e"`
					Time          float64         `json:"E"`
					EventTime     float64         `json:"T"`
					Symbol        string          `json:"s"`
					LastUID       int             `json:"U"`
					UpdateID      int             `json:"u"`
					BeforeUID     int             `json:"pu"`
					BidDepthDelta [][]interface{} `json:"b"`
					AskDepthDelta [][]interface{} `json:"a"`
				}{}
				if err := json.Unmarshal(message, &rawDepth); err != nil {
					Logger.Error("depth wsUnmarshal failed ", zap.Error(err))
					return
				}
				t, _ := timeFromUnixTimestampFloat(rawDepth.Time)

				et, _ := timeFromUnixTimestampFloat(rawDepth.EventTime)
				de := &DepthEvent{
					WSEvent: WSEvent{
						Type:   rawDepth.Type,
						Time:   t,
						Symbol: rawDepth.Symbol,
					},
					EventTime: et,
				}

				de.BeforeUID = rawDepth.BeforeUID
				de.LastUpdateID = rawDepth.LastUID
				de.UpdateID = rawDepth.UpdateID

				for _, b := range rawDepth.BidDepthDelta {
					p, _ := floatFromString(b[0])

					q, _ := floatFromString(b[1])

					de.Bids = append(de.Bids, &Order{
						Price:    p,
						Quantity: q,
					})
				}

				for _, b := range rawDepth.AskDepthDelta {
					p, _ := floatFromString(b[0])

					q, _ := floatFromString(b[1])

					de.Asks = append(de.Asks, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				dech <- de
			}
		}
	}()

	go as.exitHandler(c, done, url)
	return dech, done, nil
}

func (as *apiService) FutureKlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@kline_%s", strings.ToLower(kwr.Symbol), string(kwr.Interval))

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	kech := make(chan *KlineEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				Logger.Error(" kline websocket connect closed")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					c = ReConnectWebSocket(url)
					Logger.Error("kline websocket read failed", zap.Error(err))
					continue
				}

				if strings.Contains(string(message), "result") {
					continue
				}

				rawKline := struct {
					Type     string  `json:"e"`
					Time     float64 `json:"E"`
					Symbol   string  `json:"S"`
					OpenTime float64 `json:"t"`
					Kline    struct {
						Interval                 string  `json:"i"`
						FirstTradeID             int64   `json:"f"` // 这根K线期间第一笔成交ID
						LastTradeID              int64   `json:"L"` // 这根K线期间末一笔成交ID
						Final                    bool    `json:"x"` // 这根K线是否完结(是否已经开始下一根K线)
						OpenTime                 float64 `json:"t"`
						CloseTime                float64 `json:"T"`
						Open                     string  `json:"o"`
						High                     string  `json:"h"`
						Low                      string  `json:"l"`
						Close                    string  `json:"c"`
						Volume                   string  `json:"v"` // 这根K线期间成交量
						NumberOfTrades           int     `json:"n"` // 这根K线期间成交笔数
						QuoteAssetVolume         string  `json:"q"` // 这根K线期间成交额
						TakerBuyBaseAssetVolume  string  `json:"V"` // 主动买入的成交量
						TakerBuyQuoteAssetVolume string  `json:"Q"` // 主动买入的成交额
					} `json:"k"`
				}{}
				if err := json.Unmarshal(message, &rawKline); err != nil {
					Logger.Error("kline wsUnmarshal failed ", zap.Error(err))
					return
				}
				t, _ := timeFromUnixTimestampFloat(rawKline.Time)

				ot, _ := timeFromUnixTimestampFloat(rawKline.Kline.OpenTime)

				ct, _ := timeFromUnixTimestampFloat(rawKline.Kline.CloseTime)

				open, _ := floatFromString(rawKline.Kline.Open)

				cls, _ := floatFromString(rawKline.Kline.Close)

				high, _ := floatFromString(rawKline.Kline.High)

				low, _ := floatFromString(rawKline.Kline.Low)

				vol, _ := floatFromString(rawKline.Kline.Volume)

				qav, _ := floatFromString(rawKline.Kline.QuoteAssetVolume)

				tbbav, _ := floatFromString(rawKline.Kline.TakerBuyBaseAssetVolume)

				tbqav, _ := floatFromString(rawKline.Kline.TakerBuyQuoteAssetVolume)

				ke := &KlineEvent{
					WSEvent: WSEvent{
						Type:   rawKline.Type,
						Time:   t,
						Symbol: rawKline.Symbol,
					},
					Interval:     Interval(rawKline.Kline.Interval),
					FirstTradeID: rawKline.Kline.FirstTradeID,
					LastTradeID:  rawKline.Kline.LastTradeID,
					Final:        rawKline.Kline.Final,
					Kline: Kline{
						OpenTime:                 ot,
						CloseTime:                ct,
						Open:                     open,
						Close:                    cls,
						High:                     high,
						Low:                      low,
						Volume:                   vol,
						NumberOfTrades:           rawKline.Kline.NumberOfTrades,
						QuoteAssetVolume:         qav,
						TakerBuyBaseAssetVolume:  tbbav,
						TakerBuyQuoteAssetVolume: tbqav,
						Final:                    rawKline.Kline.Final,
					},
				}
				kech <- ke
			}
		}
	}()

	go as.exitHandler(c, done, url)
	return kech, done, nil
}

//
func (as *apiService) FutureTradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {

	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@aggTrade", strings.ToLower(twr.Symbol))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	aggtech := make(chan *AggTradeEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				Logger.Error("trade websocket connect closed")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					c = ReConnectWebSocket(url)
					Logger.Error("trade websocket read failed ", zap.Error(err))
					continue
				}
				rawAggTrade := struct {
					Type         string  `json:"e"`
					Time         float64 `json:"E"`
					Symbol       string  `json:"s"`
					TradeID      int     `json:"a"`
					Price        string  `json:"p"`
					Quantity     string  `json:"q"`
					FirstTradeID int     `json:"f"`
					LastTradeID  int     `json:"l"`
					Timestamp    float64 `json:"T"`
					IsMaker      bool    `json:"m"`
				}{}
				if err := json.Unmarshal(message, &rawAggTrade); err != nil {
					Logger.Error("trade wsUnmarshal failed ", zap.Error(err))
					return
				}
				t, _ := timeFromUnixTimestampFloat(rawAggTrade.Time)

				price, _ := floatFromString(rawAggTrade.Price)

				qty, _ := floatFromString(rawAggTrade.Quantity)

				ts, _ := timeFromUnixTimestampFloat(rawAggTrade.Timestamp)

				ae := &AggTradeEvent{
					WSEvent: WSEvent{
						Type:   rawAggTrade.Type,
						Time:   t,
						Symbol: rawAggTrade.Symbol,
					},
					AggTrade: AggTrade{
						ID:           rawAggTrade.TradeID,
						Price:        price,
						Quantity:     qty,
						FirstTradeID: rawAggTrade.FirstTradeID,
						LastTradeID:  rawAggTrade.LastTradeID,
						Timestamp:    ts,
						BuyerMaker:   rawAggTrade.IsMaker,
					},
				}
				aggtech <- ae
			}
		}
	}()

	go as.exitHandler(c, done, url)
	return aggtech, done, nil
}

func (as *apiService) FutureUserDataWebsocket(urwr UserDataWebsocketRequest) (chan *FutureAccountEvent, chan struct{}, error) {

	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s", urwr.ListenKey)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	done := make(chan struct{})
	aech := make(chan *FutureAccountEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				Logger.Error("user future data websocket connect close ")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					c = ReConnectWebSocket(url)
					Logger.Error("user data websocket read failed ", zap.Error(err))
					continue
				}

				if strings.Contains(string(message), util.ListenKeyExpired) { // listenkey 过期
					Logger.Info("websocket receive event ")
					stream, _ := as.StartFutureUserDataStream()
					url := fmt.Sprintf("wss://fstream.binance.com/ws/%s", stream.ListenKey)
					c = ReConnectWebSocket(url)
					continue
				}

				if strings.Contains(string(message), util.MARGIN_CALL) { // 追加保证金
					continue
				}

				if strings.Contains(string(message), util.ACCOUNT_UPDATE) { // 账户更新

					accUp := struct {
						Type      string  `json:"e"` // 事件类型
						EventTime float64 `json:"E"` // 事件时间
						Time      float64 `json:"T"` // 撮合时间
						Acc       struct {
							Event   string `json:"m"`
							Balance []struct {
								Symbol        string `json:"a"`
								WalletBalance string `json:"wd"`
								CurBalance    string `json:"cw"`
								BalanceChange string `json:"bc"`
							} `json:"B"`
							Property []struct {
								Symbol string `json:"s"`
								Pa     string `json:"pa"`
								EP     string `json:"ep"`
								CR     string `json:"cr"`
								UP     string `json:"up"`
								MT     string `json:"mt"`
								IW     string `json:"iw"`
								PS     string `json:"ps"`
							} `json:"P"`
						} `json:"a"`
					}{}
					if err := json.Unmarshal(message, &accUp); err != nil {
						Logger.Error("user acc data wsUnmarshal failed ", zap.Error(err))
						return
					}

					ae := &FutureAccountEvent{
						EventName: util.ACCOUNT_UPDATE,
						AE: &AccEvent{
							Type:      accUp.Type,
							EventTime: accUp.EventTime,
							Time:      accUp.Time,
						},
					}
					for _, v := range accUp.Acc.Balance {

						w, _ := floatFromString(v.WalletBalance)
						c, _ := floatFromString(v.CurBalance)
						b, _ := floatFromString(v.BalanceChange)

						ae.AE.Acc.Balance = append(ae.AE.Acc.Balance, struct {
							Symbol        string
							WalletBalance float64
							CurBalance    float64
							BalanceChange float64
						}{
							v.Symbol, w, c, b,
						})
					}

					for _, v := range accUp.Acc.Property {

						pa, _ := floatFromString(v.Pa)
						ep, _ := floatFromString(v.EP)
						cr, _ := floatFromString(v.CR)
						up, _ := floatFromString(v.UP)
						iw, _ := floatFromString(v.IW)

						ae.AE.Acc.Property = append(ae.AE.Acc.Property, struct {
							Symbol string
							Pa     float64
							EP     float64
							CR     float64
							UP     float64
							MT     string
							IW     float64
							PS     string
						}{
							v.Symbol, pa, ep, cr, up, v.MT, iw, v.PS,
						})
					}

					aech <- ae
				}

				if strings.Contains(string(message), util.ORDER_TRADE_UPDATE) { // 交易订单更新

					orderUp := struct {
						Type      string  `json:"e"` // 事件类型
						EventTime float64 `json:"E"` // 事件时间
						Time      float64 `json:"T"` // 撮合时间
						Order     struct {
							Symbol        string  `json:"s"`  // 交易对
							ClientOrderID string  `json:"c"`  // 客户端自定订单ID
							Side          string  `json:"S"`  // 订单方向
							OrderType     string  `json:"o"`  // 订单类型
							TimeInForce   string  `json:"f"`  // 有效方式
							OrigQty       string  `json:"q"`  // 订单原始数量
							Price         string  `json:"p"`  // 订单原始价格
							AvgPrice      string  `json:"ap"` // 订单平均价格
							StopPrice     string  `json:"sp"` // 条件订单触发价格，对追踪止损单无效
							NewEvent      string  `json:"x"`  // 本次事件的具体执行类型
							OrderStatus   string  `json:"X"`  // 订单的当前状态
							ID            int64   `json:"i"`  // 订单ID
							LastQty       string  `json:"l"`  // 订单末次成交量
							ExecutedQty   string  `json:"z"`  // 订单累计已成交量
							LastPrice     string  `json:"L"`  // 订单末次成交价格
							RateAssetType string  `json:"N"`  // 手续费资产类型
							RateQ         string  `json:"n"`  // 手续费数量
							Time          float64 `json:"T"`  // 成交时间
							TimeID        int     `json:"t"`  // 成交ID
							BuyEquity     string  `json:"b"`  // 买单净值
							SellEquity    string  `json:"a"`  // 卖单净值
							IsTaker       bool    `json:"m"`  // 该成交是作为挂单成交吗？
							IsReduce      bool    `json:"R"`  // 是否是只减仓单
							NowType       string  `json:"wt"` // 触发价类型
							OrigType      string  `json:"ot"` // 原始订单类型
							PositionSide  string  `json:"ps"` // 持仓方向
							IsClose       bool    `json:"cp"` // 是否为触发平仓单
							Profit        string  `json:"rp"` // 该交易实现盈亏
						} `json:"o"`
					}{}
					if err := json.Unmarshal(message, &orderUp); err != nil {
						Logger.Error("user data wsUnmarshal failed ", zap.Error(err))
						return
					}

					oe := &FutureAccountEvent{
						EventName: util.ORDER_TRADE_UPDATE,
						OE: &OrderEvent{
							Type:      orderUp.Type,
							EventTime: orderUp.EventTime,
							Time:      orderUp.Time,
						},
					}
					or, _ := floatFromString(orderUp.Order.OrigQty)
					pr, _ := floatFromString(orderUp.Order.Price)
					av, _ := floatFromString(orderUp.Order.AvgPrice)
					sp, _ := floatFromString(orderUp.Order.StopPrice)
					lq, _ := floatFromString(orderUp.Order.LastQty)
					eq, _ := floatFromString(orderUp.Order.ExecutedQty)
					lp, _ := floatFromString(orderUp.Order.LastPrice)
					rq, _ := floatFromString(orderUp.Order.RateQ)
					be, _ := floatFromString(orderUp.Order.BuyEquity)
					se, _ := floatFromString(orderUp.Order.SellEquity)
					profit, _ := floatFromString(orderUp.Order.Profit)

					t, _ := timeFromUnixTimestampFloat(orderUp.Order.Time)
					oe.OE.Order.Symbol = orderUp.Order.Symbol
					oe.OE.Order.ClientOrderID = orderUp.Order.ClientOrderID
					oe.OE.Order.Side = orderUp.Order.Side
					oe.OE.Order.OrderType = orderUp.Order.OrderType
					oe.OE.Order.TimeInForce = orderUp.Order.TimeInForce
					oe.OE.Order.OrigQty = or
					oe.OE.Order.Price = pr
					oe.OE.Order.AvgPrice = av
					oe.OE.Order.StopPrice = sp
					oe.OE.Order.LastQty = lq
					oe.OE.Order.ExecutedQty = eq
					oe.OE.Order.LastPrice = lp
					oe.OE.Order.RateQ = rq
					oe.OE.Order.BuyEquity = be
					oe.OE.Order.SellEquity = se
					oe.OE.Order.Profit = profit
					oe.OE.Order.NewEvent = EventType(orderUp.Order.NewEvent)
					oe.OE.Order.OrderStatus = OrderStatus(orderUp.Order.OrderStatus)
					oe.OE.Order.ID = int64(orderUp.Order.ID)
					oe.OE.Order.RateAssetType = orderUp.Order.RateAssetType
					oe.OE.Order.Time = t
					oe.OE.Order.IsTaker = orderUp.Order.IsTaker
					oe.OE.Order.IsReduce = orderUp.Order.IsReduce
					oe.OE.Order.IsClose = orderUp.Order.IsClose
					oe.OE.Order.NowType = OrderType(orderUp.Order.NowType)
					oe.OE.Order.OrderType = orderUp.Order.OrigType
					oe.OE.Order.PositionSide = orderUp.Order.PositionSide
					aech <- oe

				}

				if strings.Contains(string(message), util.ACCOUNT_CONFIG_UPDATE) { // 杠杆倍数 等配置更新
					continue
				}

			}
		}
	}()

	go as.future_exitHandler(c, done, urwr.ListenKey)
	return aech, done, nil
}

func (as *apiService) future_exitHandler(c *websocket.Conn, done chan struct{}, listenKey string) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	defer c.Close()

	r := req{
		Method: "SET_PROPERTY",
		Id:     5,
	}
	r.Params = append(r.Params, "combined", false)

	data, err := json.Marshal(r)
	if err != nil {
		Logger.Error("write websocket message marshal failed ", zap.Error(err))
		return
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		Logger.Error("websocket write message failed ", zap.Error(err))
		return
	}

	for {
		select {
		case _ = <-ticker.C:

			err := c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				url := fmt.Sprintf("wss://fstream.binance.com/ws/%s", listenKey)
				c = ReConnectWebSocket(url)
				Logger.Error("websocket write message failed ", zap.Error(err))
				return
			}

			as.KeepAliveFutureUserDataStream(&Stream{
				ListenKey: listenKey,
			})
		case <-as.Ctx.Done():
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			Logger.Info("websocket close connect ")
			return
		}
	}
}
