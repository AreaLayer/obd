// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rsmc.proto

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

type RsmcPaymentRequest struct {
	ChannelId            string             `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Amount               float64            `protobuf:"fixed64,2,opt,name=amount,proto3" json:"amount,omitempty"`
	RecipientInfo        *RecipientNodeInfo `protobuf:"bytes,3,opt,name=recipientInfo,proto3" json:"recipientInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *RsmcPaymentRequest) Reset()         { *m = RsmcPaymentRequest{} }
func (m *RsmcPaymentRequest) String() string { return proto.CompactTextString(m) }
func (*RsmcPaymentRequest) ProtoMessage()    {}
func (*RsmcPaymentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fecbc7e287e63a67, []int{0}
}

func (m *RsmcPaymentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RsmcPaymentRequest.Unmarshal(m, b)
}
func (m *RsmcPaymentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RsmcPaymentRequest.Marshal(b, m, deterministic)
}
func (m *RsmcPaymentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RsmcPaymentRequest.Merge(m, src)
}
func (m *RsmcPaymentRequest) XXX_Size() int {
	return xxx_messageInfo_RsmcPaymentRequest.Size(m)
}
func (m *RsmcPaymentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RsmcPaymentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RsmcPaymentRequest proto.InternalMessageInfo

func (m *RsmcPaymentRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *RsmcPaymentRequest) GetAmount() float64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *RsmcPaymentRequest) GetRecipientInfo() *RecipientNodeInfo {
	if m != nil {
		return m.RecipientInfo
	}
	return nil
}

type RsmcPaymentResponse struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	AmountA              float64  `protobuf:"fixed64,2,opt,name=amount_a,json=amountA,proto3" json:"amount_a,omitempty"`
	AmountB              float64  `protobuf:"fixed64,3,opt,name=amount_b,json=amountB,proto3" json:"amount_b,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RsmcPaymentResponse) Reset()         { *m = RsmcPaymentResponse{} }
func (m *RsmcPaymentResponse) String() string { return proto.CompactTextString(m) }
func (*RsmcPaymentResponse) ProtoMessage()    {}
func (*RsmcPaymentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fecbc7e287e63a67, []int{1}
}

func (m *RsmcPaymentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RsmcPaymentResponse.Unmarshal(m, b)
}
func (m *RsmcPaymentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RsmcPaymentResponse.Marshal(b, m, deterministic)
}
func (m *RsmcPaymentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RsmcPaymentResponse.Merge(m, src)
}
func (m *RsmcPaymentResponse) XXX_Size() int {
	return xxx_messageInfo_RsmcPaymentResponse.Size(m)
}
func (m *RsmcPaymentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RsmcPaymentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RsmcPaymentResponse proto.InternalMessageInfo

func (m *RsmcPaymentResponse) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *RsmcPaymentResponse) GetAmountA() float64 {
	if m != nil {
		return m.AmountA
	}
	return 0
}

func (m *RsmcPaymentResponse) GetAmountB() float64 {
	if m != nil {
		return m.AmountB
	}
	return 0
}

type LatestRsmcTxRequest struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LatestRsmcTxRequest) Reset()         { *m = LatestRsmcTxRequest{} }
func (m *LatestRsmcTxRequest) String() string { return proto.CompactTextString(m) }
func (*LatestRsmcTxRequest) ProtoMessage()    {}
func (*LatestRsmcTxRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fecbc7e287e63a67, []int{2}
}

func (m *LatestRsmcTxRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LatestRsmcTxRequest.Unmarshal(m, b)
}
func (m *LatestRsmcTxRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LatestRsmcTxRequest.Marshal(b, m, deterministic)
}
func (m *LatestRsmcTxRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LatestRsmcTxRequest.Merge(m, src)
}
func (m *LatestRsmcTxRequest) XXX_Size() int {
	return xxx_messageInfo_LatestRsmcTxRequest.Size(m)
}
func (m *LatestRsmcTxRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LatestRsmcTxRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LatestRsmcTxRequest proto.InternalMessageInfo

func (m *LatestRsmcTxRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

type RsmcTxResponse struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	AmountA              float64  `protobuf:"fixed64,2,opt,name=amount_a,json=amountA,proto3" json:"amount_a,omitempty"`
	AmountB              float64  `protobuf:"fixed64,3,opt,name=amount_b,json=amountB,proto3" json:"amount_b,omitempty"`
	PeerA                string   `protobuf:"bytes,4,opt,name=peer_a,json=peerA,proto3" json:"peer_a,omitempty"`
	PeerB                string   `protobuf:"bytes,5,opt,name=peer_b,json=peerB,proto3" json:"peer_b,omitempty"`
	CurrState            int32    `protobuf:"varint,6,opt,name=curr_state,json=currState,proto3" json:"curr_state,omitempty"`
	TxHash               string   `protobuf:"bytes,7,opt,name=tx_hash,json=txHash,proto3" json:"tx_hash,omitempty"`
	TxType               int32    `protobuf:"varint,8,opt,name=tx_type,json=txType,proto3" json:"tx_type,omitempty"`
	H                    string   `protobuf:"bytes,9,opt,name=h,proto3" json:"h,omitempty"`
	R                    string   `protobuf:"bytes,10,opt,name=r,proto3" json:"r,omitempty"`
	AmountHtlc           float64  `protobuf:"fixed64,11,opt,name=amount_htlc,json=amountHtlc,proto3" json:"amount_htlc,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RsmcTxResponse) Reset()         { *m = RsmcTxResponse{} }
func (m *RsmcTxResponse) String() string { return proto.CompactTextString(m) }
func (*RsmcTxResponse) ProtoMessage()    {}
func (*RsmcTxResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fecbc7e287e63a67, []int{3}
}

func (m *RsmcTxResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RsmcTxResponse.Unmarshal(m, b)
}
func (m *RsmcTxResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RsmcTxResponse.Marshal(b, m, deterministic)
}
func (m *RsmcTxResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RsmcTxResponse.Merge(m, src)
}
func (m *RsmcTxResponse) XXX_Size() int {
	return xxx_messageInfo_RsmcTxResponse.Size(m)
}
func (m *RsmcTxResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RsmcTxResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RsmcTxResponse proto.InternalMessageInfo

func (m *RsmcTxResponse) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *RsmcTxResponse) GetAmountA() float64 {
	if m != nil {
		return m.AmountA
	}
	return 0
}

func (m *RsmcTxResponse) GetAmountB() float64 {
	if m != nil {
		return m.AmountB
	}
	return 0
}

func (m *RsmcTxResponse) GetPeerA() string {
	if m != nil {
		return m.PeerA
	}
	return ""
}

func (m *RsmcTxResponse) GetPeerB() string {
	if m != nil {
		return m.PeerB
	}
	return ""
}

func (m *RsmcTxResponse) GetCurrState() int32 {
	if m != nil {
		return m.CurrState
	}
	return 0
}

func (m *RsmcTxResponse) GetTxHash() string {
	if m != nil {
		return m.TxHash
	}
	return ""
}

func (m *RsmcTxResponse) GetTxType() int32 {
	if m != nil {
		return m.TxType
	}
	return 0
}

func (m *RsmcTxResponse) GetH() string {
	if m != nil {
		return m.H
	}
	return ""
}

func (m *RsmcTxResponse) GetR() string {
	if m != nil {
		return m.R
	}
	return ""
}

func (m *RsmcTxResponse) GetAmountHtlc() float64 {
	if m != nil {
		return m.AmountHtlc
	}
	return 0
}

type TxListRequest struct {
	ChannelId            string   `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	PageSize             int32    `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageIndex            int32    `protobuf:"varint,3,opt,name=page_index,json=pageIndex,proto3" json:"page_index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxListRequest) Reset()         { *m = TxListRequest{} }
func (m *TxListRequest) String() string { return proto.CompactTextString(m) }
func (*TxListRequest) ProtoMessage()    {}
func (*TxListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_fecbc7e287e63a67, []int{4}
}

func (m *TxListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxListRequest.Unmarshal(m, b)
}
func (m *TxListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxListRequest.Marshal(b, m, deterministic)
}
func (m *TxListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxListRequest.Merge(m, src)
}
func (m *TxListRequest) XXX_Size() int {
	return xxx_messageInfo_TxListRequest.Size(m)
}
func (m *TxListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TxListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TxListRequest proto.InternalMessageInfo

func (m *TxListRequest) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *TxListRequest) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *TxListRequest) GetPageIndex() int32 {
	if m != nil {
		return m.PageIndex
	}
	return 0
}

type TxListResponse struct {
	List                 []*RsmcTxResponse `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	TotalCount           int32             `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	PageSize             int32             `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	PageIndex            int32             `protobuf:"varint,4,opt,name=page_index,json=pageIndex,proto3" json:"page_index,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *TxListResponse) Reset()         { *m = TxListResponse{} }
func (m *TxListResponse) String() string { return proto.CompactTextString(m) }
func (*TxListResponse) ProtoMessage()    {}
func (*TxListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_fecbc7e287e63a67, []int{5}
}

func (m *TxListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxListResponse.Unmarshal(m, b)
}
func (m *TxListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxListResponse.Marshal(b, m, deterministic)
}
func (m *TxListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxListResponse.Merge(m, src)
}
func (m *TxListResponse) XXX_Size() int {
	return xxx_messageInfo_TxListResponse.Size(m)
}
func (m *TxListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TxListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TxListResponse proto.InternalMessageInfo

func (m *TxListResponse) GetList() []*RsmcTxResponse {
	if m != nil {
		return m.List
	}
	return nil
}

func (m *TxListResponse) GetTotalCount() int32 {
	if m != nil {
		return m.TotalCount
	}
	return 0
}

func (m *TxListResponse) GetPageSize() int32 {
	if m != nil {
		return m.PageSize
	}
	return 0
}

func (m *TxListResponse) GetPageIndex() int32 {
	if m != nil {
		return m.PageIndex
	}
	return 0
}

func init() {
	proto.RegisterType((*RsmcPaymentRequest)(nil), "proxy.RsmcPaymentRequest")
	proto.RegisterType((*RsmcPaymentResponse)(nil), "proxy.RsmcPaymentResponse")
	proto.RegisterType((*LatestRsmcTxRequest)(nil), "proxy.LatestRsmcTxRequest")
	proto.RegisterType((*RsmcTxResponse)(nil), "proxy.RsmcTxResponse")
	proto.RegisterType((*TxListRequest)(nil), "proxy.TxListRequest")
	proto.RegisterType((*TxListResponse)(nil), "proxy.TxListResponse")
}

func init() {
	proto.RegisterFile("rsmc.proto", fileDescriptor_fecbc7e287e63a67)
}

var fileDescriptor_fecbc7e287e63a67 = []byte{
	// 498 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0x95, 0xdb, 0x26, 0x6d, 0xa6, 0xbb, 0x2b, 0xe1, 0x65, 0xc1, 0x5b, 0x84, 0xa8, 0x72, 0x0a,
	0x97, 0x1e, 0x0a, 0x67, 0x44, 0xbb, 0x1c, 0xb6, 0xd2, 0x0a, 0x21, 0x6f, 0x4f, 0x5c, 0x22, 0x37,
	0x35, 0x24, 0xab, 0x34, 0x31, 0xb6, 0x2b, 0x25, 0xfb, 0x17, 0xf8, 0x03, 0x5c, 0xf9, 0x49, 0xfc,
	0x23, 0x64, 0xe7, 0x83, 0x06, 0x55, 0xb0, 0x97, 0x3d, 0xce, 0x7b, 0x63, 0xcf, 0x9b, 0xe7, 0x97,
	0x00, 0x48, 0xb5, 0x8b, 0x66, 0x42, 0xe6, 0x3a, 0xc7, 0x8e, 0x90, 0x79, 0x51, 0x4e, 0x3c, 0x29,
	0x6a, 0xc4, 0xff, 0x8e, 0x00, 0x53, 0xb5, 0x8b, 0x3e, 0xb1, 0x72, 0xc7, 0x33, 0x4d, 0xf9, 0xb7,
	0x3d, 0x57, 0x1a, 0xbf, 0x04, 0x88, 0x62, 0x96, 0x65, 0x3c, 0x0d, 0x93, 0x2d, 0x41, 0x53, 0x14,
	0x78, 0xd4, 0xab, 0x91, 0xd5, 0x16, 0x3f, 0x03, 0x97, 0xed, 0xf2, 0x7d, 0xa6, 0x49, 0x6f, 0x8a,
	0x02, 0x44, 0xeb, 0x0a, 0xbf, 0x83, 0x53, 0xc9, 0xa3, 0x44, 0x24, 0x3c, 0xd3, 0xab, 0xec, 0x4b,
	0x4e, 0xfa, 0x53, 0x14, 0x8c, 0xe7, 0x64, 0x66, 0xe7, 0xce, 0x68, 0xc3, 0x7d, 0xcc, 0xb7, 0xdc,
	0xf0, 0xb4, 0xdb, 0xee, 0xdf, 0xc1, 0x79, 0x47, 0x8c, 0x12, 0x79, 0xa6, 0xf8, 0xff, 0xd4, 0x5c,
	0xc2, 0xa8, 0x9a, 0x1f, 0xb2, 0x5a, 0xcf, 0xb0, 0xaa, 0x17, 0x07, 0xd4, 0xc6, 0x6a, 0x69, 0xa9,
	0xa5, 0xff, 0x16, 0xce, 0x6f, 0x98, 0xe6, 0x4a, 0x9b, 0x89, 0xeb, 0xe2, 0x61, 0x9b, 0xfb, 0x3f,
	0x7b, 0x70, 0xd6, 0x1c, 0x78, 0x44, 0x75, 0xf8, 0x02, 0x5c, 0xc1, 0xb9, 0x0c, 0x19, 0x19, 0xd8,
	0x0b, 0x1d, 0x53, 0x2d, 0x5a, 0x78, 0x43, 0x9c, 0x3f, 0xf0, 0xd2, 0x4a, 0xd8, 0x4b, 0x19, 0x2a,
	0xcd, 0x34, 0x27, 0xee, 0x14, 0x05, 0x0e, 0xf5, 0x0c, 0x72, 0x6b, 0x00, 0xfc, 0x1c, 0x86, 0xba,
	0x08, 0x63, 0xa6, 0x62, 0x32, 0xb4, 0xc7, 0x5c, 0x5d, 0x5c, 0x33, 0x15, 0xd7, 0x84, 0x2e, 0x05,
	0x27, 0x23, 0x7b, 0xc8, 0xd5, 0xc5, 0xba, 0x14, 0x1c, 0x9f, 0x00, 0x8a, 0x89, 0x67, 0x7b, 0x51,
	0x6c, 0x2a, 0x49, 0xa0, 0xaa, 0x24, 0x7e, 0x05, 0xe3, 0x5a, 0x75, 0xac, 0xd3, 0x88, 0x8c, 0xad,
	0x70, 0xa8, 0xa0, 0x6b, 0x9d, 0x46, 0xfe, 0x1d, 0x9c, 0xae, 0x8b, 0x9b, 0x44, 0x3d, 0x34, 0x4d,
	0x2f, 0xc0, 0x13, 0xec, 0x2b, 0x0f, 0x55, 0x72, 0xcf, 0xad, 0x45, 0x0e, 0x1d, 0x19, 0xe0, 0x36,
	0xb9, 0xb7, 0xee, 0x5a, 0x32, 0xc9, 0xb6, 0xbc, 0xb0, 0x2e, 0x39, 0xd4, 0xb6, 0xaf, 0x0c, 0xe0,
	0xff, 0x40, 0x70, 0xd6, 0x0c, 0xab, 0xdf, 0xe3, 0x35, 0x0c, 0xd2, 0x44, 0x69, 0x82, 0xa6, 0xfd,
	0x60, 0x3c, 0xbf, 0x68, 0xb2, 0xd7, 0x79, 0x34, 0x6a, 0x5b, 0xcc, 0x2a, 0x3a, 0xd7, 0x2c, 0x0d,
	0xa3, 0x36, 0xcc, 0x0e, 0x05, 0x0b, 0x5d, 0xd9, 0x40, 0x77, 0xa4, 0xf5, 0xff, 0x29, 0x6d, 0xf0,
	0x97, 0xb4, 0xf9, 0x2f, 0x04, 0x03, 0x33, 0x15, 0x7f, 0x80, 0xf1, 0x41, 0xaa, 0xf1, 0xe5, 0x81,
	0xa2, 0xee, 0x67, 0x37, 0x99, 0x1c, 0xa3, 0xea, 0xb5, 0x16, 0x70, 0x72, 0x98, 0x57, 0xdc, 0xf4,
	0x1e, 0x09, 0xf1, 0xe4, 0xf8, 0xd2, 0xf8, 0x3d, 0x3c, 0xa9, 0xbc, 0x5a, 0x96, 0x57, 0xad, 0xfb,
	0x4f, 0xeb, 0xde, 0xce, 0x93, 0xb5, 0x37, 0x74, 0xbd, 0x5d, 0x0e, 0x3e, 0xf7, 0xc4, 0x66, 0xe3,
	0xda, 0x7f, 0xc7, 0x9b, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x3b, 0xdb, 0x65, 0x34, 0x5b, 0x04,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RsmcClient is the client API for Rsmc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RsmcClient interface {
	RsmcPayment(ctx context.Context, in *RsmcPaymentRequest, opts ...grpc.CallOption) (*RsmcPaymentResponse, error)
	LatestRsmcTx(ctx context.Context, in *LatestRsmcTxRequest, opts ...grpc.CallOption) (*RsmcTxResponse, error)
	TxListByChannelId(ctx context.Context, in *TxListRequest, opts ...grpc.CallOption) (*TxListResponse, error)
}

type rsmcClient struct {
	cc grpc.ClientConnInterface
}

func NewRsmcClient(cc grpc.ClientConnInterface) RsmcClient {
	return &rsmcClient{cc}
}

func (c *rsmcClient) RsmcPayment(ctx context.Context, in *RsmcPaymentRequest, opts ...grpc.CallOption) (*RsmcPaymentResponse, error) {
	out := new(RsmcPaymentResponse)
	err := c.cc.Invoke(ctx, "/proxy.Rsmc/RsmcPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rsmcClient) LatestRsmcTx(ctx context.Context, in *LatestRsmcTxRequest, opts ...grpc.CallOption) (*RsmcTxResponse, error) {
	out := new(RsmcTxResponse)
	err := c.cc.Invoke(ctx, "/proxy.Rsmc/LatestRsmcTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rsmcClient) TxListByChannelId(ctx context.Context, in *TxListRequest, opts ...grpc.CallOption) (*TxListResponse, error) {
	out := new(TxListResponse)
	err := c.cc.Invoke(ctx, "/proxy.Rsmc/TxListByChannelId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RsmcServer is the server API for Rsmc service.
type RsmcServer interface {
	RsmcPayment(context.Context, *RsmcPaymentRequest) (*RsmcPaymentResponse, error)
	LatestRsmcTx(context.Context, *LatestRsmcTxRequest) (*RsmcTxResponse, error)
	TxListByChannelId(context.Context, *TxListRequest) (*TxListResponse, error)
}

// UnimplementedRsmcServer can be embedded to have forward compatible implementations.
type UnimplementedRsmcServer struct {
}

func (*UnimplementedRsmcServer) RsmcPayment(ctx context.Context, req *RsmcPaymentRequest) (*RsmcPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RsmcPayment not implemented")
}
func (*UnimplementedRsmcServer) LatestRsmcTx(ctx context.Context, req *LatestRsmcTxRequest) (*RsmcTxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LatestRsmcTx not implemented")
}
func (*UnimplementedRsmcServer) TxListByChannelId(ctx context.Context, req *TxListRequest) (*TxListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TxListByChannelId not implemented")
}

func RegisterRsmcServer(s *grpc.Server, srv RsmcServer) {
	s.RegisterService(&_Rsmc_serviceDesc, srv)
}

func _Rsmc_RsmcPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RsmcPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RsmcServer).RsmcPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Rsmc/RsmcPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RsmcServer).RsmcPayment(ctx, req.(*RsmcPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rsmc_LatestRsmcTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LatestRsmcTxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RsmcServer).LatestRsmcTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Rsmc/LatestRsmcTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RsmcServer).LatestRsmcTx(ctx, req.(*LatestRsmcTxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rsmc_TxListByChannelId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RsmcServer).TxListByChannelId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proxy.Rsmc/TxListByChannelId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RsmcServer).TxListByChannelId(ctx, req.(*TxListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Rsmc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proxy.Rsmc",
	HandlerType: (*RsmcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RsmcPayment",
			Handler:    _Rsmc_RsmcPayment_Handler,
		},
		{
			MethodName: "LatestRsmcTx",
			Handler:    _Rsmc_LatestRsmcTx_Handler,
		},
		{
			MethodName: "TxListByChannelId",
			Handler:    _Rsmc_TxListByChannelId_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rsmc.proto",
}
