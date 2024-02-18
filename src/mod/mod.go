package mod

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ReqParam struct {
	URL    string
	Method string
	Query  url.Values
	Body   url.Values
	Header http.Header
	APIKEY string
}

func (r *ReqParam) SetParam(key string, value interface{}) *ReqParam {
	if r.Query == nil {
		r.Query = url.Values{}
	}
	r.Query.Set(key, fmt.Sprintf("%v", value))
	return r
}

type Kline struct {
	StartTime   time.Time
	CloseTime   time.Time
	Volume      float64
	BuyVolume   float64
	SellVolume  float64
	Quote       float64
	BuyQuote    float64
	SellQuote   float64
	TradeNumber int
	Open        float64
	Close       float64
	High        float64
	Low         float64
	Final       bool
}

type Depth struct {
	LastUpdateID int
	BeforeUID    int
	UpdateID     int
	MessageTime  time.Time
	Bids         []*Order // 买方出价
	Asks         []*Order //卖方出价
}

type Order struct {
	Price    float64
	Quantity float64
}
