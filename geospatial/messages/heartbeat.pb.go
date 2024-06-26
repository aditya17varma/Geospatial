// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.20.3
// source: proto/heartbeat.proto

package messages

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Heartbeat message
type Heartbeat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId       string                 `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	State        string                 `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
	Timestamp    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Coordinates  *Geocoordinate         `protobuf:"bytes,4,opt,name=coordinates,proto3" json:"coordinates,omitempty"`
	BatteryState *Battery               `protobuf:"bytes,5,opt,name=battery_state,json=batteryState,proto3" json:"battery_state,omitempty"`
}

func (x *Heartbeat) Reset() {
	*x = Heartbeat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Heartbeat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Heartbeat) ProtoMessage() {}

func (x *Heartbeat) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Heartbeat.ProtoReflect.Descriptor instead.
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{0}
}

func (x *Heartbeat) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *Heartbeat) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *Heartbeat) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Heartbeat) GetCoordinates() *Geocoordinate {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

func (x *Heartbeat) GetBatteryState() *Battery {
	if x != nil {
		return x.BatteryState
	}
	return nil
}

// Geocoordinates message
type Geocoordinate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *Geocoordinate) Reset() {
	*x = Geocoordinate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Geocoordinate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Geocoordinate) ProtoMessage() {}

func (x *Geocoordinate) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Geocoordinate.ProtoReflect.Descriptor instead.
func (*Geocoordinate) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{1}
}

func (x *Geocoordinate) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Geocoordinate) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

// Battery State message
type Battery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatteryLevel  int32 `protobuf:"varint,1,opt,name=battery_level,json=batteryLevel,proto3" json:"battery_level,omitempty"`
	SeekingCharge bool  `protobuf:"varint,2,opt,name=seeking_charge,json=seekingCharge,proto3" json:"seeking_charge,omitempty"`
}

func (x *Battery) Reset() {
	*x = Battery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Battery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Battery) ProtoMessage() {}

func (x *Battery) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Battery.ProtoReflect.Descriptor instead.
func (*Battery) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{2}
}

func (x *Battery) GetBatteryLevel() int32 {
	if x != nil {
		return x.BatteryLevel
	}
	return 0
}

func (x *Battery) GetSeekingCharge() bool {
	if x != nil {
		return x.SeekingCharge
	}
	return false
}

// Join Request
type JoinRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeAddr     string         `protobuf:"bytes,1,opt,name=node_addr,json=nodeAddr,proto3" json:"node_addr,omitempty"`
	NodeId       string         `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Coordinates  *Geocoordinate `protobuf:"bytes,3,opt,name=coordinates,proto3" json:"coordinates,omitempty"`
	BatteryState *Battery       `protobuf:"bytes,4,opt,name=battery_state,json=batteryState,proto3" json:"battery_state,omitempty"`
}

func (x *JoinRequest) Reset() {
	*x = JoinRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinRequest) ProtoMessage() {}

func (x *JoinRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinRequest.ProtoReflect.Descriptor instead.
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{3}
}

func (x *JoinRequest) GetNodeAddr() string {
	if x != nil {
		return x.NodeAddr
	}
	return ""
}

func (x *JoinRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *JoinRequest) GetCoordinates() *Geocoordinate {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

func (x *JoinRequest) GetBatteryState() *Battery {
	if x != nil {
		return x.BatteryState
	}
	return nil
}

// Recon Message
type ReconRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeAddr     string                 `protobuf:"bytes,1,opt,name=node_addr,json=nodeAddr,proto3" json:"node_addr,omitempty"`
	NodeId       string                 `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Coordinates  *Geocoordinate         `protobuf:"bytes,3,opt,name=coordinates,proto3" json:"coordinates,omitempty"`
	BatteryState *Battery               `protobuf:"bytes,4,opt,name=battery_state,json=batteryState,proto3" json:"battery_state,omitempty"`
	Timestamp    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ReconRequest) Reset() {
	*x = ReconRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconRequest) ProtoMessage() {}

func (x *ReconRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconRequest.ProtoReflect.Descriptor instead.
func (*ReconRequest) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{4}
}

func (x *ReconRequest) GetNodeAddr() string {
	if x != nil {
		return x.NodeAddr
	}
	return ""
}

func (x *ReconRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *ReconRequest) GetCoordinates() *Geocoordinate {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

func (x *ReconRequest) GetBatteryState() *Battery {
	if x != nil {
		return x.BatteryState
	}
	return nil
}

func (x *ReconRequest) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

// Join Response
type JoinResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accept      bool           `protobuf:"varint,1,opt,name=accept,proto3" json:"accept,omitempty"`
	NodeId      string         `protobuf:"bytes,2,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Coordinates *Geocoordinate `protobuf:"bytes,3,opt,name=coordinates,proto3" json:"coordinates,omitempty"`
}

func (x *JoinResponse) Reset() {
	*x = JoinResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinResponse) ProtoMessage() {}

func (x *JoinResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinResponse.ProtoReflect.Descriptor instead.
func (*JoinResponse) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{5}
}

func (x *JoinResponse) GetAccept() bool {
	if x != nil {
		return x.Accept
	}
	return false
}

func (x *JoinResponse) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *JoinResponse) GetCoordinates() *Geocoordinate {
	if x != nil {
		return x.Coordinates
	}
	return nil
}

// Kill Node
type KillNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kill bool `protobuf:"varint,1,opt,name=kill,proto3" json:"kill,omitempty"`
}

func (x *KillNode) Reset() {
	*x = KillNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KillNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KillNode) ProtoMessage() {}

func (x *KillNode) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KillNode.ProtoReflect.Descriptor instead.
func (*KillNode) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{6}
}

func (x *KillNode) GetKill() bool {
	if x != nil {
		return x.Kill
	}
	return false
}

// Server message
type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response string `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{7}
}

func (x *ServerMessage) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

// Error message
type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorDetail string `protobuf:"bytes,1,opt,name=error_detail,json=errorDetail,proto3" json:"error_detail,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{8}
}

func (x *Error) GetErrorDetail() string {
	if x != nil {
		return x.ErrorDetail
	}
	return ""
}

// Recon Response
type ReconResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HubHost string `protobuf:"bytes,1,opt,name=hub_host,json=hubHost,proto3" json:"hub_host,omitempty"`
	HubPort string `protobuf:"bytes,2,opt,name=hub_port,json=hubPort,proto3" json:"hub_port,omitempty"`
	NodeId  string `protobuf:"bytes,3,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Accept  bool   `protobuf:"varint,4,opt,name=accept,proto3" json:"accept,omitempty"`
}

func (x *ReconResponse) Reset() {
	*x = ReconResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReconResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReconResponse) ProtoMessage() {}

func (x *ReconResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReconResponse.ProtoReflect.Descriptor instead.
func (*ReconResponse) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{9}
}

func (x *ReconResponse) GetHubHost() string {
	if x != nil {
		return x.HubHost
	}
	return ""
}

func (x *ReconResponse) GetHubPort() string {
	if x != nil {
		return x.HubPort
	}
	return ""
}

func (x *ReconResponse) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *ReconResponse) GetAccept() bool {
	if x != nil {
		return x.Accept
	}
	return false
}

type Wrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Msg:
	//
	//	*Wrapper_ReconReq
	//	*Wrapper_ReconRes
	//	*Wrapper_JoinReq
	//	*Wrapper_JoinRes
	//	*Wrapper_KillNode
	//	*Wrapper_ServerMessage
	//	*Wrapper_HeartbeatMessage
	//	*Wrapper_ErrorMessage
	Msg isWrapper_Msg `protobuf_oneof:"msg"`
}

func (x *Wrapper) Reset() {
	*x = Wrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_heartbeat_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Wrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Wrapper) ProtoMessage() {}

func (x *Wrapper) ProtoReflect() protoreflect.Message {
	mi := &file_proto_heartbeat_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Wrapper.ProtoReflect.Descriptor instead.
func (*Wrapper) Descriptor() ([]byte, []int) {
	return file_proto_heartbeat_proto_rawDescGZIP(), []int{10}
}

func (m *Wrapper) GetMsg() isWrapper_Msg {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (x *Wrapper) GetReconReq() *ReconRequest {
	if x, ok := x.GetMsg().(*Wrapper_ReconReq); ok {
		return x.ReconReq
	}
	return nil
}

func (x *Wrapper) GetReconRes() *ReconResponse {
	if x, ok := x.GetMsg().(*Wrapper_ReconRes); ok {
		return x.ReconRes
	}
	return nil
}

func (x *Wrapper) GetJoinReq() *JoinRequest {
	if x, ok := x.GetMsg().(*Wrapper_JoinReq); ok {
		return x.JoinReq
	}
	return nil
}

func (x *Wrapper) GetJoinRes() *JoinResponse {
	if x, ok := x.GetMsg().(*Wrapper_JoinRes); ok {
		return x.JoinRes
	}
	return nil
}

func (x *Wrapper) GetKillNode() *KillNode {
	if x, ok := x.GetMsg().(*Wrapper_KillNode); ok {
		return x.KillNode
	}
	return nil
}

func (x *Wrapper) GetServerMessage() *ServerMessage {
	if x, ok := x.GetMsg().(*Wrapper_ServerMessage); ok {
		return x.ServerMessage
	}
	return nil
}

func (x *Wrapper) GetHeartbeatMessage() *Heartbeat {
	if x, ok := x.GetMsg().(*Wrapper_HeartbeatMessage); ok {
		return x.HeartbeatMessage
	}
	return nil
}

func (x *Wrapper) GetErrorMessage() *Error {
	if x, ok := x.GetMsg().(*Wrapper_ErrorMessage); ok {
		return x.ErrorMessage
	}
	return nil
}

type isWrapper_Msg interface {
	isWrapper_Msg()
}

type Wrapper_ReconReq struct {
	ReconReq *ReconRequest `protobuf:"bytes,1,opt,name=recon_req,json=reconReq,proto3,oneof"`
}

type Wrapper_ReconRes struct {
	ReconRes *ReconResponse `protobuf:"bytes,2,opt,name=recon_res,json=reconRes,proto3,oneof"`
}

type Wrapper_JoinReq struct {
	JoinReq *JoinRequest `protobuf:"bytes,3,opt,name=join_req,json=joinReq,proto3,oneof"`
}

type Wrapper_JoinRes struct {
	JoinRes *JoinResponse `protobuf:"bytes,4,opt,name=join_res,json=joinRes,proto3,oneof"`
}

type Wrapper_KillNode struct {
	KillNode *KillNode `protobuf:"bytes,5,opt,name=kill_node,json=killNode,proto3,oneof"`
}

type Wrapper_ServerMessage struct {
	ServerMessage *ServerMessage `protobuf:"bytes,6,opt,name=server_message,json=serverMessage,proto3,oneof"`
}

type Wrapper_HeartbeatMessage struct {
	HeartbeatMessage *Heartbeat `protobuf:"bytes,7,opt,name=heartbeat_message,json=heartbeatMessage,proto3,oneof"`
}

type Wrapper_ErrorMessage struct {
	ErrorMessage *Error `protobuf:"bytes,9,opt,name=error_message,json=errorMessage,proto3,oneof"`
}

func (*Wrapper_ReconReq) isWrapper_Msg() {}

func (*Wrapper_ReconRes) isWrapper_Msg() {}

func (*Wrapper_JoinReq) isWrapper_Msg() {}

func (*Wrapper_JoinRes) isWrapper_Msg() {}

func (*Wrapper_KillNode) isWrapper_Msg() {}

func (*Wrapper_ServerMessage) isWrapper_Msg() {}

func (*Wrapper_HeartbeatMessage) isWrapper_Msg() {}

func (*Wrapper_ErrorMessage) isWrapper_Msg() {}

var File_proto_heartbeat_proto protoreflect.FileDescriptor

var file_proto_heartbeat_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd5, 0x01, 0x0a, 0x09, 0x48, 0x65, 0x61,
	0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x30, 0x0a, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x47, 0x65, 0x6f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69,
	0x6e, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65,
	0x73, 0x12, 0x2d, 0x0a, 0x0d, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x42, 0x61, 0x74, 0x74, 0x65,
	0x72, 0x79, 0x52, 0x0c, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x22, 0x49, 0x0a, 0x0d, 0x47, 0x65, 0x6f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x55, 0x0a, 0x07, 0x42,
	0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72,
	0x79, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x62,
	0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x73,
	0x65, 0x65, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x68, 0x61, 0x72, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0d, 0x73, 0x65, 0x65, 0x6b, 0x69, 0x6e, 0x67, 0x43, 0x68, 0x61, 0x72,
	0x67, 0x65, 0x22, 0xa4, 0x01, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x6f, 0x64, 0x65, 0x41, 0x64, 0x64, 0x72, 0x12,
	0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x0b, 0x63, 0x6f, 0x6f, 0x72,
	0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e,
	0x47, 0x65, 0x6f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x63,
	0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x12, 0x2d, 0x0a, 0x0d, 0x62, 0x61,
	0x74, 0x74, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x08, 0x2e, 0x42, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x52, 0x0c, 0x62, 0x61, 0x74,
	0x74, 0x65, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65, 0x22, 0xdf, 0x01, 0x0a, 0x0c, 0x52, 0x65,
	0x63, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x6f,
	0x64, 0x65, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e,
	0x6f, 0x64, 0x65, 0x41, 0x64, 0x64, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x12, 0x30, 0x0a, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x47, 0x65, 0x6f, 0x63, 0x6f, 0x6f, 0x72, 0x64,
	0x69, 0x6e, 0x61, 0x74, 0x65, 0x52, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x73, 0x12, 0x2d, 0x0a, 0x0d, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x42, 0x61, 0x74, 0x74,
	0x65, 0x72, 0x79, 0x52, 0x0c, 0x62, 0x61, 0x74, 0x74, 0x65, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x71, 0x0a, 0x0c, 0x4a,
	0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61,
	0x63, 0x63, 0x65, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x61, 0x63, 0x63,
	0x65, 0x70, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x30, 0x0a, 0x0b,
	0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x47, 0x65, 0x6f, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74,
	0x65, 0x52, 0x0b, 0x63, 0x6f, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x73, 0x22, 0x1e,
	0x0a, 0x08, 0x4b, 0x69, 0x6c, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x69,
	0x6c, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6b, 0x69, 0x6c, 0x6c, 0x22, 0x2b,
	0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2a, 0x0a, 0x05, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x22, 0x76, 0x0a, 0x0d, 0x52, 0x65, 0x63, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x75, 0x62, 0x5f,
	0x68, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x75, 0x62, 0x48,
	0x6f, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x68, 0x75, 0x62, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x68, 0x75, 0x62, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x70,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x70, 0x74, 0x22,
	0x91, 0x03, 0x0a, 0x07, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x09, 0x72,
	0x65, 0x63, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x52, 0x65, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52,
	0x08, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x2d, 0x0a, 0x09, 0x72, 0x65, 0x63,
	0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x52,
	0x65, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x08,
	0x72, 0x65, 0x63, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x29, 0x0a, 0x08, 0x6a, 0x6f, 0x69, 0x6e,
	0x5f, 0x72, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x4a, 0x6f, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x07, 0x6a, 0x6f, 0x69, 0x6e,
	0x52, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x08, 0x6a, 0x6f, 0x69, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x4a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x07, 0x6a, 0x6f, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x12,
	0x28, 0x0a, 0x09, 0x6b, 0x69, 0x6c, 0x6c, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x09, 0x2e, 0x4b, 0x69, 0x6c, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52,
	0x08, 0x6b, 0x69, 0x6c, 0x6c, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x37, 0x0a, 0x0e, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x48, 0x00, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x39, 0x0a, 0x11, 0x68, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x5f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x48, 0x65, 0x61, 0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x48, 0x00, 0x52, 0x10, 0x68, 0x65, 0x61,
	0x72, 0x74, 0x62, 0x65, 0x61, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2d, 0x0a,
	0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x0c,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x42, 0x05, 0x0a, 0x03,
	0x6d, 0x73, 0x67, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_heartbeat_proto_rawDescOnce sync.Once
	file_proto_heartbeat_proto_rawDescData = file_proto_heartbeat_proto_rawDesc
)

func file_proto_heartbeat_proto_rawDescGZIP() []byte {
	file_proto_heartbeat_proto_rawDescOnce.Do(func() {
		file_proto_heartbeat_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_heartbeat_proto_rawDescData)
	})
	return file_proto_heartbeat_proto_rawDescData
}

var file_proto_heartbeat_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_heartbeat_proto_goTypes = []interface{}{
	(*Heartbeat)(nil),             // 0: Heartbeat
	(*Geocoordinate)(nil),         // 1: Geocoordinate
	(*Battery)(nil),               // 2: Battery
	(*JoinRequest)(nil),           // 3: JoinRequest
	(*ReconRequest)(nil),          // 4: ReconRequest
	(*JoinResponse)(nil),          // 5: JoinResponse
	(*KillNode)(nil),              // 6: KillNode
	(*ServerMessage)(nil),         // 7: ServerMessage
	(*Error)(nil),                 // 8: Error
	(*ReconResponse)(nil),         // 9: ReconResponse
	(*Wrapper)(nil),               // 10: Wrapper
	(*timestamppb.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_proto_heartbeat_proto_depIdxs = []int32{
	11, // 0: Heartbeat.timestamp:type_name -> google.protobuf.Timestamp
	1,  // 1: Heartbeat.coordinates:type_name -> Geocoordinate
	2,  // 2: Heartbeat.battery_state:type_name -> Battery
	1,  // 3: JoinRequest.coordinates:type_name -> Geocoordinate
	2,  // 4: JoinRequest.battery_state:type_name -> Battery
	1,  // 5: ReconRequest.coordinates:type_name -> Geocoordinate
	2,  // 6: ReconRequest.battery_state:type_name -> Battery
	11, // 7: ReconRequest.timestamp:type_name -> google.protobuf.Timestamp
	1,  // 8: JoinResponse.coordinates:type_name -> Geocoordinate
	4,  // 9: Wrapper.recon_req:type_name -> ReconRequest
	9,  // 10: Wrapper.recon_res:type_name -> ReconResponse
	3,  // 11: Wrapper.join_req:type_name -> JoinRequest
	5,  // 12: Wrapper.join_res:type_name -> JoinResponse
	6,  // 13: Wrapper.kill_node:type_name -> KillNode
	7,  // 14: Wrapper.server_message:type_name -> ServerMessage
	0,  // 15: Wrapper.heartbeat_message:type_name -> Heartbeat
	8,  // 16: Wrapper.error_message:type_name -> Error
	17, // [17:17] is the sub-list for method output_type
	17, // [17:17] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_proto_heartbeat_proto_init() }
func file_proto_heartbeat_proto_init() {
	if File_proto_heartbeat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_heartbeat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Heartbeat); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Geocoordinate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Battery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KillNode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReconResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_heartbeat_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Wrapper); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_proto_heartbeat_proto_msgTypes[10].OneofWrappers = []interface{}{
		(*Wrapper_ReconReq)(nil),
		(*Wrapper_ReconRes)(nil),
		(*Wrapper_JoinReq)(nil),
		(*Wrapper_JoinRes)(nil),
		(*Wrapper_KillNode)(nil),
		(*Wrapper_ServerMessage)(nil),
		(*Wrapper_HeartbeatMessage)(nil),
		(*Wrapper_ErrorMessage)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_heartbeat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_heartbeat_proto_goTypes,
		DependencyIndexes: file_proto_heartbeat_proto_depIdxs,
		MessageInfos:      file_proto_heartbeat_proto_msgTypes,
	}.Build()
	File_proto_heartbeat_proto = out.File
	file_proto_heartbeat_proto_rawDesc = nil
	file_proto_heartbeat_proto_goTypes = nil
	file_proto_heartbeat_proto_depIdxs = nil
}
