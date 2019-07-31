package pfcp

type RecSessionEstablishRequest struct {
    InternalMessageBase
    SessionHeaderData SessionHeader
    NodeIDData NodeID
    FSEIDData FSEID
}
