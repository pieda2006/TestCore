package pfcp

import (
    "net"
    "time"
    "log"
    "encoding/json"
)

var ownNodeInstance *OwnNode

type OwnNode struct {
    PFCPNode
    EndPointChannel map[string](chan []byte)
}

func (p *OwnNode) EndPointInitialize(endPointType uint8) (error) {
    //課題：EndPointの上限数はConfig読み込みに切り替える START
    p.EndPointChannel = make(map[string](chan []byte),10)
    //END
    p.NodeType = endPointType
    return nil
}

func (p *OwnNode) SetTransactionChan(seqNum uint32, transactionchan chan []byte){
    return
}

func (p *OwnNode) ExecuteEndPoint() (error){
    endpointfactoryinstance := GetPFCPEndPointFactoryInstance()
    endpointList := endpointfactoryinstance.GetEndPointList()
    //課題：後で返された分だけEndpointインスタンス生成を行うようにする START
    log.Println("EndPointインスタンス起動")
    for key, value := range endpointList {
        if key != p.IPaddress {
            p.EndPointChannel[value.GetNodeIPaddress()] = value.GetEndPointChan()
            value.SetNodeChan(p.MyChannel)
            go value.ExecuteEndPoint()
        }
    }
    //END

    log.Println("受信用ソケットオープン")
    log.Println("受信用ソケットIPアドレス：",p.IPaddress," 受信ポート：",p.Port)
    address, _ := net.ResolveUDPAddr("udp", p.IPaddress + ":" + p.Port)
    listener, _ := net.ListenUDP("udp", address)
    defer listener.Close()
    log.Println("受信用ソケットIPアドレス：",p.IPaddress," 受信ポート：",p.Port)
    //課題：バッファサイズはConfig可変に変更する。 START
    recevebuff := make([]byte,1500)
    //END

    for {
        log.Println("外部向けソケット刈り取り実施")
        deadline := time.Now()
        //課題：ソケット刈り取り周期はConfigで変更可能にする。 START
        listener.SetReadDeadline(deadline.Add(3 * time.Second))
        //END
        recBufferLen , distIPaddress , _ := listener.ReadFromUDP(recevebuff)
        if recBufferLen != 0 {
            log.Println("受信信号をEndPointインスタンスへ通知")
            signalnotify := new(RecSignalNotify)
            signalnotify.MsgType = SIGNAL_RECEVE_NOTIFY
            signalnotify.SignalBuffer = recevebuff
            message,_ := json.Marshal(signalnotify)
            log.Println("通知先IPaddress：",distIPaddress.IP.String())
            p.EndPointChannel[distIPaddress.IP.String()] <- message

        }
    }
}
