package generator

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// RawPacket — "сырой" пакет от внешнего источника.
type RawPacket struct {
	ID        uuid.UUID
	Timestamp time.Time
	Payload   []int
}

// Generator — эмулятор внешнего источника.
type Generator struct {
	interval time.Duration
	k        int
	out      chan<- RawPacket

	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

func NewGenerator(interval time.Duration, k int, out chan<- RawPacket) *Generator {
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	if k <= 0 {
		k = 1
	}
	return &Generator{
		interval: interval,
		k:        k,
		out:      out,
		done:     make(chan struct{}),
	}
}

func (g *Generator) Start() {
	g.ctx, g.cancel = context.WithCancel(context.Background())
	go g.run()
}

func (g *Generator) Stop() error {
	if g.cancel != nil {
		g.cancel()
	}
	<-g.done
	return nil
}

func (g *Generator) run() {
	defer close(g.done)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ticker := time.NewTicker(g.interval)
	defer ticker.Stop()

	for {
		select {
		case <-g.ctx.Done():
			close(g.out)
			return
		case t := <-ticker.C:
			p := RawPacket{
				ID:        uuid.New(),
				Timestamp: t.UTC(),
				Payload:   make([]int, g.k),
			}

			for i := range p.Payload {
				p.Payload[i] = r.Intn(1000) // 0..999
			}

			select {
			case g.out <- p:
				fmt.Println(p)
			case <-g.ctx.Done():
				close(g.out)
				return
			}
		}
	}
}
