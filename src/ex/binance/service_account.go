package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type rawExecutedOrder struct {
	Symbol        string  `json:"symbol"`
	OrderID       int     `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         string  `json:"price"`
	OrigQty       string  `json:"origQty"`
	ExecutedQty   string  `json:"executedQty"`
	Status        string  `json:"status"`
	TimeInForce   string  `json:"timeInForce"`
	Type          string  `json:"type"`
	Side          string  `json:"side"`
	StopPrice     string  `json:"stopPrice"`
	IcebergQty    string  `json:"icebergQty"`
	Time          float64 `json:"time"`
}

type rawExecutedFutureOrder struct {
	Symbol        string  `json:"symbol"`
	OrderID       int     `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         string  `json:"price"`
	OrigQty       string  `json:"origQty"`
	ExecutedQty   string  `json:"executedQty"`
	Status        string  `json:"status"`
	TimeInForce   string  `json:"timeInForce"`
	Type          string  `json:"type"`
	Side          string  `json:"side"`
	StopPrice     string  `json:"stopPrice"`
	IcebergQty    string  `json:"icebergQty"`
	Time          float64 `json:"time"`
}

func (as *apiService) NewFutureOrder(or NewFutureOrderRequest) (*FutureProcessedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)

	if or.PositionSide != "" {
		params["positionSide"] = string(or.PositionSide)
	}

	params["timeInForce"] = string(or.TimeInForce)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', 8, 64)
	params["price"] = strconv.FormatFloat(or.Price, 'f', 8, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)

	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.StopPrice != 0.0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', 8, 64)
	}

	res, err := as.request("POST", "fapi/v1/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := struct {
		Symbol        string  `json:"symbol"`
		CumQuote      string  `json:"cumQuote"`      // 成交金额
		ExecutedQty   string  `json:"executedQty"`   // 成交量
		ClientOrderId string  `json:"clientOrderId"` // 用户自定义订单号
		OrderId       int64   `json:"orderId"`       // 系统订单号
		AvgPrice      string  `json:"avgPrice"`      // 平均成交价
		OrigQty       string  `json:"origQty"`       // 原始委托数量
		Price         string  `json:"price"`         // 委托价格
		Side          string  `json:"side"`          // 买卖方向
		PositionSide  string  `json:"positionSide"`  // 持仓方向
		Status        string  `json:"status"`        // 订单状态
		StopPrice     string  `json:"stopPrice"`     // 触发价
		ClosePosition bool    `json:"closePosition"` // 是否条件全平仓
		TimeInForce   string  `json:"timeInForce"`   // 有效方法
		Type          string  `json:"type"`          // 订单类型
		OrigType      string  `json:"origType"`      // 触发前订单类型
		ActivatePrice string  `json:"activatePrice"` // 跟踪止损激活价格， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     string  `json:"priceRate"`     // 跟踪止损回调比例， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		WorkingType   string  `json:"workingType"`   // 条件价格触发类型
		PriceProtect  bool    `json:"priceProtect"`  // 是否开启条件单触发保护
		UpdateTime    float64 `json:"updateTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	cq, _ := floatFromString(rawOrder.CumQuote)
	eq, _ := floatFromString(rawOrder.ExecutedQty)
	ap, _ := floatFromString(rawOrder.AvgPrice)
	oq, _ := floatFromString(rawOrder.OrigQty)
	p, _ := floatFromString(rawOrder.Price)
	sp, _ := floatFromString(rawOrder.StopPrice)
	aep, _ := floatFromString(rawOrder.ActivatePrice)
	pr, _ := floatFromString(rawOrder.PriceRate)
	t, _ := timeFromUnixTimestampFloat(rawOrder.UpdateTime)

	return &FutureProcessedOrder{
		Symbol:        rawOrder.Symbol,
		ClientOrderId: rawOrder.ClientOrderId,
		OrderId:       rawOrder.OrderId,
		CumQuote:      cq,
		StopPrice:     sp,
		ActivatePrice: aep,
		PriceRate:     pr,
		AvgPrice:      ap,
		ExecutedQty:   eq,
		OrigQty:       oq,
		Price:         p,
		Side:          rawOrder.Side,
		PositionSide:  rawOrder.PositionSide,
		Type:          rawOrder.Type,
		UpdateTime:    t,
	}, nil
}

func (as *apiService) NewOrder(or NewOrderRequest) (*ProcessedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)
	params["timeInForce"] = string(or.TimeInForce)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', 8, 64)
	params["price"] = strconv.FormatFloat(or.Price, 'f', 8, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)

	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.StopPrice != 0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', 10, 64)
	}
	if or.IcebergQty != 0 {
		params["icebergQty"] = strconv.FormatFloat(or.IcebergQty, 'f', 10, 64)
	}

	res, err := as.request("POST", "api/v3/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := struct {
		Symbol        string  `json:"symbol"`
		OrderID       int64   `json:"orderId"`
		ClientOrderID string  `json:"clientOrderId"`
		TransactTime  float64 `json:"transactTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	t, err := timeFromUnixTimestampFloat(rawOrder.TransactTime)
	if err != nil {
		return nil, err
	}

	return &ProcessedOrder{
		Symbol:        rawOrder.Symbol,
		OrderID:       rawOrder.OrderID,
		ClientOrderID: rawOrder.ClientOrderID,
		TransactTime:  t,
	}, nil
}

func (as *apiService) NewOrderTest(or NewOrderRequest) error {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)
	params["timeInForce"] = string(or.TimeInForce)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', 8, 64)
	params["price"] = strconv.FormatFloat(or.Price, 'f', 8, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)

	fmt.Println("xxxx", params["quantity"], params["price"])
	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.StopPrice != 0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', 10, 64)
	}
	if or.IcebergQty != 0 {
		params["icebergQty"] = strconv.FormatFloat(or.IcebergQty, 'f', 10, 64)
	}

	res, err := as.request("POST", "api/v3/order/test", params, true, true)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}
	return nil
}

func (as *apiService) QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = qor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(qor.Timestamp), 10)
	if qor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(qor.OrderID, 10)
	}
	if qor.OrigClientOrderID != "" {
		params["origClientOrderId"] = qor.OrigClientOrderID
	}
	if qor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(qor.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := &rawExecutedOrder{}
	if err := json.Unmarshal(textRes, rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	eo, err := executedOrderFromRaw(rawOrder)
	if err != nil {
		return nil, err
	}
	return eo, nil
}

func (as *apiService) CancelFutureOrder(cfr CancelFutureOrderRequest) (*CanceledFutureOrder, error) {
	params := make(map[string]string)
	params["symbol"] = cfr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(cfr.Timestamp), 10)
	if cfr.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(cfr.OrderID, 10)
	}
	if cfr.OrigClientOrderID != "" {
		params["origClientOrderId"] = cfr.OrigClientOrderID
	}

	if cfr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(cfr.RecvWindow), 10)
	}

	res, err := as.request("DELETE", "fapi/v1/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.delete")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawCanceledOrder := struct {
		Symbol            string `json:"symbol"`
		OrigClientOrderID string `json:"origClientOrderId"`
		OrderID           int64  `json:"orderId"`
		ClientOrderID     string `json:"clientOrderId"`
	}{}
	if err := json.Unmarshal(textRes, &rawCanceledOrder); err != nil {
		return nil, errors.Wrap(err, "cancelOrder unmarshal failed")
	}

	return &CanceledFutureOrder{
		Symbol:            rawCanceledOrder.Symbol,
		OrigClientOrderID: rawCanceledOrder.OrigClientOrderID,
		OrderID:           rawCanceledOrder.OrderID,
	}, nil
}

func (as *apiService) CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error) {
	params := make(map[string]string)
	params["symbol"] = cor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(cor.Timestamp), 10)
	if cor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(cor.OrderID, 10)
	}
	if cor.OrigClientOrderID != "" {
		params["origClientOrderId"] = cor.OrigClientOrderID
	}
	if cor.NewClientOrderID != "" {
		params["newClientOrderId"] = cor.NewClientOrderID
	}
	if cor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(cor.RecvWindow), 10)
	}

	res, err := as.request("DELETE", "api/v3/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.delete")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawCanceledOrder := struct {
		Symbol            string `json:"symbol"`
		OrigClientOrderID string `json:"origClientOrderId"`
		OrderID           int64  `json:"orderId"`
		ClientOrderID     string `json:"clientOrderId"`
	}{}
	if err := json.Unmarshal(textRes, &rawCanceledOrder); err != nil {
		return nil, errors.Wrap(err, "cancelOrder unmarshal failed")
	}

	return &CanceledOrder{
		Symbol:            rawCanceledOrder.Symbol,
		OrigClientOrderID: rawCanceledOrder.OrigClientOrderID,
		OrderID:           rawCanceledOrder.OrderID,
		ClientOrderID:     rawCanceledOrder.ClientOrderID,
	}, nil
}

func (as *apiService) QueryOneFutureOrder(qfo QueryFutureOrderRequest) (*ExecutedFutureOrder, error) {
	params := make(map[string]string)
	params["symbol"] = qfo.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(qfo.Timestamp), 10)
	if qfo.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(qfo.RecvWindow), 10)
	}
	params["origClientOrderId"] = qfo.OrigClientOrderId

	res, err := as.request("GET", "fapi/v1/openOrder", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from openOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := struct {
		Symbol        string  `json:"symbol"`
		CumQuote      string  `json:"cumQuote"`      // 成交金额
		ExecutedQty   string  `json:"executedQty"`   // 成交量
		ClientOrderId string  `json:"clientOrderId"` // 用户自定义订单号
		OrderId       int64   `json:"orderId"`       // 系统订单号
		AvgPrice      string  `json:"avgPrice"`      // 平均成交价
		OrigQty       string  `json:"origQty"`       // 原始委托数量
		Price         string  `json:"price"`         // 委托价格
		Side          string  `json:"side"`          // 买卖方向
		PositionSide  string  `json:"positionSide"`  // 持仓方向
		Status        string  `json:"status"`        // 订单状态
		StopPrice     string  `json:"stopPrice"`     // 触发价
		ClosePosition bool    `json:"closePosition"` // 是否条件全平仓
		TimeInForce   string  `json:"timeInForce"`   // 有效方法
		Type          string  `json:"type"`          // 订单类型
		OrigType      string  `json:"origType"`      // 触发前订单类型
		ActivatePrice string  `json:"activatePrice"` // 跟踪止损激活价格， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     string  `json:"priceRate"`     // 跟踪止损回调比例， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		WorkingType   string  `json:"workingType"`   // 条件价格触发类型
		PriceProtect  bool    `json:"priceProtect"`  // 是否开启条件单触发保护
		UpdateTime    float64 `json:"updateTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrder); err != nil {
		return nil, errors.Wrap(err, "openOrders unmarshal failed")
	}

	cu, _ := floatFromString(rawOrder.CumQuote)
	ex, _ := floatFromString(rawOrder.ExecutedQty)
	or, _ := floatFromString(rawOrder.OrigQty)
	pr, _ := floatFromString(rawOrder.Price)
	st, _ := floatFromString(rawOrder.StopPrice)
	ac, _ := floatFromString(rawOrder.ActivatePrice)
	pre, _ := floatFromString(rawOrder.PriceRate)
	t, _ := timeFromUnixTimestampFloat(rawOrder.UpdateTime)

	eo := &ExecutedFutureOrder{
		Symbol:        rawOrder.Symbol,
		CumQty:        cu,
		ExecutedQty:   ex,
		OrderID:       rawOrder.OrderId,
		ClientOrderID: rawOrder.ClientOrderId,
		AvgPrice:      rawOrder.AvgPrice,
		OrigQty:       or,
		Price:         pr,
		Side:          OrderSide(rawOrder.Side),
		PositionSide:  rawOrder.PositionSide,
		Status:        OrderStatus(rawOrder.Status),
		StopPrice:     st,
		ClosePosition: rawOrder.ClosePosition,
		TimeInForce:   TimeInForce(rawOrder.TimeInForce),
		Type:          OrderType(rawOrder.Type),
		OrigType:      rawOrder.OrigType,
		ActivetePrice: ac,
		PriceRate:     pre,
		WorkingType:   rawOrder.WorkingType,
		PriceProtect:  rawOrder.PriceProtect,
		UpdateTime:    t,
	}

	return eo, nil
}

func (as *apiService) QueryAllFutureOrder(qfo QueryFutureOrderRequest) ([]*ExecutedFutureOrder, error) {
	params := make(map[string]string)
	params["symbol"] = qfo.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(qfo.Timestamp), 10)
	if qfo.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(qfo.RecvWindow), 10)
	}

	res, err := as.request("GET", "fapi/v1/openOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from openOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []struct {
		Symbol        string  `json:"symbol"`
		CumQuote      string  `json:"cumQuote"`      // 成交金额
		ExecutedQty   string  `json:"executedQty"`   // 成交量
		ClientOrderId string  `json:"clientOrderId"` // 用户自定义订单号
		OrderId       int64   `json:"orderId"`       // 系统订单号
		AvgPrice      string  `json:"avgPrice"`      // 平均成交价
		OrigQty       string  `json:"origQty"`       // 原始委托数量
		Price         string  `json:"price"`         // 委托价格
		Side          string  `json:"side"`          // 买卖方向
		PositionSide  string  `json:"positionSide"`  // 持仓方向
		Status        string  `json:"status"`        // 订单状态
		StopPrice     string  `json:"stopPrice"`     // 触发价
		ClosePosition bool    `json:"closePosition"` // 是否条件全平仓
		TimeInForce   string  `json:"timeInForce"`   // 有效方法
		Type          string  `json:"type"`          // 订单类型
		OrigType      string  `json:"origType"`      // 触发前订单类型
		ActivatePrice string  `json:"activatePrice"` // 跟踪止损激活价格， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     string  `json:"priceRate"`     // 跟踪止损回调比例， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		WorkingType   string  `json:"workingType"`   // 条件价格触发类型
		PriceProtect  bool    `json:"priceProtect"`  // 是否开启条件单触发保护
		UpdateTime    float64 `json:"updateTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "openOrders unmarshal failed")
	}

	var eoc []*ExecutedFutureOrder
	for _, rawOrder := range rawOrders {

		cu, _ := floatFromString(rawOrder.CumQuote)
		ex, _ := floatFromString(rawOrder.ExecutedQty)
		or, _ := floatFromString(rawOrder.OrigQty)
		pr, _ := floatFromString(rawOrder.Price)
		st, _ := floatFromString(rawOrder.StopPrice)
		ac, _ := floatFromString(rawOrder.ActivatePrice)
		pre, _ := floatFromString(rawOrder.PriceRate)
		t, _ := timeFromUnixTimestampFloat(rawOrder.UpdateTime)

		eo := &ExecutedFutureOrder{
			Symbol:        rawOrder.Symbol,
			CumQty:        cu,
			ExecutedQty:   ex,
			OrderID:       rawOrder.OrderId,
			ClientOrderID: rawOrder.ClientOrderId,
			AvgPrice:      rawOrder.AvgPrice,
			OrigQty:       or,
			Price:         pr,
			Side:          OrderSide(rawOrder.Side),
			PositionSide:  rawOrder.PositionSide,
			Status:        OrderStatus(rawOrder.Status),
			StopPrice:     st,
			ClosePosition: rawOrder.ClosePosition,
			TimeInForce:   TimeInForce(rawOrder.TimeInForce),
			Type:          OrderType(rawOrder.Type),
			OrigType:      rawOrder.OrigType,
			ActivetePrice: ac,
			PriceRate:     pre,
			WorkingType:   rawOrder.WorkingType,
			PriceProtect:  rawOrder.PriceProtect,
			UpdateTime:    t,
		}

		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = oor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(oor.Timestamp), 10)
	if oor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(oor.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/openOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from openOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []*rawExecutedOrder{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "openOrders unmarshal failed")
	}

	var eoc []*ExecutedOrder
	for _, rawOrder := range rawOrders {
		eo, err := executedOrderFromRaw(rawOrder)
		if err != nil {
			return nil, err
		}
		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = aor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(aor.Timestamp), 10)
	if aor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(aor.OrderID, 10)
	}
	if aor.Limit != 0 {
		params["limit"] = strconv.Itoa(aor.Limit)
	}
	if aor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(aor.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/allOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from allOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []*rawExecutedOrder{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "allOrders unmarshal failed")
	}

	var eoc []*ExecutedOrder
	for _, rawOrder := range rawOrders {
		eo, err := executedOrderFromRaw(rawOrder)
		if err != nil {
			return nil, err
		}
		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) QueryAllHistoryFutureOrders(afo AllFutureOrdersRequest) ([]*ExecutedFutureOrder, error) {
	params := make(map[string]string)
	params["symbol"] = afo.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(afo.Timestamp), 10)
	if afo.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(afo.OrderID, 10)
	}
	if afo.Limit != 0 {
		params["limit"] = strconv.Itoa(afo.Limit)
	}
	if afo.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(afo.RecvWindow), 10)
	}

	if afo.StartTime != 0 {
		params["startTime"] = strconv.Itoa(int(afo.StartTime))
	}

	if afo.EndTime != 0 {
		params["endTime"] = strconv.Itoa(int(afo.EndTime))
	}

	res, err := as.request("GET", "fapi/v1/allOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from allOrders.get")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []struct {
		Symbol        string  `json:"symbol"`
		CumQuote      string  `json:"cumQuote"`      // 成交金额
		ExecutedQty   string  `json:"executedQty"`   // 成交量
		ClientOrderId string  `json:"clientOrderId"` // 用户自定义订单号
		OrderId       int64   `json:"orderId"`       // 系统订单号
		AvgPrice      string  `json:"avgPrice"`      // 平均成交价
		OrigQty       string  `json:"origQty"`       // 原始委托数量
		Price         string  `json:"price"`         // 委托价格
		Side          string  `json:"side"`          // 买卖方向
		PositionSide  string  `json:"positionSide"`  // 持仓方向
		Status        string  `json:"status"`        // 订单状态
		StopPrice     string  `json:"stopPrice"`     // 触发价
		ClosePosition bool    `json:"closePosition"` // 是否条件全平仓
		TimeInForce   string  `json:"timeInForce"`   // 有效方法
		Type          string  `json:"type"`          // 订单类型
		OrigType      string  `json:"origType"`      // 触发前订单类型
		ActivatePrice string  `json:"activatePrice"` // 跟踪止损激活价格， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     string  `json:"priceRate"`     // 跟踪止损回调比例， 仅`TRAILING_STOP_MARKET` 订单返回此字段
		WorkingType   string  `json:"workingType"`   // 条件价格触发类型
		PriceProtect  bool    `json:"priceProtect"`  // 是否开启条件单触发保护
		UpdateTime    float64 `json:"updateTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "rawOrders unmarshal failed")
	}

	var eoc []*ExecutedFutureOrder
	for _, rawOrder := range rawOrders {

		cu, _ := floatFromString(rawOrder.CumQuote)
		ex, _ := floatFromString(rawOrder.ExecutedQty)
		or, _ := floatFromString(rawOrder.OrigQty)
		pr, _ := floatFromString(rawOrder.Price)
		st, _ := floatFromString(rawOrder.StopPrice)
		ac, _ := floatFromString(rawOrder.ActivatePrice)
		pre, _ := floatFromString(rawOrder.PriceRate)
		t, _ := timeFromUnixTimestampFloat(rawOrder.UpdateTime)

		eo := &ExecutedFutureOrder{
			Symbol:        rawOrder.Symbol,
			CumQty:        cu,
			ExecutedQty:   ex,
			OrderID:       rawOrder.OrderId,
			ClientOrderID: rawOrder.ClientOrderId,
			AvgPrice:      rawOrder.AvgPrice,
			OrigQty:       or,
			Price:         pr,
			Side:          OrderSide(rawOrder.Side),
			PositionSide:  rawOrder.PositionSide,
			Status:        OrderStatus(rawOrder.Status),
			StopPrice:     st,
			ClosePosition: rawOrder.ClosePosition,
			TimeInForce:   TimeInForce(rawOrder.TimeInForce),
			Type:          OrderType(rawOrder.Type),
			OrigType:      rawOrder.OrigType,
			ActivetePrice: ac,
			PriceRate:     pre,
			WorkingType:   rawOrder.WorkingType,
			PriceProtect:  rawOrder.PriceProtect,
			UpdateTime:    t,
		}

		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) Account(ar AccountRequest) (*Account, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(ar.Timestamp), 10)
	if ar.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ar.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/account", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from account.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawAccount := struct {
		MakerCommision   int64 `json:"makerCommision"`
		TakerCommission  int64 `json:"takerCommission"`
		BuyerCommission  int64 `json:"buyerCommission"`
		SellerCommission int64 `json:"sellerCommission"`
		CanTrade         bool  `json:"canTrade"`
		CanWithdraw      bool  `json:"canWithdraw"`
		CanDeposit       bool  `json:"canDeposit"`
		Balances         []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		}
	}{}
	if err := json.Unmarshal(textRes, &rawAccount); err != nil {
		return nil, errors.Wrap(err, "rawAccount unmarshal failed")
	}

	acc := &Account{
		MakerCommision:  rawAccount.MakerCommision,
		TakerCommision:  rawAccount.TakerCommission,
		BuyerCommision:  rawAccount.BuyerCommission,
		SellerCommision: rawAccount.SellerCommission,
		CanTrade:        rawAccount.CanTrade,
		CanWithdraw:     rawAccount.CanWithdraw,
		CanDeposit:      rawAccount.CanDeposit,
	}
	for _, b := range rawAccount.Balances {
		f, err := floatFromString(b.Free)
		if err != nil {
			return nil, err
		}
		l, err := floatFromString(b.Locked)
		if err != nil {
			return nil, err
		}
		acc.Balances = append(acc.Balances, &Balance{
			Asset:  b.Asset,
			Free:   f,
			Locked: l,
		})
	}

	return acc, nil
}

func (as *apiService) MyTrades(mtr MyTradesRequest) ([]*Trade, error) {
	params := make(map[string]string)
	params["symbol"] = mtr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(mtr.Timestamp), 10)
	if mtr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mtr.RecvWindow), 10)
	}
	if mtr.FromID != 0 {
		params["orderId"] = strconv.FormatInt(mtr.FromID, 10)
	}
	if mtr.Limit != 0 {
		params["limit"] = strconv.Itoa(mtr.Limit)
	}

	res, err := as.request("GET", "api/v3/myTrades", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from myTrades.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawTrades := []struct {
		ID              int64   `json:"id"`
		Price           string  `json:"price"`
		Qty             string  `json:"qty"`
		Commission      string  `json:"commission"`
		CommissionAsset string  `json:"commissionAsset"`
		Time            float64 `json:"time"`
		IsBuyer         bool    `json:"isBuyer"`
		IsMaker         bool    `json:"isMaker"`
		IsBestMatch     bool    `json:"isBestMatch"`
	}{}
	if err := json.Unmarshal(textRes, &rawTrades); err != nil {
		return nil, errors.Wrap(err, "rawTrades unmarshal failed")
	}

	var tc []*Trade
	for _, rt := range rawTrades {
		price, err := floatFromString(rt.Price)
		if err != nil {
			return nil, err
		}
		qty, err := floatFromString(rt.Qty)
		if err != nil {
			return nil, err
		}
		commission, err := floatFromString(rt.Commission)
		if err != nil {
			return nil, err
		}
		t, err := timeFromUnixTimestampFloat(rt.Time)
		if err != nil {
			return nil, err
		}
		tc = append(tc, &Trade{
			ID:              rt.ID,
			Price:           price,
			Qty:             qty,
			Commission:      commission,
			CommissionAsset: rt.CommissionAsset,
			Time:            t,
			IsBuyer:         rt.IsBuyer,
			IsMaker:         rt.IsMaker,
			IsBestMatch:     rt.IsBestMatch,
		})
	}
	return tc, nil
}

func (as *apiService) Withdraw(wr WithdrawRequest) (*WithdrawResult, error) {
	params := make(map[string]string)
	params["asset"] = wr.Asset
	params["address"] = wr.Address
	params["amount"] = strconv.FormatFloat(wr.Amount, 'f', 10, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(wr.Timestamp), 10)
	if wr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(wr.RecvWindow), 10)
	}
	if wr.Name != "" {
		params["name"] = wr.Name
	}

	res, err := as.request("POST", "wapi/v1/withdraw.html", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from withdraw.post")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawResult := struct {
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
	}{}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.Wrap(err, "rawTrades unmarshal failed")
	}

	return &WithdrawResult{
		Msg:     rawResult.Msg,
		Success: rawResult.Success,
	}, nil
}
func (as *apiService) DepositHistory(hr HistoryRequest) ([]*Deposit, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(hr.Timestamp), 10)
	if hr.Asset != "" {
		params["asset"] = hr.Asset
	}
	if hr.Status != nil {
		params["status"] = strconv.Itoa(*hr.Status)
	}
	if !hr.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.StartTime), 10)
	}
	if !hr.EndTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.EndTime), 10)
	}
	if hr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(hr.RecvWindow), 10)
	}

	res, err := as.request("POST", "wapi/v1/getDepositHistory.html", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from depositHistory.post")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawDepositHistory := struct {
		DepositList []struct {
			InsertTime float64 `json:"insertTime"`
			Amount     float64 `json:"amount"`
			Asset      string  `json:"asset"`
			Status     int     `json:"status"`
		}
		Success bool `json:"success"`
	}{}
	if err := json.Unmarshal(textRes, &rawDepositHistory); err != nil {
		return nil, errors.Wrap(err, "rawDepositHistory unmarshal failed")
	}

	var dc []*Deposit
	for _, d := range rawDepositHistory.DepositList {
		t, err := timeFromUnixTimestampFloat(d.InsertTime)
		if err != nil {
			return nil, err
		}
		dc = append(dc, &Deposit{
			InsertTime: t,
			Amount:     d.Amount,
			Asset:      d.Asset,
			Status:     d.Status,
		})
	}

	return dc, nil
}
func (as *apiService) WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(hr.Timestamp), 10)
	if hr.Asset != "" {
		params["asset"] = hr.Asset
	}
	if hr.Status != nil {
		params["status"] = strconv.Itoa(*hr.Status)
	}
	if !hr.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.StartTime), 10)
	}
	if !hr.EndTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.EndTime), 10)
	}
	if hr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(hr.RecvWindow), 10)
	}

	res, err := as.request("POST", "wapi/v1/getWithdrawHistory.html", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from withdrawHistory.post")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawWithdrawHistory := struct {
		WithdrawList []struct {
			Amount    float64 `json:"amount"`
			Address   string  `json:"address"`
			TxID      string  `json:"txId"`
			Asset     string  `json:"asset"`
			ApplyTime float64 `json:"insertTime"`
			Status    int     `json:"status"`
		}
		Success bool `json:"success"`
	}{}
	if err := json.Unmarshal(textRes, &rawWithdrawHistory); err != nil {
		return nil, errors.Wrap(err, "rawWithdrawHistory unmarshal failed")
	}

	var wc []*Withdrawal
	for _, w := range rawWithdrawHistory.WithdrawList {
		t, err := timeFromUnixTimestampFloat(w.ApplyTime)
		if err != nil {
			return nil, err
		}
		wc = append(wc, &Withdrawal{
			Amount:    w.Amount,
			Address:   w.Address,
			TxID:      w.TxID,
			Asset:     w.Asset,
			ApplyTime: t,
			Status:    w.Status,
		})
	}

	return wc, nil
}

func (as *apiService) FutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error) {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(fbr.Timestamp), 10)
	if fbr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(fbr.RecvWindow), 10)
	}

	res, err := as.request("GET", "fapi/v2/balance", params, true, true)
	if err != nil {
		return nil, errors.Wrap(err, "FutureBalance request failed")
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from FutureBalance")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawResult := []struct {
		AccountAlias       string `json:"accountAlias "`
		Asset              string `json:"asset"`
		Balance            string `json:"balance"`
		CrossWalletBalance string `json:"crossWalletBalance"`
		CrossUnPnl         string `json:"crossUnPnl"`
		AvailableBalance   string `json:"availableBalance"`
		MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
		MarginAvailable    bool   `json:"marginAvailable"`
		UpdateTime         int64  `json:"updateTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.Wrap(err, "futurebalance unmarshal failed")
	}

	var fbi []*FutureBalanceInfo

	for _, v := range rawResult {

		balance, err := floatFromString(v.Balance)
		if err != nil {
			return nil, err
		}

		crossWalletBalance, err := floatFromString(v.CrossWalletBalance)
		if err != nil {
			return nil, err
		}

		crossUnPnl, err := floatFromString(v.CrossUnPnl)
		if err != nil {
			return nil, err
		}

		availableBalance, err := floatFromString(v.AvailableBalance)
		if err != nil {
			return nil, err
		}

		maxWithdrawAmount, err := floatFromString(v.AvailableBalance)
		if err != nil {
			return nil, err
		}

		tm := time.Unix(0, v.UpdateTime)

		t := FutureBalanceInfo{
			AccountAlias:       v.AccountAlias,
			Asset:              v.Asset,
			Balance:            balance,
			CrossWalletBalance: crossWalletBalance,
			CrossUnPnl:         crossUnPnl,
			AvailableBalance:   availableBalance,
			MaxWithdrawAmount:  maxWithdrawAmount,
			MarginAvailable:    v.MarginAvailable,
			UpdateTime:         tm,
		}

		fbi = append(fbi, &t)

	}
	return fbi, nil
}

func (as *apiService) FutureAccount(far FutureAccountRequest) (*FutureAccountInfo, error) {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(far.Timestamp), 10)
	if far.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(far.RecvWindow), 10)
	}

	res, err := as.request("GET", "fapi/v2/account", params, true, true)
	if err != nil {
		return nil, errors.Wrap(err, "FutureAccount request failed")
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from FutureBalance")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {

		return nil, as.handleError(textRes)
	}

	rawResult := struct {
		FeeTier                     int    `json:"feeTier"`
		CanTrade                    bool   `json:"canTrade"`
		CanDeposit                  bool   `json:"canDeposit"`
		CanWithdraw                 bool   `json:"canWithdraw"`
		UpdateTime                  int    `json:"updateTime"`
		TotalInitialMargin          string `json:"totalInitialMargin"`
		TotalMaintMargin            string `json:"totalMaintMargin"`
		TotalWalletBalance          string `json:"totalWalletBalance"`
		TotalUnrealizedProfit       string `json:"totalUnrealizedProfit"`
		TotalMarginBalance          string `json:"totalMarginBalance "`
		TotalPositionInitialMargin  string `json:"totalPositionInitialMargin"`
		TotalOpenOrderInitialMargin string `json:"totalOpenOrderInitialMargin"`
		TotalCrossWalletBalance     string `json:"totalCrossWalletBalance"`
		TotalCrossUnPnl             string `json:"totalCrossUnPnl"`
		AvailableBalance            string `json:"availableBalance"`
		MaxWithdrawAmount           string `json:"maxWithdrawAmount"`
		Asset                       []struct {
			Asset                  string  `json:"asset"`
			WalletBalance          string  `json:"walletBalance"`          // 余额
			UnrealizedProfit       string  `json:"unrealizedProfit "`      // 未实现盈亏
			MarginBalance          string  `json:"marginBalance"`          // 保证金余额
			MaintMargin            string  `json:"maintMargin"`            // 维持保证金
			InitialMargin          string  `json:"initialMargin"`          // 当前所需起始保证金
			PositionInitialMargin  string  `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
			OpenOrderInitialMargin string  `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
			CrossWalletBalance     string  `json:"crossWalletBalance"`     // 全仓账户余额
			CrossUnPnl             string  `json:"crossUnPnl"`             // 全仓持仓未实现盈亏
			AvailableBalance       string  `json:"availableBalance"`       // 可用余额
			MaxWithdrawAmount      string  `json:"maxWithdrawAmount"`      // 最大可转出余额
			MarginAvailable        string  `json:"marginAvailable"`        // 是否可用作联合保证金
			UpdateTime             float64 `json:"updateTime"`
		}
		Positions []struct {
			Symbol                 string  `json:"symbol"`
			InitialMargin          string  `json:"initialMargin"`          // 当前所需起始保证金(基于最新标记价格)
			MaintMargin            string  `json:"maintMargin"`            // 维持保证金
			UnrealizedProfit       string  `json:"unrealizedProfit"`       // 持仓未实现盈亏
			PositionInitialMargin  string  `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
			OpenOrderInitialMargin string  `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
			Leverage               string  `json:"leverage"`               // 杠杆倍率
			Isolated               bool    `json:"isolated"`               // 是否是逐仓模式
			EntryPrice             string  `json:"entryPrice"`             // 持仓成本价
			MaxNotional            string  `json:"maxNotional"`            // 当前杠杆下用户可用的最大名义价值
			PositionSide           string  `json:"PositionSide"`           // 持仓方向
			PositionAmt            string  `json:"PositionAmt"`            // 持仓数量
			UpdateTime             float64 `json:"updateTime"`
		}
	}{}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.Wrap(err, "futurebalance unmarshal failed")
	}

	facc := &FutureAccountInfo{
		FeeTier:     rawResult.FeeTier,
		CanTrade:    rawResult.CanTrade,
		CanDeposit:  rawResult.CanDeposit,
		CanWithdraw: rawResult.CanWithdraw,
		UpdateTime:  rawResult.UpdateTime,
	}

	facc.TotalInitialMargin, _ = floatFromString(rawResult.TotalInitialMargin)
	facc.TotalMaintMargin, _ = floatFromString(rawResult.TotalMaintMargin)
	facc.TotalCrossWalletBalance, _ = floatFromString(rawResult.TotalCrossWalletBalance)
	facc.TotalUnrealizedProfit, _ = floatFromString(rawResult.TotalUnrealizedProfit)
	facc.TotalMarginBalance, _ = floatFromString(rawResult.TotalMarginBalance)
	facc.TotalPositionInitialMargin, _ = floatFromString(rawResult.TotalPositionInitialMargin)
	facc.TotalOpenOrderInitialMargin, _ = floatFromString(rawResult.TotalOpenOrderInitialMargin)
	facc.TotalWalletBalance, _ = floatFromString(rawResult.TotalWalletBalance)
	facc.TotalCrossUnPnl, _ = floatFromString(rawResult.TotalCrossUnPnl)
	facc.AvailableBalance, _ = floatFromString(rawResult.AvailableBalance)
	facc.MaxWithdrawAmount, _ = floatFromString(rawResult.MaxWithdrawAmount)

	for _, asset := range rawResult.Asset {

		WalletBalance, _ := floatFromString(asset.WalletBalance)
		UnrealizedProfit, _ := floatFromString(asset.UnrealizedProfit)
		MarginBalance, _ := floatFromString(asset.MarginBalance)
		InitialMargin, _ := floatFromString(asset.InitialMargin)
		PositionInitialMargin, _ := floatFromString(asset.PositionInitialMargin)
		OpenOrderInitialMargin, _ := floatFromString(asset.OpenOrderInitialMargin)
		CrossWalletBalance, _ := floatFromString(asset.CrossWalletBalance)
		CrossUnPnl, _ := floatFromString(asset.CrossUnPnl)
		AvailableBalance, _ := floatFromString(asset.AvailableBalance)
		MaxWithdrawAmount, _ := floatFromString(asset.MaxWithdrawAmount)
		MarginAvailable, _ := floatFromString(asset.MarginAvailable)

		t, err := timeFromUnixTimestampFloat(asset.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		facc.Asset = append(facc.Asset, &FutureAsset{
			Asset:                  asset.Asset,
			WalletBalance:          WalletBalance,
			UnrealizedProfit:       UnrealizedProfit,
			MarginBalance:          MarginBalance,
			InitialMargin:          InitialMargin,
			PositionInitialMargin:  PositionInitialMargin,
			OpenOrderInitialMargin: OpenOrderInitialMargin,
			CrossWalletBalance:     CrossWalletBalance,
			CrossUnPnl:             CrossUnPnl,
			AvailableBalance:       AvailableBalance,
			MaxWithdrawAmount:      MaxWithdrawAmount,
			MarginAvailable:        MarginAvailable,
			UpdateTime:             t,
		})
	}

	for _, position := range rawResult.Positions {
		InitialMargin, _ := floatFromString(position.InitialMargin)
		MaintMargin, _ := floatFromString(position.MaintMargin)
		UnrealizedProfit, _ := floatFromString(position.UnrealizedProfit)
		PositionInitialMargin, _ := floatFromString(position.PositionInitialMargin)
		OpenOrderInitialMargin, _ := floatFromString(position.OpenOrderInitialMargin)
		Leverage, _ := floatFromString(position.Leverage)
		EntryPrice, _ := floatFromString(position.EntryPrice)
		MaxNotional, _ := floatFromString(position.MaxNotional)
		PositionAmt, _ := floatFromString(position.PositionAmt)

		t, err := timeFromUnixTimestampFloat(position.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		facc.Positions = append(facc.Positions, &FuturePositions{
			Symbol:                 position.Symbol,
			InitialMargin:          InitialMargin,
			MaintMargin:            MaintMargin,
			UnrealizedProfit:       UnrealizedProfit,
			PositionInitialMargin:  PositionInitialMargin,
			OpenOrderInitialMargin: OpenOrderInitialMargin,
			Leverage:               Leverage,
			Isolated:               position.Isolated,
			EntryPrice:             EntryPrice,
			MaxNotional:            MaxNotional,
			PositionSide:           position.PositionSide,
			PositionAmt:            PositionAmt,
			UpdateTime:             t,
		})
	}

	return facc, nil
}

func (as *apiService) UserPoundage(udr UserPoundageRequest) (*UserPoundageInfo, error) {
	params := make(map[string]string)
	params["symbol"] = udr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(udr.Timestamp), 10)

	if udr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(udr.RecvWindow), 10)
	}

	res, err := as.request("GET", "fapi/v1/commissionRate", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from UserPoundage.GET")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawUserPoundage := struct {
		Symbol              string `json:"symbol"`
		MakerCommissionRate string `json:"makerCommissionRate"`
		TakerCommissionRate string `json:"takerCommissionRate"`
	}{}
	if err := json.Unmarshal(textRes, &rawUserPoundage); err != nil {
		return nil, errors.Wrap(err, "rawUserPoundage unmarshal failed")
	}

	mcr, _ := floatFromString(rawUserPoundage.MakerCommissionRate)
	tcr, _ := floatFromString(rawUserPoundage.TakerCommissionRate)

	wc := &UserPoundageInfo{
		Symbol:              rawUserPoundage.Symbol,
		MakerCommissionRate: mcr,
		TakerCommissionRate: tcr,
	}

	return wc, nil
}

func (as *apiService) AdjustLeverage(alr AdjustLeverageRequest) (*AdjustLeverageInfo, error) {
	params := make(map[string]string)
	params["symbol"] = alr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(alr.Timestamp), 10)

	if alr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(alr.RecvWindow), 10)
	}

	if alr.Leverage != 0 {
		params["leverage"] = strconv.Itoa(alr.Leverage)
	}

	res, err := as.request("POST", "/fapi/v1/leverage", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from AdjustLeverage.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawAdjustLeverage := struct {
		Leverage         int    `json:"leverage"`
		MaxNotionalValue string `json:"maxNotionalValue"`
		Symbol           string `json:"symbol"`
	}{}
	if err := json.Unmarshal(textRes, &rawAdjustLeverage); err != nil {
		return nil, errors.Wrap(err, "rawAdjustLeverage unmarshal failed")
	}

	mnl, _ := intFromString(rawAdjustLeverage.MaxNotionalValue)

	al := &AdjustLeverageInfo{
		Leverage:         rawAdjustLeverage.Leverage,
		MaxNotionalValue: mnl,
		Symbol:           rawAdjustLeverage.Symbol,
	}
	return al, nil
}

func (as *apiService) PositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error) {
	params := make(map[string]string)
	params["symbol"] = pmr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(pmr.Timestamp), 10)

	if pmr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(pmr.RecvWindow), 10)
	}
	params["positionSide"] = pmr.PositionSide
	params["amount"] = strconv.FormatFloat(pmr.Amount, 'f', 8, 64)
	if pmr.Type != 0 {
		params["type"] = strconv.Itoa(pmr.Type)
	}

	res, err := as.request("POST", "/fapi/v1/positionMargin", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from PositionMargin.POST")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawPositionMargin := struct {
		Amount string `json:"amount "`
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Type   int    `json:"type"`
	}{}
	if err := json.Unmarshal(textRes, &rawPositionMargin); err != nil {
		return nil, errors.Wrap(err, "rawPositionMargin unmarshal failed")
	}

	rpm, _ := floatFromString(rawPositionMargin.Amount)

	pm := &PositionMarginInfo{
		Amount: rpm,
		Code:   rawPositionMargin.Code,
		Msg:    rawPositionMargin.Msg,
		Type:   rawPositionMargin.Type,
	}

	return pm, nil
}

func (as *apiService) UserTradesHistory(uth UserTradesHistoryRequest) ([]*UserTradesHistoryInfo, error) {
	params := make(map[string]string)
	params["symbol"] = uth.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(uth.Timestamp), 10)

	if uth.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(uth.RecvWindow), 10)
	}
	if uth.Limit != 0 {
		params["limit"] = strconv.Itoa(uth.Limit)
	}
	if uth.FromId != 0 {
		params["fromId"] = strconv.Itoa(uth.FromId)
	}
	if uth.StartTime != 0 {
		params["startTime"] = strconv.Itoa(int(uth.StartTime))
	}

	if uth.EndTime != 0 {
		params["endTime"] = strconv.Itoa(int(uth.EndTime))
	}

	res, err := as.request("GET", "fapi/v1/userTrades", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from UserTradesHistory.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawUserTradesHistory := []struct {
		Buyer           bool    `json:"buyer"`
		Commission      string  `json:"commission"`
		CommissionAsset string  `json:"commissionAsset"`
		Id              int     `json:"id"`
		Maker           bool    `json:"maker"`
		OrderId         int     `json:"orderId"`
		Price           string  `json:"price"`
		Qty             string  `json:"qty"`
		QuoteQty        string  `json:"quoteQty"`
		RealizedPnl     string  `json:"realizedPnl"`
		Side            string  `json:"side"`
		PositionSide    string  `json:"positionSide"`
		Symbol          string  `json:"symbol"`
		Time            float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawUserTradesHistory); err != nil {
		return nil, errors.Wrap(err, "rawUserTradesHistory unmarshal failed")
	}

	var fbi []*UserTradesHistoryInfo

	for _, rt := range rawUserTradesHistory {

		t, err := timeFromUnixTimestampFloat(rt.Time)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		commission, _ := floatFromString(rt.Commission)

		price, _ := floatFromString(rt.Price)

		qty, _ := floatFromString(rt.Qty)

		quoteQty, _ := floatFromString(rt.QuoteQty)

		realizedPnl, _ := floatFromString(rt.RealizedPnl)

		fbi = append(fbi, &UserTradesHistoryInfo{
			Time:            t,
			Commission:      commission,
			Price:           price,
			Qty:             qty,
			QuoteQty:        quoteQty,
			RealizedPnl:     realizedPnl,
			Buyer:           rt.Buyer,
			CommissionAsset: rt.CommissionAsset,
			Id:              rt.Id,
			Maker:           rt.Maker,
			OrderId:         rt.OrderId,
			Side:            rt.Side,
			PositionSide:    rt.PositionSide,
			Symbol:          rt.Symbol,
		})
	}
	return fbi, nil
}

func (as *apiService) ChangeMarginType(mtr MarginTypeRequest) error {

	params := make(map[string]string)
	params["symbol"] = mtr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(mtr.Timestamp), 10)

	if mtr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mtr.RecvWindow), 10)
	}

	params["marginType"] = string(mtr.MarginType)

	res, err := as.request("POST", "fapi/v1/marginType", params, true, true)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from UserMarginType.GET")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}

	return nil

}

func (as *apiService) QueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error) {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(ups.Timestamp), 10)

	if ups.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ups.RecvWindow), 10)
	}

	res, err := as.request("GET", "fapi/v1/positionSide/dual", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from PositionMargin.POST")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawPositionMargin := struct {
		DualSidePosition bool `json:"dualSidePosition"`
	}{}
	if err := json.Unmarshal(textRes, &rawPositionMargin); err != nil {
		return nil, errors.Wrap(err, "rawPositionMargin unmarshal failed")
	}

	up := &UserPositionSideInfo{
		DualSidePosition: rawPositionMargin.DualSidePosition,
	}

	return up, nil
}

func (as *apiService) ChangeUserPositionSide(ups ChangeUserPositionSideRequest) error {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(ups.Timestamp), 10)

	params["dualSidePosition"] = string(ups.DualSidePosition)

	if ups.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ups.RecvWindow), 10)
	}

	res, err := as.request("POST", "fapi/v1/positionSide/dual", params, true, true)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from PositionMargin.POST")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}

	return nil

}

func executedOrderFromRaw(reo *rawExecutedOrder) (*ExecutedOrder, error) {
	price, err := strconv.ParseFloat(reo.Price, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.CloseTime")
	}
	origQty, err := strconv.ParseFloat(reo.OrigQty, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.OrigQty")
	}
	execQty, err := strconv.ParseFloat(reo.ExecutedQty, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.ExecutedQty")
	}
	stopPrice, err := strconv.ParseFloat(reo.StopPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.StopPrice")
	}
	icebergQty, err := strconv.ParseFloat(reo.IcebergQty, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.IcebergQty")
	}
	t, err := timeFromUnixTimestampFloat(reo.Time)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.CloseTime")
	}

	return &ExecutedOrder{
		Symbol:        reo.Symbol,
		OrderID:       reo.OrderID,
		ClientOrderID: reo.ClientOrderID,
		Price:         price,
		OrigQty:       origQty,
		ExecutedQty:   execQty,
		Status:        OrderStatus(reo.Status),
		TimeInForce:   TimeInForce(reo.TimeInForce),
		Type:          OrderType(reo.Type),
		Side:          OrderSide(reo.Side),
		StopPrice:     stopPrice,
		IcebergQty:    icebergQty,
		Time:          t,
	}, nil
}
