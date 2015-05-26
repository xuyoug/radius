package radius

import (
//"bytes"
)

// 定义radius报文基本属性

//定义radius的CODE 一个字节
type R_Code uint8

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

//定义radius的ID 一个字节
type R_Id uint8

const R_Id_MAX R_Id = 254

//定义radius的Len 两个字节
type R_Length uint16

const (
	radiusLength_MIN R_Length = 20
	radiusLength_MAX R_Length = 4096
)

//定义authenticator 16字节
type R_Authenticator []byte

const (
	Radius_Authenticator_LEN = 16
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
