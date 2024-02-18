package strategy

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	. "tinyquant/src/logger"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

type Market struct {
	Last24HHigh float64
	Las24HLow   float64

	MinuteKlineList        *MyKlineQueue // 近60条k线, 分钟线
	FifteenMinuteKlineList *MyKlineQueue
	OneHourKlineList       *MyKlineQueue
	FourHourKlineList      *MyKlineQueue
	DayKlineList           *MyKlineQueue
}

type Kline struct {
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
	CloseTime time.Time
	BuyVolume float64
}

type UpDownLink struct {
	UpLink             float64
	DownLink           float64
	MaxHigh            float64
	MinLow             float64
	MaxHighIndex       int
	MinLowIndex        int
	AvgVolume          float64
	MaxVolume          float64
	AvgPrice           float64
	HalfSampleAvgPrice float64
}

func (m *Market) InitMarket(symbol string) error {
	m.MinuteKlineList.Clear()
	m.FifteenMinuteKlineList.Clear()
	// m.OneHourKlineList.Clear()
	m.FourHourKlineList.Clear()
	// m.DayKlineList.Clear()

	//	Logger.Info("初始化 k 线 ")
	// 获取 近5天日线
	// daylist, err := Binance.GetFutureKlines(symbol, 30, binance.Day)
	// if err != nil {
	// 	Logger.Error("Get day future kline failed", zap.Error(err))
	// 	return err
	// }

	// m.Las24HLow = daylist[29].Low
	// m.Last24HHigh = daylist[29].High

	FourHourlist, err := Binance.GetFutureKlines(symbol, m.FourHourKlineList.Capacity, binance.FourHours) // 30天
	if err != nil {
		Logger.Error("Get 4/6*30 hour future kline failed", zap.Error(err))
		return err
	}

	// OneHourList, err := Binance.GetFutureKlines(symbol, m.OneHourKlineList.Capacity, binance.Hour) //24小时
	// if err != nil {
	// 	Logger.Error("Get 1/24 hour future kline failed", zap.Error(err))
	// 	return err
	// }

	FifteenMinuteKlineList, err := Binance.GetFutureKlines(symbol, m.FifteenMinuteKlineList.Capacity, binance.FifteenMinutes) // 15分钟 4小时
	if err != nil {
		Logger.Error("Get 15/4*4 minute future kline failed", zap.Error(err))
		return err
	}

	Minutelist, err := Binance.GetFutureKlines(symbol, m.MinuteKlineList.Capacity, binance.Minute) //分钟 一小时
	if err != nil {
		Logger.Error("Get 1/60 minute future kline failed", zap.Error(err))
		return err
	}

	// for _, v := range daylist {
	// 	m.DayKlineList.EnQqueu(&Kline{
	// 		Open:      v.Open,
	// 		Close:     v.Close,
	// 		High:      v.High,
	// 		Low:       v.Low,
	// 		Volume:    v.Volume,
	// 		BuyVolume: v.TakerBuyBaseAssetVolume,
	// 	})
	// }

	for _, v := range FourHourlist {
		m.FourHourKlineList.EnQqueu(&Kline{
			Open:      v.Open,
			Close:     v.Close,
			High:      v.High,
			Low:       v.Low,
			Volume:    v.Volume,
			BuyVolume: v.TakerBuyBaseAssetVolume,
			CloseTime: v.CloseTime,
		})
	}

	// for _, v := range OneHourList {
	// 	m.OneHourKlineList.EnQqueu(&Kline{
	// 		Open:      v.Open,
	// 		Close:     v.Close,
	// 		High:      v.High,
	// 		Low:       v.Low,
	// 		Volume:    v.Volume,
	// 		BuyVolume: v.TakerBuyBaseAssetVolume,
	// 	})
	// }

	for _, v := range FifteenMinuteKlineList {
		m.FifteenMinuteKlineList.EnQqueu(&Kline{
			Open:      v.Open,
			Close:     v.Close,
			High:      v.High,
			Low:       v.Low,
			Volume:    v.Volume,
			BuyVolume: v.TakerBuyBaseAssetVolume,
			CloseTime: v.CloseTime,
		})
		Logger.Sugar().Debugf("%+v", v)
	}
	for _, v := range Minutelist {
		m.MinuteKlineList.EnQqueu(&Kline{
			Open:      v.Open,
			Close:     v.Close,
			High:      v.High,
			Low:       v.Low,
			Volume:    v.Volume,
			BuyVolume: v.TakerBuyBaseAssetVolume,
			CloseTime: v.CloseTime,
		})
		Logger.Sugar().Debugf("%+v", v)
	}
	m.MinuteKlineList.UpdateUpDownLink(true)
	m.FifteenMinuteKlineList.UpdateUpDownLink(true)
	// m.OneHourKlineList.UpdateUpDownLink(true)
	m.FourHourKlineList.UpdateUpDownLink(true)
	// m.DayKlineList.UpdateUpDownLink(true)

	return nil
}

func (m *Market) UpdateKlineListTicker(symbol string) {

	ticker := time.NewTicker(1 * time.Minute)

	defer ticker.Stop()

	for {
		<-ticker.C
		//Logger.Info("开始更新 K 线")

		m.InitMarket(symbol)

	}
}

type MyKlineQueue struct {
	Data     []*Kline
	Capacity int
	Head     int //队头
	Full     bool
	*sync.RWMutex
	*UpDownLink
}

func NewQueue(len int) *MyKlineQueue {
	return &MyKlineQueue{Capacity: len, Head: -1, Data: make([]*Kline, len), UpDownLink: &UpDownLink{}, RWMutex: &sync.RWMutex{}}
}

// 插入
func (queue *MyKlineQueue) EnQqueu(val *Kline) {
	queue.Lock()
	defer queue.Unlock()
	if queue.IsEmpty() {
		queue.Head = 0
	} else if queue.Full {
		if queue.Data[queue.Head].CloseTime.Equal(val.CloseTime) {
			queue.Data[queue.Head] = val
			return
		} else if queue.Data[queue.Head].CloseTime.After(val.CloseTime) {
			return
		}
		queue.Head = (queue.Head + 1) % queue.Capacity
	} else {
		queue.Head = queue.Head + 1
		if queue.Head == queue.Capacity-2 {
			queue.Full = true
		}
	}
	queue.Data[queue.Head] = val
}

//判满
func (queue *MyKlineQueue) IsFull() bool {
	return queue.Head+1 == queue.Capacity
}

//判断空
func (queue *MyKlineQueue) IsEmpty() bool {
	return queue.Head == -1
}

func (queue *MyKlineQueue) Clear() {
	queue.Head = -1
}

func (queue *MyKlineQueue) GetUpDownLink() *UpDownLink {
	queue.RLock()
	defer queue.RUnlock()
	return &UpDownLink{
		UpLink:             queue.DownLink,
		DownLink:           queue.DownLink,
		MaxHigh:            queue.MaxHigh,
		MinLow:             queue.MinLow,
		MaxHighIndex:       queue.MaxHighIndex,
		MinLowIndex:        queue.MinLowIndex,
		AvgVolume:          queue.AvgVolume,
		MaxVolume:          queue.MaxVolume,
		AvgPrice:           queue.AvgPrice,
		HalfSampleAvgPrice: queue.HalfSampleAvgPrice,
	}
}

func (queue *MyKlineQueue) UpdateUpDownLink(lastValid bool) {
	if !queue.Full {
		return
	}
	High := 0.0
	Low := 0.0
	max_High := 0.0
	max_High_index := 0
	min_Low := 0.0
	min_Low_index := 0
	avg_High := 0.0
	avg_Low := 0.0

	Volume := 0.0
	HalfSampleVolume := 0.0
	MaxVolume := 0.0

	avgPrice := 0.0

	queue.Lock()
	if lastValid {
		j := queue.Capacity / 2
		for k := 0; k < queue.Capacity; k++ {
			kl := queue.Data[k]
			if kl.Volume > MaxVolume {
				MaxVolume = kl.Volume
			}
			High += kl.High
			Low += kl.Low
			avgPrice += kl.Close

			if kl.High > max_High {
				max_High = kl.High
				max_High_index = k
			}
			if kl.Low < min_Low {
				min_Low = kl.Low
				min_Low_index = k
			}
			Volume += kl.Volume
			if k >= j {
				HalfSampleVolume += kl.Volume
			}
		}
	} else {
		i := 0
		j := queue.Capacity / 2
		for k := queue.Head + 1; k < queue.Capacity; k++ {
			i++
			kl := queue.Data[k]
			if kl.Volume > MaxVolume {
				MaxVolume = kl.Volume
			}
			High += kl.High
			Low += kl.Low
			avgPrice += kl.Close
			if kl.High > max_High {
				max_High = kl.High
				max_High_index = k
			}
			if kl.Low < min_Low {
				min_Low = kl.Low
				min_Low_index = k
			}
			Volume += kl.Volume
			if i >= j {
				HalfSampleVolume += kl.Volume
			}
		}
		for k := 0; k < queue.Head; k++ {
			i++
			kl := queue.Data[k]
			if kl.Volume > MaxVolume {
				MaxVolume = kl.Volume
			}
			High += kl.High
			Low += kl.Low
			avgPrice += kl.Close
			if kl.High > max_High {
				max_High = kl.High
				max_High_index = k
			}
			if kl.Low < min_Low {
				min_Low = kl.Low
				min_Low_index = k
			}
			Volume += kl.Volume
			if i >= j {
				HalfSampleVolume += kl.Volume
			}
		}
	}

	var avgVolume float64
	var halfSampleAvgVolume float64
	if lastValid {
		avg_High = High / float64(queue.Capacity)
		avg_Low = Low / float64(queue.Capacity)
		avgPrice = avgPrice / float64(queue.Capacity)
		avgVolume = Volume / float64(queue.Capacity)
		halfSampleAvgVolume = HalfSampleVolume / float64(queue.Capacity/2)
	} else {
		avg_High = High / float64(queue.Capacity-1)
		avg_Low = Low / float64(queue.Capacity-1)
		avgPrice = avgPrice / float64(queue.Capacity-1)
		avgVolume = Volume / float64(queue.Capacity-1)
		halfSampleAvgVolume = HalfSampleVolume / float64((queue.Capacity/2)-1)
	}

	queue.UpLink = (avg_High + max_High) / 2 // 上行线
	queue.UpLink, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", queue.UpLink), 64)

	queue.DownLink = (avg_Low + min_Low) / 2 // 下行线
	queue.DownLink, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", queue.DownLink), 64)
	queue.MaxHigh = max_High
	queue.MinLow = min_Low
	queue.MaxHighIndex = max_High_index
	queue.MinLowIndex = min_Low_index
	queue.AvgVolume = avgVolume
	queue.MaxVolume = MaxVolume
	queue.AvgPrice = avgPrice
	queue.HalfSampleAvgPrice = halfSampleAvgVolume
	queue.Unlock()
}

func (queue *MyKlineQueue) GetNewPrice() float64 {
	if !queue.Full {
		return -1
	}
	queue.RLock()
	defer queue.RUnlock()
	return queue.Data[queue.Head].Close
}
