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
