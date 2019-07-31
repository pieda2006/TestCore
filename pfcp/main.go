package main

import (
    "runtime"
    "./pfcp"
    "log"
    "time"
    "encoding/json"
)

func main() {
    cpus := runtime.NumCPU()
    runtime.GOMAXPROCS(cpus)

    log.SetFlags(log.Ltime|log.Lshortfile)
    log.Println("処理開始")

    endPointFactory := pfcp.GetPFCPEndPointFactoryInstance()

    log.Println("Nodeインスタンス生成")
    endPointFactory.CreateEndPointList()
    endPointList := endPointFactory.GetEndPointList()

    log.Println("自ノード処理起動")
    for key, value := range endPointList {
        log.Println("OWNPORT：", value.GetNodePort())
        if key == pfcp.PFCP_OWN_NODE_IP && value.GetNodePort() == pfcp.PFCP_OWN_NODE_PORT {
            go value.ExecuteEndPoint()
        }
    }

    time.Sleep(5 * time.Second)

    var sessionChan chan []byte = nil

    if pfcp.PFCP_ENDPOINT_MODE == pfcp.ENDPOINT_CLIENT {
        log.Println("セッションへセッション確立要求送信処理開始")
        sendEstablishRequestMessage := new(pfcp.SendEstablishRequest)
        sendEstablishRequestMessage.MsgType = pfcp.SEND_SESSION_ESTABLISHMENT_REQUEST
        pfcpSessionFactoryIns := pfcp.GetPFCPSessionFactoryInstance()
        pfcpSessionIns := pfcpSessionFactoryIns.CreateSession()
        pfcpSessionIns.SetDstIPaddress(pfcp.PFCP_DST_NODE_IP)
        pfcpSessionIns.SetDstPort(pfcp.PFCP_DST_NODE_PORT)
        pfcpSessionIns.SetMyIPaddress(pfcp.PFCP_OWN_NODE_IP)
        pfcpSessionIns.SetMyPort(pfcp.PFCP_OWN_NODE_PORT)
        sessionChan = pfcpSessionIns.GetSessionChan()
        go pfcpSessionIns.ExecuteSession()
        sendMessage,_ := json.Marshal(sendEstablishRequestMessage)
        sessionChan <- sendMessage
        log.Println("セッションへセッション確立要求送信")
    }

    time.Sleep(5 * time.Second)

    if pfcp.PFCP_ENDPOINT_MODE == pfcp.ENDPOINT_CLIENT {
        log.Println("セッションへセッション確立要求送信処理開始")
        sendDeleteRequestMessage := new(pfcp.SendDeleteRequest)
        sendDeleteRequestMessage.MsgType = pfcp.SEND_SESSION_DELETE_REQUEST
        sendMessage,_ := json.Marshal(sendDeleteRequestMessage)
        sessionChan <- sendMessage
        log.Println("セッションへセッションdelete要求送信")
    }

    for {
        time.Sleep(3 * time.Second)
    }

}
