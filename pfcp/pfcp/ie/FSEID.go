package pfcp

import (
    "encoding/binary"
    "net"
)

type FSEID struct {
    InformationElement
    IPaddressType uint8
    SEID uint64
    IPv4Address string
    IPv6Address string
}

func (p *FSEID) CreateSignal(signal []byte) {
    binary.BigEndian.PutUint16(signal[0:2],p.IEType)
    binary.BigEndian.PutUint16(signal[2:4],p.IELength)
    //課題：IPv4以外の処理も追加 START
    signal[4] = p.IPaddressType
    binary.BigEndian.PutUint64(signal[5:13],p.SEID)
    copy(signal[13:17], net.ParseIP(p.IPv4Address)[12:16])
    //END
    p.CreateNextIE(signal[17:])
}

func (p *FSEID) IEInitialize() (error){
    p.IEType = FSEID_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = 13
    p.MessageLength = 0
    p.IPaddressType = 0
    p.SEID = 0
    p.IPv4Address = ""
    p.IPv6Address = ""
    return nil
}

func (p *FSEID) GetSEID() uint64 {
    return p.SEID
}

func (p *FSEID) SetSEID(seid uint64) {
    p.SEID = seid
}

func (p *FSEID) SetIPaddressType(iptype uint8) {
    p.IPaddressType = iptype
}

func (p *FSEID) SetIPv4Address(ipaddress string) {
    p.IPv4Address = ipaddress
}

func (p *FSEID) IEInitializeFromBuff(buffer []byte) (uint32){
    p.IEType = FSEID_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IPaddressType = uint8(buffer[4] & 0x3)
    p.SEID = binary.BigEndian.Uint64(buffer[5:13])
    var nextbyte uint32
    switch p.IPaddressType {
    case FSEID_IPV6 :
        //課題：IPv6の場合の処理実装 START
        //END
    case FSEID_IPV4 :
        p.IPv4Address = net.IP(buffer[13:17]).String()
        nextbyte = 17
    case FSEID_DUAL :
        //課題：IPdualの場合の処理実装 STRT
        //END
    }

    return nextbyte
}
