package packet

import (
	"encoding/binary"
	"io"
	"log"
)

var AuthPacket byte = 0x01
var DataPacket byte = 0x02
var PacketCloseConnection byte = 0x03

type Packet struct {
    Type byte
    ContentLength int64
    Data io.Reader
    dataStreamCh chan struct{}
}

func (p *Packet) GetAllData() ([]byte, error) {
    data := make([]byte, p.ContentLength)
    n, err := p.Data.Read(data)
    if err != nil {
        return nil, err
    }
    p.CloseDataStream()
    return data[:n], nil
}

func (p *Packet) CloseDataStream() {
    close(p.dataStreamCh)
}

func (p *Packet) isDataStreamClosed() bool {
    _, ok := <-p.dataStreamCh
    return !ok
}

func NewPacketReader(r io.Reader) (chan Packet){
    pktCh := make(chan Packet)

    go func () {
        var err error
        defer func(){
            if err != nil{
                if err != io.EOF {
                    log.Fatal(err)
                }
                close(pktCh)
            }
        }()

        for {
            packet := Packet{
                dataStreamCh: make(chan struct{}),
            }
            buf := make([]byte, 1)
            _, err = r.Read(buf)
            if err != nil {
                break
            }
            packet.Type = buf[0]

            err = binary.Read(r, binary.LittleEndian, &packet.ContentLength)
            if err != nil {
                break
            }

            packet.Data = io.LimitReader(r, packet.ContentLength)

            pktCh <- packet

            // Wait for data stream to be closed via reader
            packet.isDataStreamClosed()
        }
    }()

    return pktCh
}

func NewPacketWriter(w io.Writer, pktCh chan Packet) {
    for {
        pkt, ok := <-pktCh
        if !ok { 
            break 
        } 
        _, err := w.Write([]byte{pkt.Type})
        if err != nil {
            break
        }
        err = binary.Write(w, binary.LittleEndian, pkt.ContentLength)
        if err != nil {
            break
        }
        _, err = io.Copy(w, pkt.Data)
        if err != nil {
            break
        }
        pkt.CloseDataStream()
    }
}


