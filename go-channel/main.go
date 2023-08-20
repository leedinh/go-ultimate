package main

import "fmt"

type Server struct {
	users  map[string]string
	userch chan string
}

func NewServer() *Server {
	return &Server{
		users:  make(map[string]string),
		userch: make(chan string, 10000),
	}
}

func (s *Server) AddUser(user string) {
	s.users[user] = user
	fmt.Println("Add user: ", user)
}

func (s *Server) Run() {
	go s.loop()
}

func (s *Server) loop() {
	for {
		user := <-s.userch
		s.AddUser(user)
	}
}

func main() {
	userch := make(chan string)
	userch <- "Dinh"
	user := <-userch
	fmt.Println(user)
}

func sendMessage(msgch chan<- string) {
	msgch <- "Hello"
}

func readMessage(msgch <-chan string) {
	msg := <-msgch
	fmt.Println(msg)
}
