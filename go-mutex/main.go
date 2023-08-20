package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	mu     sync.RWMutex
	health int
}

func NewPlayer() *Player {
	return &Player{
		health: 100,
	}
}

func (p *Player) GetHealth() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.health
}

func (p *Player) GetHit(dmg int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.health -= dmg
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
