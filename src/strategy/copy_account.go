package strategy

import (
	. "tinyquant/src/logger"
	fb "tinyquant/src/quant/future_binance"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

type DocumentaryAccount struct {
	Account []*BinanceFutureAsset
}

func (acc *DocumentaryAccount) InitDocumentaryAccount(symbol string) error {

	if util.DocumentaryApi == nil {
		panic("Get api list failed")
	}

	if len(util.DocumentaryApi.Binance) == 0 {
		return nil
	}

	apilist := util.DocumentaryApi.Binance
	for _, cfg := range apilist {

		Acc := &BinanceFutureAsset{}
		var s string

		switch cfg.Type {
		case "usdt":
			Acc.Binance = &fb.Binance{}
			s = symbol
			// case "coin":
			// 	Acc.Binance = &coin_fb.Binance{}
			// 	s = util.COIN_ETHUSD
		}

		Acc.Quantity = util.Round(cfg.Quantity, 3)
		Acc.Name = cfg.Name
		Acc.Type = cfg.Type

		Acc.Binance.InitBinance(cfg.ApiKey, cfg.SecretKey)

		// 调整当前杠杆倍数
		err := Acc.Binance.AdjustBinanceLeverage(s, util.BinanceLeverage)
		if err != nil {
			Logger.Error("adjust leverage failed", zap.Error(err))
		}

		// 更改持仓模式

		err = Acc.Binance.ChangeBinanceUserPositionSide(binance.PosithonBothSide)
		if err != nil {
			Logger.Error("change user position side failed", zap.Error(err))
			// return err
		}

		// 改变全仓模式

		err = Acc.Binance.ChangeBinanceMarginType(s, binance.POSITION_ISOLATED) // 全仓
		if err != nil {
			Logger.Error("change margin type failed", zap.Error(err))
			//return err
		}

		ba, err := Acc.Binance.GetFutureBalance()

		if err != nil {
			Logger.Error("CopyAccount :", zap.Any("account", Acc.Name))
			Logger.Error("Get Future balance failed ", zap.Error(err))
			return err
		}

		for _, v := range ba {

			if util.ACCOUNTASSET[symbol] == v.Asset {

				Acc.Asset = v.Asset
				Acc.Balance = v.Balance                       // 总余额，包括已持仓的和当前盈利
				Acc.CrossWalletBalance = v.CrossWalletBalance // 全仓余额
				Acc.CrossUnPnl = v.CrossUnPnl                 // 未实现盈亏
				Acc.AvailableBalance = v.AvailableBalance     // 下单可用余额
				Acc.MaxWithdrawAmount = v.MaxWithdrawAmount
				Acc.MarginAvailable = v.MarginAvailable

			}
		}

		Logger.Info("CopyTrade 初始化成功 :", zap.Any("account", Acc.Name))
		acc.Account = append(acc.Account, Acc)

	}

	return nil

}
