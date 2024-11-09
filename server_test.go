package main

import (
	"chiragsoni81245/reverse-proxy/packet"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func startPingPongBackendServer(addr string) (error){
    listener, err := net.Listen("tcp", addr) 
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

func getPingPongBackend(clientAddr string) (string, error){
    backendAddr := "127.0.0.1:9000"
    err := startPingPongBackendServer(backendAddr)
    if err != nil {
        return "", err
    }
    return backendAddr, nil
}

func TestServer(t *testing.T) {
    proxy := NewProxyServer("127.0.0.1", 8000, getPingPongBackend)
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
    pktRCh := packet.NewPacketReader(proxyClient)

    for i:=0; i<3; i++ {
        proxyClient.Write(buf.Bytes())

        pkt := <-pktRCh 
        assert.Equal(t, pkt.Type, packet.DataPacket)
        assert.Equal(t, pkt.ContentLength, int64(len(data)))
        outputData, err := pkt.GetAllData()
        if err != nil {
            t.Error(err)
        }
        assert.Equal(t, outputData, data)
    }
}
