package worker_pool

import (
	"context"
	"fmt"
	"github.com/hong195/aggregator-sevice/pkg/generator"
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
}

func NewPool(n int, in <-chan generator.RawPacket) *Pool {
	if n < 1 {
		n = 1
	}
	return &Pool{
		workerCount: n,
		in:          in,
		errCh:       make(chan error, 1),
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
			case pkt, ok := <-p.in:
				if !ok {
					return
				}
				fmt.Println(pkt)
			default:
				fmt.Println("done")
				return
			}
		}
	case pkt, ok := <-p.in:
		if !ok {
			fmt.Println("done 2")
			return
		}
		fmt.Println(pkt)
	}
}

func maxVal(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, v := range xs[1:] {
		if v > m {
			m = v
		}
	}
	return m
}
