package pfcp

import (
)

type PFCPNode struct {
    IPaddress string
    Port string
    MyChannel chan []byte
    NodeType uint8
    NodeChan chan []byte
}

func (p *PFCPNode) SetNodeIPaddress(nodeip string) {
    p.IPaddress = nodeip
}
func (p *PFCPNode) SetNodePort(nodeport string) {
    p.Port = nodeport
}


func (p *PFCPNode) GetEndPointChan() (chan []byte){
    return p.MyChannel
}

func (p *PFCPNode) GetNodeIPaddress() string {
    return p.IPaddress
}

func (p *PFCPNode) GetNodePort() string {
    return p.Port
}

func (p *PFCPNode) SetNodeChan(nodechan chan []byte) {
    p.NodeChan = nodechan
}
