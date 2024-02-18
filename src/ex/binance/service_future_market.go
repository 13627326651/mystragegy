package binance

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

func (as *apiService) FutureKlines(kr KlinesRequest) ([]*Kline, error) {
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

	res, err := as.request("GET", "fapi/v1/klines", params, false, false)
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
