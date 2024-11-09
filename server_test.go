package main

import (
	"chiragsoni81245/reverse-proxy/packet"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func startPingPongBackendServer(host string, port int) (error){
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
    if err != nil {
        return err
    }

    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                break
            }

            pktRCh := packet.NewPacketReader(conn)
            packet.NewPacketWriter(conn, pktRCh)
        }
    }()

    return nil
}

func TestServer(t *testing.T) {
    err := startPingPongBackendServer("127.0.0.1", 9000)
    if err != nil {
        t.Error(err)
    }
    proxy := NewServer("127.0.0.1", 8000)
    go proxy.startAcceptConnLoop()
    time.Sleep(1*time.Second)

    proxyClient, err := net.Dial("tcp", "127.0.0.1:8000")
    if err != nil {
        t.Error(err)
        return
    }
    time.Sleep(100*time.Millisecond)

    buf := packet.GetPacketBuffer(packet.AuthPacket, []byte{})    
    proxyClient.Write(buf.Bytes())

    data := []byte{'c','h','i','r','a','g'}
    buf = packet.GetPacketBuffer(packet.DataPacket, data)    
    proxyClient.Write(buf.Bytes())

    pktRCh := packet.NewPacketReader(proxyClient)
    pkt := <-pktRCh 
    assert.Equal(t, pkt.Type, packet.DataPacket)
    assert.Equal(t, pkt.ContentLength, int64(len(data)))
    outputData, err := pkt.GetAllData()
    if err != nil {
        t.Error(err)
    }
    assert.Equal(t, outputData, data)
}
