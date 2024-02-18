package binance

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

func (as *apiService) FutureCoinKlines(kr KlinesRequest) ([]*Kline, error) {
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

	res, err := as.request("GET", "dapi/v1/klines", params, false, false)

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
			return nil, errors.Wrap(err, "cannot parse Kline.Volume") // 成交量
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
			return nil, errors.Wrap(err, "cannot parse Kline.TakerBuyBaseAssetVolume") // 主动买入成交量
		}
		tbqav, err := floatFromString(k[10])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.TakerBuyQuoteAssetVolume") // 主动买入成交额
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

func (as *apiService) OrderCoinBook(obr OrderBookRequest) (*OrderBook, error) {
	params := make(map[string]string)
	params["symbol"] = obr.Symbol
	if obr.Limit != 0 {
		params["limit"] = strconv.Itoa(obr.Limit)
	}
	res, err := as.request("GET", "dapi/v1/depth", params, false, false)
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

func (as *apiService) CoinTopLongShortPositionRatio(glsarr CoinTopLongShortPositionRatioRequest) ([]*CoinTopLongShortPositionRatioInfo, error) {
	params := make(map[string]string)
	params["pair"] = glsarr.Pair

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

	res, err := as.request("GET", "/futures/data/topLongShortPositionRatio", params, true, true)
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
		Pair           string  `json:"pair"`
		LongShortRatio string  `json:"longShortRatio"`
		LongAccount    string  `json:"longAccount"`
		ShortAccount   string  `json:"shortAccount"`
		Timestamp      float64 `json:"timestamp"`
	}{}
	if err := json.Unmarshal(textRes, &rawGlobalLongShortAccountRatio); err != nil {
		return nil, errors.Wrap(err, "rawGlobalLongShortAccountRatio unmarshal failed")
	}

	var tri []*CoinTopLongShortPositionRatioInfo

	for _, tlspri := range rawGlobalLongShortAccountRatio {
		lsr, _ := floatFromString(tlspri.LongShortRatio)
		la, _ := floatFromString(tlspri.LongAccount)
		sa, _ := floatFromString(tlspri.ShortAccount)
		t, err := timeFromUnixTimestampFloat(tlspri.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		tri = append(tri, &CoinTopLongShortPositionRatioInfo{
			Pair:           tlspri.Pair,
			LongShortRatio: lsr,
			LongAccount:    la,
			ShortAccount:   sa,
			Timestamp:      t,
		})
	}
	return tri, nil
}

func (as *apiService) CoinContractPosition(bpr CoinContractPositionRequest) ([]*CoinContractPositionInfo, error) {
	params := make(map[string]string)
	params["pair"] = bpr.Pair
	params["contractType"] = bpr.ContractType
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
		Pair                 string  `json:"pair"`
		ContractType         string  `json:"contractType"`
		SumOpenInterest      string  `json:"sumOpenInterest"`
		SumOpenInterestValue string  `json:"sumOpenInterestValue"`
		Timestamp            float64 `json:"timestamp"`
	}{}

	if err := json.Unmarshal(textRes, &rawContractPosition); err != nil {
		return nil, errors.Wrap(err, "rawContractPosition unmarshal failed")
	}

	var fbi []*CoinContractPositionInfo

	for _, rcp := range rawContractPosition {
		soi, _ := floatFromString(rcp.SumOpenInterest)
		soiv, _ := floatFromString(rcp.SumOpenInterestValue)
		t, err := timeFromUnixTimestampFloat(rcp.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		fbi = append(fbi, &CoinContractPositionInfo{
			Pair:                 rcp.Pair,
			ContractType:         rcp.ContractType,
			SumOpenInterest:      soi,
			SumOpenInterestValue: soiv,
			Timestamp:            t,
		})
	}
	return fbi, nil
}

func (as *apiService) CoinGetNewPrice(npr NewPriceRequest) ([]*CoinNewPriceInfo, error) {

	params := make(map[string]string)
	params["symbol"] = npr.Symbol

	res, err := as.request("GET", "dapi/v1/ticker/price", params, false, false)
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

	rawNewPrices := []struct {
		Symbol     string  `json:"symbol"`
		Ps         string  `json:"ps"`
		Price      string  `json:"price"`
		UpdateTime float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawNewPrices); err != nil {
		return nil, errors.Wrap(err, "rawNewPrices failed")
	}

	var tpc []*CoinNewPriceInfo

	for _, rcp := range rawNewPrices {
		price, _ := floatFromString(rcp.Price)

		t, _ := timeFromUnixTimestampFloat(rcp.UpdateTime)

		tpc = append(tpc, &CoinNewPriceInfo{
			Symbol:     rcp.Symbol,
			Ps:         rcp.Ps,
			Price:      price,
			UpdateTime: t,
		})
	}

	return tpc, nil
}

func (as *apiService) CoinPriceChangeSituation(pcsr PriceChangeSituationRequest) ([]*CoinPriceChangeSituationInfo, error) {
	params := make(map[string]string)
	params["symbol"] = pcsr.Symbol

	res, err := as.request("GET", "/dapi/v1/ticker/24hr", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from CoinPriceChangeSituation.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawPriceChangeSituation := []struct {
		Symbol             string  `json:"symbol"`
		Pair               string  `json:"pair"`
		PriceChange        string  `json:"priceChange"`
		PriceChangePercent string  `json:"priceChangePercent"`
		WeightedAvgPrice   string  `json:"weightedAvgPrice"`
		LastPrice          string  `json:"lastPrice"`
		LastQty            string  `json:"lastQty"`
		OpenPrice          string  `json:"openPrice"`
		HighPrice          string  `json:"highPrice"`
		LowPrice           string  `json:"lowPrice"`
		Volume             string  `json:"volume"`
		BaseVolume         string  `json:"baseVolume"`
		OpenTime           float64 `json:"openTime"`
		CloseTime          float64 `json:"closeTime"`
		FirstId            int     `json:"firstId"`
		LastId             int     `json:"lastId"`
		Count              int     `json:"count"`
	}{}
	if err := json.Unmarshal(textRes, &rawPriceChangeSituation); err != nil {
		return nil, errors.Wrap(err, "PriceChangeSituation unmarshal failed")
	}

	var cc []*CoinPriceChangeSituationInfo
	for _, rs := range rawPriceChangeSituation {
		pc, _ := floatFromString(rs.PriceChange)
		pcp, _ := floatFromString(rs.PriceChangePercent)
		wap, _ := floatFromString(rs.WeightedAvgPrice)
		lp, _ := floatFromString(rs.LastPrice)
		lq, _ := floatFromString(rs.LastQty)
		op, _ := floatFromString(rs.OpenPrice)
		hp, _ := floatFromString(rs.HighPrice)
		lowp, _ := floatFromString(rs.LowPrice)
		v, _ := floatFromString(rs.Volume)
		qv, _ := floatFromString(rs.BaseVolume)

		t, err := timeFromUnixTimestampFloat(rs.OpenTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		t1, err := timeFromUnixTimestampFloat(rs.CloseTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		if rs.FirstId != 0 {
			params["firstId"] = strconv.Itoa(rs.FirstId)
		}

		if rs.LastId != 0 {
			params["lastId"] = strconv.Itoa(rs.LastId)
		}

		if rs.Count != 0 {
			params["count"] = strconv.Itoa(rs.Count)
		}

		cc = append(cc, &CoinPriceChangeSituationInfo{
			Symbol:             rs.Symbol,
			Pair:               rs.Pair,
			PriceChange:        pc,
			PriceChangePercent: pcp,
			WeightedAvgPrice:   wap,
			LastPrice:          lp,
			LastQty:            lq,
			OpenPrice:          op,
			HighPrice:          hp,
			LowPrice:           lowp,
			Volume:             v,
			BaseVolume:         qv,
			OpenTime:           t,
			CloseTime:          t1,
			FirstId:            rs.FirstId,
			LastId:             rs.LastId,
			Count:              rs.Count,
		})
	}

	return cc, nil
}

func (as *apiService) CoinTakerlongshortRatio(tlr CoinTakerlongshortRatioRequest) ([]*CoinTakerlongshortRatioInfo, error) {
	params := make(map[string]string)
	params["pair"] = tlr.Pair
	params["contractType"] = tlr.ContractType
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

	res, err := as.request("GET", "/futures/data/takerBuySellVol", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from CoinTakerlongshortRatio.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawTakerlongshortRatio := []struct {
		Pair              string  `json:"pair"`
		ContractType      string  `json:"contractType"`
		TakerBuyVol       string  `json:"takerBuyVol"`
		TakerSellVol      string  `json:"takerSellVol"`
		TakerBuyVolValue  string  `json:"takerBuyVolValue"`
		TakerSellVolValue string  `json:"takerSellVolValue"`
		Timestamp         float64 `json:"timestamp"`
	}{}
	if err := json.Unmarshal(textRes, &rawTakerlongshortRatio); err != nil {
		return nil, errors.Wrap(err, "rawTakerlongshortRatio unmarshal failed")
	}

	var fbi []*CoinTakerlongshortRatioInfo

	for _, tlri := range rawTakerlongshortRatio {
		tbv, _ := floatFromString(tlri.TakerBuyVol)
		tsv, _ := floatFromString(tlri.TakerSellVol)
		tbvv, _ := floatFromString(tlri.TakerBuyVolValue)
		tsvv, _ := floatFromString(tlri.TakerSellVolValue)
		t, err := timeFromUnixTimestampFloat(tlri.Timestamp)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		fbi = append(fbi, &CoinTakerlongshortRatioInfo{
			Pair:              tlri.Pair,
			ContractType:      tlri.ContractType,
			TakerBuyVol:       tbv,
			TakerSellVol:      tsv,
			TakerBuyVolValue:  tbvv,
			TakerSellVolValue: tsvv,
			Timestamp:         t,
		})
	}
	return fbi, nil
}

func (as *apiService) CoinBestBookTicker(bbtr BestBookTickerRequest) ([]*CoinBestBookTickerInfo, error) {
	params := make(map[string]string)
	params["symbol"] = bbtr.Symbol

	res, err := as.request("GET", "/dapi/v1/ticker/bookTicker", params, true, true)
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

	rawBestBookTicker := []struct {
		Symbol   string  `json:"symbol"`
		Pair     string  `json:"pair"`
		BidPrice string  `json:"bidPrice "`
		BidQty   string  `json:"bidQty"`
		AskPrice string  `json:"askPrice"`
		AskQty   string  `json:"askQty"`
		Time     float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawBestBookTicker); err != nil {
		return nil, errors.Wrap(err, "OpenInterestNums unmarshal failed")
	}

	var cbi []*CoinBestBookTickerInfo
	for _, rawBestBookTicker := range rawBestBookTicker {

		bp, _ := floatFromString(rawBestBookTicker.BidPrice)
		bq, _ := floatFromString(rawBestBookTicker.BidQty)
		ap, _ := floatFromString(rawBestBookTicker.AskPrice)
		aq, _ := floatFromString(rawBestBookTicker.AskQty)
		t, err := timeFromUnixTimestampFloat(rawBestBookTicker.Time)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		cbi = append(cbi, &CoinBestBookTickerInfo{
			Symbol:   rawBestBookTicker.Symbol,
			Pair:     rawBestBookTicker.Pair,
			BidPrice: bp,
			BidQty:   bq,
			AskPrice: ap,
			AskQty:   aq,
			Time:     t,
		})
	}

	return cbi, nil
}

func (as *apiService) CoinOpenInterestNums(oir OpenInterestNumsRequest) (*CoinOpenInterestNumsInfo, error) {
	params := make(map[string]string)
	params["symbol"] = oir.Symbol

	res, err := as.request("GET", "/dapi/v1/openInterest", params, true, true)
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
		Pair         string  `json:"pair"`
		ContractType string  `json:"contractType"`
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

	roin := &CoinOpenInterestNumsInfo{
		Symbol:       rawOpenInterestNums.Symbol,
		Pair:         rawOpenInterestNums.Pair,
		ContractType: rawOpenInterestNums.ContractType,
		OpenInterest: oi,
		Time:         t,
	}

	return roin, nil
}

func (as *apiService) CoinGlobalLongShortAccountRatio(glsarr CoinGlobalLongShortAccountRatioRequest) ([]*GlobalLongShortAccountRatioInfo, error) {
	params := make(map[string]string)
	params["pair"] = glsarr.Pair

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

	res, err := as.request("GET", "/futures/data/globalLongShortAccountRatio", params, true, true)
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
		Pair           string  `json:"pair"`
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
			Symbol:         tlspri.Pair,
			LongShortRatio: lsr,
			LongAccount:    la,
			ShortAccount:   sa,
			Timestamp:      t,
		})
	}
	return tri, nil
}
