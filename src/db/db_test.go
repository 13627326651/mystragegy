package db

import (
	"testing"
	"time"
)

func init() {

	//	logger.InitLogger()

	//util.InitParam()
	InitMysql()
}

func Test_ConnectDB(t *testing.T) {
	InitMysql()
}

func Test_UpdateOrder(t *testing.T) {
	err := UpdateOrder(&Order{
		OrderID:       8389765512700395223,
		Symbol:        "ETHUSDT",
		IsOrderFinish: 0,
		UpdateTime:    time.Now(),
	})
	if err != nil {
		t.Error(err)
	}
}
