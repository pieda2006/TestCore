package pfcp

type SendEstablishRequest struct {
    InternalMessageBase
    SessionHeaderData SessionHeader
    NodeIDData NodeID
    FSEIDData FSEID
}
