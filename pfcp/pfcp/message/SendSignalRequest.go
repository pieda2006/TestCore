package pfcp

type SendSignalRequest struct {
    InternalMessageBase
    SignalBuffer []byte
}
