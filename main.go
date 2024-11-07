package main

import (
	"chiragsoni81245/reverse-proxy/packet"
	"errors"
	"fmt"
	"net"
)

type Server struct {
    addr string
    listener net.Listener
}

func (s *Server) startAcceptConnLoop() error {
    listener, err := net.Listen("tcp", s.addr)
    if err != nil {
        return err
    }
    s.listener = listener

    for {
        conn, err := s.listener.Accept()
        if errors.Is(err, net.ErrClosed) {
            return nil
        }
        if err != nil {
            fmt.Printf("TCP accept error: %s\n", err)
        }

        go s.handleConn(conn)
    }
}

func (s *Server) authenticate(authPacket packet.Packet) error {
    return nil
}

func (s *Server) handleConn(conn net.Conn) {
    pktRCh := packet.NewPacketReader(conn)
    authPacket, ok := <-pktRCh

    if !ok && authPacket.Type != packet.AuthPacket {
        conn.Close()
        close(pktRCh)
        return
    }

    err := s.authenticate(authPacket)
    if err != nil {
        conn.Close()
        close(pktRCh)
        return
    }

    // Connect with backend server
    backendConn, err := net.Dial("127.0.0.1", "9000")
    if err != nil {
        conn.Close()
        close(pktRCh)
        return
    }

    s.handleCommunication(conn, backendConn)
}

func (s *Server) handleCommunication(conn net.Conn, backendConn net.Conn) {
    localPacketRCh := packet.NewPacketReader(conn)
    remotePacketWCh := packet.NewPacketWriter(backendConn)     
    localPacketWCh := packet.NewPacketWriter(conn)     
    remotePacketRCh := packet.NewPacketReader(backendConn)

    for {
        select {
        case pkt := <-localPacketRCh:
            remotePacketWCh <- pkt 
        case pkt := <-remotePacketRCh:
            localPacketWCh <- pkt 
        }
    }
}

func NewServer(host string, port int) (*Server){
    return &Server{
        addr: fmt.Sprintf("%s:%d", host, port),
    }
}

func main() {
    s := NewServer("127.0.0.1", 8080)
    s.startAcceptConnLoop()
}
