package pfcp

type SendHeatBeatResponse struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecoveryTimeData RecTimeStamp
}
