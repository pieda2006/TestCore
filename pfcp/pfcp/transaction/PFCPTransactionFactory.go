package pfcp

import (
    "log"
)

var TransactionFactoryInstance *PFCPTransactionFactory

func GetPFCPTransactionFactoryInstance() *PFCPTransactionFactory {
    if TransactionFactoryInstance == nil {
        log.Println("PFCPTransactionFactoryインスタンス生成")
        TransactionFactoryInstance = new(PFCPTransactionFactory)
        TransactionFactoryInstance.TransactionFactoryInitialize()
    }
    return TransactionFactoryInstance
}

type PFCPTransactionFactory struct {
    SequenceNum uint32
}

func (p *PFCPTransactionFactory) TransactionFactoryInitialize () (error){
    p.SequenceNum = 1
    return nil
}

func (p *PFCPTransactionFactory) ConvertTransactionType(msgType uint32) uint32 {
    switch msgType {
    case HEATBEAT_REQUEST:
        return HEATBEAT_TRANSACTION
    case HEATBEAT_RESPONSE:
        return HEATBEAT_TRANSACTION
    case ASSOCIATION_SETUP_REQUEST:
        return ASSOCIATION_SETUP_TRANSACTION
    case ASSOCIATION_SETUP_RESPONSE:
        return ASSOCIATION_SETUP_TRANSACTION
    case SESSION_ESTABLISHMENT_REQUEST:
        return SESSION_ESTABLISHMENT_TRANSACTION
    case SESSION_ESTABLISHMENT_RESPONSE:
        return SESSION_ESTABLISHMENT_TRANSACTION
    case SESSION_DELETE_REQUEST:
        return SESSION_DELETE_TRANSACTION
    case SESSION_DELETE_RESPONSE:
        return SESSION_DELETE_TRANSACTION
    default:
        return 0
    }
}

func (p *PFCPTransactionFactory) CreateTransaction(msgType uint32, seqNum uint32) PFCPTransactionIF {
    var transactionIF PFCPTransactionIF
    switch p.ConvertTransactionType(msgType) {
    case HEATBEAT_TRANSACTION:
        log.Println("HeatBeatトランザクション生成")
        transaction := new(HeatBeatTransaction)
        transaction.TransactionInitialize()
        transactionIF = transaction
    case ASSOCIATION_SETUP_TRANSACTION:
        log.Println("AssociationReqトランザクション生成")
        transaction := new(AssociationReqTransaction)
        transaction.TransactionInitialize()
        transactionIF = transaction
    case SESSION_ESTABLISHMENT_TRANSACTION:
        log.Println("SessionEstablishのトランザクション生成")
        transaction := new(SessionEstablishTransaction)
        transaction.TransactionInitialize()
        transactionIF = transaction
    case SESSION_DELETE_TRANSACTION:
        log.Println("SessionDeleteのトランザクション生成")
        transaction := new(SessionDeleteTransaction)
        transaction.TransactionInitialize()
        transactionIF = transaction
    default:
        return nil
    }
    if seqNum == 0 {
        transactionIF.SetSeqNum(p.SequenceNum)
        p.SequenceNum++
    } else {
        transactionIF.SetSeqNum(seqNum)
    }

    return transactionIF

}
