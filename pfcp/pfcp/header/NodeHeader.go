package pfcp

import (
    "log"
    "encoding/binary"
)

type NodeHeader struct {
    PFCPHeader
}

func (p *NodeHeader)Initialize() {
    p.Version = 1
    p.MPflag = 0
    p.Sflag = 0
    p.MessageLength = 0
    p.HeaderLength = 4
}

func (p *NodeHeader)InitializeFromBuf(signalBuffer []byte) uint32 {
    p.Version = signalBuffer[0] >> 5
    p.MPflag = (signalBuffer[0] >> 1) & 0x01
    p.Sflag = signalBuffer[0] & 0x01
    p.MessageLength = binary.BigEndian.Uint16(signalBuffer[2:4])
    p.SequenceNumber = binary.BigEndian.Uint32(signalBuffer[4:8]) >> 8
    return 8
}

func (p *NodeHeader) CreateSignal() ([]byte){
    msgLength := p.GetLength() + 4
    log.Println("信号生成 Length: ",msgLength)
    buffer := make([]byte, msgLength)
    buffer[0] = (p.Version << 5) | (p.MPflag << 1) | p.Sflag
    buffer[1] = p.MessageType
    binary.BigEndian.PutUint16(buffer[2:4],p.MessageLength)
    binary.BigEndian.PutUint32(buffer[4:8],p.SequenceNumber << 8)
    if p.IEInterface != nil {
        p.IEInterface.CreateSignal(buffer[8:])
    }
    return buffer
}
