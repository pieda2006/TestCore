package pfcp

type SendAssociationSetupRequest struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecTimeStampData RecTimeStamp
}
