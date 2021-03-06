// Code generated by protoc-gen-go. DO NOT EDIT.
// source: yandex/cloud/iam/v1/service_account.proto

package iam

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/yandex-cloud/go-genproto/yandex/cloud/validation"
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

// A ServiceAccount resource. For more information, see [Service accounts](/docs/iam/concepts/users/service-accounts).
type ServiceAccount struct {
	// ID of the service account.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// ID of the folder that the service account belongs to.
	FolderId string `protobuf:"bytes,2,opt,name=folder_id,json=folderId,proto3" json:"folder_id,omitempty"`
	// Creation timestamp.
	CreatedAt *timestamp.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// Name of the service account.
	// The name is unique within the cloud. 3-63 characters long.
	Name string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	// Description of the service account. 0-256 characters long.
	Description          string   `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ServiceAccount) Reset()         { *m = ServiceAccount{} }
func (m *ServiceAccount) String() string { return proto.CompactTextString(m) }
func (*ServiceAccount) ProtoMessage()    {}
func (*ServiceAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_053d0ddb735dcde2, []int{0}
}

func (m *ServiceAccount) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceAccount.Unmarshal(m, b)
}
func (m *ServiceAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceAccount.Marshal(b, m, deterministic)
}
func (m *ServiceAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceAccount.Merge(m, src)
}
func (m *ServiceAccount) XXX_Size() int {
	return xxx_messageInfo_ServiceAccount.Size(m)
}
func (m *ServiceAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceAccount.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceAccount proto.InternalMessageInfo

func (m *ServiceAccount) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ServiceAccount) GetFolderId() string {
	if m != nil {
		return m.FolderId
	}
	return ""
}

func (m *ServiceAccount) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *ServiceAccount) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ServiceAccount) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func init() {
	proto.RegisterType((*ServiceAccount)(nil), "yandex.cloud.iam.v1.ServiceAccount")
}

func init() {
	proto.RegisterFile("yandex/cloud/iam/v1/service_account.proto", fileDescriptor_053d0ddb735dcde2)
}

var fileDescriptor_053d0ddb735dcde2 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xc1, 0x4b, 0xc3, 0x30,
	0x18, 0xc5, 0x69, 0x9d, 0x62, 0x33, 0xd8, 0x21, 0x5e, 0x4a, 0x45, 0x2c, 0x9e, 0xe6, 0x61, 0x09,
	0xd3, 0x93, 0x0c, 0x0f, 0xf3, 0xe6, 0x75, 0x7a, 0xf2, 0x52, 0xbe, 0x26, 0xdf, 0xea, 0x07, 0x4d,
	0x53, 0xd2, 0xb4, 0xe8, 0x3f, 0xe5, 0xdf, 0x28, 0x26, 0x0e, 0x14, 0x76, 0x0b, 0xef, 0xbd, 0xe4,
	0xfd, 0x5e, 0xd8, 0xed, 0x27, 0x74, 0x1a, 0x3f, 0xa4, 0x6a, 0xed, 0xa8, 0x25, 0x81, 0x91, 0xd3,
	0x5a, 0x0e, 0xe8, 0x26, 0x52, 0x58, 0x81, 0x52, 0x76, 0xec, 0xbc, 0xe8, 0x9d, 0xf5, 0x96, 0x5f,
	0xc4, 0xa8, 0x08, 0x51, 0x41, 0x60, 0xc4, 0xb4, 0x2e, 0xae, 0x1b, 0x6b, 0x9b, 0x16, 0x65, 0x88,
	0xd4, 0xe3, 0x5e, 0x7a, 0x32, 0x38, 0x78, 0x30, 0x7d, 0xbc, 0x55, 0x5c, 0xfd, 0x2b, 0x98, 0xa0,
	0x25, 0x0d, 0x9e, 0x6c, 0x17, 0xed, 0x9b, 0xaf, 0x84, 0x2d, 0x5e, 0x62, 0xdd, 0x36, 0xb6, 0xf1,
	0x05, 0x4b, 0x49, 0xe7, 0x49, 0x99, 0x2c, 0xb3, 0x5d, 0x4a, 0x9a, 0x5f, 0xb2, 0x6c, 0x6f, 0x5b,
	0x8d, 0xae, 0x22, 0x9d, 0xa7, 0x41, 0x3e, 0x8f, 0xc2, 0xb3, 0xe6, 0x0f, 0x8c, 0x29, 0x87, 0xe0,
	0x51, 0x57, 0xe0, 0xf3, 0x93, 0x32, 0x59, 0xce, 0xef, 0x0a, 0x11, 0xa1, 0xc4, 0x01, 0x4a, 0xbc,
	0x1e, 0xa0, 0x76, 0xd9, 0x6f, 0x7a, 0xeb, 0x39, 0x67, 0xb3, 0x0e, 0x0c, 0xe6, 0xb3, 0xf0, 0x64,
	0x38, 0xf3, 0x92, 0xcd, 0x35, 0x0e, 0xca, 0x51, 0xff, 0xc3, 0x98, 0x9f, 0x06, 0xeb, 0xaf, 0xf4,
	0xf4, 0xf8, 0xb6, 0x69, 0xc8, 0xbf, 0x8f, 0xb5, 0x50, 0xd6, 0xc8, 0x38, 0x6e, 0x15, 0xc7, 0x35,
	0x76, 0xd5, 0x60, 0x17, 0x4a, 0xe5, 0x91, 0x6f, 0xdd, 0x10, 0x98, 0xfa, 0x2c, 0xd8, 0xf7, 0xdf,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x65, 0x42, 0x1f, 0xbf, 0x78, 0x01, 0x00, 0x00,
}
