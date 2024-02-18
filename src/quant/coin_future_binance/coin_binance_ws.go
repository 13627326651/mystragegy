package future

import (
	"fmt"
	"os"
	"os/signal"
	"tinyquant/src/mod"

	"github.com/rootpd/binance"
)

// 获取 交易流
func (b *Binance) GetTradeWs() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	kech, done, err := b.CoinFutureTradeWebsocket(binance.TradeWebsocketRequest{
		Symbol: "FILUSDT",
	})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case ke := <-kech:
				fmt.Println("xxxx", ke)
			case <-done:
				break
			}
		}
	}()
	fmt.Println("canceling context")
	<-done
	fmt.Println("exit")
	<-interrupt
	fmt.Println("interrupt")
}

// 获取 深度
func (b *Binance) GetFutureDepthWs(symbol string) (chan *binance.DepthEvent, chan struct{}) {

	kech, done, err := b.CoinFutureDepthWebsocket(binance.DepthWebsocketRequest{
		Symbol: symbol,
	})
	if err != nil {
		panic(err)
	}

	return kech, done
}

func (b *Binance) GetKlineWs(symbol string, interval binance.Interval) chan *mod.Kline {

	binance_kline := make(chan *mod.Kline)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	kech, done, err := b.CoinFutureKlineWebsocket(binance.KlineWebsocketRequest{
		Symbol:   symbol,
		Interval: interval,
	})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case ke := <-kech:
				t := &mod.Kline{
					StartTime:   ke.OpenTime,
					CloseTime:   ke.CloseTime,
					Volume:      ke.Volume,
					BuyVolume:   ke.TakerBuyBaseAssetVolume,
					SellVolume:  ke.Volume - ke.TakerBuyBaseAssetVolume,
					Quote:       ke.QuoteAssetVolume,
					BuyQuote:    ke.TakerBuyQuoteAssetVolume,
					SellQuote:   ke.QuoteAssetVolume - ke.TakerBuyQuoteAssetVolume,
					TradeNumber: ke.NumberOfTrades,
					Open:        ke.Open,
					Close:       ke.Close,
					High:        ke.High,
					Low:         ke.Low,
					Final:       ke.Final,
				}

				binance_kline <- t

			case <-done:
				break
			}
		}
	}()

	return binance_kline
}
