package tcpnetwork

import (
	// "errors"
	"net"
	"sync/atomic"
	"time"
)

const (
	kServerConf_SendBufferSize int = 1024
)

type TCPNetworkConf struct {
	SendBufferSize int
}

type TCPNetwork struct {
	Conf           TCPNetworkConf
	EventQueue     chan *ConnEvent
	connId         int
	conns          map[int]*Connection
	readTimeoutSec int // accept
	headerLength   int
	shutdownFlag   int32
}

func NewTCPNetwork(eventQueueSize int, headerLength int) *TCPNetwork {
	s := &TCPNetwork{}
	//	default config
	s.Conf.SendBufferSize = kServerConf_SendBufferSize
	s.EventQueue = make(chan *ConnEvent, eventQueueSize)
	s.connId = 0
	s.conns = make(map[int]*Connection)
	s.readTimeoutSec = 0
	s.shutdownFlag = 0 //
	s.headerLength = headerLength
	return s
}

// Push implements the IEventQueue interface
func (t *TCPNetwork) Push(evt *ConnEvent) {
	if nil == t.EventQueue {
		return
	}

	//	push timeout
	select {
	case t.EventQueue <- evt:
		{

		}
	case <-time.After(time.Second * 5):
		{
			// TODO:
			// evt.Conn.close()
		}
	}
}

// Pop the event in event queue
func (t *TCPNetwork) Pop() *ConnEvent {
	evt, ok := <-t.EventQueue
	if !ok {
		//	event queue already closed
		return nil
	}

	return evt
}

// GetEventQueue get the event queue channel
func (t *TCPNetwork) GetEventQueue() <-chan *ConnEvent {
	return t.EventQueue
}

func (t *TCPNetwork) GetConn(id int) *Connection {
	if v, ok := t.conns[id]; ok {
		return v
	}
	return nil
}

// Connect the remote server
func (t *TCPNetwork) Connect(addr string) (int, error) {
	conn, err := net.Dial("tcp", addr)
	if nil != err {
		return 0, err
	}
	t.connId++
	connection := t.createConn(conn, t.connId)
	connection.run()
	connection.init()

	return connection.ConnId, nil
}

// Shutdown frees all connection and stop the listener
func (t *TCPNetwork) Shutdown() {
	if !atomic.CompareAndSwapInt32(&t.shutdownFlag, 0, 1) {
		return
	}

	//	close all connections
	t.disconnectAllConnections()
}

func (t *TCPNetwork) disconnectAllConnections() {
	for k, c := range t.conns {
		c.Close()
		delete(t.conns, k)
	}
}

func (t *TCPNetwork) createConn(c net.Conn, id int) *Connection {
	conn := newConnection(c, t.Conf.SendBufferSize, t, id)
	return conn
}

// ServeWithHandler process all events in the event queue and dispatch to the IEventHandler

// func (t *TCPNetwork) onEvent(evt *ConnEvent) {
// 	switch evt.EventType {
// 	case KConnEvent_Connected:
// 		{
// 			//	add to connection map
// 			connId := 0
// 			if kServerConn == evt.Conn.from {
// 				connId = t.connIdForServer + 1
// 				t.connIdForServer = connId
// 				t.connsForServer[connId] = evt.Conn
// 			} else {
// 				connId = t.connIdForClient + 1
// 				t.connIdForClient = connId
// 				t.connsForClient[connId] = evt.Conn
// 			}
// 			evt.Conn.connId = connId

// 			handler.OnConnected(evt)
// 		}
// 	case KConnEvent_Disconnected:
// 		{
// 			handler.OnDisconnected(evt)

// 			//	remove from connection map
// 			if kServerConn == evt.Conn.from {
// 				delete(t.connsForServer, evt.Conn.connId)
// 			} else {
// 				delete(t.connsForClient, evt.Conn.connId)
// 			}
// 		}
// 	case KConnEvent_Data:
// 		{
// 			handler.OnRecv(evt)
// 		}
// 	}
// }

type TCPNetworkServer struct {
	TCPNetwork
	listener net.Listener
}

func NewTCPNetworkServer(eventQueueSize int, headerLength int) *TCPNetworkServer {
	s := &TCPNetworkServer{}
	s.Conf.SendBufferSize = kServerConf_SendBufferSize
	s.EventQueue = make(chan *ConnEvent, eventQueueSize)
	s.connId = 0
	s.conns = make(map[int]*Connection)
	s.readTimeoutSec = 0
	s.shutdownFlag = 0
	s.headerLength = headerLength
	//	default config
	return s
}

// Listen an address to accept client connection
func (ts *TCPNetworkServer) Listen(addr string) error {
	ls, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}

	//	accept
	ts.listener = ls
	go ts.acceptRoutine()
	return err
}

func (ts *TCPNetworkServer) acceptRoutine() {
	// after accept temporary failure, enter sleep and try again
	var tempDelay time.Duration

	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			// check if the error is an temporary error
			if acceptErr, ok := err.(net.Error); ok && acceptErr.Temporary() {
				if 0 == tempDelay {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				logWarn("Accept error %s , retry after %d ms", acceptErr.Error(), tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			logError("accept routine quit.error:%s", err.Error())
			ts.listener = nil
			return
		}

		//	process conn event
		ts.connId++
		connection := ts.createConn(conn, ts.connId)
		connection.SetReadTimeoutSec(ts.readTimeoutSec)
		connection.init()
		connection.run()
	}
}

func (ts *TCPNetworkServer) Shutdown() {
	ts.TCPNetwork.Shutdown()

	if nil != ts.listener {
		ts.listener.Close()
	}
}

// func (ts *TCPNetworkServer) onEvent(evt *ConnEvent) {
// SERVE_LOOP:
// 	for {
// 		select {
// 		case evt, ok := <-ts.eventQueue:
// 			{
// 				if !ok {
// 					//	channel closed or shutdown
// 					break SERVE_LOOP
// 				}

// 				ts.TCPNetwork.onEvent(evt)
// 			}
// 		}
// 	}
// }
