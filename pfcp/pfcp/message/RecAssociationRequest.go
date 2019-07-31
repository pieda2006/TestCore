package pfcp

type RecAssociationRequest struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecoveryTimeData RecTimeStamp
}
