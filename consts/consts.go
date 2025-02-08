package consts

const (
	Long  = "long"
	Short = "short"
	Buy   = "buy"
	Sell  = "sell"
)

const (
	GateLong            = "dual_long"
	GateShort           = "dual_short"
	CloseShort          = "close_short"
	CloseLong           = "close_long"
	CloseShortPosition  = "close-short-position"
	CloseLongPosition   = "close-long-position"
	TimerIntervalSecond = 5
	ReconnectWaitSecond = 60

	RestTestUrl = "https://fx-api-testnet.gateio.ws/api/v4"

	WsBaseUrl = "wss://fx-ws.gateio.ws/v4/ws"

	WsEventUnsubscribe = "unsubscribe"
	WsEventSubscribe   = "subscribe"
	WsEventUpdate      = "update"
)
