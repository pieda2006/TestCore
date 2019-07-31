package pfcp

type RecHeatBeatRequest struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecoveryTimeData RecTimeStamp
}
