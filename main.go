package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
	msgCh      chan Message
}

type Message struct {
	from    string
	payload []byte
}

func NewServer(listnAddrs string) *Server {
	return &Server{listenAddr: listnAddrs, quitCh: make(chan struct{}), msgCh: make(chan Message, 10)}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()
	<-s.quitCh
	close(s.msgCh)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		con, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			continue
		}
		fmt.Println("New Connection Established with server:", con.RemoteAddr())
		go s.readLoop(con)

	}
}

func (s *Server) readLoop(con net.Conn) {
	defer con.Close()
	buf := make([]byte, 2048)
	for {
		n, err := con.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}
		s.msgCh <- Message{payload: buf[:n], from: con.RemoteAddr().String()}
	}
}

func main() {

	server := NewServer(":3000")
	go func() {
		for msg := range server.msgCh {
			fmt.Printf("Received Message from Connection (%s): %s\n", msg.from, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())

}
