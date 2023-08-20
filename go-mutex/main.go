package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type Player struct {
	health int32
}

func NewPlayer() *Player {
	return &Player{
		health: 100,
	}
}

func (p *Player) GetHealth() int {
	return int(atomic.LoadInt32(&p.health))
}

func (p *Player) GetHit(dmg int) {
	atomic.AddInt32(&p.health, int32(-dmg))
}

func StartUILoop(p *Player) {
	ticker := time.NewTicker(time.Second)
	for {
		fmt.Println("Player health: ", p.GetHealth())
		<-ticker.C
	}
}

func StartGameLoop(p *Player) {
	ticker := time.NewTicker(time.Millisecond * 300)
	for {
		p.GetHit(rand.Intn(40))
		if p.GetHealth() <= 0 {
			fmt.Println("Game over")
			break
		}
		<-ticker.C
	}
}

func main() {
	player := NewPlayer()
	go StartUILoop(player)
	StartGameLoop(player)
}
