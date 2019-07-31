package pfcp

type InformationElementIF interface {
    GetLength() (uint16)
    CreateSignal(signal []byte)
    SetMemberElement(memberElement InformationElementIF)
    SetNextElement(nextElement InformationElementIF)
    IEInitialize() (error)
}
