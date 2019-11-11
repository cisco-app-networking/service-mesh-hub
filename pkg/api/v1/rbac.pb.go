// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/mesh-projects/api/v1/rbac.proto

package v1

import (
	bytes "bytes"
	fmt "fmt"
	math "math"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/types"
	_ "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// if RBAC policies have been specified, service isolation will not work as expected.
// Istio, the model for this feature, has an exclude-by-default behavior. If RBAC is enabled, services can only
// do what the RBAC policy allows them to do. This means that in the special case of
type RbacStatusCode int32

const (
	// Initial value, will be replaced as soon as operator evaluates the config
	RbacStatusCode_PENDING_VERIFICATION RbacStatusCode = 0
	// Config was applied successfully
	RbacStatusCode_OK RbacStatusCode = 1
	// If a mesh does not support the provided configuration, this error code is returned.
	RbacStatusCode_ERROR_RBAC_MODE_NOT_SUPPORTED_BY_MESH RbacStatusCode = 2
	// If other policies exist, we lose the ability to have total ON vs. OFF control of RBAC. Isolation is not supported
	// in this case.
	RbacStatusCode_ERROR_CANNOT_ISOLATE_RBAC_SINCE_POLICIES_EXIST RbacStatusCode = 3
	// If the config is not accepted for any other reason, this code is returned
	RbacStatusCode_ERROR_CONFIG_NOT_ACCEPTED RbacStatusCode = 4
)

var RbacStatusCode_name = map[int32]string{
	0: "PENDING_VERIFICATION",
	1: "OK",
	2: "ERROR_RBAC_MODE_NOT_SUPPORTED_BY_MESH",
	3: "ERROR_CANNOT_ISOLATE_RBAC_SINCE_POLICIES_EXIST",
	4: "ERROR_CONFIG_NOT_ACCEPTED",
}

var RbacStatusCode_value = map[string]int32{
	"PENDING_VERIFICATION":                  0,
	"OK":                                    1,
	"ERROR_RBAC_MODE_NOT_SUPPORTED_BY_MESH": 2,
	"ERROR_CANNOT_ISOLATE_RBAC_SINCE_POLICIES_EXIST": 3,
	"ERROR_CONFIG_NOT_ACCEPTED":                      4,
}

func (x RbacStatusCode) String() string {
	return proto.EnumName(RbacStatusCode_name, int32(x))
}

func (RbacStatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_22c76505334eb8ab, []int{0}
}

// Configure RBAC properties on the mesh
type RbacMode struct {
	// Mode describes the desired RBAC behavior an optionally takes mode-specific configuration
	// implementation note: using oneof instead of enum since future modes may accept config
	//
	// Types that are valid to be assigned to Mode:
	//	*RbacMode_Unspecified_
	//	*RbacMode_Disable_
	//	*RbacMode_Enable_
	Mode isRbacMode_Mode `protobuf_oneof:"mode"`
	// Set by operator
	// - Initialized as pending.
	// - If isolation cannot be expressed, an error code corresponding to the reason is reported.
	// - If isolation can be expressed, an "OK" status code is reported.
	Status               *RbacStatus `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *RbacMode) Reset()         { *m = RbacMode{} }
func (m *RbacMode) String() string { return proto.CompactTextString(m) }
func (*RbacMode) ProtoMessage()    {}
func (*RbacMode) Descriptor() ([]byte, []int) {
	return fileDescriptor_22c76505334eb8ab, []int{0}
}
func (m *RbacMode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacMode.Unmarshal(m, b)
}
func (m *RbacMode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacMode.Marshal(b, m, deterministic)
}
func (m *RbacMode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacMode.Merge(m, src)
}
func (m *RbacMode) XXX_Size() int {
	return xxx_messageInfo_RbacMode.Size(m)
}
func (m *RbacMode) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacMode.DiscardUnknown(m)
}

var xxx_messageInfo_RbacMode proto.InternalMessageInfo

type isRbacMode_Mode interface {
	isRbacMode_Mode()
	Equal(interface{}) bool
}

type RbacMode_Unspecified_ struct {
	Unspecified *RbacMode_Unspecified `protobuf:"bytes,1,opt,name=unspecified,proto3,oneof" json:"unspecified,omitempty"`
}
type RbacMode_Disable_ struct {
	Disable *RbacMode_Disable `protobuf:"bytes,2,opt,name=disable,proto3,oneof" json:"disable,omitempty"`
}
type RbacMode_Enable_ struct {
	Enable *RbacMode_Enable `protobuf:"bytes,3,opt,name=enable,proto3,oneof" json:"enable,omitempty"`
}

func (*RbacMode_Unspecified_) isRbacMode_Mode() {}
func (*RbacMode_Disable_) isRbacMode_Mode()     {}
func (*RbacMode_Enable_) isRbacMode_Mode()      {}

func (m *RbacMode) GetMode() isRbacMode_Mode {
	if m != nil {
		return m.Mode
	}
	return nil
}

func (m *RbacMode) GetUnspecified() *RbacMode_Unspecified {
	if x, ok := m.GetMode().(*RbacMode_Unspecified_); ok {
		return x.Unspecified
	}
	return nil
}

func (m *RbacMode) GetDisable() *RbacMode_Disable {
	if x, ok := m.GetMode().(*RbacMode_Disable_); ok {
		return x.Disable
	}
	return nil
}

func (m *RbacMode) GetEnable() *RbacMode_Enable {
	if x, ok := m.GetMode().(*RbacMode_Enable_); ok {
		return x.Enable
	}
	return nil
}

func (m *RbacMode) GetStatus() *RbacStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*RbacMode) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*RbacMode_Unspecified_)(nil),
		(*RbacMode_Disable_)(nil),
		(*RbacMode_Enable_)(nil),
	}
}

// Unspecified is the default RBAC policy mode
// If a particular mesh does not support RBAC policy, this is the only allowed mode.
// Compatibility: [all]
type RbacMode_Unspecified struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RbacMode_Unspecified) Reset()         { *m = RbacMode_Unspecified{} }
func (m *RbacMode_Unspecified) String() string { return proto.CompactTextString(m) }
func (*RbacMode_Unspecified) ProtoMessage()    {}
func (*RbacMode_Unspecified) Descriptor() ([]byte, []int) {
	return fileDescriptor_22c76505334eb8ab, []int{0, 0}
}
func (m *RbacMode_Unspecified) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacMode_Unspecified.Unmarshal(m, b)
}
func (m *RbacMode_Unspecified) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacMode_Unspecified.Marshal(b, m, deterministic)
}
func (m *RbacMode_Unspecified) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacMode_Unspecified.Merge(m, src)
}
func (m *RbacMode_Unspecified) XXX_Size() int {
	return xxx_messageInfo_RbacMode_Unspecified.Size(m)
}
func (m *RbacMode_Unspecified) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacMode_Unspecified.DiscardUnknown(m)
}

var xxx_messageInfo_RbacMode_Unspecified proto.InternalMessageInfo

// Disable indicates that RBAC policies should not be enforced
// Compatibility: [only: istio]
type RbacMode_Disable struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RbacMode_Disable) Reset()         { *m = RbacMode_Disable{} }
func (m *RbacMode_Disable) String() string { return proto.CompactTextString(m) }
func (*RbacMode_Disable) ProtoMessage()    {}
func (*RbacMode_Disable) Descriptor() ([]byte, []int) {
	return fileDescriptor_22c76505334eb8ab, []int{0, 1}
}
func (m *RbacMode_Disable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacMode_Disable.Unmarshal(m, b)
}
func (m *RbacMode_Disable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacMode_Disable.Marshal(b, m, deterministic)
}
func (m *RbacMode_Disable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacMode_Disable.Merge(m, src)
}
func (m *RbacMode_Disable) XXX_Size() int {
	return xxx_messageInfo_RbacMode_Disable.Size(m)
}
func (m *RbacMode_Disable) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacMode_Disable.DiscardUnknown(m)
}

var xxx_messageInfo_RbacMode_Disable proto.InternalMessageInfo

// Enable mode tells the mesh to evaluate any policies that are defined
// Compatibility: [only: istio]
type RbacMode_Enable struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RbacMode_Enable) Reset()         { *m = RbacMode_Enable{} }
func (m *RbacMode_Enable) String() string { return proto.CompactTextString(m) }
func (*RbacMode_Enable) ProtoMessage()    {}
func (*RbacMode_Enable) Descriptor() ([]byte, []int) {
	return fileDescriptor_22c76505334eb8ab, []int{0, 2}
}
func (m *RbacMode_Enable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacMode_Enable.Unmarshal(m, b)
}
func (m *RbacMode_Enable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacMode_Enable.Marshal(b, m, deterministic)
}
func (m *RbacMode_Enable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacMode_Enable.Merge(m, src)
}
func (m *RbacMode_Enable) XXX_Size() int {
	return xxx_messageInfo_RbacMode_Enable.Size(m)
}
func (m *RbacMode_Enable) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacMode_Enable.DiscardUnknown(m)
}

var xxx_messageInfo_RbacMode_Enable proto.InternalMessageInfo

type RbacStatus struct {
	// Status code summarizing Rbac Config acceptance state
	Code RbacStatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=zephyr.solo.io.RbacStatusCode" json:"code,omitempty"`
	// As needed according to the status code, this message will surface any relevant configuration details or issues.
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RbacStatus) Reset()         { *m = RbacStatus{} }
func (m *RbacStatus) String() string { return proto.CompactTextString(m) }
func (*RbacStatus) ProtoMessage()    {}
func (*RbacStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_22c76505334eb8ab, []int{1}
}
func (m *RbacStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RbacStatus.Unmarshal(m, b)
}
func (m *RbacStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RbacStatus.Marshal(b, m, deterministic)
}
func (m *RbacStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RbacStatus.Merge(m, src)
}
func (m *RbacStatus) XXX_Size() int {
	return xxx_messageInfo_RbacStatus.Size(m)
}
func (m *RbacStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_RbacStatus.DiscardUnknown(m)
}

var xxx_messageInfo_RbacStatus proto.InternalMessageInfo

func (m *RbacStatus) GetCode() RbacStatusCode {
	if m != nil {
		return m.Code
	}
	return RbacStatusCode_PENDING_VERIFICATION
}

func (m *RbacStatus) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterEnum("zephyr.solo.io.RbacStatusCode", RbacStatusCode_name, RbacStatusCode_value)
	proto.RegisterType((*RbacMode)(nil), "zephyr.solo.io.RbacMode")
	proto.RegisterType((*RbacMode_Unspecified)(nil), "zephyr.solo.io.RbacMode.Unspecified")
	proto.RegisterType((*RbacMode_Disable)(nil), "zephyr.solo.io.RbacMode.Disable")
	proto.RegisterType((*RbacMode_Enable)(nil), "zephyr.solo.io.RbacMode.Enable")
	proto.RegisterType((*RbacStatus)(nil), "zephyr.solo.io.RbacStatus")
}

func init() {
	proto.RegisterFile("github.com/solo-io/mesh-projects/api/v1/rbac.proto", fileDescriptor_22c76505334eb8ab)
}

var fileDescriptor_22c76505334eb8ab = []byte{
	// 495 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0xdd, 0x6a, 0xdb, 0x30,
	0x14, 0xc7, 0xe3, 0xd4, 0x38, 0xed, 0x09, 0x0b, 0x41, 0xf4, 0x22, 0x33, 0x2c, 0x2b, 0x65, 0x83,
	0x7d, 0x50, 0x9b, 0x64, 0x37, 0x1b, 0xec, 0x26, 0x71, 0xd4, 0x46, 0xac, 0xb1, 0x83, 0xec, 0x8e,
	0xad, 0x37, 0xc6, 0x1f, 0xaa, 0xe3, 0x35, 0x89, 0x8c, 0xe5, 0x6c, 0x6c, 0x57, 0x7b, 0x9c, 0x5d,
	0x0e, 0xf6, 0x46, 0x7b, 0x83, 0xbd, 0xc1, 0x88, 0x6c, 0xd3, 0x16, 0x1a, 0x7a, 0x65, 0x1d, 0xeb,
	0xf7, 0xfb, 0x4b, 0x3e, 0x1c, 0xc3, 0x30, 0x49, 0x8b, 0xc5, 0x26, 0x34, 0x22, 0xbe, 0x32, 0x05,
	0x5f, 0xf2, 0x93, 0x94, 0x9b, 0x2b, 0x26, 0x16, 0x27, 0x59, 0xce, 0xbf, 0xb0, 0xa8, 0x10, 0x66,
	0x90, 0xa5, 0xe6, 0xd7, 0x81, 0x99, 0x87, 0x41, 0x64, 0x64, 0x39, 0x2f, 0x38, 0xea, 0xfc, 0x60,
	0xd9, 0xe2, 0x7b, 0x6e, 0x6c, 0x79, 0x23, 0xe5, 0xfa, 0xe0, 0x9e, 0x0c, 0xf9, 0xbc, 0x4e, 0x8b,
	0x5a, 0xaf, 0xeb, 0x32, 0x42, 0x3f, 0x4c, 0x78, 0xc2, 0xe5, 0xd2, 0xdc, 0xae, 0xaa, 0xb7, 0xfd,
	0x84, 0xf3, 0x64, 0xc9, 0x4c, 0x59, 0x85, 0x9b, 0x2b, 0x33, 0xde, 0xe4, 0x41, 0x91, 0xf2, 0xf5,
	0xae, 0xfd, 0x6f, 0x79, 0x90, 0x65, 0x2c, 0x17, 0xe5, 0xfe, 0xf1, 0x9f, 0x26, 0xec, 0xd3, 0x30,
	0x88, 0x66, 0x3c, 0x66, 0x68, 0x0a, 0xed, 0xcd, 0x5a, 0x64, 0x2c, 0x4a, 0xaf, 0x52, 0x16, 0xf7,
	0x94, 0x23, 0xe5, 0x45, 0x7b, 0xf8, 0xcc, 0xb8, 0x7b, 0x77, 0xa3, 0xc6, 0x8d, 0x8b, 0x1b, 0x76,
	0xda, 0xa0, 0xb7, 0x55, 0xf4, 0x1e, 0x5a, 0x71, 0x2a, 0x82, 0x70, 0xc9, 0x7a, 0x4d, 0x99, 0x72,
	0xb4, 0x33, 0x65, 0x52, 0x72, 0xd3, 0x06, 0xad, 0x15, 0xf4, 0x0e, 0x34, 0xb6, 0x96, 0xf2, 0x9e,
	0x94, 0x9f, 0xee, 0x94, 0xf1, 0xba, 0x72, 0x2b, 0x01, 0xbd, 0x05, 0x4d, 0x14, 0x41, 0xb1, 0x11,
	0x3d, 0x55, 0xaa, 0xfa, 0x7d, 0xaa, 0x2b, 0x89, 0xb1, 0xfa, 0xf3, 0x9f, 0xaa, 0xd0, 0x8a, 0xd7,
	0x1f, 0x41, 0xfb, 0xd6, 0x07, 0xe9, 0x07, 0xd0, 0xaa, 0x6e, 0xa6, 0xef, 0x83, 0x56, 0x9e, 0x33,
	0xd6, 0x40, 0x5d, 0xf1, 0x98, 0x1d, 0x5f, 0x02, 0xdc, 0xe4, 0xa0, 0x21, 0xa8, 0x11, 0x8f, 0x99,
	0xec, 0x57, 0x67, 0xd8, 0xdf, 0x7d, 0xa2, 0xc5, 0x63, 0x46, 0x25, 0x8b, 0x7a, 0xd0, 0x5a, 0x31,
	0x21, 0x82, 0xa4, 0x6c, 0xd0, 0x01, 0xad, 0xcb, 0x57, 0xbf, 0x15, 0xe8, 0xdc, 0x55, 0x50, 0x0f,
	0x0e, 0xe7, 0xd8, 0x9e, 0x10, 0xfb, 0xcc, 0xff, 0x88, 0x29, 0x39, 0x25, 0xd6, 0xc8, 0x23, 0x8e,
	0xdd, 0x6d, 0x20, 0x0d, 0x9a, 0xce, 0x87, 0xae, 0x82, 0x5e, 0xc2, 0x73, 0x4c, 0xa9, 0x43, 0x7d,
	0x3a, 0x1e, 0x59, 0xfe, 0xcc, 0x99, 0x60, 0xdf, 0x76, 0x3c, 0xdf, 0xbd, 0x98, 0xcf, 0x1d, 0xea,
	0xe1, 0x89, 0x3f, 0xfe, 0xec, 0xcf, 0xb0, 0x3b, 0xed, 0x36, 0xd1, 0x10, 0x8c, 0x12, 0xb5, 0x46,
	0xf6, 0x96, 0x21, 0xae, 0x73, 0x3e, 0xf2, 0x70, 0x69, 0xba, 0xc4, 0xb6, 0xb0, 0x3f, 0x77, 0xce,
	0x89, 0x45, 0xb0, 0xeb, 0xe3, 0x4f, 0xc4, 0xf5, 0xba, 0x7b, 0xe8, 0x09, 0x3c, 0xae, 0x1c, 0xc7,
	0x3e, 0x25, 0x67, 0x32, 0x7b, 0x64, 0x59, 0x78, 0xee, 0xe1, 0x49, 0x57, 0x1d, 0x0f, 0x7e, 0xfd,
	0xed, 0x2b, 0x97, 0xaf, 0x1f, 0xfc, 0x2f, 0xb2, 0xeb, 0xa4, 0x1a, 0xee, 0x50, 0x93, 0xe3, 0xf7,
	0xe6, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf6, 0x42, 0xbb, 0xe6, 0x4d, 0x03, 0x00, 0x00,
}

func (this *RbacMode) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode)
	if !ok {
		that2, ok := that.(RbacMode)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if that1.Mode == nil {
		if this.Mode != nil {
			return false
		}
	} else if this.Mode == nil {
		return false
	} else if !this.Mode.Equal(that1.Mode) {
		return false
	}
	if !this.Status.Equal(that1.Status) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *RbacMode_Unspecified_) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode_Unspecified_)
	if !ok {
		that2, ok := that.(RbacMode_Unspecified_)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Unspecified.Equal(that1.Unspecified) {
		return false
	}
	return true
}
func (this *RbacMode_Disable_) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode_Disable_)
	if !ok {
		that2, ok := that.(RbacMode_Disable_)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Disable.Equal(that1.Disable) {
		return false
	}
	return true
}
func (this *RbacMode_Enable_) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode_Enable_)
	if !ok {
		that2, ok := that.(RbacMode_Enable_)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Enable.Equal(that1.Enable) {
		return false
	}
	return true
}
func (this *RbacMode_Unspecified) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode_Unspecified)
	if !ok {
		that2, ok := that.(RbacMode_Unspecified)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *RbacMode_Disable) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode_Disable)
	if !ok {
		that2, ok := that.(RbacMode_Disable)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *RbacMode_Enable) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacMode_Enable)
	if !ok {
		that2, ok := that.(RbacMode_Enable)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *RbacStatus) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RbacStatus)
	if !ok {
		that2, ok := that.(RbacStatus)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Code != that1.Code {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
