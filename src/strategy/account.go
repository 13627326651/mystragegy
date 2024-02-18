package strategy

import (
	"os"
	"sync"
	"time"
	. "tinyquant/src/logger"
	quant "tinyquant/src/quant"
	fb "tinyquant/src/quant/future_binance"
	"tinyquant/src/util"

	"github.com/rootpd/binance"
	"go.uber.org/zap"
)

type BinanceFutureAsset struct {
	*sync.RWMutex
	Symbol   string
	Binance  quant.Binance
	Name     string
	Type     string
	Quantity float64

	Asset              string  // 资产
	Balance            float64 // 总余额
	CrossWalletBalance float64 // 全仓余额
	CrossUnPnl         float64 // 全仓持仓未实现盈亏
	AvailableBalance   float64 // 下单可用余额
	MaxWithdrawAmount  float64 // 最大可转出余额
	MarginAvailable    bool    // 是否可用作联合保证金
	//PositionQty        float64 // 持仓量
}

func (acc *BinanceFutureAsset) InitAccount(symbol string) error {
	acc.Symbol = symbol
	if f, ok := os.LookupEnv("futures"); ok {

		if f == "usdt" {
			Logger.Info("tinyquant u本位 启动")
			Binance = &fb.Binance{}
		}

		// if f == "coin" {
		// Logger.Info("tinyquant 币本位 启动")
		// Binance = &coin_fb.Binance{}
		// }

	} else {
		panic("Get Binance API failed")
	}

	//客户端初始化
	Binance.InitBinance(util.BINANCE_API_KEY, util.BINANCE_SECRET_KEY)

	// 调整当前杠杆倍数 403?
	err := Binance.AdjustBinanceLeverage(symbol, util.BinanceLeverage)
	if err != nil {
		Logger.Error("adjust leverage failed", zap.Error(err))
		// return err
	}

	// 更改持仓模式
	err = Binance.ChangeBinanceUserPositionSide(binance.PosithonBothSide)
	if err != nil {
		Logger.Error("change user position side failed", zap.Error(err))
		// return err
	}

	// 改变全仓模式
	err = Binance.ChangeBinanceMarginType(symbol, binance.POSITION_CROSSED) // 全仓
	if err != nil {
		Logger.Error("change margin type failed", zap.Error(err))
		//return err
	}
	if err = acc.LoadAccount(); err != nil {
		return err
	}
	// acc.ReloadAccount()
	return nil
}

func (acc *BinanceFutureAsset) ReloadAccount() {
	timer := time.NewTimer(1 * time.Second)
	go func() {
		for {
			select {
			case <-timer.C:
				acc.LoadAccount()
				timer.Reset(1 * time.Second)
			}
		}
	}()
}

func (acc *BinanceFutureAsset) LoadAccount() error {
	// 获取账户余额
	ba, err := Binance.GetFutureBalance()
	if err != nil {
		Logger.Error("Get Future balance failed ", zap.Error(err))
		return err
	}
	acc.Lock()
	defer acc.Unlock()
	for _, v := range ba {
		if util.ACCOUNTASSET[acc.Symbol] == v.Asset {
			// Logger.Info("当前资产 :", zap.Any("asset", v))
			acc.Asset = v.Asset
			acc.Balance = v.Balance                       // 总余额，包括已持仓的和当前盈利
			acc.CrossWalletBalance = v.CrossWalletBalance // 全仓余额
			acc.CrossUnPnl = v.CrossUnPnl                 // 未实现盈亏
			acc.AvailableBalance = v.AvailableBalance     // 下单可用余额
			acc.MaxWithdrawAmount = v.MaxWithdrawAmount
			acc.MarginAvailable = v.MarginAvailable
		}
	}
	return nil
}
