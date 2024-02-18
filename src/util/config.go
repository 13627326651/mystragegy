package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitParam(follow bool) {
	InitConfig(follow)
	InitLogParam()
	// InitMysqlParams()
	InitQuantParam()
	// InitRedisParams()
	InitApiKey()
}

const (
	BinanceLeverage = 100 // 币安 U本位杠杆
)

type ORDER_TYPE_CONTROL int

const (
	BOTH = ORDER_TYPE_CONTROL(0)
	BUY  = ORDER_TYPE_CONTROL(1)
	SELL = ORDER_TYPE_CONTROL(2)
	NONE = ORDER_TYPE_CONTROL(3)
)

type ORIGIN_ORDER_STATUS int

const (
	COMMON          = ORIGIN_ORDER_STATUS(0)     // 手动单
	PIN             = ORIGIN_ORDER_STATUS(1)     // 插针单
	CLOSECOMMON     = ORIGIN_ORDER_STATUS(10)    // 普通平仓单
	PINCLOSECOMMON  = ORIGIN_ORDER_STATUS(100)   // 插针平仓单
	LOSSCLOSECOMMON = ORIGIN_ORDER_STATUS(1000)  // 止损平仓单
	FLOW            = ORIGIN_ORDER_STATUS(10000) // 顺势单
)

type ORIGIN_ORDER_FLAG int

const (
	UNKNNOW     = ORIGIN_ORDER_FLAG(0) // 手动
	ADDPOSITION = ORIGIN_ORDER_FLAG(1) // 开仓
	DELPOSITION = ORIGIN_ORDER_FLAG(2) // 平仓
)

const (
	ETHUSDT = "ETHUSDT"
	ETHBUSD = "ETHBUSD"
)

var ACCOUNTASSET = map[string]string{
	ETHUSDT: "USDT",
	ETHBUSD: "BUSD",
}

const COIN_ETHUSD = "ETHUSD_PERP"

const OrderType = "order_type"
const CopyOrderID = "copy_id"
const PositionParams = "pp"
const TryBuyCount = "tbc"
const TrySellCount = "tsc"
const AllTryBuyCount = "atbc"
const AllTrySellCount = "atsc"

var Order_Precision = map[string]int{
	"FILUSDT": 8,
	"BTCUSDT": 8,
	"ETHUSDT": 6,
}

// 账户信息推送事件
const (
	ListenKeyExpired      = "listenKeyExpired"
	MARGIN_CALL           = "MARGIN_CALL"
	ACCOUNT_UPDATE        = "ACCOUNT_UPDATE"
	ORDER_TRADE_UPDATE    = "ORDER_TRADE_UPDATE"
	ACCOUNT_CONFIG_UPDATE = "ACCOUNT_CONFIG_UPDATE"
)

var (
	BINANCE_API_KEY    string
	BINANCE_SECRET_KEY string
)

var (
	WXROBOTURL string
)

var (
	MysqlHost   string
	MysqlUser   string
	MysqlPass   string
	MysqlDBName string
)

var (
	Quantity                    float64
	Profits                     float64
	VolumeIncrease              float64
	VolumeIncreaseForClose      float64
	SpringPrice                 float64
	PlaceTest                   bool
	TerracedPrice               []float64 //连续开单T度
	CancelCloseOrderLevel       float64   //取消平仓单
	CreatCloseOrderLevel        float64   //创建平仓单
	IncreaseQuantityLevel       float64   //增加开仓
	DoubleCreatOrderLevel       float64   //增加开仓仓位满足Profits的倍数
	ContinuousOrderValidityTime int64
	SupportLevel                float64
	PressureLevel               float64
)

var (
	Console      bool
	File         bool
	Path         string
	FileLevel    string
	ConsoleLevel string
)

var (
	RedisHost string
	RedisPass string
)

type Api struct {
	Binance []*BinanceConfig
}

type BinanceConfig struct {
	ApiKey    string
	SecretKey string
	Quantity  float64
	Type      string
	Name      string
}

var DocumentaryApi *Api

func InitConfig(documentary bool) {

	viper.SetConfigName("config")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("use config : ", viper.AllSettings())

	if documentary {
		vjson := viper.New()
		vjson.SetConfigName("api")
		vjson.SetConfigType("json")
		vjson.AddConfigPath(".")

		if err := vjson.ReadInConfig(); err != nil {
			panic("read api config failed")
		}

		err = vjson.Unmarshal(&documentary)
		if err != nil {
			panic("Unmarshal config failed")
		}
	}
}

func InitMysqlParams() {

	MysqlHost = viper.GetString("mysql.Host")
	if MysqlHost == "" {
		panic("Get Host failed ")
	}

	MysqlUser = viper.GetString("mysql.User")
	if MysqlUser == "" {
		panic("Get User failed ")
	}

	MysqlPass = viper.GetString("mysql.Pass")
	if MysqlPass == "" {
		panic("Get Pass failed ")
	}

	MysqlDBName = viper.GetString("mysql.DBName")
	if MysqlDBName == "" {
		panic("Get DBName failed ")
	}
}

func InitRedisParams() {
	RedisHost = viper.GetString("redis.Host")
	if RedisHost == "" {
		panic("Get Host failed ")
	}

	viper.SetDefault("redis.Pass", "unknow")
	RedisPass = viper.GetString("redis.Pass")
	if RedisPass == "unknow" {
		panic("Get Pass failed ")
	}
}

func InitQuantParam() {
	viper.SetDefault("quant.Quantity", 0.01)
	Quantity = viper.GetFloat64("quant.Quantity")
	viper.SetDefault("quant.Profits", 0.0125)
	Profits = viper.GetFloat64("quant.Profits")
	viper.SetDefault("quant.VolumeIncrease", 5.0)
	VolumeIncrease = viper.GetFloat64("quant.VolumeIncrease")
	viper.SetDefault("quant.VolumeIncreaseForClose", 4.0)
	VolumeIncreaseForClose = viper.GetFloat64("quant.VolumeIncreaseForClose")
	viper.SetDefault("quant.VolumeIncrease", 0.0025)
	SpringPrice = viper.GetFloat64("quant.SpringPrice")
	viper.SetDefault("quant.PlaceTest", true)
	PlaceTest = viper.GetBool("quant.PlaceTest")
	TerracedPrice = make([]float64, 0)
	TerracedPrice = append(TerracedPrice, viper.GetFloat64("quant.TerracedPrice0"))
	TerracedPrice = append(TerracedPrice, viper.GetFloat64("quant.TerracedPrice1"))
	TerracedPrice = append(TerracedPrice, viper.GetFloat64("quant.TerracedPrice2"))
	TerracedPrice = append(TerracedPrice, viper.GetFloat64("quant.TerracedPrice3"))
	TerracedPrice = append(TerracedPrice, viper.GetFloat64("quant.TerracedPrice4"))
	viper.SetDefault("quant.CancelCloseOrderLevel", 0.03)
	CancelCloseOrderLevel = viper.GetFloat64("quant.CancelCloseOrderLevel")
	viper.SetDefault("quant.CreatCloseOrderLevel", 0.02)
	CreatCloseOrderLevel = viper.GetFloat64("quant.CreatCloseOrderLevel")
	viper.SetDefault("quant.IncreaseQuantityLevel", 0.04)
	IncreaseQuantityLevel = viper.GetFloat64("quant.IncreaseQuantityLevel")
	viper.SetDefault("quant.DoubleCreatOrderLevel", 2.0)
	DoubleCreatOrderLevel = viper.GetFloat64("quant.DoubleCreatOrderLevel")
	viper.SetDefault("quant.ContinuousOrderValidityTime", 10)
	ContinuousOrderValidityTime = viper.GetInt64("quant.ContinuousOrderValidityTime")

	viper.SetDefault("quant.PressureLevel", 1360)
	PressureLevel = viper.GetFloat64("quant.PressureLevel")

	viper.SetDefault("quant.SupportLevel", 1310)
	SupportLevel = viper.GetFloat64("quant.SupportLevel")

}

func InitLogParam() {
	viper.SetDefault("log.Console", true)
	Console = viper.GetBool("log.Console")
	viper.SetDefault("log.File", true)
	File = viper.GetBool("log.File")
	viper.SetDefault("log.Path", "./log/zap.log")
	Path = viper.GetString("log.Path")
	viper.SetDefault("log.FileLevel", "debug")
	FileLevel = viper.GetString("log.FileLevel")
	viper.SetDefault("log.ConsoleLevel", "debug")
	ConsoleLevel = viper.GetString("log.ConsoleLevel")

	fmt.Printf("%v %v %v %v %v\n", Console, File, Path, FileLevel, ConsoleLevel)
}

func InitApiKey() {

	BINANCE_API_KEY = viper.GetString("system.ApiKey")
	if BINANCE_API_KEY == "" {
		panic("Get ApiKey  failed ")
	}
	BINANCE_SECRET_KEY = viper.GetString("system.SecretKey")
	if BINANCE_SECRET_KEY == "" {
		panic("Get secretKey failed ")
	}

	WXROBOTURL = viper.GetString("system.robot")
	if WXROBOTURL == "" {
		panic("get robot url failed ")
	}
}
