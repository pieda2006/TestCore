package pfcp

type PFCPNodeIF interface {
    SetNodeIPaddress(nodeip string)
    SetNodePort(nodeport string)
    ExecuteEndPoint() (error)
    EndPointInitialize(endPointType uint8) (error)
    GetEndPointChan() (chan []byte)
    GetNodeIPaddress() string
    GetNodePort() string
    SetNodeChan(chan []byte)
    SetTransactionChan(seqNum uint32, transactionchan chan []byte)
}
