package ws

import (
	"github.com/go-mogu/gate-sdk/model"
)

type Client struct {
	gateBaseWsClient *GateBaseWsClient
}

func (p *Client) Init(uri string, listener OnReceive, errorListener OnReceive) *Client {
	p.gateBaseWsClient = new(GateBaseWsClient).Init(uri)
	p.gateBaseWsClient.SetListener(listener, errorListener)
	p.gateBaseWsClient.ConnectWebSocket()

	return p

}

func (p *Client) Connect() *Client {
	p.gateBaseWsClient.Connect()
	return p
}

func (p *Client) UnSubscribe(req *model.WsBaseReq) {
	delete(p.gateBaseWsClient.ScribeMap, req.Channel)
	p.SendMessageByType(req)
}

func (p *Client) Subscribe(req *model.WsBaseReq, listener OnReceive) {
	p.gateBaseWsClient.AllScribeSet[req] = struct{}{}
	p.gateBaseWsClient.ScribeMap[req.Channel] = listener
	p.gateBaseWsClient.SendByType(req)
}

func (p *Client) SendMessageByType(req *model.WsBaseReq) {
	p.gateBaseWsClient.SendByType(req)
}
