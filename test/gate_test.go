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
