package worker_pool

import (
	"context"
	"fmt"
	"github.com/hong195/aggregator-sevice/internal/usecase"
	"github.com/hong195/aggregator-sevice/internal/usecase/command"
	"github.com/hong195/aggregator-sevice/pkg/generator"
	"github.com/hong195/aggregator-sevice/pkg/logger"
	"sync"
	"time"
)

const defaultStopTimeout = 3 * time.Second

type Pool struct {
	workerCount int
	in          <-chan generator.RawPacket
	wg          sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc

	errCh chan error
	u     *usecase.UseCases
	l     logger.Interface
}

func NewPool(n int, in <-chan generator.RawPacket, u *usecase.UseCases, l logger.Interface) *Pool {
	if n < 1 {
		n = 1
	}
	return &Pool{
		workerCount: n,
		in:          in,
		errCh:       make(chan error, 1),
		u:           u,
		l:           l,
	}
}

func (p *Pool) Start() {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.wg.Add(p.workerCount)

	for i := 1; i <= p.workerCount; i++ {
		go p.worker(i)
	}
}

func (p *Pool) Shutdown() error {

	if p.cancel != nil {
		p.cancel()
	}

	p.wg.Wait()
	return nil
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()
	select {
	case <-p.ctx.Done():
		// быстро дренируем буфер и выходим
		for {
			select {
			case _, ok := <-p.in:
				if !ok {
					return
				}
			default:
				fmt.Println("done")
				return
			}
		}
	case pkt, ok := <-p.in:
		if !ok {
			return
		}

		fmt.Println(pkt)

		packetToStore := command.NewStoreDataPacket(pkt.ID.String(), pkt.Timestamp.UnixMilli(), pkt.Payload)

		err := p.u.Commands.StoreDataPacket.Handle(p.ctx, packetToStore)
		if err != nil {
			p.l.Error(err, "worker pool - run - StoreDataPacket")
			return
		}
	}
}
