package pfcp

type SendAssociationSetupResponse struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecTimeStampData RecTimeStamp
    CauseData Cause
}
