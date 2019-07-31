package pfcp

type RecSessionEstablishResponse struct {
    InternalMessageBase
    SessionHeaderData SessionHeader
    NodeIDData NodeID
    FSEIDData FSEID
    CauseData Cause
}
