// Code generated by protoc-gen-go. DO NOT EDIT.
// source: node.proto

package im

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

type NodeType int32

const (
	NodeType_Super  NodeType = 0
	NodeType_Master NodeType = 1
)

var NodeType_name = map[int32]string{
	0: "Super",
	1: "Master",
}

var NodeType_value = map[string]int32{
	"Super":  0,
	"Master": 1,
}

func (x NodeType) String() string {
	return proto.EnumName(NodeType_name, int32(x))
}

func (NodeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

type Node struct {
	Type                 NodeType `protobuf:"varint,1,opt,name=Type,proto3,enum=im.NodeType" json:"Type,omitempty"`
	Votes                uint64   `protobuf:"varint,2,opt,name=Votes,proto3" json:"Votes,omitempty"`
	PeerID               string   `protobuf:"bytes,3,opt,name=PeerID,proto3" json:"PeerID,omitempty"`
	Owner                []byte   `protobuf:"bytes,5,opt,name=Owner,proto3" json:"Owner,omitempty"`
	Sig                  []byte   `protobuf:"bytes,4,opt,name=Sig,proto3" json:"Sig,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c843d59d2d938e7, []int{0}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetType() NodeType {
	if m != nil {
		return m.Type
	}
	return NodeType_Super
}

func (m *Node) GetVotes() uint64 {
	if m != nil {
		return m.Votes
	}
	return 0
}

func (m *Node) GetPeerID() string {
	if m != nil {
		return m.PeerID
	}
	return ""
}

func (m *Node) GetOwner() []byte {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *Node) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func init() {
	proto.RegisterEnum("im.NodeType", NodeType_name, NodeType_value)
	proto.RegisterType((*Node)(nil), "im.Node")
}

func init() { proto.RegisterFile("node.proto", fileDescriptor_0c843d59d2d938e7) }

var fileDescriptor_0c843d59d2d938e7 = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xca, 0xcb, 0x4f, 0x49,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xca, 0xcc, 0x55, 0xaa, 0xe3, 0x62, 0xf1, 0xcb,
	0x4f, 0x49, 0x15, 0x52, 0xe0, 0x62, 0x09, 0xa9, 0x2c, 0x48, 0x95, 0x60, 0x54, 0x60, 0xd4, 0xe0,
	0x33, 0xe2, 0xd1, 0xcb, 0xcc, 0xd5, 0x03, 0x89, 0x83, 0xc4, 0x82, 0xc0, 0x32, 0x42, 0x22, 0x5c,
	0xac, 0x61, 0xf9, 0x25, 0xa9, 0xc5, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x10, 0x8e, 0x90,
	0x18, 0x17, 0x5b, 0x40, 0x6a, 0x6a, 0x91, 0xa7, 0x8b, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x67, 0x10,
	0x94, 0x07, 0x52, 0xed, 0x5f, 0x9e, 0x97, 0x5a, 0x24, 0xc1, 0xaa, 0xc0, 0xa8, 0xc1, 0x13, 0x04,
	0xe1, 0x08, 0x09, 0x70, 0x31, 0x07, 0x67, 0xa6, 0x4b, 0xb0, 0x80, 0xc5, 0x40, 0x4c, 0x2d, 0x45,
	0x2e, 0x0e, 0x98, 0x3d, 0x42, 0x9c, 0x5c, 0xac, 0xc1, 0xa5, 0x05, 0xa9, 0x45, 0x02, 0x0c, 0x42,
	0x5c, 0x5c, 0x6c, 0xbe, 0x89, 0xc5, 0x25, 0xa9, 0x45, 0x02, 0x8c, 0x49, 0x6c, 0x60, 0xd7, 0x1a,
	0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x91, 0x0f, 0x22, 0x43, 0xbb, 0x00, 0x00, 0x00,
}
