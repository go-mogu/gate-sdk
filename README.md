### Installation
```shell
go get -u github.com/go-mogu/gate-sdk
```
### Getting Started
Please follow the installation instruction and execute the following Go code:

#### restApi
```go
package test

import (
	"context"
	"github.com/antihax/optional"
	"github.com/gateio/gateapi-go/v6"
	"testing"
)

func TestGate(t *testing.T) {
	client := gateapi.NewAPIClient(gateapi.NewConfiguration())
	// uncomment the next line if your are testing against testnet
	//client.ChangeBasePath("https://fx-api-testnet.gateio.ws/api/v4")
	ctx := context.WithValue(context.Background(),
		gateapi.ContextGateAPIV4,
		gateapi.GateAPIV4{
			Key:    "YOU_KEY",
			Secret: "YOU_SECRET",
		},
	)
	candlesticks, _, err := client.FuturesApi.ListFuturesCandlesticks(ctx, "usdt", "SOL_USDT", &gateapi.ListFuturesCandlesticksOpts{
		From:     optional.Int64{},
		To:       optional.Int64{},
		Limit:    optional.Int32{},
		Interval: optional.NewString("15m"),
	})
	if err != nil {
		panic(err)
	}
	t.Log(candlesticks)
}

```

#### wsApi
```go
package test

import (
	"github.com/go-mogu/gate-sdk/consts"
	"github.com/go-mogu/gate-sdk/model"
	"github.com/go-mogu/gate-sdk/ws"
	"net/url"
	"testing"
)

const (
	Key    = "YOUR_API_KEY"
	Secret = "YOUR_API_SECRETY"
)

func TestGateWsKline(t *testing.T) {
	u := url.URL{Scheme: "wss", Host: "fx-ws.gateio.ws", Path: "/v4/ws/usdt"}
	client := new(ws.Client).Init(u.String(), func(message model.WsBaseRes) {
		t.Log("default:", message)
	}, func(message model.WsBaseRes) {
		t.Log("error:", message)
	})

	subReq := model.NewWsBaseReq("futures.candlesticks", consts.WsEventSubscribe, []string{"15m", "BTC_USDT"})
	client.Subscribe(subReq, func(message model.WsBaseRes) {
		t.Log(message)
	})
	positionsMsg := model.NewWsBaseReq("futures.positions", consts.WsEventSubscribe, []string{"USERID", "BTC_USDT"})
	positionsMsg.Sign(Key, Secret)
	client.Subscribe(positionsMsg, func(message model.WsBaseRes) {
		t.Log(message)
	})

	client.Connect()
}

```