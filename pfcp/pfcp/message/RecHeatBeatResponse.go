package pfcp

type RecHeatBeatResponse struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecoveryTimeData RecTimeStamp
}
