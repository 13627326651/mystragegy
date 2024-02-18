package db

import (
	"encoding/json"
	"fmt"
	. "tinyquant/src/logger"
	"tinyquant/src/util"

	"go.uber.org/zap"
)

type UserParams struct {
	BuyVolume      float64 `json:"bv"`
	SellVolume     float64 `json:"sv"`
	BuyStartPrice  float64 `json:"bs"`
	SellStartPrice float64 `json:"ss"`
	BuyAllPrice    float64 `json:"bp"`
	SellAllPrice   float64 `json:"sp"`
}

func InsertOrderType(order_id int64, order_type util.ORIGIN_ORDER_STATUS) error {

	tx := GetRedisClient()

	id := fmt.Sprintf("%d", order_id)
	err := tx.HSet(util.OrderType, id, int(order_type)).Err()
	if err != nil {
		Logger.Error("insert order type  failed", zap.Error(err))
		return err
	}

	return nil
}

func DelOrderType(order_id int64) error {

	tx := GetRedisClient()

	id := fmt.Sprintf("%d", order_id)
	err := tx.HDel(util.OrderType, id).Err()
	if err != nil {
		Logger.Error("Del order type  failed", zap.Error(err))
		return err
	}

	return nil

}

func GetOrderType(order_id int64) (util.ORIGIN_ORDER_STATUS, error) {

	tx := GetRedisClient()

	id := fmt.Sprintf("%d", order_id)
	val, err := tx.HGet(util.OrderType, id).Int()
	if err != nil {
		Logger.Error("Get order type  failed", zap.Error(err))
		return util.ORIGIN_ORDER_STATUS(val), err
	}
	return util.ORIGIN_ORDER_STATUS(val), nil
}

func InsertCopyOrderID(order_id int64, copy_id map[string]int64) error {

	tx := GetRedisClient()

	id := fmt.Sprintf("%d", order_id)

	ret, err := json.Marshal(copy_id)
	if err != nil {
		Logger.Error("json marshal failed ", zap.Error(err))
		return err
	}
	err = tx.HSet(util.CopyOrderID, id, string(ret)).Err()
	if err != nil {
		Logger.Error("insert order type  failed", zap.Error(err))
		return err
	}

	return nil
}

func GetCopyOrderID(order_id int64) (map[string]int64, error) {

	tx := GetRedisClient()

	id := fmt.Sprintf("%d", order_id)

	val, err := tx.HGet(util.CopyOrderID, id).Bytes()
	if err != nil {
		Logger.Error("Get Copy id failed", zap.Error(err))
		return nil, err
	}

	ret := make(map[string]int64)

	err = json.Unmarshal(val, &ret)
	if err != nil {
		Logger.Error("json unmarshal failed ", zap.Error(err))
		return nil, err
	}

	return ret, nil
}

func DelCopyOrder(order_id int64) error {

	tx := GetRedisClient()

	id := fmt.Sprintf("%d", order_id)
	err := tx.HDel(util.CopyOrderID, id).Err()
	if err != nil {
		Logger.Error("Del Copy id  failed", zap.Error(err))
		return err
	}

	return nil

}

func InsertPositionParams(up *UserParams) error {

	tx := GetRedisClient()

	ret, err := json.Marshal(up)
	if err != nil {
		Logger.Error("json marshal failed ", zap.Error(err))
		return err
	}
	err = tx.Set(util.PositionParams, string(ret), 0).Err()
	if err != nil {
		Logger.Error("insert order type  failed", zap.Error(err))
		return err
	}

	return nil
}

func GetPositionParams() (*UserParams, error) {

	tx := GetRedisClient()

	val, err := tx.Get(util.PositionParams).Bytes()
	if err != nil {
		Logger.Error("Get Copy id failed", zap.Error(err))
		return nil, err
	}

	ret := new(UserParams)

	err = json.Unmarshal(val, &ret)
	if err != nil {
		Logger.Error("json unmarshal failed ", zap.Error(err))
		return nil, err
	}

	return ret, nil
}

func InsertTryBuyCount(tryBuyCount int) error {

	tx := GetRedisClient()

	err := tx.Set(util.TryBuyCount, tryBuyCount, 0).Err()
	if err != nil {
		Logger.Error("insert Try Buy Count failed", zap.Error(err))
		return err
	}
	return nil
}

func GetTryBuyCount() (int, error) {

	tx := GetRedisClient()

	val, err := tx.Get(util.TryBuyCount).Int()
	if err != nil {
		Logger.Error("Get Try Buy Count failed", zap.Error(err))
		return 0, err
	}

	return val, nil

}

func InsertTrySellCount(trySellCount int) error {

	tx := GetRedisClient()

	err := tx.Set(util.TrySellCount, trySellCount, 0).Err()
	if err != nil {
		Logger.Error("insert Try Sell Count failed", zap.Error(err))
		return err
	}
	return nil
}

func GetTrySellCount() (int, error) {

	tx := GetRedisClient()

	val, err := tx.Get(util.TrySellCount).Int()
	if err != nil {
		Logger.Error("Get Sell Buy Count failed", zap.Error(err))
		return 0, err
	}

	return val, nil

}

func InsertAllTryBuyCount(AllTryBuyCount int) error {

	tx := GetRedisClient()

	err := tx.Set(util.AllTryBuyCount, AllTryBuyCount, 0).Err()
	if err != nil {
		Logger.Error("insert All Try Buy Count failed", zap.Error(err))
		return err
	}
	return nil
}

func GetAllTryBuyCount() (int, error) {

	tx := GetRedisClient()

	val, err := tx.Get(util.AllTryBuyCount).Int()
	if err != nil {
		Logger.Error("Get All Try Buy Count failed", zap.Error(err))
		return 0, err
	}

	return val, nil

}

func InsertAllTrySellCount(AllTrySellCount int) error {

	tx := GetRedisClient()

	err := tx.Set(util.AllTrySellCount, AllTrySellCount, 0).Err()
	if err != nil {
		Logger.Error("insert All Try Sell Count failed", zap.Error(err))
		return err
	}
	return nil
}

func GetAllTrySellCount() (int, error) {

	tx := GetRedisClient()

	val, err := tx.Get(util.AllTrySellCount).Int()
	if err != nil {
		Logger.Error("Get All Try Sell Count failed", zap.Error(err))
		return 0, err
	}

	return val, nil

}
