package util_test

import (
	"fmt"
	"testing"
	"tinyquant/src/util"
)

func Test_SendCardMsg(t *testing.T) {

	msg := fmt.Sprintf("订单类型 : 开仓 \n订单品种 : %s  \n订单方向 : %s  \n成交价格 : %f  \n成交数量 :  %f \n手续费 %f ", "ETHUSDT", "BUY", 3899.24, 0.01, 0.07312312)

	util.SendOrderMsg(msg)

}
