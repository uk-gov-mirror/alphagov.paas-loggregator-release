package reliability

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type WorkerClient struct {
	addr   string
	runner Runner
}

func NewWorkerClient(addr string, r Runner) *WorkerClient {
	return &WorkerClient{
		addr:   addr,
		runner: r,
	}
}

func (w *WorkerClient) Run(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.Dial(w.addr, nil)
	if err != nil {
		return err
	}
	log.Println("connected to control server")

	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		defer cancel()

		for {
			var test Test
			err := conn.ReadJSON(&test)
			if err != nil {
				break
			}

			log.Println("test received from control server")
			go w.runner.Run(&test)
		}
	}()

	<-ctx.Done()

	return conn.Close()
}
