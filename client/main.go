package main

import (
	"bufio"
	"gamedef"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
	"github.com/davyxu/golog"
	"os"
	"strings"
)

var log = golog.New("main")

func ReadConsole(callback func(string)) {

	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		text = strings.TrimRight(text, "\n\r ")

		text = strings.TrimLeft(text, " ")

		callback(text)
	}
}

func main() {
	queue := cellnet.NewEventQueue()

	peer := socket.NewConnector(queue).Start("127.0.0.1:8801")
	peer.SetName("client")

	session := 1

	cellnet.RegisterMessage(peer, "gamedef.EchoAck", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.EchoAck)
		log.Infof("sid%d say: %s", msg.Errorcode, msg.Content)
	})

	queue.StartLoop()

	ReadConsole(func(str string) {
		session = session + 1
		peer.(socket.Connector).DefaultSession().Send(&gamedef.EchoReq{
			Session: session,
			Content: str,
		})
	})
}
