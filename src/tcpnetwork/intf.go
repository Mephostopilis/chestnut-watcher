package tcpnetwork

type IEventQueue interface {
	Push(*ConnEvent)
	Pop() *ConnEvent
}

type IEventHandler interface {
	OnConnected(evt *ConnEvent)
	OnDisconnected(evt *ConnEvent)
	OnRecv(evt *ConnEvent)
}

type IUnpacker interface {
	Unpack(*Connection, []byte) ([]byte, error)
}
