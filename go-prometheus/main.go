package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto" //nolint:depguard
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	msgCounter prometheus.Counter
	msgLatency prometheus.Histogram
}

func newPromMetrics() *PromMetrics {
	msgCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "actor_messages_total",
		Help: "count of actor messages received",
	})
	msgLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "actor_message_latency",
		Help:    "latency of actor messages received",
		Buckets: prometheus.LinearBuckets(0, 5, 10),
	})
	return &PromMetrics{
		msgCounter: msgCounter,
		msgLatency: msgLatency,
	}
}

func (p *PromMetrics) WithMetrics() func(actor.ReceiveFunc) actor.ReceiveFunc {
	return func(next actor.ReceiveFunc) actor.ReceiveFunc {
		return func(c *actor.Context) {
			fmt.Println("intercepting", c.Message())
			start := time.Now()
			p.msgCounter.Inc()
			next(c)
			ms := time.Since(start).Microseconds()
			p.msgLatency.Observe(float64(ms))
		}
	}
}

type Message struct {
	data string
}

type Foo struct {
}

func (f *Foo) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		fmt.Println("Foo started")
	case actor.Stopped:
		fmt.Println("Foo stopping")
	case Message:
		fmt.Println("Foo received message: ", msg.data)
	}
}

func newFoo() actor.Receiver {
	return &Foo{}
}

func main() {
	go func() {
		http.ListenAndServe(":2222", promhttp.Handler())
	}()

	e := actor.NewEngine()
	pM := newPromMetrics()
	pid := e.Spawn(newFoo, "foo", actor.WithMiddleware(pM.WithMetrics()))
	e.Send(pid, Message{data: "Hello"})

	i := 0
	for {
		e.Send(pid, Message{data: fmt.Sprintf("Hello %d", i)})
		time.Sleep(2 * time.Second)
		i++
	}
}
