package future

import (
	"time"

	"github.com/rootpd/binance"
)

// 获取账户订单推送
func (b *Binance) GetAccountWs() (chan *binance.FutureAccountEvent, chan struct{}) {

	stream, err := b.StartFutureUserDataStream()
	if err != nil {
		panic(err)
	}

	kech, done, err := b.FutureUserDataWebsocket(binance.UserDataWebsocketRequest{
		ListenKey: stream.ListenKey,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(50 * time.Minute)

		defer ticker.Stop()

		for {
			<-ticker.C

			b.KeepAliveFutureUserDataStream(stream)
		}
	}()

	return kech, done
}
