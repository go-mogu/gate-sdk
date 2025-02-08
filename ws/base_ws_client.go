package ws

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-mogu/gate-sdk/consts"
	"github.com/go-mogu/gate-sdk/model"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron"
	"log"
	"sync"
	"time"
)

type GateBaseWsClient struct {
	ctx              context.Context
	NeedLogin        bool
	Connection       bool
	LoginStatus      bool
	Url              string
	Listener         OnReceive
	ErrorListener    OnReceive
	Ticker           *time.Ticker
	SendMutex        *sync.Mutex
	WebSocketClient  *websocket.Conn
	LastReceivedTime time.Time
	AllScribeSet     map[*model.WsBaseReq]struct{}
	ScribeMap        map[string]OnReceive
	c                *cron.Cron
}

func (p *GateBaseWsClient) Init(uri string) *GateBaseWsClient {
	p.ctx = context.Background()
	p.Connection = false
	p.Url = uri
	p.AllScribeSet = make(map[*model.WsBaseReq]struct{})
	p.ScribeMap = make(map[string]OnReceive)
	p.SendMutex = &sync.Mutex{}
	p.Ticker = time.NewTicker(consts.TimerIntervalSecond * time.Second)
	return p
}

func (p *GateBaseWsClient) SetListener(msgListener OnReceive, errorListener OnReceive) {
	p.Listener = msgListener
	p.ErrorListener = errorListener
}

func (p *GateBaseWsClient) Connect() {
	p.tickerLoop()
}

func (p *GateBaseWsClient) ConnectWebSocket() {
	var err error
	log.Default().Println("WebSocket connecting...")
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
	p.WebSocketClient, _, err = websocket.DefaultDialer.Dial(p.Url, nil)
	if err != nil {
		fmt.Printf("WebSocket connected error: %s\n", err)
		return
	}
	log.Default().Println("WebSocket connected")
	p.Connection = true
	p.LastReceivedTime = time.Now()
	p.StartReadLoop()
	p.ExecutePing()
}

func (p *GateBaseWsClient) StartReadLoop() {
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				log.Default().Println(exception)
			}
		}()
		p.ReadLoop()
	}()
}

func (p *GateBaseWsClient) ExecutePing() {
	c := cron.New()
	_ = c.AddFunc("*/15 * * * * *", p.ping)
	c.Start()
	p.c = c
}
func (p *GateBaseWsClient) ping() {
	ping := &model.WsBaseReq{
		Channel: "futures.ping",
		Event:   "",
		Payload: make([]string, 0),
	}
	p.SendByType(ping)
}

func (p *GateBaseWsClient) SendByType(req *model.WsBaseReq) {
	req.Time = time.Now().Unix()
	msgByte, _ := json.Marshal(req)
	p.Send(msgByte)
}

func (p *GateBaseWsClient) Send(data []byte) {
	if p.WebSocketClient == nil {
		log.Default().Println("WebSocket sent error: no connection available")
		return
	}
	p.SendMutex.Lock()
	err := p.WebSocketClient.WriteMessage(websocket.TextMessage, data)
	p.SendMutex.Unlock()
	if err != nil {
		log.Default().Printf("WebSocket sent error: data=%s, error=%s", data, err)
	}
}

func (p *GateBaseWsClient) tickerLoop() {
	log.Default().Println("tickerLoop started")
	for {
		select {
		case <-p.Ticker.C:
			elapsedSecond := time.Now().Sub(p.LastReceivedTime).Seconds()

			if elapsedSecond > consts.ReconnectWaitSecond {
				log.Default().Println("WebSocket reconnect...")
				p.disconnectWebSocket()
				p.ConnectWebSocket()
				time.Sleep(5 * time.Second)
				for req, _ := range p.AllScribeSet {
					p.SendByType(req)
				}
			}
		}
	}
}

func (p *GateBaseWsClient) disconnectWebSocket() {
	if p.WebSocketClient == nil {
		return
	}

	log.Default().Println("WebSocket disconnecting...")
	err := p.WebSocketClient.Close()
	if err != nil {
		log.Default().Printf("WebSocket disconnect error: %s\n", err)
		return
	}

	log.Default().Println("WebSocket disconnected")
}

func (p *GateBaseWsClient) ReadLoop() {
	for {

		if p.WebSocketClient == nil {
			log.Default().Println("Read error: no connection available")
			//time.Sleep(TimerIntervalSecond * time.Second)
			continue
		}

		_, buf, err := p.WebSocketClient.ReadMessage()
		if err != nil {
			log.Default().Printf("Read error: %s", err)
			continue
		}
		p.LastReceivedTime = time.Now()
		var message model.WsBaseRes
		err = json.Unmarshal(buf, &message)
		if err != nil {
			log.Default().Printf("Read error: %s", err)
			continue
		}

		if message.Channel == "futures.pong" {
			//g.Log().Debug(p.ctx, "Keep connected:", message)
			continue
		}

		listener := p.GetListener(message.Channel)
		listener(message)
	}

}

func (p *GateBaseWsClient) GetListener(channel string) OnReceive {

	v, e := p.ScribeMap[channel]

	if !e {
		return p.Listener
	}
	return v
}

type OnReceive func(message model.WsBaseRes)
