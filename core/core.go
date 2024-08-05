package core

import (
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Printf("new connection to %s\n", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	buf := make([]byte, 2048)
	msg := fmt.Sprintf("%s: ", conn.RemoteAddr())
	conn.Write([]byte(msg))
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Close connection with %s\n", conn.RemoteAddr())
			<-s.quitCh
		}

		fmt.Printf("%s: %s", conn.RemoteAddr(), string(buf[:n]))
		receivedMsg := fmt.Sprintf("Received --> %s", string(buf[:n]))
		conn.Write([]byte(receivedMsg))
		conn.Write([]byte(msg))
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln
	s.acceptLoop()

	return nil
}

// func main() {
// 	fmt.Println("Server started")
// 	server := NewServer(":8000")
// 	log.Fatal(server.Start())
// }
