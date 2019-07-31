package pfcp

import (
    "log"
)

var EndPointFactoryInstance *PFCPEndPointFactory

func GetPFCPEndPointFactoryInstance() *PFCPEndPointFactory {
    if EndPointFactoryInstance == nil {
        EndPointFactoryInstance = new(PFCPEndPointFactory)
        EndPointFactoryInstance.EndPointFactoryInitialize()
    }
    return EndPointFactoryInstance
}


type PFCPEndPointFactory struct {
    EndPointList map[string](PFCPNodeIF)
}

func (p *PFCPEndPointFactory) EndPointFactoryInitialize() (error) {
    //課題：最大ノード数はConfig値で保持するようにする START
    p.EndPointList = make(map[string](PFCPNodeIF),PFCP_MAX_ENDPOINT_NODE)
    //END
    return nil

}
func (p *PFCPEndPointFactory) CreateEndPointList() (error) {

    log.Println("自サーバインスタンス生成")
    pfcpownnode := new(OwnNode)
    pfcpownnode.EndPointInitialize(ENDPOINT_OWN)
    pfcpownnode.SetNodeIPaddress(PFCP_OWN_NODE_IP)
    pfcpownnode.SetNodePort(PFCP_OWN_NODE_PORT)
    p.EndPointList[PFCP_OWN_NODE_IP] = pfcpownnode

    log.Println("EndPointインスタンス生成")
    pfcpendpointnode := new(PFCPEndPoint)
    pfcpendpointnode.EndPointInitialize(PFCP_ENDPOINT_MODE)
    pfcpendpointnode.SetNodeIPaddress(PFCP_DST_NODE_IP)
    pfcpendpointnode.SetNodePort(PFCP_DST_NODE_PORT)
    pfcpendpointnode.SetMyIPaddress(PFCP_OWN_NODE_IP)
    pfcpendpointnode.SetMyPort(PFCP_OWN_NODE_PORT)
    p.EndPointList[PFCP_DST_NODE_IP] = pfcpendpointnode

    return nil
}

func (p *PFCPEndPointFactory) GetEndPointIF(ipaddress string) PFCPNodeIF{
    return p.EndPointList[ipaddress]
}

func (p *PFCPEndPointFactory) GetEndPointList() map[string](PFCPNodeIF){

    return p.EndPointList
}
