package strategy

import (
	"sync"
	"time"
	. "tinyquant/src/logger"

	"go.uber.org/zap"
)

var OB *OrderBook

type OrderBook struct {
	LastUpdateID   int
	MessageTime    time.Time
	ChangeInfoTime time.Time
	Bids           OrderSlice // 买方出价
	Asks           OrderSlice // 卖方出价
}

type Order struct {
	Price    float64
	Quantity float64
}

type OrderBookMap struct {
	LastUpdateID   int
	MessageTime    time.Time
	ChangeInfoTime time.Time
	Bids           map[float64]float64
	Asks           map[float64]float64
	mutex          sync.RWMutex
}

func (o *OrderBookMap) InitOrderBook(symbol string, limit int) error {

	ob, err := Binance.GetDepth(symbol, limit)
	if err != nil {
		Logger.Error("Get Depth failed ", zap.Error(err))
		return err
	}

	o.LastUpdateID = ob.LastUpdateID
	o.MessageTime = ob.MessageTime

	o.mutex.Lock()
	defer o.mutex.Unlock()

	for _, v := range ob.Asks {
		o.Asks[v.Price] = v.Quantity
	}

	for _, v := range ob.Bids {
		o.Bids[v.Price] = v.Quantity
	}

	return nil
	/*obTmp := &OrderBook{}

	obTmp.LastUpdateID = ob.LastUpdateID
	obTmp.MessageTime = ob.MessageTime

	for _, v := range ob.Asks {
		obTmp.Asks = append(obTmp.Asks, Order{
			Price:    v.Price,
			Quantity: v.Quantity,
		})
	}

	for _, v := range ob.Bids {
		obTmp.Bids = append(obTmp.Bids, Order{
			Price:    v.Price,
			Quantity: v.Quantity,
		})
	}

	sort.Sort(OrderSlice(obTmp.Asks))
	sort.Sort(OrderSlice(obTmp.Bids))

	OB = obTmp*/
}

func DepthMaptoSlice() {

}

func (o *OrderBookMap) DepthUpdate(symbol string) {

	depth_chan, depth_done := Binance.GetFutureDepthWs(symbol)

	o.Asks = map[float64]float64{}
	o.Bids = map[float64]float64{}

	update_flag := false

	go func() {
		time.Sleep(1000 * time.Millisecond)
		o.InitOrderBook(symbol, 1000)
	}()

	for {
		select {

		case depth := <-depth_chan:

			if update_flag {

				o.mutex.Lock()
				for _, v := range depth.Asks {
					if v.Quantity == 0.0 {
						delete(o.Asks, v.Price)
					}

					o.Asks[v.Price] = v.Quantity

				}

				for _, v := range depth.Bids {
					if v.Quantity == 0.0 {
						delete(o.Bids, v.Price)
					}
					o.Bids[v.Price] = v.Quantity

				}
				o.mutex.Unlock()

				continue
			}
			if depth.LastUpdateID == 0 || depth.UpdateID == 0 || o.LastUpdateID == 0 {
				continue
			}
			if depth.LastUpdateID <= o.LastUpdateID && depth.UpdateID >= o.LastUpdateID {
				update_flag = true
			}

		case <-depth_done:

		default:
			break
		}
	}

}

func (o *OrderBookMap) GetOrderBookRate() float64 {

	o.mutex.RLock()
	defer o.mutex.RUnlock()

	var ask float64
	var bid float64

	for _, v := range o.Asks {
		ask += v
	}
	for _, v := range o.Bids {
		bid += v
	}

	return ask / bid

}

type OrderSlice []Order

func (a OrderSlice) Len() int {
	return len(a)
}
func (a OrderSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a OrderSlice) Less(i, j int) bool {
	return a[i].Price < a[j].Price
}
