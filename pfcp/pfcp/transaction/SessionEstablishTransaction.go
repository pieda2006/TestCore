package pfcp

import (
    "encoding/json"
    "log"
)

type SessionEstablishTransaction struct {
    PFCPTransaction
}


func (p *SessionEstablishTransaction) ExecuteTransaction() {
    for {
        message, ok := <-p.MyChannel
        if !ok {
            return
        }

        var messageBase InternalMessageBase
        json.Unmarshal(message, &messageBase)

        log.Println("Recieve MessageType: ",messageBase.MsgType)

        switch messageBase.MsgType {
        case SIGNAL_RECEVE_NOTIFY:
            var recSignalNotifyMessage RecSignalNotify
            json.Unmarshal(message, &recSignalNotifyMessage)
            messageType := p.GetMsgType(recSignalNotifyMessage.SignalBuffer)
            log.Println("Message Typeの取得 : ",messageType)
            switch messageType {
            case SESSION_ESTABLISHMENT_REQUEST:
                switch p.TransactionStatus {
                case INITIALIZE_STATUS :
                    p.TransactionStatus = RECEAVE_REQUEST
                    recSessionEstablishRequestMessage := new(RecSessionEstablishRequest)
                    nextByte := recSessionEstablishRequestMessage.SessionHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                    bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                    for nextByte < bufferlen {
                        ieType := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:])
                        log.Println("bufferlen: ",bufferlen," nextByte: ",nextByte," IEtype: ",ieType)
                        switch ieType {
                        case NOADID_IE:
                            nextByte = recSessionEstablishRequestMessage.NodeIDData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        case FSEID_IE:
                            nextByte = recSessionEstablishRequestMessage.FSEIDData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        default:
                        }
                    }
                    log.Println("Sessionを生成")
                    sessionFactoryInstance := GetPFCPSessionFactoryInstance()
                    pfcpSessionIns := sessionFactoryInstance.CreateSession()
                    pfcpSessionIns.SetDstIPaddress(p.DstIPaddress)
                    pfcpSessionIns.SetDstPort(p.DstPort)
                    pfcpSessionIns.SetMyIPaddress(p.MyIPaddress)
                    pfcpSessionIns.SetMyPort(p.MyPort)
                    pfcpSessionIns.SetTransactionChan(p.SequenceNum, p.MyChannel)
                    p.SessionChannel = pfcpSessionIns.GetSessionChan()
                    go pfcpSessionIns.ExecuteSession()
                    log.Println("Establish RequestをSessionへ通知")
                    recSessionEstablishRequestMessage.MsgType = REC_SESSION_ESTABLISHMENT_REQUEST
                    sendMessage,_ := json.Marshal(recSessionEstablishRequestMessage)
                    p.SessionChannel <- sendMessage
                    p.TransactionStatus = WAITE_ANS_REQ
              case SENDING_REQUEST :
              case WAITE_ANSWER :
              case RECEAVE_ANSWER :
              case SENDING_ANSWER :
              case WAITE_ANS_REQ :
              case RECEAVE_REQUEST :
              case FINALIZE_STATUS :
              default :
              }
            case SESSION_ESTABLISHMENT_RESPONSE:
                log.Println("SESSION_ESTABLISHMENT_RESPONSE受信")
                log.Println("TransactionStatus: ",p.TransactionStatus)
                switch p.TransactionStatus {
                case INITIALIZE_STATUS :
                case SENDING_REQUEST :
                case WAITE_ANSWER :
                    p.TransactionStatus = RECEAVE_ANSWER
                    recSessionEstablishResponseMessage := new(RecSessionEstablishResponse)
                    nextByte := recSessionEstablishResponseMessage.SessionHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                    bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                    log.Println("IEtype: ",p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:])," nexByte: ",nextByte," BufferLen: ",bufferlen)
                    for nextByte < bufferlen {
                        ieType := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:])
                        log.Println("IEtype: ",ieType," nexByte: ",nextByte," BufferLen: ",bufferlen)
                        switch ieType {
                        case NOADID_IE:
                            nextByte = recSessionEstablishResponseMessage.NodeIDData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        case FSEID_IE:
                            nextByte = recSessionEstablishResponseMessage.FSEIDData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        case CAUSE_IE:
                            nextByte = recSessionEstablishResponseMessage.CauseData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        default:
                            //課題：デコード対象外のIEを読み飛ばす処理を入れる。
                            //END
                        }
                    }
                    log.Println("Establish RequestをSessionへ通知")
                    recSessionEstablishResponseMessage.MsgType = REC_SESSION_ESTABLISHMENT_RESPONSE
                    sendMessage,_ := json.Marshal(recSessionEstablishResponseMessage)
                    p.SessionChannel <- sendMessage
                    p.TransactionStatus = FINALIZE_STATUS
                    log.Println("トランザクション終了処理")
                    p.SendFinalizeMessage()
                    return
                case RECEAVE_ANSWER :
                case SENDING_ANSWER :
                case WAITE_ANS_REQ :
                case RECEAVE_REQUEST :
                case FINALIZE_STATUS :
                default :
                }
            default:
            }
        case SEND_SESSION_ESTABLISHMENT_REQUEST :
            switch p.TransactionStatus {
            case INITIALIZE_STATUS :
                p.TransactionStatus = SENDING_REQUEST
                log.Println("SESSION_ESTABLISHMENT_REQUESTE送信要求受信")
                var sendEstablishRequestMessage SendEstablishRequest
                json.Unmarshal(message, &sendEstablishRequestMessage)
                log.Println("送信先のEndPoint取得")
                endpointFactoryIns := GetPFCPEndPointFactoryInstance()
                endpointIF := endpointFactoryIns.GetEndPointIF(p.DstIPaddress)
                p.EndPointChannel = endpointIF.GetEndPointChan()
                endpointIF.SetTransactionChan(p.SequenceNum, p.MyChannel)
                log.Println("ヘッダ情報設定")
                sendEstablishRequestMessage.SessionHeaderData.SetMsgType(SESSION_ESTABLISHMENT_REQUEST)
                sendEstablishRequestMessage.SessionHeaderData.SetSequenceNum(p.SequenceNum)
                sendEstablishRequestMessage.SessionHeaderData.SetInformationElement(&sendEstablishRequestMessage.FSEIDData)
                sendEstablishRequestMessage.SessionHeaderData.SetInformationElement(&sendEstablishRequestMessage.NodeIDData)
                log.Println("EndPointへ応答信号送信要求受信")
                sendSignalMessage := new(SendSignalRequest)
                sendSignalMessage.MsgType = SIGNAL_SEND_REQUEST
                sendSignalMessage.SignalBuffer = sendEstablishRequestMessage.SessionHeaderData.CreateSignal()
                sendMessage,_ := json.Marshal(sendSignalMessage)
                p.EndPointChannel <- sendMessage
                p.TransactionStatus = WAITE_ANSWER
                log.Println("状態遷移： ",p.TransactionStatus)
            case SENDING_REQUEST :
            case WAITE_ANSWER :
            case RECEAVE_ANSWER :
            case SENDING_ANSWER :
            case WAITE_ANS_REQ :
            case RECEAVE_REQUEST :
            case FINALIZE_STATUS :
            default :
            }
        case SEND_SESSION_ESTABLISHMENT_RESPONSE :
            switch p.TransactionStatus {
            case INITIALIZE_STATUS :
            case SENDING_REQUEST :
            case WAITE_ANSWER :
            case RECEAVE_ANSWER :
            case SENDING_ANSWER :
            case WAITE_ANS_REQ :
                p.TransactionStatus = SENDING_ANSWER
                log.Println("SESSION_ESTABLISHMENT_RESPONSE送信要求受信")
                var sendEstablishResponseMessage SendEstablishResponse
                json.Unmarshal(message, &sendEstablishResponseMessage)
                log.Println("ヘッダ情報設定")
                sendEstablishResponseMessage.SessionHeaderData.SetMsgType(SESSION_ESTABLISHMENT_RESPONSE)
                sendEstablishResponseMessage.SessionHeaderData.SetSequenceNum(p.SequenceNum)
                sendEstablishResponseMessage.SessionHeaderData.SetInformationElement(&sendEstablishResponseMessage.CauseData)
                sendEstablishResponseMessage.SessionHeaderData.SetInformationElement(&sendEstablishResponseMessage.NodeIDData)
                sendEstablishResponseMessage.SessionHeaderData.SetInformationElement(&sendEstablishResponseMessage.FSEIDData)
                log.Println("EndPointへ応答信号送信要求受信")
                sendSignalMessage := new(SendSignalRequest)
                sendSignalMessage.MsgType = SIGNAL_SEND_REQUEST
                sendSignalMessage.SignalBuffer = sendEstablishResponseMessage.SessionHeaderData.CreateSignal()
                sendMessage,_ := json.Marshal(sendSignalMessage)
                p.EndPointChannel <- sendMessage
                p.TransactionStatus = FINALIZE_STATUS
                log.Println("トランザクション終了処理:セッションへ終了通知")
                p.SendFinalizeMessage()
                return
            case RECEAVE_REQUEST :
            case FINALIZE_STATUS :
            default :
            }
        default:
        }
    }
}
