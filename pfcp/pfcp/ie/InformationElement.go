package pfcp

import (
    "log"
)

type InformationElement struct {
    IEType uint16
    MemberElement InformationElementIF
    NextElement InformationElementIF
    IELength uint16
    MessageLength uint16
}


func (p *InformationElement) GetLength() (uint16) {
    if p.MessageLength == 0 {
        p.MessageLength = p.IELength + 4 //4は固定分
        if p.MemberElement != nil {
            p.MessageLength = p.MessageLength + p.MemberElement.GetLength()
        }
        if p.NextElement != nil {
            p.MessageLength = p.MessageLength + p.NextElement.GetLength()
        }
    }
    log.Println("IE: ",p.IEType," Length: ",p.MessageLength)
    return p.MessageLength
}
func (p *InformationElement) SetMemberElement(memberElement InformationElementIF) {
    p.MemberElement = memberElement
}
func (p *InformationElement) SetNextElement(nextElement InformationElementIF) {
    p.NextElement = nextElement
}
func (p *InformationElement) CreateNextIE(signal []byte) {
    if p.NextElement != nil {
        p.NextElement.CreateSignal(signal)
    }
}

