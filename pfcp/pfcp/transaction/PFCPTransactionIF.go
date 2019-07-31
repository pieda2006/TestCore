package pfcp

type PFCPTransactionIF interface {
    SetSeqNum(seqNum uint32)
    SetDstIPaddress(dstIPaddress string)
    SetDstPort(dstPort string)
    SetMyIPaddress(myIPaddress string)
    SetMyPort(myPort string)
    SetEndPointChanel(endPintChannel chan []byte)
    GetTransactionChanel() (chan []byte)
    TransactionInitialize() (error)
    ExecuteTransaction()
    GetSeqNum() uint32
    SetSessionChan(sessionchan chan []byte)
}
