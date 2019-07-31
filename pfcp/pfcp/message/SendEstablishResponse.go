package pfcp

type SendEstablishResponse struct {
    InternalMessageBase
    SessionHeaderData SessionHeader
    NodeIDData NodeID
    FSEIDData FSEID
    CauseData Cause
}
