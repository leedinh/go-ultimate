package main

import (
	"log"
	"net"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type server struct {
	listenAddr string
	ln         net.Listener
	sessions   map[*actor.PID]net.Conn
}

type session struct {
	cnn net.Conn
}

type ConnAdd struct {
	pid *actor.PID
	cnn net.Conn
}

type ConnRemove struct {
	pid *actor.PID
}

func newSession(cnn net.Conn) actor.Producer {
	return func() actor.Receiver {
		return &session{cnn: cnn}
	}
}

func (s *session) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
	case actor.Started:
		log.Println("session started", s.cnn.RemoteAddr())
		go s.readLoop(c)
	case actor.Stopped:
		_ = msg
	case []byte:
		log.Println("session received", string(msg))
		s.cnn.Write(msg)
	}
}

func (s *session) readLoop(c *actor.Context) {
	buf := make([]byte, 1024)
	for {
		n, err := s.cnn.Read(buf)
		if err != nil {
			log.Println("read error", err)
			break
		}
		msg := buf[:n]
		c.Send(c.PID(), msg)
	}

	// Read loop ended, close connection
	c.Send(c.Parent(), &ConnRemove{pid: c.PID()})
}

func newServer(listenAddr string) actor.Producer {
	return func() actor.Receiver {
		return &server{
			listenAddr: listenAddr,
			sessions:   make(map[*actor.PID]net.Conn)}
	}

}

func (s *server) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		ln, err := net.Listen("tcp", s.listenAddr)
		if err != nil {
			panic(err)
		}
		s.ln = ln
	case actor.Started:
		log.Println("server started", s.listenAddr)
		go s.acceptLoop(c)
	case actor.Stopped:
		_ = msg
	case *ConnAdd:
		log.Println("conn added", msg.pid, msg.cnn.RemoteAddr())
		s.sessions[msg.pid] = msg.cnn
	case *ConnRemove:
		log.Println("conn removed", msg.pid, s.sessions[msg.pid].RemoteAddr())
		delete(s.sessions, msg.pid)
	}
}

func (s *server) acceptLoop(c *actor.Context) {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("accept error", err)
			break
		}
		pid := c.SpawnChild(newSession(conn), "session", actor.WithTags(conn.RemoteAddr().String()))
		c.Send(c.PID(), &ConnAdd{pid: pid, cnn: conn})
	}
}

func main() {
	e := actor.NewEngine()
	e.Spawn(newServer(":6000"), "server")

	time.Sleep(100 * time.Second)
	<-make(chan struct{})
}
