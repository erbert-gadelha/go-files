// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: grpcarquivo.proto

package grpcarquivo

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Content       string                 `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Request) Reset() {
	*x = Request{}
	mi := &file_grpcarquivo_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_grpcarquivo_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_grpcarquivo_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Lines         int32                  `protobuf:"varint,1,opt,name=lines,proto3" json:"lines,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file_grpcarquivo_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_grpcarquivo_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_grpcarquivo_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetLines() int32 {
	if x != nil {
		return x.Lines
	}
	return 0
}

var File_grpcarquivo_proto protoreflect.FileDescriptor

var file_grpcarquivo_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x67, 0x72, 0x70, 0x63, 0x61, 0x72, 0x71, 0x75, 0x69, 0x76, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x67, 0x72, 0x70, 0x63, 0x61, 0x72, 0x71, 0x75, 0x69, 0x76, 0x6f,
	0x22, 0x23, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x20, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x6c, 0x69, 0x6e, 0x65, 0x73, 0x32, 0x4b, 0x0a, 0x0e, 0x41, 0x72, 0x71, 0x75, 0x69,
	0x76, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x4c, 0x69, 0x6e, 0x65, 0x73, 0x12, 0x14, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x61, 0x72,
	0x71, 0x75, 0x69, 0x76, 0x6f, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x61, 0x72, 0x71, 0x75, 0x69, 0x76, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x61, 0x72,
	0x71, 0x75, 0x69, 0x76, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_grpcarquivo_proto_rawDescOnce sync.Once
	file_grpcarquivo_proto_rawDescData []byte
)

func file_grpcarquivo_proto_rawDescGZIP() []byte {
	file_grpcarquivo_proto_rawDescOnce.Do(func() {
		file_grpcarquivo_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_grpcarquivo_proto_rawDesc), len(file_grpcarquivo_proto_rawDesc)))
	})
	return file_grpcarquivo_proto_rawDescData
}

var file_grpcarquivo_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_grpcarquivo_proto_goTypes = []any{
	(*Request)(nil),  // 0: grpcarquivo.Request
	(*Response)(nil), // 1: grpcarquivo.Response
}
var file_grpcarquivo_proto_depIdxs = []int32{
	0, // 0: grpcarquivo.ArquivoService.CountLines:input_type -> grpcarquivo.Request
	1, // 1: grpcarquivo.ArquivoService.CountLines:output_type -> grpcarquivo.Response
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_grpcarquivo_proto_init() }
func file_grpcarquivo_proto_init() {
	if File_grpcarquivo_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_grpcarquivo_proto_rawDesc), len(file_grpcarquivo_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpcarquivo_proto_goTypes,
		DependencyIndexes: file_grpcarquivo_proto_depIdxs,
		MessageInfos:      file_grpcarquivo_proto_msgTypes,
	}.Build()
	File_grpcarquivo_proto = out.File
	file_grpcarquivo_proto_goTypes = nil
	file_grpcarquivo_proto_depIdxs = nil
}
