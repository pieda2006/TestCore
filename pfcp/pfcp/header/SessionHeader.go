package pfcp

import (
    "log"
    "encoding/binary"
)

type SessionHeader struct {
    PFCPHeader
    SEID uint64
    MessagePriority uint8
}

func (p *SessionHeader)Initialize() {
    log.Println("初期化処理実施")
    p.Version = 1
    p.MPflag = 1
    p.Sflag = 1
    p.SEID = 0
    p.MessagePriority = 0
    p.SequenceNumber = 0
    p.MessageLength = 0
    p.HeaderLength = 12
}

func (p *SessionHeader)InitializeFromBuf(signalBuffer []byte) uint32 {
    log.Println("受信信号から信号インスタンス生成")
    p.Version = signalBuffer[0] >> 5
    p.MPflag = (signalBuffer[0] >> 1) & 0x01
    p.Sflag = signalBuffer[0] & 0x01
    p.MessageType = signalBuffer[1]
    p.MessageLength = binary.BigEndian.Uint16(signalBuffer[2:4])
    p.SEID = binary.BigEndian.Uint64(signalBuffer[4:12])
    p.SequenceNumber = binary.BigEndian.Uint32(signalBuffer[12:16]) >> 8
    p.MessagePriority = signalBuffer[15] >> 4
    return 16
}

func (p *SessionHeader) CreateSignal() ([]byte){
    msgLength := p.GetLength() + 4
    log.Println("信号生成 Length: ",msgLength)
    buffer := make([]byte, msgLength)
    buffer[0] = (p.Version << 5) | (p.MPflag << 1) | p.Sflag
    buffer[1] = p.MessageType
    binary.BigEndian.PutUint16(buffer[2:4],p.MessageLength)
    binary.BigEndian.PutUint64(buffer[4:12],p.SEID)
    binary.BigEndian.PutUint32(buffer[12:16],p.SequenceNumber << 8)
    buffer[15] = p.MessagePriority << 4
    if p.IEInterface != nil {
        p.IEInterface.CreateSignal(buffer[16:])
    }
    return buffer
}

func (p *SessionHeader) SetSEID(seid uint64){
    p.SEID = seid
}

func (p *SessionHeader) GetSEID() uint64{
    return p.SEID
}

func (p *SessionHeader) GetMessagePriority() uint8{
    return p.MessagePriority
}

func (p *SessionHeader) SetMessagePriority(priority uint8){
    p.MessagePriority = priority
}
