package pfcp

import (
    "encoding/json"
    "log"
)

type AssociationReqTransaction struct {
    PFCPTransaction
}

func (p *AssociationReqTransaction) ExecuteTransaction() {
    for {
        message, ok := <-p.MyChannel
        if !ok {
            return
        }

        var messageBase InternalMessageBase
        json.Unmarshal(message, &messageBase)

        switch messageBase.MsgType {
        case SEND_ASSOCIATION_REQUEST:
            switch p.TransactionStatus {
            case INITIALIZE_STATUS :
              log.Println("SEND_ASSOCIATION_REQUEST受信")
              p.TransactionStatus = SENDING_REQUEST
              var sendAssociationSetupRequestMessage SendAssociationSetupRequest
              json.Unmarshal(message, &sendAssociationSetupRequestMessage)
              log.Println("Header生成")
              sendAssociationSetupRequestMessage.NodeHeaderData.Initialize()
              sendAssociationSetupRequestMessage.NodeHeaderData.SetMsgType(ASSOCIATION_SETUP_REQUEST)
              sendAssociationSetupRequestMessage.NodeHeaderData.SetSequenceNum(p.SequenceNum)
              log.Println("RecTimeStamp生成")
              sendAssociationSetupRequestMessage.NodeHeaderData.SetInformationElement(&sendAssociationSetupRequestMessage.RecTimeStampData)
              signalSendReq := new(SendSignalRequest)
              signalSendReq.MsgType = SIGNAL_SEND_REQUEST
              signalSendReq.SignalBuffer = sendAssociationSetupRequestMessage.NodeHeaderData.CreateSignal()
              sendMessage,_ := json.Marshal(signalSendReq)
              log.Println("ASSOCIATION_SETUP_REQUEST送信要求を通知")
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
        case SEND_ASSOCIATION_RESPONSE:
          switch p.TransactionStatus {
          case INITIALIZE_STATUS :
          case SENDING_REQUEST :
          case WAITE_ANSWER :
          case RECEAVE_ANSWER :
          case SENDING_ANSWER :
          case WAITE_ANS_REQ :
              log.Println("SEND_ASSOCIATION_RESPONSE受信")
              p.TransactionStatus = SENDING_ANSWER
              var sendAssociationSetupResponseMessage SendAssociationSetupResponse
              json.Unmarshal(message, &sendAssociationSetupResponseMessage)
              log.Println("Header生成")
              sendAssociationSetupResponseMessage.NodeHeaderData.Initialize()
              sendAssociationSetupResponseMessage.NodeHeaderData.SetMsgType(ASSOCIATION_SETUP_RESPONSE)
              sendAssociationSetupResponseMessage.NodeHeaderData.SetSequenceNum(p.SequenceNum)
              sendAssociationSetupResponseMessage.NodeHeaderData.SetInformationElement(&sendAssociationSetupResponseMessage.RecTimeStampData)
              sendAssociationSetupResponseMessage.NodeHeaderData.SetInformationElement(&sendAssociationSetupResponseMessage.CauseData)
              signalSendReq := new(SendSignalRequest)
              signalSendReq.MsgType = SIGNAL_SEND_REQUEST
              signalSendReq.SignalBuffer = sendAssociationSetupResponseMessage.NodeHeaderData.CreateSignal()
              sendMessage,_ := json.Marshal(signalSendReq)
              log.Println("ASSOCIATION_SETUP_RESPONSE送信要求を通知")
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
            messageType := p.GetMsgType(recSignalNotifyMessage.SignalBuffer)
            log.Println("信号受信 MsgType: ",messageType)
            switch messageType {
            case ASSOCIATION_SETUP_REQUEST:
              switch p.TransactionStatus {
              case INITIALIZE_STATUS :
                log.Println("ASSOCIATION_SETUP_REQUEST受信")
                p.TransactionStatus = RECEAVE_REQUEST
                recAssociationRequest := new(RecAssociationRequest)
                recAssociationRequest.MsgType = REC_ASSOCIATION_REQUEST
                nextByte := recAssociationRequest.NodeHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                for nextByte < bufferlen {
                    log.Println("NextByte: ",nextByte," BufferLen: ",bufferlen)
                    ietype := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:(nextByte+2)])
                    switch ietype {
                    case RECOVERY_TIME_STAMP_IE :
                        nextByte = recAssociationRequest.RecoveryTimeData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                    default:
                    }
                }
                log.Println("ASSOCIATION_SETUP_REQUESTをEndPointへ通知")
                sendMessage,_ := json.Marshal(recAssociationRequest)
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

            case ASSOCIATION_SETUP_RESPONSE:
                  switch p.TransactionStatus {
                  case INITIALIZE_STATUS :
                  case SENDING_REQUEST :
                  case WAITE_ANSWER :
                    log.Println("ASSOCIATION_SETUP_RESPONSE受信")
                      p.TransactionStatus = RECEAVE_ANSWER
                      recAssociationRequest := new(RecAssociationResponse)
                      recAssociationRequest.MsgType = REC_ASSOCIATION_RESPONSE
                      nextByte := recAssociationRequest.NodeHeaderData.InitializeFromBuf(recSignalNotifyMessage.SignalBuffer)
                      bufferlen := uint32(p.GetMsgLength(recSignalNotifyMessage.SignalBuffer) + 4)
                      for nextByte < bufferlen {
                          ietype := p.GetIEType(recSignalNotifyMessage.SignalBuffer[nextByte:(nextByte+2)])
                          switch ietype {
                          case RECOVERY_TIME_STAMP_IE :
                              nextByte = recAssociationRequest.RecoveryTimeData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                          case CAUSE_IE :
                              nextByte = recAssociationRequest.CauseData.IEInitializeFromBuff(recSignalNotifyMessage.SignalBuffer[nextByte:]) + nextByte
                          default:
                          }
                      }
                      log.Println("ASSOCIATION_SETUP_RESPONSEをEndPointへ通知")
                      sendMessage,_ := json.Marshal(recAssociationRequest)
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
