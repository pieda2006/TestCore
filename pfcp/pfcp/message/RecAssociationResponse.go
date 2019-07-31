package pfcp

type RecAssociationResponse struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecoveryTimeData RecTimeStamp
    CauseData Cause
}
