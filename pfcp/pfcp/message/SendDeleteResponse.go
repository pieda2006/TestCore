package pfcp

type SendDeleteResponse struct {
    InternalMessageBase
    SessionHeaderData SessionHeader
    CauseData Cause
}
