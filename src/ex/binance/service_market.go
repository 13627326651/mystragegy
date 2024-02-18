package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func (as *apiService) Ping() error {
	params := make(map[string]string)
	response, err := as.request("GET", "api/v1/ping", params, false, false)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", response.StatusCode)
	return nil
}

func (as *apiService) Time() (time.Time, error) {
	params := make(map[string]string)
	res, err := as.request("GET", "api/v1/time", params, false, false)
	if err != nil {
		return time.Time{}, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to read response from Time")
	}
	defer res.Body.Close()
	var rawTime struct {
		ServerTime string `json:"serverTime"`
	}
	if err := json.Unmarshal(textRes, &rawTime); err != nil {
		return time.Time{}, errors.Wrap(err, "timeResponse unmarshal failed")
	}
	t, err := timeFromUnixTimestampFloat(rawTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (as *apiService) ExchangeInfo() (string, error) {
	params := make(map[string]string)
	params["symbol"] = "ETHUSDT"
	res, err := as.request("GET", "fapi/v1/exchangeInfo", params, false, false)
	if err != nil {
		return "", err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read response from Time")
	}
	defer res.Body.Close()

	return string(textRes), nil
}

func (as *apiService) NewPrice(nb OrderNewPriceRequest) (*NewPrice, error) {

	params := make(map[string]string)
	params["symbol"] = nb.Symbol

	res, err := as.request("GET", "api/v3/ticker/price", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Time")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}
	fmt.Println("xxxxx", string(textRes))
	pr := struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}{}
	if err := json.Unmarshal(textRes, &pr); err != nil {
		return nil, errors.Wrap(err, "new price unmarshal failed")
	}

	pri := &NewPrice{
		Symbol: pr.Symbol,
		Price:  pr.Price,
	}

	return pri, nil
}

func (as *apiService) OrderBook(obr OrderBookRequest) (*OrderBook, error) {
	params := make(map[string]string)
	params["symbol"] = obr.Symbol
	if obr.Limit != 0 {
		params["limit"] = strconv.Itoa(obr.Limit)
	}
	res, err := as.request("GET", "fapi/v1/depth", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Time")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawBook := &struct {
		LastUpdateID   int             `json:"lastUpdateId"`
		MessageTime    float64         `json:"E"`
		ChangeInfoTime float64         `json:"T"`
		Bids           [][]interface{} `json:"bids"`
		Asks           [][]interface{} `json:"asks"`
	}{}
	if err := json.Unmarshal(textRes, rawBook); err != nil {
		return nil, errors.Wrap(err, "timeResponse unmarshal failed")
	}

	m, _ := timeFromUnixTimestampFloat(rawBook.MessageTime)

	ob := &OrderBook{
		LastUpdateID: rawBook.LastUpdateID,
		MessageTime:  m,
	}
	extractOrder := func(rawPrice, rawQuantity interface{}) (*Order, error) {
		price, err := floatFromString(rawPrice)
		if err != nil {
			return nil, err
		}
		quantity, err := floatFromString(rawQuantity)
		if err != nil {
			return nil, err
		}
		return &Order{
			Price:    price,
			Quantity: quantity,
		}, nil
	}
	for _, bid := range rawBook.Bids {
		order, err := extractOrder(bid[0], bid[1])
		if err != nil {
			return nil, err
		}
		ob.Bids = append(ob.Bids, order)
	}
	for _, ask := range rawBook.Asks {
		order, err := extractOrder(ask[0], ask[1])
		if err != nil {
			return nil, err
		}
		ob.Asks = append(ob.Asks, order)
	}

	return ob, nil
}

func (as *apiService) AggTrades(atr AggTradesRequest) ([]*AggTrade, error) {
	params := make(map[string]string)
	params["symbol"] = atr.Symbol
	if atr.FromID != 0 {
		params["fromId"] = strconv.FormatInt(atr.FromID, 10)
	}
	if atr.StartTime != 0 {
		params["startTime"] = strconv.FormatInt(atr.StartTime, 10)
	}
	if atr.EndTime != 0 {
		params["endTime"] = strconv.FormatInt(atr.EndTime, 10)
	}
	if atr.Limit != 0 {
		params["limit"] = strconv.Itoa(atr.Limit)
	}

	res, err := as.request("GET", "api/v1/aggTrades", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from AggTrades")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawAggTrades := []struct {
		ID             int    `json:"a"`
		Price          string `json:"p"`
		Quantity       string `json:"q"`
		FirstTradeID   int    `json:"f"`
		LastTradeID    int    `json:"l"`
		Timestamp      int64  `json:"T"`
		BuyerMaker     bool   `json:"m"`
		BestPriceMatch bool   `json:"M"`
	}{}
	if err := json.Unmarshal(textRes, &rawAggTrades); err != nil {
		return nil, errors.Wrap(err, "aggTrades unmarshal failed")
	}
	aggTrades := []*AggTrade{}
	for _, rawTrade := range rawAggTrades {
		price, err := floatFromString(rawTrade.Price)
		if err != nil {
			return nil, err
		}
		quantity, err := floatFromString(rawTrade.Quantity)
		if err != nil {
			return nil, err
		}
		t := time.Unix(0, rawTrade.Timestamp*int64(time.Millisecond))

		aggTrades = append(aggTrades, &AggTrade{
			ID:             rawTrade.ID,
			Price:          price,
			Quantity:       quantity,
			FirstTradeID:   rawTrade.FirstTradeID,
			LastTradeID:    rawTrade.LastTradeID,
			Timestamp:      t,
			BuyerMaker:     rawTrade.BuyerMaker,
			BestPriceMatch: rawTrade.BestPriceMatch,
		})
	}
	return aggTrades, nil
}

func (as *apiService) Klines(kr KlinesRequest) ([]*Kline, error) {
	params := make(map[string]string)
	params["symbol"] = kr.Symbol
	params["interval"] = string(kr.Interval)
	if kr.Limit != 0 {
		params["limit"] = strconv.Itoa(kr.Limit)
	}
	if kr.StartTime != 0 {
		params["startTime"] = strconv.FormatInt(kr.StartTime, 10)
	}
	if kr.EndTime != 0 {
		params["endTime"] = strconv.FormatInt(kr.EndTime, 10)
	}

	res, err := as.request("GET", "api/v1/klines", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Klines")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawKlines := [][]interface{}{}
	if err := json.Unmarshal(textRes, &rawKlines); err != nil {
		return nil, errors.Wrap(err, "rawKlines unmarshal failed")
	}
	klines := []*Kline{}
	for _, k := range rawKlines {
		ot, err := timeFromUnixTimestampFloat(k[0])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.OpenTime")
		}
		open, err := floatFromString(k[1])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Open")
		}
		high, err := floatFromString(k[2])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.High")
		}
		low, err := floatFromString(k[3])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Low")
		}
		cls, err := floatFromString(k[4])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Close")
		}
		volume, err := floatFromString(k[5])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Volume")
		}
		ct, err := timeFromUnixTimestampFloat(k[6])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.CloseTime")
		}
		qav, err := floatFromString(k[7])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.QuoteAssetVolume")
		}
		not, ok := k[8].(float64)
		if !ok {
			return nil, errors.Wrap(err, "cannot parse Kline.NumberOfTrades")
		}
		tbbav, err := floatFromString(k[9])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.TakerBuyBaseAssetVolume")
		}
		tbqav, err := floatFromString(k[10])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.TakerBuyQuoteAssetVolume")
		}
		klines = append(klines, &Kline{
			OpenTime:                 ot,
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    cls,
			Volume:                   volume,
			CloseTime:                ct,
			QuoteAssetVolume:         qav,
			NumberOfTrades:           int(not),
			TakerBuyBaseAssetVolume:  tbbav,
			TakerBuyQuoteAssetVolume: tbqav,
		})
	}
	return klines, nil
}

func (as *apiService) Ticker24(tr TickerRequest) (*Ticker24, error) {
	params := make(map[string]string)
	params["symbol"] = tr.Symbol

	res, err := as.request("GET", "api/v1/ticker/24hr", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawTicker24 := struct {
		PriceChange        string  `json:"priceChange"`
		PriceChangePercent string  `json:"priceChangePercent"`
		WeightedAvgPrice   string  `json:"weightedAvgPrice"`
		PrevClosePrice     string  `json:"prevClosePrice"`
		LastPrice          string  `json:"lastPrice"`
		BidPrice           string  `json:"bidPrice"`
		AskPrice           string  `json:"askPrice"`
		OpenPrice          string  `json:"openPrice"`
		HighPrice          string  `json:"highPrice"`
		LowPrice           string  `json:"lowPrice"`
		Volume             string  `json:"volume"`
		OpenTime           float64 `json:"openTime"`
		CloseTime          float64 `json:"closeTime"`
		FirstID            int
		LastID             int
		Count              int
	}{}
	if err := json.Unmarshal(textRes, &rawTicker24); err != nil {
		return nil, errors.Wrap(err, "rawTicker24 unmarshal failed")
	}

	pc, err := strconv.ParseFloat(rawTicker24.PriceChange, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.PriceChange")
	}
	pcPercent, err := strconv.ParseFloat(rawTicker24.PriceChangePercent, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.PriceChangePercent")
	}
	wap, err := strconv.ParseFloat(rawTicker24.WeightedAvgPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.WeightedAvgPrice")
	}
	pcp, err := strconv.ParseFloat(rawTicker24.PrevClosePrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.PrevClosePrice")
	}
	lastPrice, err := strconv.ParseFloat(rawTicker24.LastPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.LastPrice")
	}
	bp, err := strconv.ParseFloat(rawTicker24.BidPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.BidPrice")
	}
	ap, err := strconv.ParseFloat(rawTicker24.AskPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.AskPrice")
	}
	op, err := strconv.ParseFloat(rawTicker24.OpenPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.OpenPrice")
	}
	hp, err := strconv.ParseFloat(rawTicker24.HighPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.HighPrice")
	}
	lowPrice, err := strconv.ParseFloat(rawTicker24.LowPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.LowPrice")
	}
	vol, err := strconv.ParseFloat(rawTicker24.Volume, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.Volume")
	}
	ot, err := timeFromUnixTimestampFloat(rawTicker24.OpenTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.OpenTime")
	}
	ct, err := timeFromUnixTimestampFloat(rawTicker24.CloseTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.CloseTime")
	}
	t24 := &Ticker24{
		PriceChange:        pc,
		PriceChangePercent: pcPercent,
		WeightedAvgPrice:   wap,
		PrevClosePrice:     pcp,
		LastPrice:          lastPrice,
		BidPrice:           bp,
		AskPrice:           ap,
		OpenPrice:          op,
		HighPrice:          hp,
		LowPrice:           lowPrice,
		Volume:             vol,
		OpenTime:           ot,
		CloseTime:          ct,
		FirstID:            rawTicker24.FirstID,
		LastID:             rawTicker24.LastID,
		Count:              rawTicker24.Count,
	}
	return t24, nil
}

func (as *apiService) TickerAllPrices() ([]*PriceTicker, error) {
	params := make(map[string]string)

	res, err := as.request("GET", "api/v1/ticker/allPrices", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawTickerAllPrices := []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}{}
	if err := json.Unmarshal(textRes, &rawTickerAllPrices); err != nil {
		return nil, errors.Wrap(err, "rawTickerAllPrices unmarshal failed")
	}

	var tpc []*PriceTicker
	for _, rawTickerPrice := range rawTickerAllPrices {
		p, err := strconv.ParseFloat(rawTickerPrice.Price, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerAllPrices.Price")
		}
		tpc = append(tpc, &PriceTicker{
			Symbol: rawTickerPrice.Symbol,
			Price:  p,
		})
	}
	return tpc, nil
}

func (as *apiService) TickerAllBooks() ([]*BookTicker, error) {
	params := make(map[string]string)

	res, err := as.request("GET", "api/v1/ticker/allBookTickers", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/allBookTickers")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawBookTickers := []struct {
		Symbol   string `json:"symbol"`
		BidPrice string `json:"bidPrice"`
		BidQty   string `json:"bidQty"`
		AskPrice string `json:"askPrice"`
		AskQty   string `json:"askQty"`
	}{}
	if err := json.Unmarshal(textRes, &rawBookTickers); err != nil {
		return nil, errors.Wrap(err, "rawBookTickers unmarshal failed")
	}

	var btc []*BookTicker
	for _, rawBookTicker := range rawBookTickers {
		bp, err := strconv.ParseFloat(rawBookTicker.BidPrice, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.BidPrice")
		}
		bqty, err := strconv.ParseFloat(rawBookTicker.BidQty, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.BidQty")
		}
		ap, err := strconv.ParseFloat(rawBookTicker.AskPrice, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.AskPrice")
		}
		aqty, err := strconv.ParseFloat(rawBookTicker.AskQty, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.AskQty")
		}
		btc = append(btc, &BookTicker{
			Symbol:   rawBookTicker.Symbol,
			BidPrice: bp,
			BidQty:   bqty,
			AskPrice: ap,
			AskQty:   aqty,
		})
	}
	return btc, nil
}

func (as *apiService) PremiumAndFundsRate(pfrr PremiumAndFundsRateRequest) (*PremiumAndFundsRateInfo, error) {

	params := make(map[string]string)
	params["symbol"] = pfrr.Symbol

	res, err := as.request("GET", "/fapi/v1/premiumIndex", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from PremiumAndFundsRate.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawPremiumAndFundsRate := struct {
		Symbol               string  `json:"symbol"`
		MarkPrice            string  `json:"markPrice"`
		IndexPrice           string  `json:"indexPrice"`
		EstimatedSettlePrice string  `json:"estimatedSettlePrice"`
		LastFundingRate      string  `json:"lastFundingRate"`
		NextFundingTime      float64 `json:"nextFundingTime"`
		InterestRate         string  `json:"interestRate"`
		Time                 float64 `json:"time"`
	}{}

	if err := json.Unmarshal(textRes, &rawPremiumAndFundsRate); err != nil {
		return nil, errors.Wrap(err, "PremiumAndFundsRate unmarshal failed")
	}

	mp, _ := floatFromString(rawPremiumAndFundsRate.MarkPrice)
	ip, _ := floatFromString(rawPremiumAndFundsRate.IndexPrice)
	esp, _ := floatFromString(rawPremiumAndFundsRate.EstimatedSettlePrice)
	lfr, _ := floatFromString(rawPremiumAndFundsRate.LastFundingRate)
	ir, _ := floatFromString(rawPremiumAndFundsRate.InterestRate)
	t, err := timeFromUnixTimestampFloat(rawPremiumAndFundsRate.Time)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Time")
	}
	t1, err := timeFromUnixTimestampFloat(rawPremiumAndFundsRate.NextFundingTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Time")
	}

	pf := &PremiumAndFundsRateInfo{
		Symbol:               rawPremiumAndFundsRate.Symbol,
		MarkPrice:            mp,
		IndexPrice:           ip,
		EstimatedSettlePrice: esp,
		LastFundingRate:      lfr,
		NextFundingTime:      t1,
		InterestRate:         ir,
		Time:                 t,
	}

	return pf, nil
}

func (as *apiService) PriceChangeSituation(pcsr PriceChangeSituationRequest) (*PriceChangeSituationInfo, error) {
	params := make(map[string]string)
	params["symbol"] = pcsr.Symbol

	res, err := as.request("GET", "/fapi/v1/ticker/24hr", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from PriceChangeSituation.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawPriceChangeSituation := struct {
		Symbol             string  `json:"symbol"`
		PriceChange        string  `json:"priceChange"`
		PriceChangePercent string  `json:"priceChangePercent"`
		WeightedAvgPrice   string  `json:"weightedAvgPrice"`
		LastPrice          string  `json:"lastPrice"`
		LastQty            string  `json:"lastQty"`
		OpenPrice          string  `json:"openPrice"`
		HighPrice          string  `json:"highPrice"`
		LowPrice           string  `json:"lowPrice"`
		Volume             string  `json:"volume"`
		QuoteVolume        string  `json:"quoteVolume"`
		OpenTime           float64 `json:"openTime"`
		CloseTime          float64 `json:"closeTime"`
		FirstId            int     `json:"firstId"`
		LastId             int     `json:"lastId"`
		Count              int     `json:"count"`
	}{}
	if err := json.Unmarshal(textRes, &rawPriceChangeSituation); err != nil {
		return nil, errors.Wrap(err, "PriceChangeSituation unmarshal failed")
	}

	pc, _ := floatFromString(rawPriceChangeSituation.PriceChange)
	pcp, _ := floatFromString(rawPriceChangeSituation.PriceChangePercent)
	wap, _ := floatFromString(rawPriceChangeSituation.WeightedAvgPrice)
	lp, _ := floatFromString(rawPriceChangeSituation.LastPrice)
	lq, _ := floatFromString(rawPriceChangeSituation.LastQty)
	op, _ := floatFromString(rawPriceChangeSituation.OpenPrice)
	hp, _ := floatFromString(rawPriceChangeSituation.HighPrice)
	lowp, _ := floatFromString(rawPriceChangeSituation.LowPrice)
	v, _ := floatFromString(rawPriceChangeSituation.Volume)
	qv, _ := floatFromString(rawPriceChangeSituation.QuoteVolume)

	t, err := timeFromUnixTimestampFloat(rawPriceChangeSituation.OpenTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Time")
	}

	t1, err := timeFromUnixTimestampFloat(rawPriceChangeSituation.CloseTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Time")
	}

	if rawPriceChangeSituation.FirstId != 0 {
		params["firstId"] = strconv.Itoa(rawPriceChangeSituation.FirstId)
	}

	if rawPriceChangeSituation.LastId != 0 {
		params["lastId"] = strconv.Itoa(rawPriceChangeSituation.LastId)
	}

	if rawPriceChangeSituation.Count != 0 {
		params["count"] = strconv.Itoa(rawPriceChangeSituation.Count)
	}

	rpcs := &PriceChangeSituationInfo{
		Symbol:             rawPriceChangeSituation.Symbol,
		PriceChange:        pc,
		PriceChangePercent: pcp,
		WeightedAvgPrice:   wap,
		LastPrice:          lp,
		LastQty:            lq,
		OpenPrice:          op,
		HighPrice:          hp,
		LowPrice:           lowp,
		Volume:             v,
		QuoteVolume:        qv,
		OpenTime:           t,
		CloseTime:          t1,
		FirstId:            rawPriceChangeSituation.FirstId,
		LastId:             rawPriceChangeSituation.LastId,
		Count:              rawPriceChangeSituation.Count,
	}
	return rpcs, nil
}

func (as *apiService) OpenInterestNums(oir OpenInterestNumsRequest) (*OpenInterestNumsInfo, error) {
	params := make(map[string]string)
	params["symbol"] = oir.Symbol

	res, err := as.request("GET", "/fapi/v1/openInterest", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from OpenInterestNums.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOpenInterestNums := struct {
		OpenInterest string  `json:"openInterest "`
		Symbol       string  `json:"symbol"`
		Time         float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawOpenInterestNums); err != nil {
		return nil, errors.Wrap(err, "OpenInterestNums unmarshal failed")
	}

	oi, _ := floatFromString(rawOpenInterestNums.OpenInterest)
	t, err := timeFromUnixTimestampFloat(rawOpenInterestNums.Time)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Time")
	}

	roin := &OpenInterestNumsInfo{
		Symbol:       rawOpenInterestNums.Symbol,
		OpenInterest: oi,
		Time:         t,
	}

	return roin, nil
}

func (as *apiService) BestBookTicker(bbtr BestBookTickerRequest) (*BestBookTickerInfo, error) {
	params := make(map[string]string)
	params["symbol"] = bbtr.Symbol

	res, err := as.request("GET", "/fapi/v1/ticker/bookTicker", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from BestBookTicker.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawBestBookTicker := struct {
		Symbol   string  `json:"symbol"`
		BidPrice string  `json:"bidPrice "`
		BidQty   string  `json:"bidQty"`
		AskPrice string  `json:"askPrice"`
		AskQty   string  `json:"askQty"`
		Time     float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawBestBookTicker); err != nil {
		return nil, errors.Wrap(err, "OpenInterestNums unmarshal failed")
	}

	bp, _ := floatFromString(rawBestBookTicker.BidPrice)
	bq, _ := floatFromString(rawBestBookTicker.BidQty)
	ap, _ := floatFromString(rawBestBookTicker.AskPrice)
	aq, _ := floatFromString(rawBestBookTicker.AskQty)
	t, err := timeFromUnixTimestampFloat(rawBestBookTicker.Time)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Time")
	}

	rbbt := &BestBookTickerInfo{
		Symbol:   rawBestBookTicker.Symbol,
		BidPrice: bp,
		BidQty:   bq,
		AskPrice: ap,
		AskQty:   aq,
		Time:     t,
	}
	return rbbt, nil
}

func (as *apiService) ContractPosition(bpr ContractPositionRequest) ([]*ContractPositionInfo, error) {
	params := make(map[string]string)
	params["symbol"] = bpr.Symbol

	params["period"] = bpr.Period

	if bpr.Limit != 0 {
		params["limit"] = strconv.Itoa(bpr.Limit)
	}
	if bpr.StartTime != 0 {
		params["startTime"] = strconv.Itoa(int(bpr.StartTime))
	}

	if bpr.EndTime != 0 {
		params["endTime"] = strconv.Itoa(int(bpr.EndTime))
	}

	res, err := as.request("GET", "/futures/data/openInterestHist", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from ContractPosition.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawContractPosition := []struct {
		Symbol               string  `json:"symbol"`
		SumOpenInterest      string  `json:"sumOpenInterest"`
		SumOpenInterestValue string  `json:"sumOpenInterestValue"`
		Timestamp            float64 `json:"timestamp"`
	}{}

	if err := json.Unmarshal(textRes, &rawContractPosition); err != nil {
		return nil, errors.Wrap(err, "rawContractPosition unmarshal failed")
	}

	var fbi []*ContractPositionInfo

	for _, rcp := range rawContractPosition {
		soi, _ := floatFromString(rcp.SumOpenInterest)
		soiv, _ := floatFromString(rcp.SumOpenInterestValue)
		t, err := timeFromUnixTimestampFloat(rcp.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		fbi = append(fbi, &ContractPositionInfo{
			Symbol:               rcp.Symbol,
			SumOpenInterest:      soi,
			SumOpenInterestValue: soiv,
			Timestamp:            t,
		})
	}
	return fbi, nil
}

func (as *apiService) TopLongShortPositionRatio(tspr TopLongShortPositionRatioRequest) ([]*TopLongShortPositionRatioInfo, error) {
	params := make(map[string]string)
	params["symbol"] = tspr.Symbol

	params["period"] = tspr.Period

	if tspr.Limit != 0 {
		params["limit"] = strconv.Itoa(tspr.Limit)
	}
	if tspr.StartTime != 0 {
		params["startTime"] = strconv.Itoa(int(tspr.StartTime))
	}

	if tspr.EndTime != 0 {
		params["endTime"] = strconv.Itoa(int(tspr.EndTime))
	}

	res, err := as.request("GET", "/futures/data/topLongShortPositionRatio", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from ContractPosition.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawTopLongShortPositionRatio := []struct {
		Symbol         string  `json:"symbol"`
		LongShortRatio string  `json:"longShortRatio"`
		LongAccount    string  `json:"longAccount"`
		ShortAccount   string  `json:"shortAccount"`
		Timestamp      float64 `json:"timestamp"`
	}{}
	if err := json.Unmarshal(textRes, &rawTopLongShortPositionRatio); err != nil {
		return nil, errors.Wrap(err, "rawTopLongShortPositionRatio unmarshal failed")
	}

	var tri []*TopLongShortPositionRatioInfo

	for _, tlspri := range rawTopLongShortPositionRatio {
		lsr, _ := floatFromString(tlspri.LongShortRatio)
		la, _ := floatFromString(tlspri.LongAccount)
		sa, _ := floatFromString(tlspri.ShortAccount)
		t, err := timeFromUnixTimestampFloat(tlspri.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		tri = append(tri, &TopLongShortPositionRatioInfo{
			Symbol:         tlspri.Symbol,
			LongShortRatio: lsr,
			LongAccount:    la,
			ShortAccount:   sa,
			Timestamp:      t,
		})
	}
	return tri, nil
}

func (as *apiService) GlobalLongShortAccountRatio(glsarr GlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error) {
	params := make(map[string]string)
	params["symbol"] = glsarr.Symbol

	params["period"] = glsarr.Period

	if glsarr.Limit != 0 {
		params["limit"] = strconv.Itoa(glsarr.Limit)
	}
	if glsarr.StartTime != 0 {
		params["startTime"] = strconv.Itoa(int(glsarr.StartTime))
	}

	if glsarr.EndTime != 0 {
		params["endTime"] = strconv.Itoa(int(glsarr.EndTime))
	}

	res, err := as.request("GET", "/futures/data/globalLongShortAccountRatio", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from GlobalLongShortAccountRatio.GET")
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawGlobalLongShortAccountRatio := []struct {
		Symbol         string  `json:"symbol"`
		LongShortRatio string  `json:"longShortRatio"`
		LongAccount    string  `json:"longAccount"`
		ShortAccount   string  `json:"shortAccount"`
		Timestamp      float64 `json:"timestamp"`
	}{}
	if err := json.Unmarshal(textRes, &rawGlobalLongShortAccountRatio); err != nil {
		return nil, errors.Wrap(err, "rawGlobalLongShortAccountRatio unmarshal failed")
	}

	var tri []*GlobalLongShortAccountRatioInfo

	for _, tlspri := range rawGlobalLongShortAccountRatio {
		lsr, _ := floatFromString(tlspri.LongShortRatio)
		la, _ := floatFromString(tlspri.LongAccount)
		sa, _ := floatFromString(tlspri.ShortAccount)
		t, err := timeFromUnixTimestampFloat(tlspri.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		tri = append(tri, &GlobalLongShortAccountRatioInfo{
			Symbol:         tlspri.Symbol,
			LongShortRatio: lsr,
			LongAccount:    la,
			ShortAccount:   sa,
			Timestamp:      t,
		})
	}
	return tri, nil
}

func (as *apiService) TakerlongshortRatio(tlr TakerlongshortRatioRequest) ([]*TakerlongshortRatioInfo, error) {
	params := make(map[string]string)
	params["symbol"] = tlr.Symbol

	params["period"] = tlr.Period

	if tlr.Limit != 0 {
		params["limit"] = strconv.Itoa(tlr.Limit)
	}
	if tlr.StartTime != 0 {
		params["startTime"] = strconv.Itoa(int(tlr.StartTime))
	}

	if tlr.EndTime != 0 {
		params["endTime"] = strconv.Itoa(int(tlr.EndTime))
	}

	res, err := as.request("GET", "/futures/data/takerlongshortRatio", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from TakerlongshortRatio.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawTakerlongshortRatio := []struct {
		BuySellRatio string  `json:"BuySellRatio"`
		BuyVol       string  `json:"BuyVol"`
		SellVol      string  `json:"SellVol"`
		Timestamp    float64 `json:"timestamp"`
	}{}
	if err := json.Unmarshal(textRes, &rawTakerlongshortRatio); err != nil {
		return nil, errors.Wrap(err, "rawTakerlongshortRatio unmarshal failed")
	}

	var fbi []*TakerlongshortRatioInfo

	for _, tlri := range rawTakerlongshortRatio {
		bsr, _ := floatFromString(tlri.BuySellRatio)
		bv, _ := floatFromString(tlri.BuyVol)
		sv, _ := floatFromString(tlri.SellVol)
		t, err := timeFromUnixTimestampFloat(tlri.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		fbi = append(fbi, &TakerlongshortRatioInfo{
			BuySellRatio: bsr,
			BuyVol:       bv,
			SellVol:      sv,
			Timestamp:    t,
		})
	}
	return fbi, nil
}

func (as *apiService) GetNewPrice(npr NewPriceRequest) (*NewPriceInfo, error) {

	params := make(map[string]string)
	params["symbol"] = npr.Symbol

	res, err := as.request("GET", "fapi/v1/ticker/price", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawNewPrices := struct {
		Symbol     string  `json:"symbol"`
		Price      string  `json:"price"`
		UpdateTime float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawNewPrices); err != nil {
		return nil, errors.Wrap(err, "rawNewPrices failed")
	}

	price, _ := floatFromString(rawNewPrices.Price)

	t, _ := timeFromUnixTimestampFloat(rawNewPrices.UpdateTime)

	tpc := &NewPriceInfo{
		Symbol:     rawNewPrices.Symbol,
		Price:      price,
		UpdateTime: t,
	}
	return tpc, nil
}
