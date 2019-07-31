package pfcp

type SendHeatBeatRequest struct {
    InternalMessageBase
    NodeHeaderData NodeHeader
    RecTimeStampData RecTimeStamp
}
