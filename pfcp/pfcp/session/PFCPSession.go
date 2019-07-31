package pfcp

import (
    "log"
    "encoding/json"
)

const (
    INITIALIZE_SESSION_STATUS = iota
    SEND_ESTABLISH_NOTIFY = iota
    SENDING_ESTABLISH_REQ = iota
    SENDED_ESTABLISH_REQ = iota
    SEND_ESTABLISH_ANS = iota
    REC_ESTABLISH_REQ = iota
    REC_ESTABLISH_ANS = iota
    ESTABLISH_STATUS = iota
    SENDING_DELETE_REQ = iota
    SENDED_DELETE_REQ = iota
    REC_DELETE_ANS = iota
    REC_DELETE_REQ = iota
    SEND_DELETE_NOTIFY = iota
    SEND_DELETE_ANS = iota
    FINALIZE_SESSION_STATUS = iota
)

type PFCPSession struct {
    SessionStatus uint8
    MySEID uint64
    DstSEID uint64
    MyChannel chan []byte
    TransactionChannel map[uint32](chan []byte)
    DstIPaddress string
    DstPort string
    MyIPaddress string
    MyPort string
}

func (p *PFCPSession) SessionInitialize() (error) {
    p.SessionStatus = INITIALIZE_SESSION_STATUS
    //課題：バッファはConfigで設定可能にする START
    p.MyChannel = make(chan []byte, 10)
    p.TransactionChannel = make(map[uint32](chan []byte),10)
    //END
    return nil
}

func (p *PFCPSession) SendFinalizeMessage() {
    sessionFinalizeNotifyMessage := new(SessionFinalizeNotify)
    sessionFinalizeNotifyMessage.MsgType = SESSION_FINALIZE_NOTIFY
    sendMessage,_ := json.Marshal(sessionFinalizeNotifyMessage)
    for key, value := range p.TransactionChannel {
        log.Println("key: ",key)
        value <- sendMessage
    }
}

func (p *PFCPSession) SetMySEID(seid uint64) {
    p.MySEID = seid
}


func (p *PFCPSession) SetDstIPaddress(ipaddr string) {
    p.DstIPaddress = ipaddr
}

func (p *PFCPSession) SetDstPort(port string) {
    p.DstPort = port
}

func (p *PFCPSession) SetMyIPaddress(ipaddr string) {
    p.MyIPaddress = ipaddr
}

func (p *PFCPSession) SetMyPort(port string) {
    p.MyPort = port
}

func (p *PFCPSession) SetTransactionChan(seqnum uint32, transactionChan chan []byte) {
    p.TransactionChannel[seqnum] = transactionChan
}

func (p *PFCPSession) GetSessionChan() (chan []byte) {
    return p.MyChannel
}

func (p *PFCPSession) CreateTransaction(msgType uint32) PFCPTransactionIF {
    log.Println("トランザクション生成 MessageType: ",msgType)
    transactionFactry := GetPFCPTransactionFactoryInstance()
    msgTransaction := transactionFactry.CreateTransaction(msgType, 0)
    msgTransaction.SetDstIPaddress(p.DstIPaddress)
    msgTransaction.SetDstPort(p.DstPort)
    msgTransaction.SetMyIPaddress(p.MyIPaddress)
    msgTransaction.SetMyPort(p.MyPort)
    transactionChan := msgTransaction.GetTransactionChanel()
    p.TransactionChannel[msgTransaction.GetSeqNum()] = transactionChan
    msgTransaction.SetSessionChan(p.MyChannel)
    log.Println("トランザクション起動")
    go msgTransaction.ExecuteTransaction()
    return msgTransaction
}
func (p *PFCPSession) ExecuteSession() {
    for {
        log.Println("Session処理開始")
        message, ok := <-p.MyChannel
        if !ok {
            return
        }
        log.Println("イベント受信")
        var messageBase InternalMessageBase
        json.Unmarshal(message, &messageBase)
        switch messageBase.MsgType {
        case SEND_SESSION_ESTABLISHMENT_REQUEST:
            switch p.SessionStatus {
            case INITIALIZE_SESSION_STATUS:
                p.SessionStatus = SENDING_ESTABLISH_REQ
                log.Println("Session Establish Requestの送信要求受信")
                var sendEstablishRequestMessage SendEstablishRequest
                json.Unmarshal(message, &sendEstablishRequestMessage)
                log.Println("ヘッダ情報設定")
                sendEstablishRequestMessage.SessionHeaderData.Initialize()
                sendEstablishRequestMessage.SessionHeaderData.SetSEID(0)
                log.Println("IE情報設定")
                sendEstablishRequestMessage.FSEIDData.IEInitialize()
                sendEstablishRequestMessage.FSEIDData.SetSEID(p.MySEID)
                sendEstablishRequestMessage.FSEIDData.SetIPaddressType(FSEID_IPV4)
                sendEstablishRequestMessage.FSEIDData.SetIPv4Address(p.MyIPaddress)
                sendEstablishRequestMessage.NodeIDData.IEInitialize()
                sendEstablishRequestMessage.NodeIDData.SetNodeIDType(NODETYPE_IPV4)
                sendEstablishRequestMessage.NodeIDData.SetNodeIDValueIPv4(p.MyIPaddress)
                log.Println("トランザクション生成")
                transactionIF := p.CreateTransaction(SESSION_ESTABLISHMENT_REQUEST)
                log.Println("Transactionへ通知")
                sendMessage,_ := json.Marshal(sendEstablishRequestMessage)
                transactionIF.GetTransactionChanel() <- sendMessage
                p.SessionStatus = SENDED_ESTABLISH_REQ
            case SEND_ESTABLISH_NOTIFY:
            case SENDING_ESTABLISH_REQ:
            case SENDED_ESTABLISH_REQ:
            case SEND_ESTABLISH_ANS:
            case REC_ESTABLISH_REQ:
            case REC_ESTABLISH_ANS:
            case ESTABLISH_STATUS:
            case SENDING_DELETE_REQ:
            case SENDED_DELETE_REQ:
            case REC_DELETE_ANS:
            case REC_DELETE_REQ:
            case SEND_DELETE_NOTIFY:
            case SEND_DELETE_ANS:
            case FINALIZE_SESSION_STATUS:
            default:
            }
        case REC_SESSION_ESTABLISHMENT_RESPONSE:
            switch p.SessionStatus {
            case INITIALIZE_SESSION_STATUS:
            case SEND_ESTABLISH_NOTIFY:
            case SENDING_ESTABLISH_REQ:
            case SENDED_ESTABLISH_REQ:
                log.Println("Session Establish Responseの受信")
                p.SessionStatus = REC_ESTABLISH_ANS
                var recSessionEstablishResponseMessage RecSessionEstablishResponse
                json.Unmarshal(message, &recSessionEstablishResponseMessage)
                log.Println("対向装置から通知されたSEIDを保持")
                p.DstSEID = recSessionEstablishResponseMessage.FSEIDData.GetSEID()
                //課題：上位のインスタンス通知 START
                //END
                p.SessionStatus = ESTABLISH_STATUS
            case SEND_ESTABLISH_ANS:
            case REC_ESTABLISH_REQ:
            case REC_ESTABLISH_ANS:
            case ESTABLISH_STATUS:
            case SENDING_DELETE_REQ:
            case SENDED_DELETE_REQ:
            case REC_DELETE_ANS:
            case REC_DELETE_REQ:
            case SEND_DELETE_NOTIFY:
            case SEND_DELETE_ANS:
            case FINALIZE_SESSION_STATUS:
            default:
            }
        case REC_SESSION_ESTABLISHMENT_REQUEST:
            switch p.SessionStatus {
            case INITIALIZE_SESSION_STATUS:
                p.SessionStatus = REC_ESTABLISH_REQ
                log.Println("対向からSession Establish Request受信")
                var recSessionEstablishRequestMessage RecSessionEstablishRequest
                json.Unmarshal(message, &recSessionEstablishRequestMessage)

                p.DstSEID = recSessionEstablishRequestMessage.FSEIDData.GetSEID()
                log.Println("対向装置のSEIDを保持 SEID: ",p.DstSEID)
                p.SessionStatus = SEND_ESTABLISH_NOTIFY
                //課題：PDR処理実装 START
                //処理中断が入るかも。
                //END
                p.SessionStatus = SEND_ESTABLISH_ANS
                log.Println("Session Establish Response送信")
                sendEstablishResponseMessage := new(SendEstablishResponse)
                sendEstablishResponseMessage.MsgType = SEND_SESSION_ESTABLISHMENT_RESPONSE
                log.Println("送信メッセージ初期化処理")
                sendEstablishResponseMessage.SessionHeaderData.Initialize()
                sendEstablishResponseMessage.NodeIDData.IEInitialize()
                sendEstablishResponseMessage.CauseData.IEInitialize()
                sendEstablishResponseMessage.FSEIDData.IEInitialize()
                log.Println("Header情報生成")
                sendEstablishResponseMessage.SessionHeaderData.SetSEID(p.DstSEID)
                sendEstablishResponseMessage.SessionHeaderData.SetMessagePriority(recSessionEstablishRequestMessage.SessionHeaderData.GetMessagePriority())
                log.Println("IE情報生成")
                sendEstablishResponseMessage.NodeIDData.SetNodeIDType(NODETYPE_IPV4)
                sendEstablishResponseMessage.NodeIDData.SetNodeIDValueIPv4(p.MyIPaddress)
                sendEstablishResponseMessage.CauseData.SetCause(REQUEST_ACCEPTED)
                sendEstablishResponseMessage.FSEIDData.SetSEID(p.MySEID)
                sendEstablishResponseMessage.FSEIDData.SetIPaddressType(FSEID_IPV4)
                sendEstablishResponseMessage.FSEIDData.SetIPv4Address(p.MyIPaddress)
                log.Println("Transactionへ通知")
                sendMessage,_ := json.Marshal(sendEstablishResponseMessage)
                log.Println("SeqNum: ",recSessionEstablishRequestMessage.SessionHeaderData.GetSequenceNum())
                p.TransactionChannel[recSessionEstablishRequestMessage.SessionHeaderData.GetSequenceNum()] <- sendMessage
                p.SessionStatus = ESTABLISH_STATUS
            case SEND_ESTABLISH_NOTIFY:
            case SENDING_ESTABLISH_REQ:
            case SENDED_ESTABLISH_REQ:
            case SEND_ESTABLISH_ANS:
            case REC_ESTABLISH_REQ:
            case REC_ESTABLISH_ANS:
            case ESTABLISH_STATUS:
            case SENDING_DELETE_REQ:
            case SENDED_DELETE_REQ:
            case REC_DELETE_ANS:
            case REC_DELETE_REQ:
            case SEND_DELETE_NOTIFY:
            case SEND_DELETE_ANS:
            case FINALIZE_SESSION_STATUS:
            default:
            }
        case SEND_SESSION_DELETE_REQUEST:
            switch p.SessionStatus {
            case INITIALIZE_SESSION_STATUS:
            case SEND_ESTABLISH_NOTIFY:
            case SENDING_ESTABLISH_REQ:
            case SENDED_ESTABLISH_REQ:
            case SEND_ESTABLISH_ANS:
            case REC_ESTABLISH_REQ:
            case REC_ESTABLISH_ANS:
            case ESTABLISH_STATUS:
                p.SessionStatus = SENDING_DELETE_REQ
                log.Println("Session Delete Requestの送信要求受信")
                var ｓendDeleteRequestMessage SendDeleteRequest
                json.Unmarshal(message, &ｓendDeleteRequestMessage)
                log.Println("ヘッダ情報設定")
                ｓendDeleteRequestMessage.SessionHeaderData.Initialize()
                ｓendDeleteRequestMessage.SessionHeaderData.SetSEID(p.DstSEID)
                log.Println("トランザクション生成")
                transactionIF := p.CreateTransaction(SESSION_DELETE_REQUEST)
                log.Println("Transactionへ通知")
                sendMessage,_ := json.Marshal(ｓendDeleteRequestMessage)
                transactionIF.GetTransactionChanel() <- sendMessage
                p.SessionStatus = SENDED_DELETE_REQ
            case SENDING_DELETE_REQ:
            case SENDED_DELETE_REQ:
            case REC_DELETE_ANS:
            case REC_DELETE_REQ:
            case SEND_DELETE_NOTIFY:
            case SEND_DELETE_ANS:
            case FINALIZE_SESSION_STATUS:
            default:
            }
        case REC_SESSION_DELETE_REQUEST:
          switch p.SessionStatus {
          case INITIALIZE_SESSION_STATUS:
          case SEND_ESTABLISH_NOTIFY:
          case SENDING_ESTABLISH_REQ:
          case SENDED_ESTABLISH_REQ:
          case SEND_ESTABLISH_ANS:
          case REC_ESTABLISH_REQ:
          case REC_ESTABLISH_ANS:
          case ESTABLISH_STATUS:
              p.SessionStatus = REC_DELETE_REQ
              log.Println("対向からSession Delete Request受信")
              var recSessionDeleteRequestMessage RecSessionDeleteRequest
              json.Unmarshal(message, &recSessionDeleteRequestMessage)
              log.Println("削除対象のSEIDと一致するかチェック")
              if p.MySEID == recSessionDeleteRequestMessage.SessionHeaderData.GetSEID() {
                  p.SessionStatus = SEND_DELETE_NOTIFY
                  //課題：PDR処理実装 START
                  //処理中断が入るかも。
                  //END
                  p.SessionStatus = SEND_DELETE_ANS
                  log.Println("Session Delete Response送信")
                  sendDeleteResponseMessage := new(SendDeleteResponse)
                  sendDeleteResponseMessage.MsgType = SEND_SESSION_DELETE_RESPONSE
                  log.Println("送信メッセージ初期化処理")
                  sendDeleteResponseMessage.SessionHeaderData.Initialize()
                  sendDeleteResponseMessage.CauseData.IEInitialize()
                  log.Println("Header情報生成")
                  sendDeleteResponseMessage.SessionHeaderData.SetSEID(p.DstSEID)
                  sendDeleteResponseMessage.SessionHeaderData.SetMessagePriority(recSessionDeleteRequestMessage.SessionHeaderData.GetMessagePriority())
                  log.Println("IE情報生成")
                  sendDeleteResponseMessage.CauseData.SetCause(REQUEST_ACCEPTED)
                  log.Println("Transactionへ通知")
                  sendMessage,_ := json.Marshal(sendDeleteResponseMessage)
                  log.Println("SeqNum: ",recSessionDeleteRequestMessage.SessionHeaderData.GetSequenceNum())
                  p.TransactionChannel[recSessionDeleteRequestMessage.SessionHeaderData.GetSequenceNum()] <- sendMessage
                  p.SessionStatus = FINALIZE_SESSION_STATUS
                  p.SendFinalizeMessage()
                  return
              }
          case SENDING_DELETE_REQ:
          case SENDED_DELETE_REQ:
          case REC_DELETE_ANS:
          case REC_DELETE_REQ:
          case SEND_DELETE_NOTIFY:
          case SEND_DELETE_ANS:
          case FINALIZE_SESSION_STATUS:
          default:
          }
        case REC_SESSION_DELETE_RESPONSE:
          switch p.SessionStatus {
          case INITIALIZE_SESSION_STATUS:
          case SEND_ESTABLISH_NOTIFY:
          case SENDING_ESTABLISH_REQ:
          case SENDED_ESTABLISH_REQ:
          case SEND_ESTABLISH_ANS:
          case REC_ESTABLISH_REQ:
          case REC_ESTABLISH_ANS:
          case ESTABLISH_STATUS:
          case SENDING_DELETE_REQ:
          case SENDED_DELETE_REQ:
            p.SessionStatus = REC_DELETE_ANS
            //課題：上位のインスタンス通知 START
            //END
            p.SessionStatus = FINALIZE_SESSION_STATUS
            p.SendFinalizeMessage()
            return
          case REC_DELETE_ANS:
          case REC_DELETE_REQ:
          case SEND_DELETE_NOTIFY:
          case SEND_DELETE_ANS:
          case FINALIZE_SESSION_STATUS:
          default:
          }
        case TRANSACTION_FINALIZE_NOTIFY:
            log.Println("トランザクション終了通知受信")
            var ｔransactionFinalizeNotifyMessage TransactionFinalizeNotify
            json.Unmarshal(message, &ｔransactionFinalizeNotifyMessage)
            delete(p.TransactionChannel,ｔransactionFinalizeNotifyMessage.SeqNum)
        default:
        }
    }
}
