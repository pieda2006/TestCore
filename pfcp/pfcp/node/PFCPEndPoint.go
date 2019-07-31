package pfcp

import (
  "net"
  "time"
  "log"
  "encoding/json"
  "encoding/binary"
)

type PFCPEndPoint struct {
  PFCPNode
  EndPointStatus uint8
  EndPointType uint8
  //トランザクション上限数はConfig定義に後で変える。
  TransactionChannel map[uint32](chan []byte)
  RecTime int64
  EndPointRecTime uint64
  MyIPaddress string
  MyPort string
}

const (
    ENDPOINT_UNCONECT = iota
    ENDPOINT_WAITE_ANSWER = iota
    ENDPOINT_CONECT = iota
    REC_ASSOCIATION_REQ = iota
)

func (p *PFCPEndPoint) EndPointInitialize(endPointType uint8) (error) {
    p.EndPointStatus = ENDPOINT_UNCONECT
    p.EndPointType = endPointType
    //課題：Transactionの最大数はConfig設定に切り替える。 START
    p.MyChannel = make(chan []byte, 10)
    p.TransactionChannel = make(map[uint32](chan []byte),10)
    //END
    current := time.Now()
    p.RecTime = current.Unix()
    p.EndPointRecTime = 0
    return nil
}

func (p *PFCPEndPoint) SetMyIPaddress(ipaddress string) {
    p.MyIPaddress = ipaddress
}

func (p *PFCPEndPoint) SetMyPort(port string) {
    p.MyPort = port
}

func (p *PFCPEndPoint) SetTransactionChan(seqNum uint32, transactionchan chan []byte){
    p.TransactionChannel[seqNum] = transactionchan
}

func (p *PFCPEndPoint) CreateTransaction(msgType uint32, seqNum uint32) PFCPTransactionIF {
    log.Println("トランザクション生成 MessageType: ",msgType,"シーケンス番号: ",seqNum)
    transactionFactry := GetPFCPTransactionFactoryInstance()
    msgTransaction := transactionFactry.CreateTransaction(msgType, seqNum)
    msgTransaction.SetDstIPaddress(p.IPaddress)
    msgTransaction.SetDstPort(p.Port)
    msgTransaction.SetMyIPaddress(p.MyIPaddress)
    msgTransaction.SetMyPort(p.MyPort)
    transactionChan := msgTransaction.GetTransactionChanel()
    p.TransactionChannel[msgTransaction.GetSeqNum()] = transactionChan
    msgTransaction.SetEndPointChanel(p.MyChannel)
    log.Println("トランザクション起動")
    go msgTransaction.ExecuteTransaction()
    return msgTransaction
}

func (p *PFCPEndPoint) ExecuteEndPoint() (error){
    log.Println("送信用UDPコネクション生成")
    conn, _ := net.Dial("udp", p.IPaddress + ":" + p.Port)
    log.Println("送信先IPaddress：",p.IPaddress," 宛先ポート：",p.Port)
    defer conn.Close()
    if p.EndPointType == ENDPOINT_CLIENT {
        if p.EndPointStatus == ENDPOINT_UNCONECT {
           log.Println("Association Requestトランザクション生成")
           associationReqTransaction := p.CreateTransaction(ASSOCIATION_SETUP_REQUEST, 0)
           log.Println("Association Request送信要求通知")
           associationRequest := new(SendAssociationSetupRequest)
           associationRequest.MsgType = SEND_ASSOCIATION_REQUEST
           associationRequest.RecTimeStampData.IEInitialize()
           associationRequest.RecTimeStampData.SetTimeStamp(uint32(p.RecTime))
           message,_ := json.Marshal(associationRequest)
           associationReqTransaction.GetTransactionChanel() <- message
           p.EndPointStatus = ENDPOINT_WAITE_ANSWER
        }
    }

    log.Println("EndPointインスタンスの受信処理開始")
    for {
        log.Println("EndPointインスタンスのチャネル受信")

        recMessage, ok := <-p.MyChannel
        if !ok {
            return nil
        }
        var messageBase InternalMessageBase
        json.Unmarshal(recMessage, &messageBase)

        switch messageBase.MsgType {
        case SIGNAL_RECEVE_NOTIFY:
            log.Println("SIGNAL_RECEVE_NOTIFY受信")

            var recSignalNotifyMessage RecSignalNotify
            json.Unmarshal(recMessage, &recSignalNotifyMessage)

            messageType := uint32(recSignalNotifyMessage.SignalBuffer[1])
            log.Println("受信信号messageType： ", messageType)

            if (messageType == ASSOCIATION_SETUP_REQUEST && p.EndPointStatus == ENDPOINT_UNCONECT) ||
               (messageType == ASSOCIATION_SETUP_RESPONSE && p.EndPointStatus == ENDPOINT_WAITE_ANSWER) ||
               p.EndPointStatus == ENDPOINT_CONECT {


                var seqNum uint32
                if messageType < 50 {
                    seqNum = binary.BigEndian.Uint32(recSignalNotifyMessage.SignalBuffer[4:8]) >> 8
                } else {
                    seqNum = binary.BigEndian.Uint32(recSignalNotifyMessage.SignalBuffer[12:16]) >> 8
                }

                log.Println("受信信号SeqNum: ",seqNum)

                transactionChan, ok := p.TransactionChannel[seqNum]
                if ok == false {
                    log.Println("トランザクションが見つからなかったため新規生成")
                    transactionIF := p.CreateTransaction(messageType, seqNum)
                    transactionChan = transactionIF.GetTransactionChanel()
                }
                log.Println("トランザクションへ受信信号を通知 MsgType:",recSignalNotifyMessage.MsgType)
                transactionChan <- recMessage
            }
        case SIGNAL_SEND_REQUEST:
            log.Println("SIGNAL_SEND_REQUEST受信")
            var sendSignalRequestMessage SendSignalRequest
            json.Unmarshal(recMessage, &sendSignalRequestMessage)
            messageType := uint32(sendSignalRequestMessage.SignalBuffer[1])
            if (messageType == ASSOCIATION_SETUP_REQUEST && p.EndPointStatus != ENDPOINT_CONECT) ||
               (messageType == ASSOCIATION_SETUP_RESPONSE && p.EndPointStatus != ENDPOINT_CONECT) ||
                p.EndPointStatus == ENDPOINT_CONECT {
                log.Println("対向のEndPointへ送信")
                conn.Write(sendSignalRequestMessage.SignalBuffer)
            }
        case REC_ASSOCIATION_REQUEST:
            switch p.EndPointStatus {
            case ENDPOINT_UNCONECT :
                log.Println("REC_ASSOCIATION_REQUEST受信")
                p.EndPointStatus = REC_ASSOCIATION_REQ
                var recAssociationRequestMessage RecAssociationRequest
                json.Unmarshal(recMessage, &recAssociationRequestMessage)
                log.Println("対向装置の再開時刻を保持しておく")
                p.EndPointRecTime = uint64(recAssociationRequestMessage.RecoveryTimeData.GetTimeStamp())
                log.Println("SEND_ASSOCIATION_RESPONSEをトランザクションへ通知")
                transactionChan, ok := p.TransactionChannel[recAssociationRequestMessage.NodeHeaderData.GetSequenceNum()]
                if ok == false {
                    //課題：チャネル取得失敗時の処理を入れる START
                    //END
                }
                associationSetupResponse := new(SendAssociationSetupResponse)
                log.Println("送信信号初期化")
                associationSetupResponse.MsgType = SEND_ASSOCIATION_RESPONSE
                associationSetupResponse.RecTimeStampData.IEInitialize()
                associationSetupResponse.RecTimeStampData.SetTimeStamp(uint32(p.RecTime))
                associationSetupResponse.CauseData.IEInitialize()
                associationSetupResponse.CauseData.SetCause(REQUEST_ACCEPTED)
                message,_ := json.Marshal(associationSetupResponse)
                transactionChan <- message
                log.Println("コネクション状態を接続中に変更")
                p.EndPointStatus = ENDPOINT_CONECT
            case ENDPOINT_WAITE_ANSWER :
            case ENDPOINT_CONECT :
            case REC_ASSOCIATION_REQ :
            }
        case REC_ASSOCIATION_RESPONSE:
          switch p.EndPointStatus {
          case ENDPOINT_UNCONECT :
          case ENDPOINT_WAITE_ANSWER :
              log.Println("REC_ASSOCIATION_RESPONSE受信")
              var recAssociationResponseMessage RecAssociationResponse
              json.Unmarshal(recMessage, &recAssociationResponseMessage)
              log.Println("対向装置の再開時刻を保持しておく")
              p.EndPointRecTime = uint64(recAssociationResponseMessage.RecoveryTimeData.GetTimeStamp())
              log.Println("コネクション状態を接続中に変更")
              p.EndPointStatus = ENDPOINT_CONECT
          case ENDPOINT_CONECT :
          case REC_ASSOCIATION_REQ :
          }
        case REC_HEATBEAT_REQUEST:
            log.Println("REC_HEATBEAT_REQUEST受信")
            var recHeatBeatRequestMessage RecHeatBeatRequest
            json.Unmarshal(recMessage, &recHeatBeatRequestMessage)
            log.Println("受信した再開時刻をチェック")
            if p.EndPointRecTime != uint64(recHeatBeatRequestMessage.RecoveryTimeData.GetTimeStamp()) {
                //課題：再開時刻が違うのでコネクション切断処理を入れる START
                //END
            }
            log.Println("SEND_HEATBEAT_RESPONSEをトランザクションへ通知")
            transactionChan, ok := p.TransactionChannel[recHeatBeatRequestMessage.NodeHeaderData.GetSequenceNum()]
            if ok == false {
                //課題：チャネル取得失敗時の処理を入れる START
                //END
            }
            heatbeatResponse := new(SendHeatBeatResponse)
            heatbeatResponse.MsgType = SEND_HEATBEAT_RESPONSE
            heatbeatResponse.RecoveryTimeData.IEInitialize()
            heatbeatResponse.RecoveryTimeData.SetTimeStamp(uint32(p.RecTime))
            message,_ := json.Marshal(heatbeatResponse)
            transactionChan <- message
          case REC_HEATBEAT_RESPONSE:
              log.Println("REC_HEATBEAT_RESPONSE受信")
              var recHeatBeatResponseMessage RecHeatBeatResponse
              json.Unmarshal(recMessage, &recHeatBeatResponseMessage)
              log.Println("受信した再開時刻をチェック")
              if p.EndPointRecTime != uint64(recHeatBeatResponseMessage.RecoveryTimeData.GetTimeStamp()) {
                  //課題：再開時刻が違うのでコネクション切断処理を入れる START
                  //END
              }
        case TRANSACTION_FINALIZE_NOTIFY:
            log.Println("トランザクション終了通知受信")
            var ｔransactionFinalizeNotifyMessage TransactionFinalizeNotify
            json.Unmarshal(recMessage, &ｔransactionFinalizeNotifyMessage)
            delete(p.TransactionChannel,ｔransactionFinalizeNotifyMessage.SeqNum)
        default:
        }
    }
    return nil
}
