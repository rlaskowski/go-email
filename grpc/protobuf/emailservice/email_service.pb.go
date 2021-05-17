// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: grpc/protobuf/emailservice/email_service.proto

package emailservice

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

type IncomingMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Address  *Address   `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	Subject  string     `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	Date     string     `protobuf:"bytes,4,opt,name=date,proto3" json:"date,omitempty"`
	Contents []*Content `protobuf:"bytes,5,rep,name=contents,proto3" json:"contents,omitempty"`
	Files    []*File    `protobuf:"bytes,6,rep,name=files,proto3" json:"files,omitempty"`
}

func (x *IncomingMessage) Reset() {
	*x = IncomingMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncomingMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncomingMessage) ProtoMessage() {}

func (x *IncomingMessage) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncomingMessage.ProtoReflect.Descriptor instead.
func (*IncomingMessage) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{0}
}

func (x *IncomingMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *IncomingMessage) GetAddress() *Address {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *IncomingMessage) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *IncomingMessage) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *IncomingMessage) GetContents() []*Content {
	if x != nil {
		return x.Contents
	}
	return nil
}

func (x *IncomingMessage) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

type Stat struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key           string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	MessageId     int64  `protobuf:"varint,2,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	MessageNumber int64  `protobuf:"varint,3,opt,name=message_number,json=messageNumber,proto3" json:"message_number,omitempty"`
}

func (x *Stat) Reset() {
	*x = Stat{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stat) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stat) ProtoMessage() {}

func (x *Stat) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stat.ProtoReflect.Descriptor instead.
func (*Stat) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{1}
}

func (x *Stat) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Stat) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *Stat) GetMessageNumber() int64 {
	if x != nil {
		return x.MessageNumber
	}
	return 0
}

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{2}
}

func (x *Address) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Address) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type Content struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HtmlType bool   `protobuf:"varint,1,opt,name=HtmlType,proto3" json:"HtmlType,omitempty"`
	Data     []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Content) Reset() {
	*x = Content{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Content) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Content) ProtoMessage() {}

func (x *Content) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Content.ProtoReflect.Descriptor instead.
func (*Content) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{3}
}

func (x *Content) GetHtmlType() bool {
	if x != nil {
		return x.HtmlType
	}
	return false
}

func (x *Content) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type File struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *File) Reset() {
	*x = File{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *File) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*File) ProtoMessage() {}

func (x *File) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use File.ProtoReflect.Descriptor instead.
func (*File) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{4}
}

func (x *File) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *File) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type StatRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *StatRequest) Reset() {
	*x = StatRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatRequest) ProtoMessage() {}

func (x *StatRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatRequest.ProtoReflect.Descriptor instead.
func (*StatRequest) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{5}
}

func (x *StatRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type IncomingMsgRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key           string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	MessageNumber int64  `protobuf:"varint,2,opt,name=message_number,json=messageNumber,proto3" json:"message_number,omitempty"`
}

func (x *IncomingMsgRequest) Reset() {
	*x = IncomingMsgRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncomingMsgRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncomingMsgRequest) ProtoMessage() {}

func (x *IncomingMsgRequest) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncomingMsgRequest.ProtoReflect.Descriptor instead.
func (*IncomingMsgRequest) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{6}
}

func (x *IncomingMsgRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *IncomingMsgRequest) GetMessageNumber() int64 {
	if x != nil {
		return x.MessageNumber
	}
	return 0
}

type IncomingMsgResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message []byte `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *IncomingMsgResponse) Reset() {
	*x = IncomingMsgResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncomingMsgResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncomingMsgResponse) ProtoMessage() {}

func (x *IncomingMsgResponse) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_protobuf_emailservice_email_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncomingMsgResponse.ProtoReflect.Descriptor instead.
func (*IncomingMsgResponse) Descriptor() ([]byte, []int) {
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP(), []int{7}
}

func (x *IncomingMsgResponse) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

var File_grpc_protobuf_emailservice_email_service_proto protoreflect.FileDescriptor

var file_grpc_protobuf_emailservice_email_service_proto_rawDesc = []byte{
	0x0a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0xdd,
	0x01, 0x0a, 0x0f, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x2f, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x65, 0x12, 0x31, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x06, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x5e,
	0x0a, 0x04, 0x53, 0x74, 0x61, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x37,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x39, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x48, 0x74, 0x6d, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x48, 0x74, 0x6d, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x2e, 0x0a, 0x04, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x1f, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x4d, 0x0a, 0x12, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x4d,
	0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x22, 0x2f, 0x0a, 0x13, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x4d, 0x73,
	0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x32, 0xf8, 0x01, 0x0a, 0x0c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x12, 0x19, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12,
	0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x22, 0x00, 0x30, 0x01, 0x12, 0x59, 0x0a, 0x0e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x20, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67,
	0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x65, 0x6d, 0x61,
	0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69,
	0x6e, 0x67, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30,
	0x01, 0x12, 0x4b, 0x0a, 0x0c, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x12, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x1a, 0x21, 0x2e, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x4d, 0x73, 0x67,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x1c,
	0x5a, 0x1a, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x61, 0x69, 0x6c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_protobuf_emailservice_email_service_proto_rawDescOnce sync.Once
	file_grpc_protobuf_emailservice_email_service_proto_rawDescData = file_grpc_protobuf_emailservice_email_service_proto_rawDesc
)

func file_grpc_protobuf_emailservice_email_service_proto_rawDescGZIP() []byte {
	file_grpc_protobuf_emailservice_email_service_proto_rawDescOnce.Do(func() {
		file_grpc_protobuf_emailservice_email_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_protobuf_emailservice_email_service_proto_rawDescData)
	})
	return file_grpc_protobuf_emailservice_email_service_proto_rawDescData
}

var file_grpc_protobuf_emailservice_email_service_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_grpc_protobuf_emailservice_email_service_proto_goTypes = []interface{}{
	(*IncomingMessage)(nil),     // 0: emailservice.IncomingMessage
	(*Stat)(nil),                // 1: emailservice.Stat
	(*Address)(nil),             // 2: emailservice.Address
	(*Content)(nil),             // 3: emailservice.Content
	(*File)(nil),                // 4: emailservice.File
	(*StatRequest)(nil),         // 5: emailservice.StatRequest
	(*IncomingMsgRequest)(nil),  // 6: emailservice.IncomingMsgRequest
	(*IncomingMsgResponse)(nil), // 7: emailservice.IncomingMsgResponse
}
var file_grpc_protobuf_emailservice_email_service_proto_depIdxs = []int32{
	2, // 0: emailservice.IncomingMessage.address:type_name -> emailservice.Address
	3, // 1: emailservice.IncomingMessage.contents:type_name -> emailservice.Content
	4, // 2: emailservice.IncomingMessage.files:type_name -> emailservice.File
	5, // 3: emailservice.EmailService.MessageStat:input_type -> emailservice.StatRequest
	6, // 4: emailservice.EmailService.ReceiveMessage:input_type -> emailservice.IncomingMsgRequest
	1, // 5: emailservice.EmailService.RouteMessage:input_type -> emailservice.Stat
	1, // 6: emailservice.EmailService.MessageStat:output_type -> emailservice.Stat
	7, // 7: emailservice.EmailService.ReceiveMessage:output_type -> emailservice.IncomingMsgResponse
	7, // 8: emailservice.EmailService.RouteMessage:output_type -> emailservice.IncomingMsgResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_grpc_protobuf_emailservice_email_service_proto_init() }
func file_grpc_protobuf_emailservice_email_service_proto_init() {
	if File_grpc_protobuf_emailservice_email_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IncomingMessage); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stat); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Content); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*File); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatRequest); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IncomingMsgRequest); i {
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
		file_grpc_protobuf_emailservice_email_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IncomingMsgResponse); i {
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
			RawDescriptor: file_grpc_protobuf_emailservice_email_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_grpc_protobuf_emailservice_email_service_proto_goTypes,
		DependencyIndexes: file_grpc_protobuf_emailservice_email_service_proto_depIdxs,
		MessageInfos:      file_grpc_protobuf_emailservice_email_service_proto_msgTypes,
	}.Build()
	File_grpc_protobuf_emailservice_email_service_proto = out.File
	file_grpc_protobuf_emailservice_email_service_proto_rawDesc = nil
	file_grpc_protobuf_emailservice_email_service_proto_goTypes = nil
	file_grpc_protobuf_emailservice_email_service_proto_depIdxs = nil
}