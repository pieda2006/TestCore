package pfcp

import (
  "encoding/binary"
  "encoding/json"
)

//Client
const (
    INITIALIZE_STATUS = iota
    SENDING_REQUEST = iota
    WAITE_ANSWER = iota
    RECEAVE_ANSWER = iota
    SENDING_ANSWER = iota
    WAITE_ANS_REQ = iota
    RECEAVE_REQUEST = iota
    FINALIZE_STATUS = iota
)

type PFCPTransaction struct {
    TransactionStatus uint8
    SequenceNum uint32
    MyChannel chan []byte
    EndPointChannel chan []byte
    SessionChannel chan []byte
    DstIPaddress string
    DstPort string
    MyIPaddress string
    MyPort string
}

func (p *PFCPTransaction) TransactionInitialize() (error) {
    p.TransactionStatus = INITIALIZE_STATUS
    //課題：受信バッファはConfig定義可能なように変更する START
    p.MyChannel = make(chan []byte, 10)
    //END
    return nil
}

func (p *PFCPTransaction) GetMsgType(buffer []byte) uint8 {
    return uint8(buffer[1])
}

func (p *PFCPTransaction) GetIEType(buffer []byte) uint16 {
    return binary.BigEndian.Uint16(buffer[0:2])
}
func (p *PFCPTransaction) GetMsgLength(buffer []byte) uint16 {
    return binary.BigEndian.Uint16(buffer[2:4])
}
func (p *PFCPTransaction) SetSeqNum(seqNum uint32) {
    p.SequenceNum = seqNum
}

func (p *PFCPTransaction) SetDstIPaddress(dstIPaddress string) {
    p.DstIPaddress = dstIPaddress
}

func (p *PFCPTransaction) SetDstPort(dstPort string) {
    p.DstPort = dstPort
}

func (p *PFCPTransaction) SetMyIPaddress(myIPaddress string) {
    p.MyIPaddress = myIPaddress
}

func (p *PFCPTransaction) SetMyPort(myPort string) {
    p.MyPort = myPort
}

func (p *PFCPTransaction) SetSessionChan(sessionchan chan []byte){
    p.SessionChannel = sessionchan
}

func (p *PFCPTransaction) GetTransactionChanel() (chan []byte) {
    return p.MyChannel
}

func (p *PFCPTransaction) GetSeqNum() uint32 {
    return p.SequenceNum
}

func  (p *PFCPTransaction) SetEndPointChanel(endPointChan chan []byte) {
    p.EndPointChannel = endPointChan
}

func (p *PFCPTransaction) SendFinalizeMessage() {
    transactionFinalizeNotifyMessage := new(TransactionFinalizeNotify)
    transactionFinalizeNotifyMessage.MsgType = TRANSACTION_FINALIZE_NOTIFY
    transactionFinalizeNotifyMessage.SeqNum = p.SequenceNum
    sendMessage,_ := json.Marshal(transactionFinalizeNotifyMessage)
    if p.SessionChannel != nil {
        p.SessionChannel <- sendMessage
    }
    if p.EndPointChannel != nil {
        p.EndPointChannel <- sendMessage
    }
}
