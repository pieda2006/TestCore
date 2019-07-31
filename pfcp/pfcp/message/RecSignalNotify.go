package pfcp

type RecSignalNotify struct {
    InternalMessageBase
    SignalBuffer []byte
}
