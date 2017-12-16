package main

import (
	"log"
	"sync"

	"bufio"
	"os"
	"syscall"

	"os/signal"

	"sync/atomic"

	"tcpnetwork"
)

var (
	kServerAddress  string = "localhost:14444"
	serverConnected int32
	stopFlag        int32
)

type EchoStreamProtocolHandler struct {
}

func NewEchoStreamProtocolHander() {
	return &EchoStreamProtocolHandler{}
}

func (e *EchoStreamProtocolHandler) Serialize(body interface{}) ([]byte, error) {
}

func (e *EchoStreamProtocolHandler) Unserialize(bin []byte) (interface{}, error) {
}

// echo server routine
func echoServer() (*tcpnetwork.TCPNetwork, error) {
	var err error
	server := tcpnetwork.NewTCPNetworkServer(1024, tcpnetwork.KStreamProtocol2HeaderLength)
	err = server.Listen(kServerAddress)
	if nil != err {
		return nil, err
	}

	return server, nil
}

func routineEchoServer(server *tcpnetwork.TCPNetwork, wg *sync.WaitGroup, stopCh chan struct{}) {
	defer func() {
		log.Println("server done")
		wg.Done()
	}()

	for {
		select {
		case evt, ok := <-server.GetEventQueue():
			{
				if !ok {
					return
				}

				switch evt.EventType {
				case tcpnetwork.KConnEvent_Connected:
					{
						log.Println("Client ", evt.Conn.GetRemoteAddress(), " connected")
					}
				case tcpnetwork.KConnEvent_Close:
					{
						log.Println("Client ", evt.Conn.GetRemoteAddress(), " disconnected")
					}
				case tcpnetwork.KConnEvent_Data:
					{
						log.Println("data" + evt.Data)
						// evt.Conn.Send(evt.Data, 0)
					}
				}
			}
		case <-stopCh:
			{
				return
			}
		}
	}
}

func main() {

	// create server
	server, err := echoServer()
	if nil != err {
		log.Println(err)
		return
	}

	stopCh := make(chan struct{})

	// process event
	var wg sync.WaitGroup
	wg.Add(1)
	go routineEchoServer(server, &wg, stopCh)

	// input event
	wg.Add(1)
	go routineInput(&wg, clientConn)

	// wait
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

MAINLOOP:
	for {
		select {
		case <-sc:
			{
				//	app cancelled by user , do clean up work
				log.Println("Terminating ...")
				break MAINLOOP
			}
		}
	}

	atomic.StoreInt32(&stopFlag, 1)
	log.Println("Press enter to exit")
	close(stopCh)
	wg.Wait()
	// clientConn.Close()
	server.Shutdown()
}
