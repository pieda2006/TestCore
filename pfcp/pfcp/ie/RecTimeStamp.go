package pfcp

import (
    "encoding/binary"
)

type RecTimeStamp struct {
    InformationElement
    RecoveryTime uint32
}

func (p *RecTimeStamp) CreateSignal(signal []byte) {
    binary.BigEndian.PutUint16(signal[0:2],p.IEType)
    binary.BigEndian.PutUint16(signal[2:4],p.IELength)
    binary.BigEndian.PutUint32(signal[4:8],p.RecoveryTime)
    p.CreateNextIE(signal[8:])
}

func (p *RecTimeStamp) SetTimeStamp(recTime uint32) {
    p.RecoveryTime = recTime
}

func (p *RecTimeStamp) IEInitialize() (error){
    p.IEType = RECOVERY_TIME_STAMP_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = 4
    p.MessageLength = 0
    return nil
}

func (p *RecTimeStamp) IEInitializeFromBuff(buffer []byte) (uint32){
    p.IEType = RECOVERY_TIME_STAMP_IE
    p.MemberElement = nil
    p.NextElement = nil
    p.IELength = binary.BigEndian.Uint16(buffer[2:4])
    p.RecoveryTime = binary.BigEndian.Uint32(buffer[4:8])
    p.MessageLength = 0
    return 8
}

func (p *RecTimeStamp) GetTimeStamp() uint32 {
    return p.RecoveryTime
}
