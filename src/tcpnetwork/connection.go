package tcpnetwork

import (
	// "fmt"
	// "io"
	"net"
	"sync/atomic"
	"time"

	"runtime/debug"
)

const (
	kConnStatus_None = iota
	kConnStatus_Connected
	kConnStatus_Disconnected
)

// All connection event
const (
	KConnEvent_None = iota
	KConnEvent_Connected
	KConnEvent_Disconnected
	KConnEvent_Data
	KConnEvent_Close
	KConnEvent_Total
)

const (
	kConnConf_DefaultSendTimeoutSec = 5
	kConnConf_MaxReadBufferLength   = 0xffff // 0xffff
)

// Send method flag
const (
	KConnFlag_CopySendBuffer = 1 << iota // do not copy the send buffer
	KConnFlag_NoHeader                   // do not append stream header
)

// send task
type sendTask struct {
	data []byte
	flag int64
}

// ConnEvent represents a event occurs on a connection, such as connected, disconnected or data arrived
type ConnEvent struct {
	ConnId    int
	EventType int
	Tag       int
	Session   int
	Msg       interface{}
}

func newConnEvent(et int, connId int, d interface{}) *ConnEvent {
	return &ConnEvent{
		EventType: et,
		ConnId:    connId,
		Msg:       d,
	}
}

//	Sync event callback
//	If return true, this event will not be sent to event channel
//	If return false, this event will be sent to event channel again
type FuncSyncExecute func(*ConnEvent) bool

// Connection is a wrap for net.Conn and process read and write task of the conn
// When event occurs, it will call the eventQueue to dispatch event
type Connection struct {
	ConnId              int
	StreamProtocol      *StreamProtocol
	network             *TCPNetwork
	conn                net.Conn
	status              int32
	sendMsgQueue        chan *sendTask
	sendTimeoutSec      int
	maxReadBufferLength int
	userdata            interface{}
	readTimeoutSec      int
	fnSyncExecute       FuncSyncExecute
	disableSend         int32
	localAddr           string
	remoteAddr          string
}

func newConnection(c net.Conn, sendBufferSize int, t *TCPNetwork, id int) *Connection {
	return &Connection{
		ConnId:              id,
		StreamProtocol:      NewStreamProtocol(KStreamProtocol2HeaderLength),
		network:             t,
		conn:                c,
		status:              kConnStatus_None,
		sendMsgQueue:        make(chan *sendTask, sendBufferSize),
		sendTimeoutSec:      kConnConf_DefaultSendTimeoutSec,
		maxReadBufferLength: kConnConf_MaxReadBufferLength,
	}
}

func (c *Connection) init() {
	c.localAddr = c.conn.LocalAddr().String()
	c.remoteAddr = c.conn.RemoteAddr().String()
}

//	directly close, packages in queue will not be sent
func (c *Connection) close() {
	//	set the disconnected status, use atomic operation to avoid close twice
	if atomic.CompareAndSwapInt32(&c.status, kConnStatus_Connected, kConnStatus_Disconnected) {
		c.conn.Close()
	}
}

// Close the connection, routine safe, send task in the queue will be sent before closing the connection
func (c *Connection) Close() {
	if atomic.LoadInt32(&c.status) != kConnStatus_Connected {
		return
	}

	select {
	case c.sendMsgQueue <- nil:
		{
			//	nothing
		}
	case <-time.After(time.Duration(c.sendTimeoutSec) * time.Second):
		{
			//	timeout, close the connection
			c.close()
		}
	}

	//	disable send
	atomic.StoreInt32(&c.disableSend, 1)
}

//	When don't need conection to send any thing, free it, DO NOT call it on multi routines
func (c *Connection) Free() {
	if nil != c.sendMsgQueue {
		close(c.sendMsgQueue)
		c.sendMsgQueue = nil
	}
}

func (c *Connection) syncExecuteEvent(evt *ConnEvent) bool {
	if nil == c.fnSyncExecute {
		return false
	}

	return c.fnSyncExecute(evt)
}

func (c *Connection) pushEvent(et int, d []byte) {
	//	this is for sync execute
	evt := newConnEvent(et, c.ConnId, d)
	if c.syncExecuteEvent(evt) {
		return
	}

	if nil == c.network.EventQueue {
		panic("Nil event queue")
		return
	}
	c.network.Push(evt)
}

// SetSyncExecuteFunc , you can set a callback that you can synchoronously process the event in every connection's event routine
// If the callback function return true, the event will not be dispatched
func (c *Connection) SetSyncExecuteFunc(fn FuncSyncExecute) FuncSyncExecute {
	prevFn := c.fnSyncExecute
	c.fnSyncExecute = fn
	return prevFn
}

// GetStatus get the connection's status
func (c *Connection) GetStatus() int32 {
	return c.status
}

func (c *Connection) setStatus(stat int) {
	c.status = int32(stat)
}

// GetUserdata get the userdata you set
func (c *Connection) GetUserdata() interface{} {
	return c.userdata
}

// SetUserdata set the userdata you need
func (c *Connection) SetUserdata(ud interface{}) {
	c.userdata = ud
}

// SetReadTimeoutSec set the read deadline for the connection
func (c *Connection) SetReadTimeoutSec(sec int) {
	c.readTimeoutSec = sec
}

//  GetReadTimeoutSec get the read deadline for the connection
func (c *Connection) GetReadTimeoutSec() int {
	return c.readTimeoutSec
}

// GetRemoteAddress return the remote address of the connection
func (c *Connection) GetRemoteAddress() string {
	return c.remoteAddr
}

// GetLocalAddress return the local address of the connection
func (c *Connection) GetLocalAddress() string {
	return c.localAddr
}

func (c *Connection) sendRaw(task *sendTask) error {
	if atomic.LoadInt32(&c.disableSend) != 0 {
		return kErrConnIsClosed
	}
	if atomic.LoadInt32(&c.status) != kConnStatus_Connected {
		return kErrConnIsClosed
	}

	select {
	case c.sendMsgQueue <- task:
		{
			//	nothing
		}
	case <-time.After(time.Duration(c.sendTimeoutSec) * time.Second):
		{
			//	timeout, close the connection
			logError("Send to peer %s timeout, close connection", c.GetRemoteAddress())
			c.close()
			return kErrConnSendTimeout
		}
	}

	return nil
}

// Send the buffer
func (c *Connection) Send(o interface{}) error {

	bytes, err := c.StreamProtocol.Serialize(o)
	if nil != err {
		return err
	}
	task := &sendTask{
		data: bytes,
		flag: 0,
	}

	// //	copy send buffer
	// if 0 != f&KConnFlag_CopySendBuffer {
	// 	msgCopy := make([]byte, len(bytes))
	// 	copy(msgCopy, bytes)
	// 	task.data = msgCopy
	// }

	return c.sendRaw(task)
}

// ApplyReadDealine
func (c *Connection) ApplyReadDeadline() {
	if 0 != c.readTimeoutSec {
		c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.readTimeoutSec) * time.Second))
	}
}

// ResetReadDeadline
func (c *Connection) ResetReadDeadline() {
	c.conn.SetReadDeadline(time.Time{})
}

//	run a routine to process the connection
func (c *Connection) run() {
	go c.routineMain()
}

func (c *Connection) routineMain() {
	defer func() {
		//	routine end
		e := recover()
		if e != nil {
			logFatal("Read routine panic %v, stack:", e)
			stackInfo := debug.Stack()
			logFatal(string(stackInfo))
		}

		//	close the connection
		logWarn("Read routine %s closed", c.GetRemoteAddress())
		c.close()

		//	free channel
		//	FIXED : consumers need free it, not producer

		//	post event
		c.pushEvent(KConnEvent_Disconnected, nil)
	}()

	if nil == c.StreamProtocol {
		panic("Nil stream protocol")
		return
	}

	//	connected
	c.pushEvent(KConnEvent_Connected, nil)
	atomic.StoreInt32(&c.status, kConnStatus_Connected)

	go c.routineSend()
	err := c.routineRead()
	if nil != err {
		logError("Read routine quit with error %v", err)
	}
}

func (c *Connection) routineSend() error {
	var err error

	defer func() {
		if nil != err {
			logError("Send routine quit with error %v", err)
		}
		e := recover()
		if nil != e {
			//	panic
			logFatal("Send routine panic %v, stack:", e)
			stackInfo := debug.Stack()
			logFatal(string(stackInfo))
		}
	}()

	for {
		select {
		case evt, ok := <-c.sendMsgQueue:
			{
				if !ok {
					//	channel closed, quit
					return nil
				}

				if nil == evt {
					c.close()
					return nil
				}

				if 0 == evt.flag&KConnFlag_NoHeader {
					var bytes []byte
					bytes, err = c.StreamProtocol.Serialize(evt.data)
					if nil != err {
						//	write header first
						if len(bytes) != 0 {
							_, err = c.conn.Write(bytes)
							if err != nil {
								return err
							}
						}
					} else {
						//	invalid packet
						panic("Failed to serialize header")
						break
					}
				}

				_, err = c.conn.Write(evt.data)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *Connection) routineRead() error {
	//	default buffer
	buf := make([]byte, c.maxReadBufferLength)
	var msg []byte
	var err error

	for {

		// if nil == c.unpacker {
		// 	msg, err = c.unpack(buf)
		// } else {
		// 	msg, err = c.unpacker.Unpack(c, buf)
		// }
		// if err != nil {
		// 	return err
		// }
		// if nil == msg {
		// 	// Empty pstream
		// 	continue
		// }

		// //	only push event when the connection is connected
		// if atomic.LoadInt32(&c.status) == kConnStatus_Connected {
		// 	//	try to unserialize to pb. do it in each routine to reduce the pressure of worker routine
		// 	pb := unserializePb(msg)
		// 	if nil != pb {
		// 		c.pushPbEvent(pb)
		// 	} else {
		// 		c.pushEvent(KConnEvent_Data, msg)
		// 	}
		// }
	}

	return nil
}

// func (c *Connection) unpack(buf []byte) ([]byte, error) {
// 	var err error
// 	//	read head
// 	c.ApplyReadDeadline()
// 	headerLength := int(c.streamProtocol.GetHeaderLength())
// 	if headerLength > len(buf) {
// 		return nil, fmt.Errorf("Header length %d > buffer length %d", headerLength, len(buf))
// 	}
// 	headBuf := buf[:headerLength]
// 	if _, err = io.ReadFull(c.conn, headBuf); nil != err {
// 		return nil, err
// 	}

// 	//	check length
// 	packetLength := c.streamProtocol.UnserializeHeader(headBuf)
// 	if packetLength > uint32(c.maxReadBufferLength) ||
// 		packetLength < c.streamProtocol.GetHeaderLength() {
// 		return nil, fmt.Errorf("Invalid stream length %d", packetLength)
// 	}
// 	if packetLength == c.streamProtocol.GetHeaderLength() {
// 		// Empty stream ?
// 		logFatal("Invalid stream length equal to header length")
// 		return nil, nil
// 	}

// 	//	read body
// 	c.ApplyReadDeadline()
// 	bodyLength := packetLength - c.streamProtocol.GetHeaderLength()
// 	if _, err = io.ReadFull(c.conn, buf[:bodyLength]); nil != err {
// 		return nil, err
// 	}

// 	//	ok
// 	msg := make([]byte, bodyLength)
// 	copy(msg, buf[:bodyLength])
// 	c.ResetReadDeadline()

// 	return msg, nil
// }
