package main

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
)

type Handler struct {
	state uint
}
type SetState struct {
	value uint
}

type ResetState struct{}

func newHandler() actor.Receiver {
	return &Handler{}
}

func (h *Handler) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		h.state = 10
		fmt.Println("handler Initialized")
	case actor.Started:
		fmt.Println("handler started")
	case SetState:
		h.state = msg.value
		fmt.Println("handler receive new state", h.state)
	case ResetState:
		h.state = 0
		fmt.Println("handler receive reset state", h.state)
	case actor.Stopped:
		fmt.Println("handler stopped")
		_ = msg
	}
}

func main() {
	e := actor.NewEngine()
	pid := e.Spawn(newHandler, "handler")
	fmt.Println("pid ->", pid)
	for x := 0; x < 10; x++ {
		e.Send(pid, SetState{value: uint(x)})
	}
	e.Send(pid, ResetState{})
}
