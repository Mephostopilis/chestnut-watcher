package main

import (
	// "log"
	// "fmt"
	"time"
	"server"
	// "syncdb"
	// "model"
	// "chat/gamedef"
	// "github.com/davyxu/cellnet"
	// "github.com/davyxu/cellnet/socket"
	"github.com/davyxu/cellnet/timer"
	// "github.com/davyxu/golog"
)

func main() {

	context := server.NewContext()
	context.Startup()
	
	queue := context.Queue()
	queue.StartLoop()

	// var count int
	timer.NewLoop(queue, time.Millisecond*100, func(ctx *timer.Loop) {
		log := context.Log
		log.Debugln("tick 100 ms")
		context.Update()
		// count++

		// if count >= 10 {
		// 	// signal.Done(1)
		// 	ctx.Stop()
		// }
	}, nil).Start()

	queue.Wait()

	context.Cleanup()

	// timer.After(queue, 500*time.Millisecond, func() {
	// 	log.Debugln("after 100 ms")

	// 	// signal.Done(1)
	// })

	// peer.Stop()
}
