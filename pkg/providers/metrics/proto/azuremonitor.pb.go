// Copyright (c) technicianted. All rights reserved.
// Licensed under the MIT License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.0--rc2
// source: azuremonitor.proto

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

// Metric aggregation type as supported by the metric.
type AzureMonitorMetricConfig_Aggregation int32

const (
	AzureMonitorMetricConfig_None          AzureMonitorMetricConfig_Aggregation = 0
	AzureMonitorMetricConfig_Average       AzureMonitorMetricConfig_Aggregation = 1
	AzureMonitorMetricConfig_Maximum       AzureMonitorMetricConfig_Aggregation = 2
	AzureMonitorMetricConfig_Minimum       AzureMonitorMetricConfig_Aggregation = 3
	AzureMonitorMetricConfig_Count         AzureMonitorMetricConfig_Aggregation = 4
	AzureMonitorMetricConfig_Total         AzureMonitorMetricConfig_Aggregation = 5
	AzureMonitorMetricConfig_RatePerMinute AzureMonitorMetricConfig_Aggregation = 6
)

// Enum value maps for AzureMonitorMetricConfig_Aggregation.
var (
	AzureMonitorMetricConfig_Aggregation_name = map[int32]string{
		0: "None",
		1: "Average",
		2: "Maximum",
		3: "Minimum",
		4: "Count",
		5: "Total",
		6: "RatePerMinute",
	}
	AzureMonitorMetricConfig_Aggregation_value = map[string]int32{
		"None":          0,
		"Average":       1,
		"Maximum":       2,
		"Minimum":       3,
		"Count":         4,
		"Total":         5,
		"RatePerMinute": 6,
	}
)

func (x AzureMonitorMetricConfig_Aggregation) Enum() *AzureMonitorMetricConfig_Aggregation {
	p := new(AzureMonitorMetricConfig_Aggregation)
	*p = x
	return p
}

func (x AzureMonitorMetricConfig_Aggregation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AzureMonitorMetricConfig_Aggregation) Descriptor() protoreflect.EnumDescriptor {
	return file_azuremonitor_proto_enumTypes[0].Descriptor()
}

func (AzureMonitorMetricConfig_Aggregation) Type() protoreflect.EnumType {
	return &file_azuremonitor_proto_enumTypes[0]
}

func (x AzureMonitorMetricConfig_Aggregation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AzureMonitorMetricConfig_Aggregation.Descriptor instead.
func (AzureMonitorMetricConfig_Aggregation) EnumDescriptor() ([]byte, []int) {
	return file_azuremonitor_proto_rawDescGZIP(), []int{0, 0}
}

type AzureMonitorMetricConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Target Azure resource URI.
	ResourceURI string `protobuf:"bytes,1,opt,name=resourceURI,proto3" json:"resourceURI,omitempty"`
	// Metric Azure namespace.
	MetricNamespace string `protobuf:"bytes,2,opt,name=metricNamespace,proto3" json:"metricNamespace,omitempty"`
	// Aggeragtion type for this metric. Must be supported by the metric.
	Aggregation AzureMonitorMetricConfig_Aggregation `protobuf:"varint,3,opt,name=aggregation,proto3,enum=k9sautoscaler.providers.metrics.proto.AzureMonitorMetricConfig_Aggregation" json:"aggregation,omitempty"`
	// Filter values using expressions.
	Filter *string `protobuf:"bytes,4,opt,name=filter,proto3,oneof" json:"filter,omitempty"`
}

func (x *AzureMonitorMetricConfig) Reset() {
	*x = AzureMonitorMetricConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_azuremonitor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AzureMonitorMetricConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AzureMonitorMetricConfig) ProtoMessage() {}

func (x *AzureMonitorMetricConfig) ProtoReflect() protoreflect.Message {
	mi := &file_azuremonitor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AzureMonitorMetricConfig.ProtoReflect.Descriptor instead.
func (*AzureMonitorMetricConfig) Descriptor() ([]byte, []int) {
	return file_azuremonitor_proto_rawDescGZIP(), []int{0}
}

func (x *AzureMonitorMetricConfig) GetResourceURI() string {
	if x != nil {
		return x.ResourceURI
	}
	return ""
}

func (x *AzureMonitorMetricConfig) GetMetricNamespace() string {
	if x != nil {
		return x.MetricNamespace
	}
	return ""
}

func (x *AzureMonitorMetricConfig) GetAggregation() AzureMonitorMetricConfig_Aggregation {
	if x != nil {
		return x.Aggregation
	}
	return AzureMonitorMetricConfig_None
}

func (x *AzureMonitorMetricConfig) GetFilter() string {
	if x != nil && x.Filter != nil {
		return *x.Filter
	}
	return ""
}

// Configuration for Azure Monitor based metrics provider.
// Authentication is handled using default Azure credential mechanism.
// See: https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication
// It is important that the metrics query returns exactly 1 time series to be usable
// in autoscaling.
type AzureMonitorConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AzureMonitorConfig) Reset() {
	*x = AzureMonitorConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_azuremonitor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AzureMonitorConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AzureMonitorConfig) ProtoMessage() {}

func (x *AzureMonitorConfig) ProtoReflect() protoreflect.Message {
	mi := &file_azuremonitor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AzureMonitorConfig.ProtoReflect.Descriptor instead.
func (*AzureMonitorConfig) Descriptor() ([]byte, []int) {
	return file_azuremonitor_proto_rawDescGZIP(), []int{1}
}

var File_azuremonitor_proto protoreflect.FileDescriptor

var file_azuremonitor_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x7a, 0x75, 0x72, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x25, 0x6b, 0x39, 0x73, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61,
	0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x2e, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe6, 0x02, 0x0a, 0x18,
	0x41, 0x7a, 0x75, 0x72, 0x65, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x55, 0x52, 0x49, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x55, 0x52, 0x49, 0x12, 0x28, 0x0a, 0x0f, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x12, 0x6d, 0x0a, 0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x4b, 0x2e, 0x6b, 0x39, 0x73, 0x61,
	0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x73, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x41, 0x7a, 0x75, 0x72, 0x65, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x88, 0x01, 0x01,
	0x22, 0x67, 0x0a, 0x0b, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x76, 0x65,
	0x72, 0x61, 0x67, 0x65, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x61, 0x78, 0x69, 0x6d, 0x75,
	0x6d, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x4d, 0x69, 0x6e, 0x69, 0x6d, 0x75, 0x6d, 0x10, 0x03,
	0x12, 0x09, 0x0a, 0x05, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x54,
	0x6f, 0x74, 0x61, 0x6c, 0x10, 0x05, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x61, 0x74, 0x65, 0x50, 0x65,
	0x72, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x10, 0x06, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x22, 0x14, 0x0a, 0x12, 0x41, 0x7a, 0x75, 0x72, 0x65, 0x4d, 0x6f, 0x6e,
	0x69, 0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x32, 0x5a, 0x30, 0x6b, 0x39,
	0x73, 0x2d, 0x61, 0x75, 0x74, 0x6f, 0x73, 0x63, 0x61, 0x6c, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_azuremonitor_proto_rawDescOnce sync.Once
	file_azuremonitor_proto_rawDescData = file_azuremonitor_proto_rawDesc
)

func file_azuremonitor_proto_rawDescGZIP() []byte {
	file_azuremonitor_proto_rawDescOnce.Do(func() {
		file_azuremonitor_proto_rawDescData = protoimpl.X.CompressGZIP(file_azuremonitor_proto_rawDescData)
	})
	return file_azuremonitor_proto_rawDescData
}

var file_azuremonitor_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_azuremonitor_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_azuremonitor_proto_goTypes = []interface{}{
	(AzureMonitorMetricConfig_Aggregation)(0), // 0: k9sautoscaler.providers.metrics.proto.AzureMonitorMetricConfig.Aggregation
	(*AzureMonitorMetricConfig)(nil),          // 1: k9sautoscaler.providers.metrics.proto.AzureMonitorMetricConfig
	(*AzureMonitorConfig)(nil),                // 2: k9sautoscaler.providers.metrics.proto.AzureMonitorConfig
}
var file_azuremonitor_proto_depIdxs = []int32{
	0, // 0: k9sautoscaler.providers.metrics.proto.AzureMonitorMetricConfig.aggregation:type_name -> k9sautoscaler.providers.metrics.proto.AzureMonitorMetricConfig.Aggregation
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_azuremonitor_proto_init() }
func file_azuremonitor_proto_init() {
	if File_azuremonitor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_azuremonitor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AzureMonitorMetricConfig); i {
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
		file_azuremonitor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AzureMonitorConfig); i {
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
	file_azuremonitor_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_azuremonitor_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_azuremonitor_proto_goTypes,
		DependencyIndexes: file_azuremonitor_proto_depIdxs,
		EnumInfos:         file_azuremonitor_proto_enumTypes,
		MessageInfos:      file_azuremonitor_proto_msgTypes,
	}.Build()
	File_azuremonitor_proto = out.File
	file_azuremonitor_proto_rawDesc = nil
	file_azuremonitor_proto_goTypes = nil
	file_azuremonitor_proto_depIdxs = nil
}
