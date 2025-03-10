// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: core/relay/proto/alive.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AliveStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

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
	// app origin name string
	AppOrigin string `protobuf:"bytes,7,opt,name=app_origin,json=appOrigin,proto3" json:"app_origin,omitempty"`
	// ai model hash string
	ModelHash string `protobuf:"bytes,8,opt,name=model_hash,json=modelHash,proto3" json:"model_hash,omitempty"`
	// mac addr
	Mac string `protobuf:"bytes,9,opt,name=mac,proto3" json:"mac,omitempty"`
	// memory info
	MemInfo string `protobuf:"bytes,10,opt,name=memInfo,proto3" json:"memInfo,omitempty"`
	// cpu info
	CpuInfo string `protobuf:"bytes,11,opt,name=cpu_info,json=cpuInfo,proto3" json:"cpu_info,omitempty"`
	// average e power value
	AveragePower float32 `protobuf:"fixed32,12,opt,name=average_power,json=averagePower,proto3" json:"average_power,omitempty"`
	// gpu info
	GpuInfo string `protobuf:"bytes,13,opt,name=gpu_info,json=gpuInfo,proto3" json:"gpu_info,omitempty"`
	// version
	Version string `protobuf:"bytes,14,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *AliveStatus) Reset() {
	*x = AliveStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_relay_proto_alive_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AliveStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AliveStatus) ProtoMessage() {}

func (x *AliveStatus) ProtoReflect() protoreflect.Message {
	mi := &file_core_relay_proto_alive_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AliveStatus.ProtoReflect.Descriptor instead.
func (*AliveStatus) Descriptor() ([]byte, []int) {
	return file_core_relay_proto_alive_proto_rawDescGZIP(), []int{0}
}

func (x *AliveStatus) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AliveStatus) GetStartupTime() uint64 {
	if x != nil {
		return x.StartupTime
	}
	return 0
}

func (x *AliveStatus) GetUptime() uint64 {
	if x != nil {
		return x.Uptime
	}
	return 0
}

func (x *AliveStatus) GetGuageHeight() uint64 {
	if x != nil {
		return x.GuageHeight
	}
	return 0
}

func (x *AliveStatus) GetGuageMax() uint64 {
	if x != nil {
		return x.GuageMax
	}
	return 0
}

func (x *AliveStatus) GetRelay() string {
	if x != nil {
		return x.Relay
	}
	return ""
}

func (x *AliveStatus) GetAppOrigin() string {
	if x != nil {
		return x.AppOrigin
	}
	return ""
}

func (x *AliveStatus) GetModelHash() string {
	if x != nil {
		return x.ModelHash
	}
	return ""
}

func (x *AliveStatus) GetMac() string {
	if x != nil {
		return x.Mac
	}
	return ""
}

func (x *AliveStatus) GetMemInfo() string {
	if x != nil {
		return x.MemInfo
	}
	return ""
}

func (x *AliveStatus) GetCpuInfo() string {
	if x != nil {
		return x.CpuInfo
	}
	return ""
}

func (x *AliveStatus) GetAveragePower() float32 {
	if x != nil {
		return x.AveragePower
	}
	return 0
}

func (x *AliveStatus) GetGpuInfo() string {
	if x != nil {
		return x.GpuInfo
	}
	return ""
}

func (x *AliveStatus) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type AliveStatusResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success   bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Discovery string `protobuf:"bytes,2,opt,name=discovery,proto3" json:"discovery,omitempty"`
}

func (x *AliveStatusResp) Reset() {
	*x = AliveStatusResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_relay_proto_alive_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AliveStatusResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AliveStatusResp) ProtoMessage() {}

func (x *AliveStatusResp) ProtoReflect() protoreflect.Message {
	mi := &file_core_relay_proto_alive_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AliveStatusResp.ProtoReflect.Descriptor instead.
func (*AliveStatusResp) Descriptor() ([]byte, []int) {
	return file_core_relay_proto_alive_proto_rawDescGZIP(), []int{1}
}

func (x *AliveStatusResp) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AliveStatusResp) GetDiscovery() string {
	if x != nil {
		return x.Discovery
	}
	return ""
}

var File_core_relay_proto_alive_proto protoreflect.FileDescriptor

var file_core_relay_proto_alive_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x61, 0x6c, 0x69, 0x76, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x76, 0x31, 0x22, 0x91, 0x03, 0x0a, 0x0b, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x75,
	0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x75, 0x70, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x70, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x67, 0x75, 0x61, 0x67, 0x65, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x67, 0x75, 0x61, 0x67, 0x65, 0x48, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x75, 0x61, 0x67, 0x65, 0x5f, 0x6d, 0x61,
	0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x67, 0x75, 0x61, 0x67, 0x65, 0x4d, 0x61,
	0x78, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x70, 0x70, 0x5f, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x70, 0x70,
	0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x48, 0x61, 0x73, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x63, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6d, 0x61, 0x63, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x49, 0x6e,
	0x66, 0x6f, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x70, 0x75, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x23, 0x0a, 0x0d,
	0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x0c, 0x61, 0x76, 0x65, 0x72, 0x61, 0x67, 0x65, 0x50, 0x6f, 0x77, 0x65,
	0x72, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x70, 0x75, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x67, 0x70, 0x75, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x49, 0x0a, 0x0f, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72,
	0x79, 0x32, 0x36, 0x0a, 0x05, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x2d, 0x0a, 0x05, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x12, 0x0f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x1a, 0x13, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x42, 0x0e, 0x5a, 0x0c, 0x2f, 0x72, 0x65,
	0x6c, 0x61, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_core_relay_proto_alive_proto_rawDescOnce sync.Once
	file_core_relay_proto_alive_proto_rawDescData = file_core_relay_proto_alive_proto_rawDesc
)

func file_core_relay_proto_alive_proto_rawDescGZIP() []byte {
	file_core_relay_proto_alive_proto_rawDescOnce.Do(func() {
		file_core_relay_proto_alive_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_relay_proto_alive_proto_rawDescData)
	})
	return file_core_relay_proto_alive_proto_rawDescData
}

var file_core_relay_proto_alive_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_core_relay_proto_alive_proto_goTypes = []interface{}{
	(*AliveStatus)(nil),     // 0: v1.AliveStatus
	(*AliveStatusResp)(nil), // 1: v1.AliveStatusResp
}
var file_core_relay_proto_alive_proto_depIdxs = []int32{
	0, // 0: v1.Alive.Hello:input_type -> v1.AliveStatus
	1, // 1: v1.Alive.Hello:output_type -> v1.AliveStatusResp
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_core_relay_proto_alive_proto_init() }
func file_core_relay_proto_alive_proto_init() {
	if File_core_relay_proto_alive_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_relay_proto_alive_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AliveStatus); i {
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
		file_core_relay_proto_alive_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AliveStatusResp); i {
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
			RawDescriptor: file_core_relay_proto_alive_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_core_relay_proto_alive_proto_goTypes,
		DependencyIndexes: file_core_relay_proto_alive_proto_depIdxs,
		MessageInfos:      file_core_relay_proto_alive_proto_msgTypes,
	}.Build()
	File_core_relay_proto_alive_proto = out.File
	file_core_relay_proto_alive_proto_rawDesc = nil
	file_core_relay_proto_alive_proto_goTypes = nil
	file_core_relay_proto_alive_proto_depIdxs = nil
}
