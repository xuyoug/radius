package radius

import (
//"bytes"
)

// 定义radius报文基本属性

//定义radius的CODE 一个字节
type R_Code uint8

//定义常见的radiuscode
const (
	CodeAccessRequest R_Code = iota + 1
	CodeAccessAccept
	CodeAccessReject
	CodeAccountingRequest
	CodeAccountingRespons
	CodeAccessChallenge R_Code = 11
	CodeStatusServer    R_Code = 12
	CodeStatusClient    R_Code = 13
	CodeReserved        R_Code = 255
)

//定义radius的ID
//一个字节
type R_Id uint8

//定义radius的Id最大值
const R_Id_MAX R_Id = 254

//定义radius的Len 两个字节
type R_Length uint16

//定义radiusLength的最大最小值
const (
	R_Length_MIN R_Length = 20
	R_Length_MAX R_Length = 4096
)

//定义authenticator 16字节
type R_Authenticator []byte

//定义Authenticator长度
const (
	R_Authenticator_LEN = 16
)

//定义radius的基础数据结构
type Radius struct {
	R_Code
	R_Id
	R_Length
	R_Authenticator
	AttributeList
}

//func NewRadius()(*Radius,error)
