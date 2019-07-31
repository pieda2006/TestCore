package pfcp

type RecSessionDeleteResponse struct {
    InternalMessageBase
    SessionHeaderData SessionHeader
    CauseData Cause
}
