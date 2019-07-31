package pfcp

import (
    "encoding/json"
    "log"
)

type HeatBeatTransaction struct {
    PFCPTransaction
}

func (p *HeatBeatTransaction) ExecuteTransaction() {
    for {
        message, ok := <-p.MyChannel
        if !ok {
            return
        }

        var messageBase InternalMessageBase
        json.Unmarshal(message, &messageBase)

        switch messageBase.MsgType {
        case SEND_HEATBEAT_REQUEST:
            switch p.TransactionStatus {
            case INITIALIZE_STATUS :
                p.TransactionStatus = SENDING_REQUEST
                var sendHeatBeatRequestMessage SendHeatBeatRequest
                json.Unmarshal(message, &sendHeatBeatRequestMessage)
                sendHeatBeatRequestMessage.NodeHeaderData.Initialize()
                sendHeatBeatRequestMessage.NodeHeaderData.SetMsgType(HEATBEAT_REQUEST)
                sendHeatBeatRequestMessage.NodeHeaderData.SetSequenceNum(p.SequenceNum)
                sendHeatBeatRequestMessage.NodeHeaderData.SetInformationElement(&sendHeatBeatRequestMessage.RecTimeStampData)
                signalSendReq := new(SendSignalRequest)
                signalSendReq.MsgType = SIGNAL_SEND_REQUEST
                signalSendReq.SignalBuffer = sendHeatBeatRequestMessage.NodeHeaderData.CreateSignal()
                sendMessage,_ := json.Marshal(signalSendReq)
                p.EndPointChannel <- sendMessage
                p.TransactionStatus = WAITE_ANSWER
            case SENDING_REQUEST :
            case WAITE_ANSWER :
            case RECEAVE_ANSWER :
            case SENDING_ANSWER :
            case WAITE_ANS_REQ :
            case RECEAVE_REQUEST :
            case FINALIZE_STATUS :
            default :
            }


        case SEND_HEATBEAT_RESPONSE:
              switch p.TransactionStatus {
              case INITIALIZE_STATUS :
              case SENDING_REQUEST :
              case WAITE_ANSWER :
              case RECEAVE_ANSWER :
              case SENDING_ANSWER :
              case WAITE_ANS_REQ :
                  p.TransactionStatus = SENDING_ANSWER
                  var sendHeatBeatResponseMessage SendHeatBeatResponse
                  json.Unmarshal(message, &sendHeatBeatResponseMessage)
                  sendHeatBeatResponseMessage.NodeHeaderData.Initialize()
                  sendHeatBeatResponseMessage.NodeHeaderData.SetMsgType(HEATBEAT_RESPONSE)
                  sendHeatBeatResponseMessage.NodeHeaderData.SetSequenceNum(p.SequenceNum)
                  sendHeatBeatResponseMessage.NodeHeaderData.SetInformationElement(&sendHeatBeatResponseMessage.RecoveryTimeData)
                  signalSendReq := new(SendSignalRequest)
                  signalSendReq.MsgType = SIGNAL_SEND_REQUEST
                  signalSendReq.SignalBuffer = sendHeatBeatResponseMessage.NodeHeaderData.CreateSignal()
                  sendMessage,_ := json.Marshal(signalSendReq)
                  p.EndPointChannel <- sendMessage
                  p.TransactionStatus = FINALIZE_STATUS
                  log.Println("トランザクション終了処理")
                  p.SendFinalizeMessage()
                  return
              case RECEAVE_REQUEST :
              case FINALIZE_STATUS :
              default :
              }
        case SIGNAL_RECEVE_NOTIFY:
            var recSignalNotifyMessage RecSignalNotify
            json.Unmarshal(message, &recSignalNotifyMessage)
            //MsgTypeh判定
            messageType := uint8(recSignalNotifyMessage.SignalBuffer[1])
            switch messageType {
            case HEATBEAT_REQUEST:
                switch p.TransactionStatus {
                case INITIALIZE_STATUS :
                    p.TransactionStatus = RECEAVE_REQUEST
                    recHeatBeatRequestMessage := new(RecHeatBeatRequest)
                    nextByte := recHeatBeatRequestMessage.NodeHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                    bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                    for nextByte < bufferlen {
                        ietype := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:(nextByte+2)])
                        switch ietype {
                        case RECOVERY_TIME_STAMP_IE :
                            nextByte = recHeatBeatRequestMessage.RecoveryTimeData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        default:
                        }
                    }
                    log.Println("EndPointへ送信")
                    sendMessage,_ := json.Marshal(recHeatBeatRequestMessage)
                    p.EndPointChannel <- sendMessage
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
            case HEATBEAT_RESPONSE:
                switch p.TransactionStatus {
                case INITIALIZE_STATUS :
                case SENDING_REQUEST :
                case WAITE_ANSWER :
                    p.TransactionStatus = RECEAVE_ANSWER
                    recHeatBeatResponseMessage := new(RecHeatBeatResponse)
                    nextByte := recHeatBeatResponseMessage.NodeHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                    bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                    for nextByte < bufferlen {
                        ietype := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:(nextByte+2)])
                        switch ietype {
                        case RECOVERY_TIME_STAMP_IE :
                            nextByte = recHeatBeatResponseMessage.RecoveryTimeData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                        default:
                        }
                    }
                    log.Println("EndPointへ送信")
                    sendMessage,_ := json.Marshal(recHeatBeatResponseMessage)
                    p.EndPointChannel <- sendMessage
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

        default:
        }
    }
}
