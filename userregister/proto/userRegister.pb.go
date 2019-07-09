// Code generated by protoc-gen-go. DO NOT EDIT.
// source: userRegister.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Request struct {
	MobileNum            string   `protobuf:"bytes,1,opt,name=mobileNum,proto3" json:"mobileNum,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Password             string   `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd0216664a6472cd, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetMobileNum() string {
	if m != nil {
		return m.MobileNum
	}
	return ""
}

func (m *Request) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Request) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Response struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd0216664a6472cd, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "go.micro.api.register.Request")
	proto.RegisterType((*Response)(nil), "go.micro.api.register.Response")
}

func init() { proto.RegisterFile("userRegister.proto", fileDescriptor_dd0216664a6472cd) }

var fileDescriptor_dd0216664a6472cd = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x8f, 0xc1, 0x0a, 0x82, 0x40,
	0x10, 0x86, 0x33, 0xa3, 0x74, 0xe8, 0x10, 0x03, 0x81, 0x88, 0x54, 0x78, 0xea, 0xb4, 0x87, 0x7a,
	0x8f, 0xa0, 0x85, 0xe8, 0xd4, 0x41, 0x6b, 0x90, 0x85, 0xd6, 0xdd, 0x76, 0x94, 0x5e, 0x3f, 0x5c,
	0xad, 0x2e, 0x75, 0xfb, 0xe7, 0xff, 0x60, 0xe6, 0x1b, 0xc0, 0x96, 0xc9, 0x49, 0xaa, 0x14, 0x37,
	0xe4, 0x84, 0x75, 0xa6, 0x31, 0xb8, 0xac, 0x8c, 0xd0, 0xea, 0xea, 0x8c, 0x28, 0xac, 0x12, 0x6e,
	0x80, 0xf9, 0x19, 0x66, 0x92, 0x1e, 0x2d, 0x71, 0x83, 0x19, 0xc4, 0xda, 0x94, 0xea, 0x4e, 0x87,
	0x56, 0x27, 0xc1, 0x26, 0xd8, 0xc6, 0xf2, 0x5b, 0x20, 0xc2, 0xa4, 0x2e, 0x34, 0x25, 0x63, 0x0f,
	0x7c, 0xc6, 0x14, 0x22, 0x5b, 0x30, 0x3f, 0x8d, 0xbb, 0x25, 0xa1, 0xef, 0x3f, 0x73, 0x9e, 0x41,
	0x24, 0x89, 0xad, 0xa9, 0x99, 0x70, 0x01, 0xa1, 0xe6, 0x6a, 0xd8, 0xd9, 0xc5, 0xdd, 0xa5, 0xa3,
	0xbd, 0x02, 0x1e, 0x61, 0xfe, 0xce, 0x27, 0x26, 0x87, 0x2b, 0xf1, 0x53, 0x55, 0x0c, 0x9e, 0xe9,
	0xfa, 0x2f, 0xef, 0xcf, 0xe5, 0xa3, 0x72, 0xea, 0x7f, 0xde, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x2b, 0x64, 0xca, 0x1e, 0x09, 0x01, 0x00, 0x00,
}
