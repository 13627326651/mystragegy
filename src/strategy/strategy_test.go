package strategy_test

import (
	"tinyquant/src/db"
	"tinyquant/src/logger"
	fb "tinyquant/src/quant/future_binance"
	"tinyquant/src/util"
)

var Binance fb.Binance

func init() {

	util.InitParam(false)

	db.InitMysql()

	logger.InitLogger()
	Binance = fb.Binance{}
	Binance.InitBinance(util.BINANCE_API_KEY, util.BINANCE_SECRET_KEY)

}
