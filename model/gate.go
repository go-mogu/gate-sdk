package model

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
)

type WsBaseReq struct {
	Time    int64       `json:"time"`
	Channel string      `json:"channel"`
	Event   string      `json:"event"`
	Payload []string    `json:"payload"`
	Auth    *GateWsAuth `json:"auth"`
}

type GateWsAuth struct {
	Method string `json:"method"`
	KEY    string `json:"KEY"`
	SIGN   string `json:"SIGN"`
}

func NewWsBaseReq(channel, event string, payload []string) *WsBaseReq {
	return &WsBaseReq{
		Channel: channel,
		Event:   event,
		Payload: payload,
	}
}

func sign(secret, channel, event string, t int64) string {
	message := fmt.Sprintf("channel=%s&event=%s&time=%d", channel, event, t)
	h2 := hmac.New(sha512.New, []byte(secret))
	_, _ = io.WriteString(h2, message)
	return hex.EncodeToString(h2.Sum(nil))
}

func (req *WsBaseReq) Sign(key, secret string) {
	signStr := sign(secret, req.Channel, req.Event, req.Time)
	req.Auth = &GateWsAuth{
		Method: "api_key",
		KEY:    key,
		SIGN:   signStr,
	}
}

type WsBaseRes struct {
	Time    int         `json:"time"`
	TimeMs  int64       `json:"time_ms"`
	Channel string      `json:"channel"`
	Event   string      `json:"event"`
	Result  interface{} `json:"result"`
}

type Candles struct {
	A float64 `json:"a"`
	C float64 `json:"c"`
	H float64 `json:"h"`
	L float64 `json:"l"`
	N string  `json:"n"`
	O float64 `json:"o"`
	T int64   `json:"t"`
	V int     `json:"v"`
}
