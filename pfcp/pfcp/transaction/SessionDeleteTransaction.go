package pfcp

import (
    "encoding/json"
    "log"
)

type SessionDeleteTransaction struct {
    PFCPTransaction
}


func (p *SessionDeleteTransaction) ExecuteTransaction() {
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
            case SESSION_DELETE_REQUEST:
                switch p.TransactionStatus {
                case INITIALIZE_STATUS :
                    p.TransactionStatus = RECEAVE_REQUEST
                    recSessionDeleteRequestMessage := new(RecSessionDeleteRequest)
                    recSessionDeleteRequestMessage.SessionHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                    sessionFactoryIns := GetPFCPSessionFactoryInstance()
                    log.Println("セッションを取得 SEID:",recSessionDeleteRequestMessage.SessionHeaderData.GetSEID())
                    sessionInstance := sessionFactoryIns.GetSession(recSessionDeleteRequestMessage.SessionHeaderData.GetSEID())
                    sessionInstance.SetTransactionChan(p.SequenceNum, p.MyChannel)
                    p.SessionChannel = sessionInstance.GetSessionChan()
                    log.Println("Delete RequestをSessionへ通知")
                    recSessionDeleteRequestMessage.MsgType = REC_SESSION_DELETE_REQUEST
                    sendMessage,_ := json.Marshal(recSessionDeleteRequestMessage)
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
            case SESSION_DELETE_RESPONSE:
                log.Println("SESSION_ESTABLISHMENT_RESPONSE受信")
                log.Println("TransactionStatus: ",p.TransactionStatus)
                switch p.TransactionStatus {
                case INITIALIZE_STATUS :
                case SENDING_REQUEST :
                case WAITE_ANSWER :
                    p.TransactionStatus = RECEAVE_ANSWER
                    recSessionDeleteResponseMessage := new(RecSessionDeleteResponse)
                    nextByte := recSessionDeleteResponseMessage.SessionHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                    bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                    log.Println("IEtype: ",p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:])," nexByte: ",nextByte," BufferLen: ",bufferlen)
                    for nextByte < bufferlen {
                        ieType := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:])
                        log.Println("IEtype: ",ieType," nexByte: ",nextByte," BufferLen: ",bufferlen)
                        switch ieType {
                        case CAUSE_IE:
                            nextByte = recSessionDeleteResponseMessage.CauseData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        default:
                            //課題：デコード対象外のIEを読み飛ばす処理を入れる。
                            //END
                        }
                    }
                    log.Println("Delete ResponseをSessionへ通知")
                    recSessionDeleteResponseMessage.MsgType = REC_SESSION_DELETE_RESPONSE
                    sendMessage,_ := json.Marshal(recSessionDeleteResponseMessage)
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
        case SEND_SESSION_DELETE_REQUEST :
            switch p.TransactionStatus {
            case INITIALIZE_STATUS :
                p.TransactionStatus = SENDING_REQUEST
                log.Println("SESSION_DELETE_REQUESTE送信要求受信")
                var sendDeleteRequestMessage SendDeleteRequest
                json.Unmarshal(message, &sendDeleteRequestMessage)
                log.Println("送信先のEndPoint取得")
                endpointFactoryIns := GetPFCPEndPointFactoryInstance()
                endpointIF := endpointFactoryIns.GetEndPointIF(p.DstIPaddress)
                p.EndPointChannel = endpointIF.GetEndPointChan()
                endpointIF.SetTransactionChan(p.SequenceNum, p.MyChannel)
                log.Println("ヘッダ情報設定")
                sendDeleteRequestMessage.SessionHeaderData.SetMsgType(SESSION_DELETE_REQUEST)
                sendDeleteRequestMessage.SessionHeaderData.SetSequenceNum(p.SequenceNum)
                log.Println("EndPointへ応答信号送信要求受信")
                sendSignalMessage := new(SendSignalRequest)
                sendSignalMessage.MsgType = SIGNAL_SEND_REQUEST
                sendSignalMessage.SignalBuffer = sendDeleteRequestMessage.SessionHeaderData.CreateSignal()
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
        case SEND_SESSION_DELETE_RESPONSE :
            switch p.TransactionStatus {
            case INITIALIZE_STATUS :
            case SENDING_REQUEST :
            case WAITE_ANSWER :
            case RECEAVE_ANSWER :
            case SENDING_ANSWER :
            case WAITE_ANS_REQ :
                p.TransactionStatus = SENDING_ANSWER
                log.Println("SESSION_DELETE_RESPONSE送信要求受信")
                var sendDeleteResponseMessage SendDeleteResponse
                json.Unmarshal(message, &sendDeleteResponseMessage)
                log.Println("ヘッダ情報設定")
                sendDeleteResponseMessage.SessionHeaderData.SetMsgType(SESSION_DELETE_RESPONSE)
                sendDeleteResponseMessage.SessionHeaderData.SetSequenceNum(p.SequenceNum)
                sendDeleteResponseMessage.SessionHeaderData.SetInformationElement(&sendDeleteResponseMessage.CauseData)
                log.Println("EndPointへ応答信号送信要求受信")
                sendSignalMessage := new(SendSignalRequest)
                sendSignalMessage.MsgType = SIGNAL_SEND_REQUEST
                sendSignalMessage.SignalBuffer = sendDeleteResponseMessage.SessionHeaderData.CreateSignal()
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
        case SESSION_FINALIZE_NOTIFY:
            log.Println("セッション終了通知受信")
            p.SessionChannel = nil
        default:
        }
    }
}
