package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	. "tinyquant/src/logger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (as *apiService) DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://fstream.binance.com/ws/%s@depth", strings.ToLower(dwr.Symbol))
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
					Logger.Error("read message failed", zap.Error(err))
					return
				}
				rawDepth := struct {
					Type          string          `json:"e"`
					Time          float64         `json:"E"`
					Symbol        string          `json:"s"`
					UpdateID      int             `json:"u"`
					BidDepthDelta [][]interface{} `json:"b"`
					AskDepthDelta [][]interface{} `json:"a"`
				}{}
				if err := json.Unmarshal(message, &rawDepth); err != nil {
					Logger.Error("depth wsUnmarshal failed ", zap.Error(err))
					return
				}
				t, _ := timeFromUnixTimestampFloat(rawDepth.Time)

				de := &DepthEvent{
					WSEvent: WSEvent{
						Type:   rawDepth.Type,
						Time:   t,
						Symbol: rawDepth.Symbol,
					},
				}
				for _, b := range rawDepth.BidDepthDelta {
					p, _ := floatFromString(b[0])

					q, _ := floatFromString(b[1])

					de.Bids = append(de.Bids, &Order{
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

func (as *apiService) KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
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
					Logger.Error("kline websocket read failed", zap.Error(err))
					return
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
func (as *apiService) TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
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
					Logger.Error("trade websocket read failed ", zap.Error(err))
					return
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

func (as *apiService) UserDataWebsocket(urwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error) {
	strUrl := fmt.Sprintf("wss://fstream.binance.com/ws/%s", urwr.ListenKey)

	c, _, err := websocket.DefaultDialer.Dial(strUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	aech := make(chan *AccountEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				Logger.Error("user data websocket connect close ")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					Logger.Error("user data websocket read failed ", zap.Error(err))
					return
				}

				rawAccount := struct {
					Type            string  `json:"e"`
					Time            float64 `json:"E"`
					OpenTime        float64 `json:"t"`
					MakerCommision  int64   `json:"m"`
					TakerCommision  int64   `json:"t"`
					BuyerCommision  int64   `json:"b"`
					SellerCommision int64   `json:"s"`
					CanTrade        bool    `json:"T"`
					CanWithdraw     bool    `json:"W"`
					CanDeposit      bool    `json:"D"`
					Balances        []struct {
						Asset            string `json:"a"`
						AvailableBalance string `json:"f"`
						Locked           string `json:"l"`
					} `json:"B"`
				}{}
				if err := json.Unmarshal(message, &rawAccount); err != nil {
					Logger.Error("user data wsUnmarshal failed ", zap.Error(err))
					return
				}
				t, _ := timeFromUnixTimestampFloat(rawAccount.Time)

				ae := &AccountEvent{
					WSEvent: WSEvent{
						Type: rawAccount.Type,
						Time: t,
					},
					Account: Account{
						MakerCommision:  rawAccount.MakerCommision,
						TakerCommision:  rawAccount.TakerCommision,
						BuyerCommision:  rawAccount.BuyerCommision,
						SellerCommision: rawAccount.SellerCommision,
						CanTrade:        rawAccount.CanTrade,
						CanWithdraw:     rawAccount.CanWithdraw,
						CanDeposit:      rawAccount.CanDeposit,
					},
				}
				for _, b := range rawAccount.Balances {
					free, _ := floatFromString(b.AvailableBalance)

					locked, _ := floatFromString(b.Locked)

					ae.Balances = append(ae.Balances, &Balance{
						Asset:  b.Asset,
						Free:   free,
						Locked: locked,
					})
				}
				aech <- ae
			}
		}
	}()

	go as.exitHandler(c, done, strUrl)
	return aech, done, nil
}

type req struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	Id     int           `json:"id"`
}

func (as *apiService) exitHandler(c *websocket.Conn, done chan struct{}, url string) {
	ticker := time.NewTicker(10 * time.Minute)
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
		c = ReConnectWebSocket(url)
		Logger.Error("websocket write message failed ", zap.Error(err))
		return
	}

	for {
		select {
		case _ = <-ticker.C:

			err = c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				c = ReConnectWebSocket(url)
				Logger.Error("websocket write message failed ", zap.Error(err))
				continue
			}
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
