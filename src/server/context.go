package server

import (
	"errors"
	"syncdb"
	"gamedef"
	"github.com/davyxu/golog"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/socket"
	// "github.com/davyxu/cellnet/timer"
	// "github.com/davyxu/golog"
	"github.com/astaxie/beego/orm"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type Context struct {
	Log *golog.Logger
	queue cellnet.EventQueue
	client *redis.Client
	o orm.Ormer
	peer cellnet.Peer
	request map[string] func (args interface{})
}

func NewContext() *Context {

	log := golog.New("server")

	client := syncdb.NewRedisClient()
	o := syncdb.NewOrm()
	
	queue := cellnet.NewEventQueue()

	peer := socket.NewAcceptor(queue).Start("127.0.0.1:8801")
	peer.SetName("client")

	context := &Context {
		Log : log,
		client : client,
		o : o,
		queue : queue,
		peer : peer,
		request : make(map[string]func (args interface{})),
	}
	return context
}

func (context *Context) Queue() cellnet.EventQueue {
	return context.queue
}

func (context *Context)  Startup() error {
	cellnet.RegisterMessage(context.peer, "gamedef.EchoReq", func(ev *cellnet.Event) {
		msg := ev.Msg.(*gamedef.EchoReq)
		Echo(context, msg)
	})
	return errors.New("h")
}

func (context *Context) Update() error {
	context.Log.Infoln("hello")
	return errors.New("h")	
}

func (context *Context) Cleanup() error {
	context.peer.Stop()
	return errors.New("h")
}