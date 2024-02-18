module tinyquant

go 1.12

require (
	github.com/Chronokeeper/anyxml v0.0.0-20160530174208-54457d8e98c6 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet v2.1.2+incompatible // indirect
	github.com/agrison/go-tablib v0.0.0-20160310143025-4930582c22ee // indirect
	github.com/agrison/mxj v0.0.0-20160310142625-1269f8afb3b4 // indirect
	github.com/bndr/gotabulate v1.1.2 // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-redis/redis/v7 v7.4.1
	github.com/go-sql-driver/mysql v1.4.0
	github.com/huobirdcenter/huobi_golang v0.0.0-20210226095227-8a30a95b6d0d
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robfig/cron v1.2.0
	github.com/rootpd/binance v0.0.0-20171024115603-c656b55bcff4
	github.com/spf13/viper v1.7.0
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/urfave/cli/v2 v2.4.0
	github.com/xormplus/builder v0.0.0-20200331055651-240ff40009be // indirect
	github.com/xormplus/xorm v0.0.0-20210822100304-4e1d4fcc1e67
	go.uber.org/zap v1.21.0
	gopkg.in/flosch/pongo2.v3 v3.0.0-20141028000813-5e81b817a0c4 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/rootpd/binance => ./src/ex/binance

//replace github.com/huobirdcenter/huobi_golang => ./src/ex/huobi_golang

//replace github.com/nntaoli-project/goex => ./src/ex/goex
