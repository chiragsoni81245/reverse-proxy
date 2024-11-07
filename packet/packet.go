package packet

import (
	"encoding/binary"
	"io"
	"log"
)

var AuthPacket byte = 0x01
var PacketCloseConnection byte = 0x02

type Packet struct {
    Type byte
    ContentLength int64
    Data []byte
}

func NewPacketReader(r io.Reader) (chan Packet){
    pktCh := make(chan Packet)

    go func () {
        for {
            packet := Packet{}
            buf := make([]byte, 1)
            _, err := r.Read(buf)
            if err != nil {
                log.Fatal(err)
                close(pktCh)
                break
            }
            packet.Type = buf[0]

            binary.Read(r, binary.LittleEndian, &packet.ContentLength)

            buf = make([]byte, packet.ContentLength)
            var n int
            n, err = r.Read(buf)
            if err != nil {
                log.Fatal(err)
                close(pktCh)
                break
            }
            packet.Data = buf[:n]

            pktCh <- packet
        }
    }()

    return pktCh
}

func NewPacketWriter(w io.Writer) (chan Packet){
    pktCh := make(chan Packet)

    go func () {
        for {
            pkt, ok := <-pktCh
            if !ok { break } 
            _, err := w.Write([]byte{pkt.Type})
            if err != nil {
                close(pktCh)
            }
            err = binary.Write(w, binary.LittleEndian, pkt.ContentLength)
            if err != nil {
                close(pktCh)
            }
            _, err = w.Write(pkt.Data)
            if err != nil {
                close(pktCh)
            }
        }
    }()

    return pktCh
}


