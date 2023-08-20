package main

import "testing"

func TestGame(t *testing.T) {
	player := NewPlayer()
	go StartUILoop(player)
	StartGameLoop(player)
}
