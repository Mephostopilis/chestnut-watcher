package main

import (
	"github.com/davyxu/cellnet/timer"
	"server"
	"time"
)

func main() {

	context := server.NewContext()
	context.Startup()

	queue := context.Queue()
	queue.StartLoop()

	// var count int
	timer.NewLoop(queue, time.Millisecond*100, func(ctx *timer.Loop) {

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
