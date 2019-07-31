package pfcp

import (
    "encoding/binary"
    "net"
)

type NodeID struct {
    InformationElement
    NodeIDType uint8
    NodeIDValueIPv4 string
    NodeIDValueIPv6 string
    NodeIDValueFQDN string
}

func (p *NodeID) CreateSignal(signal []byte) {
    binary.BigEndian.PutUint16(signal[0:2],p.IEType)
    binary.BigEndian.PutUint16(signal[2:4],p.IELength)
    //課題：IPv4以外の処理も追加 START
    signal[4] = p.NodeIDType
    copy(signal[5:9],net.ParseIP(p.NodeIDValueIPv4)[12:16])
    //END
    p.CreateNextIE(signal[9:])
}

func (p *NodeID) IEInitialize() (error){
    p.IEType = NOADID_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = 5
    p.MessageLength = 0
    p.NodeIDType = NODETYPE_IPV4
    p.NodeIDValueIPv4 = ""
    p.NodeIDValueIPv6 = ""
    p.NodeIDValueFQDN = ""
    return nil
}

func (p *NodeID) IEInitializeFromBuff(buffer []byte) (uint32){
    p.IEType = NOADID_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = binary.BigEndian.Uint16(buffer[2:4])
    p.NodeIDType = buffer[4] & 0xF
    var nextbyte uint32 = 9
    switch p.NodeIDType {
        case NODETYPE_IPV4 :
            p.NodeIDValueIPv4 = net.IP(buffer[5:9]).String()
        case NODETYPE_IPV6 :
            //課題：IPv6の場合のパース処理追加 START
            //END
        case NODETYPE_FQDN :
            //課題：FQDNの場合のパース処理追加 START
            //END
        default :
    }
    return nextbyte
}

func (p *NodeID) SetNodeIDType(nodetype uint8) {
    p.NodeIDType = nodetype
}

func (p *NodeID) SetNodeIDValueIPv4(ipaddress string) {
    p.NodeIDValueIPv4 = ipaddress
}

func (p *NodeID) SetNodeIDValueIPv6(ipaddress string) {
    p.NodeIDValueIPv6 = ipaddress
}

func (p *NodeID) SetNodeIDValueFQDN(fqdn string) {
    p.NodeIDValueFQDN = fqdn
}
