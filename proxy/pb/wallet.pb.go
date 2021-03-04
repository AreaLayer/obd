// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wallet.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type LoginRequest struct {
	Mnemonic             string   `protobuf:"bytes,1,opt,name=mnemonic,proto3" json:"mnemonic,omitempty"`
	LoginToken           string   `protobuf:"bytes,2,opt,name=login_token,json=loginToken,proto3" json:"login_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginRequest) Reset()         { *m = LoginRequest{} }
func (m *LoginRequest) String() string { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()    {}
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{0}
}

func (m *LoginRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginRequest.Unmarshal(m, b)
}
func (m *LoginRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginRequest.Marshal(b, m, deterministic)
}
func (m *LoginRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginRequest.Merge(m, src)
}
func (m *LoginRequest) XXX_Size() int {
	return xxx_messageInfo_LoginRequest.Size(m)
}
func (m *LoginRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LoginRequest proto.InternalMessageInfo

func (m *LoginRequest) GetMnemonic() string {
	if m != nil {
		return m.Mnemonic
	}
	return ""
}

func (m *LoginRequest) GetLoginToken() string {
	if m != nil {
		return m.LoginToken
	}
	return ""
}

type LoginResponse struct {
	UserPeerId           string   `protobuf:"bytes,1,opt,name=user_peerId,json=userPeerId,proto3" json:"user_peerId,omitempty"`
	NodePeerId           string   `protobuf:"bytes,2,opt,name=node_peerId,json=nodePeerId,proto3" json:"node_peerId,omitempty"`
	NodeAddress          string   `protobuf:"bytes,3,opt,name=node_address,json=nodeAddress,proto3" json:"node_address,omitempty"`
	HtlcFeeRate          float64  `protobuf:"fixed64,4,opt,name=htlc_fee_rate,json=htlcFeeRate,proto3" json:"htlc_fee_rate,omitempty"`
	HtlcMaxFee           float64  `protobuf:"fixed64,5,opt,name=htlc_max_fee,json=htlcMaxFee,proto3" json:"htlc_max_fee,omitempty"`
	ChainNodeType        string   `protobuf:"bytes,6,opt,name=chain_node_type,json=chainNodeType,proto3" json:"chain_node_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResponse) Reset()         { *m = LoginResponse{} }
func (m *LoginResponse) String() string { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()    {}
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{1}
}

func (m *LoginResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResponse.Unmarshal(m, b)
}
func (m *LoginResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResponse.Marshal(b, m, deterministic)
}
func (m *LoginResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResponse.Merge(m, src)
}
func (m *LoginResponse) XXX_Size() int {
	return xxx_messageInfo_LoginResponse.Size(m)
}
func (m *LoginResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResponse proto.InternalMessageInfo

func (m *LoginResponse) GetUserPeerId() string {
	if m != nil {
		return m.UserPeerId
	}
	return ""
}

func (m *LoginResponse) GetNodePeerId() string {
	if m != nil {
		return m.NodePeerId
	}
	return ""
}

func (m *LoginResponse) GetNodeAddress() string {
	if m != nil {
		return m.NodeAddress
	}
	return ""
}

func (m *LoginResponse) GetHtlcFeeRate() float64 {
	if m != nil {
		return m.HtlcFeeRate
	}
	return 0
}

func (m *LoginResponse) GetHtlcMaxFee() float64 {
	if m != nil {
		return m.HtlcMaxFee
	}
	return 0
}

func (m *LoginResponse) GetChainNodeType() string {
	if m != nil {
		return m.ChainNodeType
	}
	return ""
}

type ChangePasswordRequest struct {
	CurrentPassword      string   `protobuf:"bytes,1,opt,name=current_password,json=currentPassword,proto3" json:"current_password,omitempty"`
	NewPassword          string   `protobuf:"bytes,2,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangePasswordRequest) Reset()         { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string { return proto.CompactTextString(m) }
func (*ChangePasswordRequest) ProtoMessage()    {}
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{2}
}

func (m *ChangePasswordRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangePasswordRequest.Unmarshal(m, b)
}
func (m *ChangePasswordRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangePasswordRequest.Marshal(b, m, deterministic)
}
func (m *ChangePasswordRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangePasswordRequest.Merge(m, src)
}
func (m *ChangePasswordRequest) XXX_Size() int {
	return xxx_messageInfo_ChangePasswordRequest.Size(m)
}
func (m *ChangePasswordRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangePasswordRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ChangePasswordRequest proto.InternalMessageInfo

func (m *ChangePasswordRequest) GetCurrentPassword() string {
	if m != nil {
		return m.CurrentPassword
	}
	return ""
}

func (m *ChangePasswordRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

type ChangePasswordResponse struct {
	Result               string   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangePasswordResponse) Reset()         { *m = ChangePasswordResponse{} }
func (m *ChangePasswordResponse) String() string { return proto.CompactTextString(m) }
func (*ChangePasswordResponse) ProtoMessage()    {}
func (*ChangePasswordResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{3}
}

func (m *ChangePasswordResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangePasswordResponse.Unmarshal(m, b)
}
func (m *ChangePasswordResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangePasswordResponse.Marshal(b, m, deterministic)
}
func (m *ChangePasswordResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangePasswordResponse.Merge(m, src)
}
func (m *ChangePasswordResponse) XXX_Size() int {
	return xxx_messageInfo_ChangePasswordResponse.Size(m)
}
func (m *ChangePasswordResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangePasswordResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ChangePasswordResponse proto.InternalMessageInfo

func (m *ChangePasswordResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

type LogoutRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutRequest) Reset()         { *m = LogoutRequest{} }
func (m *LogoutRequest) String() string { return proto.CompactTextString(m) }
func (*LogoutRequest) ProtoMessage()    {}
func (*LogoutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{4}
}

func (m *LogoutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutRequest.Unmarshal(m, b)
}
func (m *LogoutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutRequest.Marshal(b, m, deterministic)
}
func (m *LogoutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutRequest.Merge(m, src)
}
func (m *LogoutRequest) XXX_Size() int {
	return xxx_messageInfo_LogoutRequest.Size(m)
}
func (m *LogoutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutRequest proto.InternalMessageInfo

type LogoutResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogoutResponse) Reset()         { *m = LogoutResponse{} }
func (m *LogoutResponse) String() string { return proto.CompactTextString(m) }
func (*LogoutResponse) ProtoMessage()    {}
func (*LogoutResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{5}
}

func (m *LogoutResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogoutResponse.Unmarshal(m, b)
}
func (m *LogoutResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogoutResponse.Marshal(b, m, deterministic)
}
func (m *LogoutResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogoutResponse.Merge(m, src)
}
func (m *LogoutResponse) XXX_Size() int {
	return xxx_messageInfo_LogoutResponse.Size(m)
}
func (m *LogoutResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogoutResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogoutResponse proto.InternalMessageInfo

type GenSeedRequest struct {
	AezeedPassphrase     []byte   `protobuf:"bytes,1,opt,name=aezeed_passphrase,json=aezeedPassphrase,proto3" json:"aezeed_passphrase,omitempty"`
	SeedEntropy          []byte   `protobuf:"bytes,2,opt,name=seed_entropy,json=seedEntropy,proto3" json:"seed_entropy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenSeedRequest) Reset()         { *m = GenSeedRequest{} }
func (m *GenSeedRequest) String() string { return proto.CompactTextString(m) }
func (*GenSeedRequest) ProtoMessage()    {}
func (*GenSeedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{6}
}

func (m *GenSeedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenSeedRequest.Unmarshal(m, b)
}
func (m *GenSeedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenSeedRequest.Marshal(b, m, deterministic)
}
func (m *GenSeedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenSeedRequest.Merge(m, src)
}
func (m *GenSeedRequest) XXX_Size() int {
	return xxx_messageInfo_GenSeedRequest.Size(m)
}
func (m *GenSeedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenSeedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenSeedRequest proto.InternalMessageInfo

func (m *GenSeedRequest) GetAezeedPassphrase() []byte {
	if m != nil {
		return m.AezeedPassphrase
	}
	return nil
}

func (m *GenSeedRequest) GetSeedEntropy() []byte {
	if m != nil {
		return m.SeedEntropy
	}
	return nil
}

type GenSeedResponse struct {
	CipherSeedMnemonic   string   `protobuf:"bytes,1,opt,name=cipher_seed_mnemonic,json=cipherSeedMnemonic,proto3" json:"cipher_seed_mnemonic,omitempty"`
	EncipheredSeed       string   `protobuf:"bytes,2,opt,name=enciphered_seed,json=encipheredSeed,proto3" json:"enciphered_seed,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenSeedResponse) Reset()         { *m = GenSeedResponse{} }
func (m *GenSeedResponse) String() string { return proto.CompactTextString(m) }
func (*GenSeedResponse) ProtoMessage()    {}
func (*GenSeedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{7}
}

func (m *GenSeedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenSeedResponse.Unmarshal(m, b)
}
func (m *GenSeedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenSeedResponse.Marshal(b, m, deterministic)
}
func (m *GenSeedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenSeedResponse.Merge(m, src)
}
func (m *GenSeedResponse) XXX_Size() int {
	return xxx_messageInfo_GenSeedResponse.Size(m)
}
func (m *GenSeedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenSeedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenSeedResponse proto.InternalMessageInfo

func (m *GenSeedResponse) GetCipherSeedMnemonic() string {
	if m != nil {
		return m.CipherSeedMnemonic
	}
	return ""
}

func (m *GenSeedResponse) GetEncipheredSeed() string {
	if m != nil {
		return m.EncipheredSeed
	}
	return ""
}

type EstimateFeeRequest struct {
	ConfTarget           int32    `protobuf:"varint,1,opt,name=conf_target,json=confTarget,proto3" json:"conf_target,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EstimateFeeRequest) Reset()         { *m = EstimateFeeRequest{} }
func (m *EstimateFeeRequest) String() string { return proto.CompactTextString(m) }
func (*EstimateFeeRequest) ProtoMessage()    {}
func (*EstimateFeeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{8}
}

func (m *EstimateFeeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EstimateFeeRequest.Unmarshal(m, b)
}
func (m *EstimateFeeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EstimateFeeRequest.Marshal(b, m, deterministic)
}
func (m *EstimateFeeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EstimateFeeRequest.Merge(m, src)
}
func (m *EstimateFeeRequest) XXX_Size() int {
	return xxx_messageInfo_EstimateFeeRequest.Size(m)
}
func (m *EstimateFeeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_EstimateFeeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_EstimateFeeRequest proto.InternalMessageInfo

func (m *EstimateFeeRequest) GetConfTarget() int32 {
	if m != nil {
		return m.ConfTarget
	}
	return 0
}

type EstimateFeeResponse struct {
	SatPerKw             int64    `protobuf:"varint,1,opt,name=sat_per_kw,json=satPerKw,proto3" json:"sat_per_kw,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EstimateFeeResponse) Reset()         { *m = EstimateFeeResponse{} }
func (m *EstimateFeeResponse) String() string { return proto.CompactTextString(m) }
func (*EstimateFeeResponse) ProtoMessage()    {}
func (*EstimateFeeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{9}
}

func (m *EstimateFeeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EstimateFeeResponse.Unmarshal(m, b)
}
func (m *EstimateFeeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EstimateFeeResponse.Marshal(b, m, deterministic)
}
func (m *EstimateFeeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EstimateFeeResponse.Merge(m, src)
}
func (m *EstimateFeeResponse) XXX_Size() int {
	return xxx_messageInfo_EstimateFeeResponse.Size(m)
}
func (m *EstimateFeeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EstimateFeeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EstimateFeeResponse proto.InternalMessageInfo

func (m *EstimateFeeResponse) GetSatPerKw() int64 {
	if m != nil {
		return m.SatPerKw
	}
	return 0
}

type AddrRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddrRequest) Reset()         { *m = AddrRequest{} }
func (m *AddrRequest) String() string { return proto.CompactTextString(m) }
func (*AddrRequest) ProtoMessage()    {}
func (*AddrRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{10}
}

func (m *AddrRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddrRequest.Unmarshal(m, b)
}
func (m *AddrRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddrRequest.Marshal(b, m, deterministic)
}
func (m *AddrRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddrRequest.Merge(m, src)
}
func (m *AddrRequest) XXX_Size() int {
	return xxx_messageInfo_AddrRequest.Size(m)
}
func (m *AddrRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddrRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddrRequest proto.InternalMessageInfo

type AddrResponse struct {
	//
	//The address encoded using a bech32 format.
	Addr                 string   `protobuf:"bytes,1,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddrResponse) Reset()         { *m = AddrResponse{} }
func (m *AddrResponse) String() string { return proto.CompactTextString(m) }
func (*AddrResponse) ProtoMessage()    {}
func (*AddrResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b88fd140af4deb6f, []int{11}
}

func (m *AddrResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddrResponse.Unmarshal(m, b)
}
func (m *AddrResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddrResponse.Marshal(b, m, deterministic)
}
func (m *AddrResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddrResponse.Merge(m, src)
}
func (m *AddrResponse) XXX_Size() int {
	return xxx_messageInfo_AddrResponse.Size(m)
}
func (m *AddrResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddrResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddrResponse proto.InternalMessageInfo

func (m *AddrResponse) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func init() {
	proto.RegisterType((*LoginRequest)(nil), "proxy.LoginRequest")
	proto.RegisterType((*LoginResponse)(nil), "proxy.LoginResponse")
	proto.RegisterType((*ChangePasswordRequest)(nil), "proxy.ChangePasswordRequest")
	proto.RegisterType((*ChangePasswordResponse)(nil), "proxy.ChangePasswordResponse")
	proto.RegisterType((*LogoutRequest)(nil), "proxy.LogoutRequest")
	proto.RegisterType((*LogoutResponse)(nil), "proxy.LogoutResponse")
	proto.RegisterType((*GenSeedRequest)(nil), "proxy.GenSeedRequest")
	proto.RegisterType((*GenSeedResponse)(nil), "proxy.GenSeedResponse")
	proto.RegisterType((*EstimateFeeRequest)(nil), "proxy.EstimateFeeRequest")
	proto.RegisterType((*EstimateFeeResponse)(nil), "proxy.EstimateFeeResponse")
	proto.RegisterType((*AddrRequest)(nil), "proxy.AddrRequest")
	proto.RegisterType((*AddrResponse)(nil), "proxy.AddrResponse")
}

func init() {
	proto.RegisterFile("wallet.proto", fileDescriptor_b88fd140af4deb6f)
}

var fileDescriptor_b88fd140af4deb6f = []byte{
	// 620 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x94, 0x5d, 0x4f, 0x1b, 0x3b,
	0x10, 0x86, 0x95, 0x40, 0x72, 0x38, 0x93, 0x4d, 0xc2, 0x31, 0x1f, 0xca, 0x59, 0x71, 0x04, 0xc7,
	0x17, 0x2d, 0x55, 0x25, 0x84, 0x40, 0x48, 0xbd, 0xed, 0x07, 0x54, 0x15, 0x05, 0x45, 0x5b, 0xa4,
	0x4a, 0xbd, 0xd9, 0x9a, 0xec, 0x40, 0x22, 0x12, 0x7b, 0x6b, 0x3b, 0x4a, 0xd2, 0x9f, 0xdb, 0x3f,
	0xd1, 0xdb, 0xca, 0xb3, 0xb3, 0x21, 0x49, 0xe9, 0x5d, 0xf6, 0x99, 0xf7, 0x1d, 0xcf, 0x8c, 0xc7,
	0x81, 0x68, 0xa2, 0x86, 0x43, 0xf4, 0x47, 0xb9, 0x35, 0xde, 0x88, 0x5a, 0x6e, 0xcd, 0x74, 0x26,
	0x2f, 0x21, 0xfa, 0x68, 0xee, 0x07, 0x3a, 0xc1, 0x6f, 0x63, 0x74, 0x5e, 0xc4, 0xb0, 0x31, 0xd2,
	0x38, 0x32, 0x7a, 0xd0, 0xeb, 0x54, 0x0e, 0x2a, 0x87, 0x7f, 0x27, 0xf3, 0x6f, 0xb1, 0x0f, 0x8d,
	0x61, 0xd0, 0xa6, 0xde, 0x3c, 0xa0, 0xee, 0x54, 0x29, 0x0c, 0x84, 0x6e, 0x02, 0x91, 0x3f, 0x2a,
	0xd0, 0xe4, 0x6c, 0x2e, 0x37, 0xda, 0x61, 0xb0, 0x8c, 0x1d, 0xda, 0x34, 0x47, 0xb4, 0x1f, 0x32,
	0xce, 0x08, 0x01, 0x75, 0x89, 0x04, 0x81, 0x36, 0x19, 0x96, 0x02, 0xce, 0x19, 0x10, 0x0b, 0xfe,
	0x87, 0x88, 0x04, 0x2a, 0xcb, 0x2c, 0x3a, 0xd7, 0x59, 0x23, 0x05, 0x99, 0x5e, 0x17, 0x48, 0x48,
	0x68, 0xf6, 0xfd, 0xb0, 0x97, 0xde, 0x21, 0xa6, 0x56, 0x79, 0xec, 0xac, 0x1f, 0x54, 0x0e, 0x2b,
	0x49, 0x23, 0xc0, 0x0b, 0xc4, 0x44, 0x79, 0x14, 0x07, 0x10, 0x91, 0x66, 0xa4, 0xa6, 0x41, 0xd7,
	0xa9, 0x91, 0x04, 0x02, 0xbb, 0x52, 0xd3, 0x0b, 0x44, 0xf1, 0x0c, 0xda, 0xbd, 0xbe, 0x1a, 0xe8,
	0x94, 0x8e, 0xf3, 0xb3, 0x1c, 0x3b, 0x75, 0x3a, 0xab, 0x49, 0xf8, 0xda, 0x64, 0x78, 0x33, 0xcb,
	0x51, 0x22, 0xec, 0xbc, 0xed, 0x2b, 0x7d, 0x8f, 0x5d, 0xe5, 0xdc, 0xc4, 0xd8, 0xac, 0x1c, 0xdd,
	0x0b, 0xd8, 0xec, 0x8d, 0xad, 0x45, 0xed, 0xd3, 0x9c, 0x43, 0xdc, 0x70, 0x9b, 0x79, 0xe9, 0xa0,
	0xa6, 0x70, 0xf2, 0x28, 0xab, 0x72, 0x53, 0x38, 0x29, 0x25, 0xf2, 0x18, 0x76, 0x57, 0x8f, 0xe1,
	0x99, 0xee, 0x42, 0xdd, 0xa2, 0x1b, 0x0f, 0x3d, 0x67, 0xe7, 0x2f, 0xd9, 0xa6, 0xe1, 0x9b, 0xb1,
	0xe7, 0x82, 0xe4, 0x26, 0xb4, 0x4a, 0x50, 0x58, 0xe5, 0x57, 0x68, 0xbd, 0x47, 0xfd, 0x09, 0x71,
	0x5e, 0xf4, 0x4b, 0xf8, 0x47, 0xe1, 0x77, 0xc4, 0x8c, 0x8a, 0xc9, 0xfb, 0x56, 0x39, 0xa4, 0xbc,
	0x51, 0xb2, 0x59, 0x04, 0xba, 0x73, 0x1e, 0xca, 0x76, 0x41, 0x8a, 0xda, 0x5b, 0x93, 0xcf, 0xa8,
	0xec, 0x28, 0x69, 0x04, 0x76, 0x5e, 0x20, 0x39, 0x84, 0xf6, 0xfc, 0x04, 0xae, 0xf7, 0x18, 0xb6,
	0x7b, 0x83, 0xbc, 0x8f, 0x36, 0x25, 0xf3, 0xca, 0x7a, 0x89, 0x22, 0x16, 0x1c, 0x57, 0xe5, 0xa2,
	0x3d, 0x87, 0x36, 0xea, 0x82, 0x63, 0x46, 0x2e, 0x9e, 0x50, 0xeb, 0x11, 0x07, 0x83, 0x3c, 0x03,
	0x71, 0xee, 0xfc, 0x60, 0xa4, 0x3c, 0x86, 0x8b, 0xe6, 0x9e, 0xf6, 0xa1, 0xd1, 0x33, 0xfa, 0x2e,
	0xf5, 0xca, 0xde, 0x63, 0x31, 0xa5, 0x5a, 0x02, 0x01, 0xdd, 0x10, 0x91, 0xa7, 0xb0, 0xb5, 0x64,
	0xe3, 0x42, 0xf7, 0x00, 0x9c, 0xf2, 0x69, 0x8e, 0x36, 0x7d, 0x98, 0x90, 0x6d, 0x2d, 0xd9, 0x70,
	0xca, 0x77, 0xd1, 0x5e, 0x4e, 0x64, 0x13, 0x1a, 0x61, 0xe1, 0xca, 0xe1, 0x4a, 0x88, 0x8a, 0x4f,
	0x36, 0x0b, 0x58, 0x0f, 0x2b, 0xca, 0x5d, 0xd1, 0xef, 0x93, 0x9f, 0x55, 0xa8, 0x7f, 0xa6, 0x47,
	0x27, 0x5e, 0xc1, 0x5f, 0x3c, 0x17, 0xb1, 0x73, 0x44, 0x4f, 0xef, 0x68, 0xf9, 0x26, 0xe2, 0xdd,
	0x55, 0xcc, 0x89, 0x4f, 0xa0, 0x46, 0x6f, 0x4a, 0x6c, 0xb1, 0x60, 0xf1, 0xbd, 0xc6, 0xdb, 0xcb,
	0x90, 0x3d, 0x57, 0xd0, 0x5a, 0x5e, 0x1e, 0xb1, 0xc7, 0xba, 0x27, 0x57, 0x37, 0xfe, 0xef, 0x0f,
	0x51, 0x4e, 0x77, 0x06, 0xf5, 0x62, 0x91, 0xc4, 0xc2, 0x71, 0x8f, 0x8b, 0x16, 0xef, 0xac, 0x50,
	0xb6, 0xbd, 0x83, 0xc6, 0xc2, 0x98, 0xc5, 0xbf, 0xac, 0xfa, 0xfd, 0xc6, 0xe2, 0xf8, 0xa9, 0x10,
	0x67, 0x39, 0x85, 0x8d, 0x6b, 0x9c, 0xfa, 0x30, 0x6c, 0x21, 0x58, 0xb7, 0x70, 0x11, 0xf1, 0xd6,
	0x12, 0x2b, 0x4c, 0x6f, 0xd6, 0xbf, 0x54, 0xf3, 0xdb, 0xdb, 0x3a, 0xfd, 0xd5, 0x9d, 0xfe, 0x0a,
	0x00, 0x00, 0xff, 0xff, 0x18, 0xc0, 0x75, 0x35, 0xfa, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WalletClient is the client API for Wallet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WalletClient interface {
	GenSeed(ctx context.Context, in *GenSeedRequest, opts ...grpc.CallOption) (*GenSeedResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
	EstimateFee(ctx context.Context, in *EstimateFeeRequest, opts ...grpc.CallOption) (*EstimateFeeResponse, error)
	//
	//NextAddr returns the next unused address within the wallet.
	NextAddr(ctx context.Context, in *AddrRequest, opts ...grpc.CallOption) (*AddrResponse, error)
}

type walletClient struct {
	cc grpc.ClientConnInterface
}

func NewWalletClient(cc grpc.ClientConnInterface) WalletClient {
	return &walletClient{cc}
}

func (c *walletClient) GenSeed(ctx context.Context, in *GenSeedRequest, opts ...grpc.CallOption) (*GenSeedResponse, error) {
	out := new(GenSeedResponse)
	err := c.cc.Invoke(ctx, "/proxy.Wallet/GenSeed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/proxy.Wallet/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) {
	out := new(ChangePasswordResponse)
	err := c.cc.Invoke(ctx, "/proxy.Wallet/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, "/proxy.Wallet/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) EstimateFee(ctx context.Context, in *EstimateFeeRequest, opts ...grpc.CallOption) (*EstimateFeeResponse, error) {
	out := new(EstimateFeeResponse)
	err := c.cc.Invoke(ctx, "/proxy.Wallet/EstimateFee", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *walletClient) NextAddr(ctx context.Context, in *AddrRequest, opts ...grpc.CallOption) (*AddrResponse, error) {
	out := new(AddrResponse)
	err := c.cc.Invoke(ctx, "/proxy.Wallet/NextAddr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WalletServer is the server API for Wallet service.
type WalletServer interface {
	GenSeed(context.Context, *GenSeedRequest) (*GenSeedResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	EstimateFee(context.Context, *EstimateFeeRequest) (*EstimateFeeResponse, error)
	//
	//NextAddr returns the next unused address within the wallet.
	NextAddr(context.Context, *AddrRequest) (*AddrResponse, error)
}

// UnimplementedWalletServer can be embedded to have forward compatible implementations.
type UnimplementedWalletServer struct {
}

func (*UnimplementedWalletServer) GenSeed(ctx context.Context, req *GenSeedRequest) (*GenSeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenSeed not implemented")
}
func (*UnimplementedWalletServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (*UnimplementedWalletServer) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (*UnimplementedWalletServer) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (*UnimplementedWalletServer) EstimateFee(ctx context.Context, req *EstimateFeeRequest) (*EstimateFeeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EstimateFee not implemented")
}
func (*UnimplementedWalletServer) NextAddr(ctx context.Context, req *AddrRequest) (*AddrResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NextAddr not implemented")
}

func RegisterWalletServer(s *grpc.Server, srv WalletServer) {
	s.RegisterService(&_Wallet_serviceDesc, srv)
}

func _Wallet_GenSeed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenSeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).GenSeed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Wallet/GenSeed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).GenSeed(ctx, req.(*GenSeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Wallet/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Wallet/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Wallet/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_EstimateFee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EstimateFeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).EstimateFee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Wallet/EstimateFee",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).EstimateFee(ctx, req.(*EstimateFeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Wallet_NextAddr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddrRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WalletServer).NextAddr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Wallet/NextAddr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WalletServer).NextAddr(ctx, req.(*AddrRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Wallet_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proxy.Wallet",
	HandlerType: (*WalletServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenSeed",
			Handler:    _Wallet_GenSeed_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Wallet_Login_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _Wallet_ChangePassword_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Wallet_Logout_Handler,
		},
		{
			MethodName: "EstimateFee",
			Handler:    _Wallet_EstimateFee_Handler,
		},
		{
			MethodName: "NextAddr",
			Handler:    _Wallet_NextAddr_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wallet.proto",
}