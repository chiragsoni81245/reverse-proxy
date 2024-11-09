package packet

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPacketReader(t *testing.T) {
    data := []byte{'t', 'e', 's', 't'}
    buf := GetPacketBuffer(AuthPacket, data)
    packetCh := NewPacketReader(buf)
    pkt := <-packetCh
    assert.Equal(t, pkt.Type, AuthPacket)
    assert.Equal(t, pkt.ContentLength, int64(4))
    outputData, err := pkt.GetAllData()
    if err != nil {
        t.Error(err)
    }
    assert.Equal(t, outputData, data)
    pkt, isOpen  :=  <-packetCh
    assert.False(t, isOpen)
}

func TestNewPacketWriter(t *testing.T) {
    data := []byte{'t', 'e', 's', 't'}
    pktBuf := GetPacketBuffer(AuthPacket, data)
    actualPacketBytes := pktBuf.Bytes()
    pktReaderCh := NewPacketReader(pktBuf)

    buf := new(bytes.Buffer)
    NewPacketWriter(buf, pktReaderCh)
    assert.Equal(t, actualPacketBytes, buf.Bytes())
}
