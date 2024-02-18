package future

import "github.com/rootpd/binance"

// 获取账户订单推送
func (b *Binance) GetAccountWs() (chan *binance.FutureAccountEvent, chan struct{}) {

	stream, err := b.StartCoinFutureUserDataStream()
	if err != nil {
		panic(err)
	}

	kech, done, err := b.CoinFutureUserDataWebsocket(binance.UserDataWebsocketRequest{
		ListenKey: stream.ListenKey,
	})
	if err != nil {
		panic(err)
	}

	return kech, done
}
