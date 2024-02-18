package huobi_test

import (
	"testing"
	"tinyquant/src/quant/huobi"
)

func Test_GetAccountInfo(t *testing.T) {

	huobi.GetAccountInfo()
}

func Test_GetAccountBalance(t *testing.T) {
	huobi.GetAccountBalance()
}

func Test_GetAllSymbolTickers(t *testing.T) {
	huobi.GetAllSymbolTickers()
}

func Test_GetSymbolKline(t *testing.T) {
	huobi.GetSymbolKline()
}

func Test_GetKlineWs(t *testing.T) {
	huobi.GetKlineWs()
}
