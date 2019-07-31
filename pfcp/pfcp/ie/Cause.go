package pfcp

import (
    "encoding/binary"
)

type Cause struct {
    InformationElement
    Cause uint8
}

func (p *Cause) CreateSignal(signal []byte) {
    binary.BigEndian.PutUint16(signal[0:2],p.IEType)
    binary.BigEndian.PutUint16(signal[2:4],p.IELength)
    signal[4] = p.Cause
    p.CreateNextIE(signal[5:])
}

func (p *Cause) SetCause(cause uint8) {
    p.Cause = cause
}

func (p *Cause) IEInitialize() (error){
    p.IEType = CAUSE_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = 1
    p.MessageLength = 0
    return nil
}

func (p *Cause) IEInitializeFromBuff(buffer []byte) (uint32){
    p.IEType = CAUSE_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = binary.BigEndian.Uint16(buffer[2:4])
    p.Cause = buffer[4]
    p.MessageLength = 0
    return 5
}

func (p *Cause) GetCause() uint8 {
    return p.Cause
}
