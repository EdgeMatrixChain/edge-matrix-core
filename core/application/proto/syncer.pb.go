// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: core/application/proto/syncer.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GetDataRequest is a request for GetData
type GetDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The hash of Data to sync
	DataHash string `protobuf:"bytes,1,opt,name=dataHash,proto3" json:"dataHash,omitempty"`
}

func (x *GetDataRequest) Reset() {
	*x = GetDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_application_proto_syncer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataRequest) ProtoMessage() {}

func (x *GetDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_application_proto_syncer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataRequest.ProtoReflect.Descriptor instead.
func (*GetDataRequest) Descriptor() ([]byte, []int) {
	return file_core_application_proto_syncer_proto_rawDescGZIP(), []int{0}
}

func (x *GetDataRequest) GetDataHash() string {
	if x != nil {
		return x.DataHash
	}
	return ""
}

// PostPeerStatusRequest is a request for post poc
type PostPeerStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The PeerStatus
	NodeId string `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}

func (x *PostPeerStatusRequest) Reset() {
	*x = PostPeerStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_application_proto_syncer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostPeerStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostPeerStatusRequest) ProtoMessage() {}

func (x *PostPeerStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_application_proto_syncer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostPeerStatusRequest.ProtoReflect.Descriptor instead.
func (*PostPeerStatusRequest) Descriptor() ([]byte, []int) {
	return file_core_application_proto_syncer_proto_rawDescGZIP(), []int{1}
}

func (x *PostPeerStatusRequest) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

// Data contains app data as []byte
type Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// map[string]*[]byte
	Data map[string][]byte `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Data) Reset() {
	*x = Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_application_proto_syncer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Data) ProtoMessage() {}

func (x *Data) ProtoReflect() protoreflect.Message {
	mi := &file_core_application_proto_syncer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Data.ProtoReflect.Descriptor instead.
func (*Data) Descriptor() ([]byte, []int) {
	return file_core_application_proto_syncer_proto_rawDescGZIP(), []int{2}
}

func (x *Data) GetData() map[string][]byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_application_proto_syncer_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_core_application_proto_syncer_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_core_application_proto_syncer_proto_rawDescGZIP(), []int{3}
}

func (x *Result) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

// AppStatus contains app peer status
type AppStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// app name
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// app startup time
	StartupTime uint64 `protobuf:"varint,2,opt,name=startup_time,json=startupTime,proto3" json:"startup_time,omitempty"`
	// app uptime
	Uptime uint64 `protobuf:"varint,3,opt,name=uptime,proto3" json:"uptime,omitempty"`
	// amount of slots currently occupying the app
	GuageHeight uint64 `protobuf:"varint,4,opt,name=guage_height,json=guageHeight,proto3" json:"guage_height,omitempty"`
	// max limit
	GuageMax uint64 `protobuf:"varint,5,opt,name=guage_max,json=guageMax,proto3" json:"guage_max,omitempty"`
	// relay addr string
	Relay string `protobuf:"bytes,6,opt,name=relay,proto3" json:"relay,omitempty"`
	// PeerID
	NodeId string `protobuf:"bytes,7,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	// Addr
	Addr string `protobuf:"bytes,8,opt,name=addr,proto3" json:"addr,omitempty"`
	// App origin name
	AppOrigin string `protobuf:"bytes,9,opt,name=app_origin,json=appOrigin,proto3" json:"app_origin,omitempty"`
	// ai model hash string
	ModelHash string `protobuf:"bytes,10,opt,name=model_hash,json=modelHash,proto3" json:"model_hash,omitempty"`
	// mac addr
	Mac string `protobuf:"bytes,11,opt,name=mac,proto3" json:"mac,omitempty"`
	// memory info
	MemInfo string `protobuf:"bytes,12,opt,name=memInfo,proto3" json:"memInfo,omitempty"`
	// cpu info
	CpuInfo string `protobuf:"bytes,13,opt,name=cpu_info,json=cpuInfo,proto3" json:"cpu_info,omitempty"`
	// average e-power
	AveragePower float32 `protobuf:"fixed32,14,opt,name=average_power,json=averagePower,proto3" json:"average_power,omitempty"`
	// gpu info
	GpuInfo string `protobuf:"bytes,15,opt,name=gpu_info,json=gpuInfo,proto3" json:"gpu_info,omitempty"`
	// version
	Version string `protobuf:"bytes,16,opt,name=version,proto3" json:"version,omitempty"`
	// relay proxy port
	RelayProxyPort uint64 `protobuf:"varint,17,opt,name=relay_proxy_port,json=relayProxyPort,proto3" json:"relay_proxy_port,omitempty"`
}

func (x *AppStatus) Reset() {
	*x = AppStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_application_proto_syncer_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AppStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppStatus) ProtoMessage() {}

func (x *AppStatus) ProtoReflect() protoreflect.Message {
	mi := &file_core_application_proto_syncer_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppStatus.ProtoReflect.Descriptor instead.
func (*AppStatus) Descriptor() ([]byte, []int) {
	return file_core_application_proto_syncer_proto_rawDescGZIP(), []int{4}
}

func (x *AppStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AppStatus) GetStartupTime() uint64 {
	if x != nil {
		return x.StartupTime
	}
	return 0
}

func (x *AppStatus) GetUptime() uint64 {
	if x != nil {
		return x.Uptime
	}
	return 0
}

func (x *AppStatus) GetGuageHeight() uint64 {
	if x != nil {
		return x.GuageHeight
	}
	return 0
}

func (x *AppStatus) GetGuageMax() uint64 {
	if x != nil {
		return x.GuageMax
	}
	return 0
}

func (x *AppStatus) GetRelay() string {
	if x != nil {
		return x.Relay
	}
	return ""
}

func (x *AppStatus) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *AppStatus) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *AppStatus) GetAppOrigin() string {
	if x != nil {
		return x.AppOrigin
	}
	return ""
}

func (x *AppStatus) GetModelHash() string {
	if x != nil {
		return x.ModelHash
	}
	return ""
}

func (x *AppStatus) GetMac() string {
	if x != nil {
		return x.Mac
	}
	return ""
}

func (x *AppStatus) GetMemInfo() string {
	if x != nil {
		return x.MemInfo
	}
	return ""
}

func (x *AppStatus) GetCpuInfo() string {
	if x != nil {
		return x.CpuInfo
	}
	return ""
}

func (x *AppStatus) GetAveragePower() float32 {
	if x != nil {
		return x.AveragePower
	}
	return 0
}

func (x *AppStatus) GetGpuInfo() string {
	if x != nil {
		return x.GpuInfo
	}
	return ""
}

func (x *AppStatus) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *AppStatus) GetRelayProxyPort() uint64 {
	if x != nil {
		return x.RelayProxyPort
	}
	return 0
}

var File_core_application_proto_syncer_proto protoreflect.FileDescriptor

var file_core_application_proto_syncer_proto_rawDesc = []byte{
	0x0a, 0x23, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x65, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61,
	0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61,
	0x48, 0x61, 0x73, 0x68, 0x22, 0x30, 0x0a, 0x15, 0x50, 0x6f, 0x73, 0x74, 0x50, 0x65, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x67, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x26,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x1c, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xe6, 0x03,
	0x0a, 0x09, 0x41, 0x70, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x75, 0x70, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0b, 0x67, 0x75, 0x61, 0x67, 0x65, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x67, 0x75, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x61, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x08, 0x67, 0x75, 0x61, 0x67, 0x65, 0x4d, 0x61, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65,
	0x6c, 0x61, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x6c, 0x61, 0x79,
	0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x1d, 0x0a,
	0x0a, 0x61, 0x70, 0x70, 0x5f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x61, 0x70, 0x70, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x1d, 0x0a, 0x0a,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6d,
	0x61, 0x63, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x61, 0x63, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x70, 0x75, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x6f,
	0x77, 0x65, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x61, 0x76, 0x65, 0x72, 0x61,
	0x67, 0x65, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x70, 0x75, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x70, 0x75, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x10,
	0x72, 0x65, 0x6c, 0x61, 0x79, 0x5f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x5f, 0x70, 0x6f, 0x72, 0x74,
	0x18, 0x11, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x50, 0x72, 0x6f,
	0x78, 0x79, 0x50, 0x6f, 0x72, 0x74, 0x32, 0xa2, 0x01, 0x0a, 0x07, 0x53, 0x79, 0x6e, 0x63, 0x41,
	0x70, 0x70, 0x12, 0x38, 0x0a, 0x0d, 0x50, 0x6f, 0x73, 0x74, 0x41, 0x70, 0x70, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x19, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x50, 0x65, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x30, 0x01, 0x12, 0x29, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x08, 0x2e, 0x76, 0x31,
	0x2e, 0x44, 0x61, 0x74, 0x61, 0x30, 0x01, 0x12, 0x32, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0d, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x70, 0x70, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x14, 0x5a, 0x12, 0x2f,
	0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_application_proto_syncer_proto_rawDescOnce sync.Once
	file_core_application_proto_syncer_proto_rawDescData = file_core_application_proto_syncer_proto_rawDesc
)

func file_core_application_proto_syncer_proto_rawDescGZIP() []byte {
	file_core_application_proto_syncer_proto_rawDescOnce.Do(func() {
		file_core_application_proto_syncer_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_application_proto_syncer_proto_rawDescData)
	})
	return file_core_application_proto_syncer_proto_rawDescData
}

var file_core_application_proto_syncer_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_core_application_proto_syncer_proto_goTypes = []interface{}{
	(*GetDataRequest)(nil),        // 0: v1.GetDataRequest
	(*PostPeerStatusRequest)(nil), // 1: v1.PostPeerStatusRequest
	(*Data)(nil),                  // 2: v1.Data
	(*Result)(nil),                // 3: v1.Result
	(*AppStatus)(nil),             // 4: v1.AppStatus
	nil,                           // 5: v1.Data.DataEntry
	(*emptypb.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_core_application_proto_syncer_proto_depIdxs = []int32{
	5, // 0: v1.Data.data:type_name -> v1.Data.DataEntry
	1, // 1: v1.SyncApp.PostAppStatus:input_type -> v1.PostPeerStatusRequest
	0, // 2: v1.SyncApp.GetData:input_type -> v1.GetDataRequest
	6, // 3: v1.SyncApp.GetStatus:input_type -> google.protobuf.Empty
	3, // 4: v1.SyncApp.PostAppStatus:output_type -> v1.Result
	2, // 5: v1.SyncApp.GetData:output_type -> v1.Data
	4, // 6: v1.SyncApp.GetStatus:output_type -> v1.AppStatus
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_core_application_proto_syncer_proto_init() }
func file_core_application_proto_syncer_proto_init() {
	if File_core_application_proto_syncer_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_application_proto_syncer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDataRequest); i {
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
		file_core_application_proto_syncer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostPeerStatusRequest); i {
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
		file_core_application_proto_syncer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Data); i {
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
		file_core_application_proto_syncer_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
		file_core_application_proto_syncer_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AppStatus); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_application_proto_syncer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_core_application_proto_syncer_proto_goTypes,
		DependencyIndexes: file_core_application_proto_syncer_proto_depIdxs,
		MessageInfos:      file_core_application_proto_syncer_proto_msgTypes,
	}.Build()
	File_core_application_proto_syncer_proto = out.File
	file_core_application_proto_syncer_proto_rawDesc = nil
	file_core_application_proto_syncer_proto_goTypes = nil
	file_core_application_proto_syncer_proto_depIdxs = nil
}
