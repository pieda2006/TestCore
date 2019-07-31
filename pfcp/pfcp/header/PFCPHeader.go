package pfcp

import (
)

type PFCPHeader struct {
    Version uint8
    MPflag uint8
    Sflag uint8
    MessageType uint8
    MessageLength uint16
    SequenceNumber uint32
    IEInterface InformationElementIF
    HeaderLength uint16
}

func (p *PFCPHeader) GetLength() (uint16) {
    if p.MessageLength == 0 {
        if p.IEInterface != nil {
            p.MessageLength = p.HeaderLength + p.IEInterface.GetLength()
        } else {
            p.MessageLength = p.HeaderLength
        }
    }
    return p.MessageLength
}

func (p *PFCPHeader) SetSequenceNum(seqNum uint32) {
    p.SequenceNumber = seqNum
}

func (p *PFCPHeader) SetInformationElement(infoElement InformationElementIF) {
    if p.IEInterface != nil {
        infoElement.SetNextElement(p.IEInterface)
    }
    p.IEInterface = infoElement
}

func (p *PFCPHeader) SetMsgType (msgType uint8) {
    p.MessageType = msgType
}

func  (p *PFCPHeader) GetSequenceNum () uint32 {
    return p.SequenceNumber
}
