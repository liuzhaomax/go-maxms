// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
//     protoc-gen-go v1.26.0
//     protoc        v4.25.1
// source: user.proto

package pb

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

type UserIDReq struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    UserID string `protobuf:"bytes,1,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *UserIDReq) Reset() {
    *x = UserIDReq{}
    if protoimpl.UnsafeEnabled {
        mi := &file_user_proto_msgTypes[0]
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        ms.StoreMessageInfo(mi)
    }
}

func (x *UserIDReq) String() string {
    return protoimpl.X.MessageStringOf(x)
}

func (*UserIDReq) ProtoMessage() {}

func (x *UserIDReq) ProtoReflect() protoreflect.Message {
    mi := &file_user_proto_msgTypes[0]
    if protoimpl.UnsafeEnabled && x != nil {
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        if ms.LoadMessageInfo() == nil {
            ms.StoreMessageInfo(mi)
        }
        return ms
    }
    return mi.MessageOf(x)
}

// Deprecated: Use UserIDReq.ProtoReflect.Descriptor instead.
func (*UserIDReq) Descriptor() ([]byte, []int) {
    return file_user_proto_rawDescGZIP(), []int{0}
}

func (x *UserIDReq) GetUserID() string {
    if x != nil {
        return x.UserID
    }
    return ""
}

type Status struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Code int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
    Desc string `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *Status) Reset() {
    *x = Status{}
    if protoimpl.UnsafeEnabled {
        mi := &file_user_proto_msgTypes[1]
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        ms.StoreMessageInfo(mi)
    }
}

func (x *Status) String() string {
    return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
    mi := &file_user_proto_msgTypes[1]
    if protoimpl.UnsafeEnabled && x != nil {
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        if ms.LoadMessageInfo() == nil {
            ms.StoreMessageInfo(mi)
        }
        return ms
    }
    return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
    return file_user_proto_rawDescGZIP(), []int{1}
}

func (x *Status) GetCode() int32 {
    if x != nil {
        return x.Code
    }
    return 0
}

func (x *Status) GetDesc() string {
    if x != nil {
        return x.Desc
    }
    return ""
}

type UserResData struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Id     int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
    UserID string `protobuf:"bytes,2,opt,name=userID,proto3" json:"userID,omitempty"`
    Mobile string `protobuf:"bytes,3,opt,name=mobile,proto3" json:"mobile,omitempty"`
}

func (x *UserResData) Reset() {
    *x = UserResData{}
    if protoimpl.UnsafeEnabled {
        mi := &file_user_proto_msgTypes[2]
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        ms.StoreMessageInfo(mi)
    }
}

func (x *UserResData) String() string {
    return protoimpl.X.MessageStringOf(x)
}

func (*UserResData) ProtoMessage() {}

func (x *UserResData) ProtoReflect() protoreflect.Message {
    mi := &file_user_proto_msgTypes[2]
    if protoimpl.UnsafeEnabled && x != nil {
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        if ms.LoadMessageInfo() == nil {
            ms.StoreMessageInfo(mi)
        }
        return ms
    }
    return mi.MessageOf(x)
}

// Deprecated: Use UserResData.ProtoReflect.Descriptor instead.
func (*UserResData) Descriptor() ([]byte, []int) {
    return file_user_proto_rawDescGZIP(), []int{2}
}

func (x *UserResData) GetId() int32 {
    if x != nil {
        return x.Id
    }
    return 0
}

func (x *UserResData) GetUserID() string {
    if x != nil {
        return x.UserID
    }
    return ""
}

func (x *UserResData) GetMobile() string {
    if x != nil {
        return x.Mobile
    }
    return ""
}

type UserRes struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Status *Status      `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
    Data   *UserResData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *UserRes) Reset() {
    *x = UserRes{}
    if protoimpl.UnsafeEnabled {
        mi := &file_user_proto_msgTypes[3]
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        ms.StoreMessageInfo(mi)
    }
}

func (x *UserRes) String() string {
    return protoimpl.X.MessageStringOf(x)
}

func (*UserRes) ProtoMessage() {}

func (x *UserRes) ProtoReflect() protoreflect.Message {
    mi := &file_user_proto_msgTypes[3]
    if protoimpl.UnsafeEnabled && x != nil {
        ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
        if ms.LoadMessageInfo() == nil {
            ms.StoreMessageInfo(mi)
        }
        return ms
    }
    return mi.MessageOf(x)
}

// Deprecated: Use UserRes.ProtoReflect.Descriptor instead.
func (*UserRes) Descriptor() ([]byte, []int) {
    return file_user_proto_rawDescGZIP(), []int{3}
}

func (x *UserRes) GetStatus() *Status {
    if x != nil {
        return x.Status
    }
    return nil
}

func (x *UserRes) GetData() *UserResData {
    if x != nil {
        return x.Data
    }
    return nil
}

var File_user_proto protoreflect.FileDescriptor

var file_user_proto_rawDesc = []byte{
    0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x23, 0x0a, 0x09,
    0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65,
    0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
    0x44, 0x22, 0x30, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63,
    0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
    0x12, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64,
    0x65, 0x73, 0x63, 0x22, 0x4d, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x44, 0x61,
    0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02,
    0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
    0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x6f,
    0x62, 0x69, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x6f, 0x62, 0x69,
    0x6c, 0x65, 0x22, 0x4c, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x12, 0x1f, 0x0a,
    0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e,
    0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x20,
    0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x55,
    0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
    0x32, 0x36, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
    0x27, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
    0x49, 0x44, 0x12, 0x0a, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x52, 0x65, 0x71, 0x1a, 0x08,
    0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62,
    0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
    file_user_proto_rawDescOnce sync.Once
    file_user_proto_rawDescData = file_user_proto_rawDesc
)

func file_user_proto_rawDescGZIP() []byte {
    file_user_proto_rawDescOnce.Do(func() {
        file_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_proto_rawDescData)
    })
    return file_user_proto_rawDescData
}

var file_user_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_user_proto_goTypes = []interface{}{
    (*UserIDReq)(nil),   // 0: UserIDReq
    (*Status)(nil),      // 1: Status
    (*UserResData)(nil), // 2: UserResData
    (*UserRes)(nil),     // 3: UserRes
}
var file_user_proto_depIdxs = []int32{
    1, // 0: UserRes.status:type_name -> Status
    2, // 1: UserRes.data:type_name -> UserResData
    0, // 2: UserService.GetUserByUserID:input_type -> UserIDReq
    3, // 3: UserService.GetUserByUserID:output_type -> UserRes
    3, // [3:4] is the sub-list for method output_type
    2, // [2:3] is the sub-list for method input_type
    2, // [2:2] is the sub-list for extension type_name
    2, // [2:2] is the sub-list for extension extendee
    0, // [0:2] is the sub-list for field type_name
}

func init() { file_user_proto_init() }
func file_user_proto_init() {
    if File_user_proto != nil {
        return
    }
    if !protoimpl.UnsafeEnabled {
        file_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
            switch v := v.(*UserIDReq); i {
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
        file_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
            switch v := v.(*Status); i {
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
        file_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
            switch v := v.(*UserResData); i {
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
        file_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
            switch v := v.(*UserRes); i {
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
            RawDescriptor: file_user_proto_rawDesc,
            NumEnums:      0,
            NumMessages:   4,
            NumExtensions: 0,
            NumServices:   1,
        },
        GoTypes:           file_user_proto_goTypes,
        DependencyIndexes: file_user_proto_depIdxs,
        MessageInfos:      file_user_proto_msgTypes,
    }.Build()
    File_user_proto = out.File
    file_user_proto_rawDesc = nil
    file_user_proto_goTypes = nil
    file_user_proto_depIdxs = nil
}
