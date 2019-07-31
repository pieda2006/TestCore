package pfcp

import (
    "log"
)

var SessionFactoryInstance *PFCPSessionFactory

func GetPFCPSessionFactoryInstance() *PFCPSessionFactory {
    if SessionFactoryInstance == nil {
        log.Println("PFCPSessionFactoryインスタンス生成")
        SessionFactoryInstance = new(PFCPSessionFactory)
        SessionFactoryInstance.SessionFactoryInitialize()
    }
    return SessionFactoryInstance
}

type PFCPSessionFactory struct {
    SeIdNum uint64
    PFCPSessionMap map[uint64](*PFCPSession)
}

func (p *PFCPSessionFactory) SessionFactoryInitialize () (error){
    p.SeIdNum = 1
    p.PFCPSessionMap = make(map[uint64](*PFCPSession),100)
    return nil
}

func (p *PFCPSessionFactory) CreateSession() *PFCPSession {
    pfcpsession := new(PFCPSession)
    pfcpsession.SessionInitialize()
    pfcpsession.SetMySEID(p.SeIdNum)
    p.PFCPSessionMap[p.SeIdNum] = pfcpsession
    p.SeIdNum++
    return pfcpsession
}

func (p *PFCPSessionFactory) GetSession(seid uint64) *PFCPSession {
  return p.PFCPSessionMap[seid]
}
