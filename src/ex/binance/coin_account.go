package binance

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func (as *apiService) CoinChangeMarginType(mtr MarginTypeRequest) error {

	params := make(map[string]string)
	params["symbol"] = mtr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(mtr.Timestamp), 10)

	if mtr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mtr.RecvWindow), 10)
	}

	params["marginType"] = string(mtr.MarginType)

	res, err := as.request("POST", "dapi/v1/marginType", params, true, true)
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

func (as *apiService) CoinChangeUserPositionSide(ups ChangeUserPositionSideRequest) error {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(ups.Timestamp), 10)

	params["dualSidePosition"] = string(ups.DualSidePosition)

	if ups.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ups.RecvWindow), 10)
	}

	res, err := as.request("POST", "dapi/v1/positionSide/dual", params, true, true)
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

func (as *apiService) CoinQueryUserPositionSide(ups UserPositionSideRequest) (*UserPositionSideInfo, error) {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(ups.Timestamp), 10)

	if ups.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ups.RecvWindow), 10)
	}

	res, err := as.request("GET", "dapi/v1/positionSide/dual", params, true, true)
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

func (as *apiService) CoinUserPoundage(udr UserPoundageRequest) (*UserPoundageInfo, error) {
	params := make(map[string]string)
	params["symbol"] = udr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(udr.Timestamp), 10)

	if udr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(udr.RecvWindow), 10)
	}

	res, err := as.request("GET", "dapi/v1/commissionRate", params, true, true)
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

func (as *apiService) CoinQueryFutureOrder(qfo CoinQueryFutureOrderRequest) ([]*ExecutedFutureOrder, error) {
	params := make(map[string]string)
	params["symbol"] = qfo.Symbol
	params["pair"] = qfo.Pair
	params["timestamp"] = strconv.FormatInt(unixMillis(qfo.Timestamp), 10)
	if qfo.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(qfo.RecvWindow), 10)
	}

	res, err := as.request("GET", "dapi/v1/openOrders", params, true, true)
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
		CumBase       string  `json:"cumBase"` // 成交金额
		Pair          string  `json:"pair"`
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

		//cu, _ := floatFromString(rawOrder.CumBase)
		ex, _ := floatFromString(rawOrder.ExecutedQty)
		or, _ := floatFromString(rawOrder.OrigQty)
		pr, _ := floatFromString(rawOrder.Price)
		st, _ := floatFromString(rawOrder.StopPrice)
		ac, _ := floatFromString(rawOrder.ActivatePrice)
		pre, _ := floatFromString(rawOrder.PriceRate)
		t, _ := timeFromUnixTimestampFloat(rawOrder.UpdateTime)

		eo := &ExecutedFutureOrder{
			Symbol: rawOrder.Symbol,
			//CumBase:       cu,
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
			//Pair:          rawOrder.Pair,
		}

		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) CoinFutureBalance(fbr FutureBalanceRequest) ([]*FutureBalanceInfo, error) {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(fbr.Timestamp), 10)
	if fbr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(fbr.RecvWindow), 10)
	}

	res, err := as.request("GET", "dapi/v1/balance", params, true, true)
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
		WithdrawAvailable  string `json:"withdrawAvailable"`
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

		//withdrawAvailable, err := floatFromString(v.WithdrawAvailable)
		//if err != nil {
		//	return nil, err
		//}

		tm := time.Unix(0, v.UpdateTime)

		t := FutureBalanceInfo{
			AccountAlias:       v.AccountAlias,
			Asset:              v.Asset,
			Balance:            balance,
			CrossWalletBalance: crossWalletBalance,
			CrossUnPnl:         crossUnPnl,
			AvailableBalance:   availableBalance,
			//WithdrawAvailable:  withdrawAvailable,
			UpdateTime: tm,
		}

		fbi = append(fbi, &t)

	}
	return fbi, nil
}

func (as *apiService) CoinFutureAccount(far FutureAccountRequest) (*CoinFutureAccountInfo, error) {
	params := make(map[string]string)

	params["timestamp"] = strconv.FormatInt(unixMillis(far.Timestamp), 10)
	if far.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(far.RecvWindow), 10)
	}

	res, err := as.request("GET", "/dapi/v1/account", params, true, true)
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
		FeeTier     int  `json:"feeTier"`
		CanTrade    bool `json:"canTrade"`
		CanDeposit  bool `json:"canDeposit"`
		CanWithdraw bool `json:"canWithdraw"`
		UpdateTime  int  `json:"updateTime"`
		Asset       []struct {
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
			UpdateTime             float64 `json:"updateTime"`
		}
		Positions []struct {
			Symbol                 string  `json:"symbol"`
			InitialMargin          string  `json:"initialMargin"`          // 当前所需起始保证金(基于最新标记价格)
			MaintMargin            string  `json:"maintMargin"`            // 维持保证金
			UnrealizedProfit       string  `json:"unrealizedProfit "`      // 持仓未实现盈亏
			PositionInitialMargin  string  `json:"positionInitialMargin"`  // 持仓所需起始保证金(基于最新标记价格)
			OpenOrderInitialMargin string  `json:"openOrderInitialMargin"` // 当前挂单所需起始保证金(基于最新标记价格)
			Leverage               string  `json:"leverage"`               // 杠杆倍率
			Isolated               bool    `json:"isolated"`               // 是否是逐仓模式
			EntryPrice             string  `json:"entryPrice"`             // 持仓成本价            // 当前杠杆下用户可用的最大名义价值
			PositionSide           string  `json:"PositionSide"`           // 持仓方向
			PositionAmt            string  `json:"PositionAmt"`            // 持仓数量
			UpdateTime             float64 `json:"updateTime"`
			MaxQty                 string  `json:"maxQty"`
		}
	}{}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.Wrap(err, "futurebalance unmarshal failed")
	}

	facc := &CoinFutureAccountInfo{
		FeeTier:     rawResult.FeeTier,
		CanTrade:    rawResult.CanTrade,
		CanDeposit:  rawResult.CanDeposit,
		CanWithdraw: rawResult.CanWithdraw,
		UpdateTime:  rawResult.UpdateTime,
	}

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
		MaintMargin, _ := floatFromString(asset.MaintMargin)

		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}

		facc.Asset = append(facc.Asset, &CoinFutureAsset{
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
			MaintMargin:            MaintMargin,
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

		PositionAmt, _ := floatFromString(position.PositionAmt)
		maxQty, _ := floatFromString(position.MaxQty)
		t, err := timeFromUnixTimestampFloat(position.UpdateTime)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		facc.Positions = append(facc.Positions, &CoinFuturePositions{
			Symbol:                 position.Symbol,
			InitialMargin:          InitialMargin,
			MaintMargin:            MaintMargin,
			UnrealizedProfit:       UnrealizedProfit,
			PositionInitialMargin:  PositionInitialMargin,
			OpenOrderInitialMargin: OpenOrderInitialMargin,
			Leverage:               Leverage,
			Isolated:               position.Isolated,
			EntryPrice:             EntryPrice,
			MaxQty:                 maxQty,
			PositionSide:           position.PositionSide,
			PositionAmt:            PositionAmt,
			UpdateTime:             t,
		})
	}

	return facc, nil
}

func (as *apiService) CoinNewFutureOrder(or NewFutureOrderRequest) (*FutureProcessedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)

	if or.PositionSide != "" {
		params["positionSide"] = string(or.PositionSide)
	}

	params["timeInForce"] = string(or.TimeInForce)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', 0, 64)
	params["price"] = strconv.FormatFloat(or.Price, 'f', 2, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)

	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.StopPrice != 0.0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', 2, 64)
	}

	res, err := as.request("POST", "dapi/v1/order", params, true, true)
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
		CumBase       string  `json:"cumBase"`       // 成交金额
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

	//cq, _ := floatFromString(rawOrder.CumBase)
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
		//CumBase:       cq,
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

func (as *apiService) CoinCancelFutureOrder(cfr CancelFutureOrderRequest) (*CanceledFutureOrder, error) {
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

	res, err := as.request("DELETE", "dapi/v1/order", params, true, true)
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
		AvgPrice          string `json:"avgPrice"`
	}{}
	if err := json.Unmarshal(textRes, &rawCanceledOrder); err != nil {
		return nil, errors.Wrap(err, "cancelOrder unmarshal failed")
	}

	//ap, _ := floatFromString(rawCanceledOrder.AvgPrice)
	return &CanceledFutureOrder{
		Symbol:            rawCanceledOrder.Symbol,
		OrigClientOrderID: rawCanceledOrder.OrigClientOrderID,
		OrderID:           rawCanceledOrder.OrderID,
		//AvgPrice:          ap,
	}, nil
}

func (as *apiService) CoinAdjustLeverage(alr AdjustLeverageRequest) (*CoinAdjustLeverageInfo, error) {
	params := make(map[string]string)
	params["symbol"] = alr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(alr.Timestamp), 10)

	if alr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(alr.RecvWindow), 10)
	}

	if alr.Leverage != 0 {
		params["leverage"] = strconv.Itoa(alr.Leverage)
	}

	res, err := as.request("POST", "/dapi/v1/leverage", params, true, true)
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
		Leverage int    `json:"leverage"`
		MaxQty   string `json:"maxQty"`
		Symbol   string `json:"symbol"`
	}{}
	if err := json.Unmarshal(textRes, &rawAdjustLeverage); err != nil {
		return nil, errors.Wrap(err, "rawAdjustLeverage unmarshal failed")
	}

	mnl, _ := intFromString(rawAdjustLeverage.MaxQty)

	al := &CoinAdjustLeverageInfo{
		Leverage: rawAdjustLeverage.Leverage,
		maxQty:   mnl,
		Symbol:   rawAdjustLeverage.Symbol,
	}
	return al, nil
}

func (as *apiService) CoinPositionMargin(pmr PositionMarginRequest) (*PositionMarginInfo, error) {
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

	res, err := as.request("POST", "/dapi/v1/positionMargin", params, true, true)
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

func (as *apiService) CoinUserTradesHistory(uth CoinUserTradesHistoryRequest) ([]*CoinUserTradesHistoryInfo, error) {
	params := make(map[string]string)
	params["symbol"] = uth.Symbol
	params["pair"] = uth.Pair
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

	res, err := as.request("GET", "dapi/v1/userTrades", params, true, true)
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
		Pair            string  `json:"pair"`
		MarginAsset     string  `json:"marginAsset"`
		Commission      string  `json:"commission"`
		CommissionAsset string  `json:"commissionAsset"`
		Id              int     `json:"id"`
		Maker           bool    `json:"maker"`
		OrderId         int     `json:"orderId"`
		Price           string  `json:"price"`
		Qty             string  `json:"qty"`
		BaseQty         string  `json:"baseQty"`
		RealizedPnl     string  `json:"realizedPnl"`
		Side            string  `json:"side"`
		PositionSide    string  `json:"positionSide"`
		Symbol          string  `json:"symbol"`
		Time            float64 `json:"time"`
	}{}
	if err := json.Unmarshal(textRes, &rawUserTradesHistory); err != nil {
		return nil, errors.Wrap(err, "rawUserTradesHistory unmarshal failed")
	}

	var fbi []*CoinUserTradesHistoryInfo

	for _, rt := range rawUserTradesHistory {

		t, err := timeFromUnixTimestampFloat(rt.Time)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Time")
		}
		commission, _ := floatFromString(rt.Commission)

		price, _ := floatFromString(rt.Price)

		qty, _ := floatFromString(rt.Qty)

		baseQty, _ := floatFromString(rt.BaseQty)

		realizedPnl, _ := floatFromString(rt.RealizedPnl)

		fbi = append(fbi, &CoinUserTradesHistoryInfo{
			Time:            t,
			Commission:      commission,
			MarginAsset:     rt.MarginAsset,
			Pair:            rt.Pair,
			Price:           price,
			Qty:             qty,
			BaseQty:         baseQty,
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

func (as *apiService) CoinAllFutureOrders(afo CoinAllFutureOrdersRequest) ([]*CoinHistoryExecutedFutureOrder, error) {
	params := make(map[string]string)
	params["symbol"] = afo.Symbol
	params["pair"] = afo.Pair
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

	res, err := as.request("GET", "dapi/v1/allOrders", params, true, true)
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
		Pair          string  `json:"pair"`
		CumBase       string  `json:"cumBase"`       // 成交金额
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
		UpdateTime    float64 `json:"updateTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "rawOrders unmarshal failed")
	}

	var eoc []*CoinHistoryExecutedFutureOrder
	for _, rawOrder := range rawOrders {

		cu, _ := floatFromString(rawOrder.CumBase)
		ex, _ := floatFromString(rawOrder.ExecutedQty)
		or, _ := floatFromString(rawOrder.OrigQty)
		pr, _ := floatFromString(rawOrder.Price)
		st, _ := floatFromString(rawOrder.StopPrice)
		ac, _ := floatFromString(rawOrder.ActivatePrice)
		pre, _ := floatFromString(rawOrder.PriceRate)
		t, _ := timeFromUnixTimestampFloat(rawOrder.UpdateTime)

		eo := &CoinHistoryExecutedFutureOrder{
			Symbol:        rawOrder.Symbol,
			CumBase:       cu,
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

			UpdateTime: t,
		}

		eoc = append(eoc, eo)
	}

	return eoc, nil
}
